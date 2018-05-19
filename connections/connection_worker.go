package connections

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/player"
)

const (
	chanOutputMessageBuffer      = 128
	chanReadMessagesBuffer       = 128
	chanDecodeMessageBuffer      = 128
	chanEncodeMessageBuffer      = 128
	chanProxyInputMessageBuffer  = 64
	chanListenInputMessageBuffer = 32

	sendInputMessageTimeout = time.Millisecond * 50

	sendCloseConnectionTimeout = time.Second

	readMessageLimit = 1024
)

type ConnectionWorker struct {
	conn   *websocket.Conn
	logger logrus.FieldLogger

	chStop      <-chan struct{}
	chsInput    []chan InputMessage
	chsInputMux *sync.RWMutex

	flagStarted bool
}

func NewConnectionWorker(conn *websocket.Conn, logger logrus.FieldLogger) *ConnectionWorker {
	return &ConnectionWorker{
		conn:        conn,
		logger:      logger,
		chsInput:    make([]chan InputMessage, 0),
		chsInputMux: &sync.RWMutex{},
	}
}

type ErrStartConnectionWorker string

func (e ErrStartConnectionWorker) Error() string {
	return "error start connection worker: " + string(e)
}

func (cw *ConnectionWorker) Start(game *game.Game, broadcast *GroupBroadcast) error {
	if cw.flagStarted {
		return ErrStartConnectionWorker("connection worker already started")
	}

	cw.flagStarted = true

	cw.conn.SetCloseHandler(cw.handleCloseConnection)
	cw.conn.SetReadLimit(readMessageLimit)

	// Input
	chInputBytes, chStop := cw.read()
	chInputMessages := cw.decode(chInputBytes, chStop)
	cw.broadcastInputMessage(chInputMessages, chStop)

	// Output
	// TODO: Create buffer const.
	chOutputMessages := cw.listenGameEvents(game.Events(chStop, 32), chStop)
	// TODO: Create buffer const.
	chOutputMessagesBroadcast := broadcast.OutputMessages(chStop, 32)
	chOutputBytes := cw.encode(chStop, chOutputMessages, chOutputMessagesBroadcast)
	cw.write(chOutputBytes, chStop)

	player := player.NewPlayer(game)
	player.Start(chStop)

	cw.chStop = chStop

	<-chStop

	cw.stopInputs()

	return nil
}

func (cw *ConnectionWorker) handleCloseConnection(code int, text string) error {
	message := websocket.FormatCloseMessage(code, "")
	cw.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(sendCloseConnectionTimeout))
	return nil
}

func (cw *ConnectionWorker) stopInputs() {
	cw.chsInputMux.Lock()
	defer cw.chsInputMux.Unlock()

	for _, ch := range cw.chsInput {
		close(ch)
	}

	cw.chsInput = cw.chsInput[:0]
}

func (cw *ConnectionWorker) read() (<-chan []byte, <-chan struct{}) {
	chout := make(chan []byte, chanReadMessagesBuffer)
	chstop := make(chan struct{}, 0)

	go func() {
		defer close(chout)
		defer close(chstop)

		for {
			messageType, data, err := cw.conn.ReadMessage()
			if err != nil {
				// TODO: Handle error?
				return
			}

			if websocket.TextMessage != messageType {
				// TODO: Handle case - unexpected message type?
				continue
			}

			chout <- data
		}
	}()

	return chout, chstop
}

func (cw *ConnectionWorker) decode(chin <-chan []byte, stop <-chan struct{}) <-chan InputMessage {
	chout := make(chan InputMessage, chanDecodeMessageBuffer)

	go func() {
		defer close(chout)

		var decoder = ffjson.NewDecoder()

		for {
			select {
			case data := <-chin:
				var inputMessage *InputMessage
				if err := decoder.Decode(data, &inputMessage); err != nil {
					// TODO: Handler error.
				} else {
					chout <- *inputMessage
				}
			case <-stop:
				return
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) broadcastInputMessage(chin <-chan InputMessage, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case inputMessage := <-chin:
				cw.chsInputMux.RLock()
				for _, ch := range cw.chsInput {
					select {
					case ch <- inputMessage:
					case <-stop:
						return
					}
				}
				cw.chsInputMux.RUnlock()
			case <-stop:
				return
			}
		}
	}()
}

func (cw *ConnectionWorker) Input(stop <-chan struct{}) <-chan InputMessage {
	// TODO: Create param buffer.

	chProxy := make(chan InputMessage, chanProxyInputMessageBuffer)

	cw.chsInputMux.Lock()
	cw.chsInput = append(cw.chsInput, chProxy)
	cw.chsInputMux.Unlock()

	chout := make(chan InputMessage, chanListenInputMessageBuffer)

	go func() {
		defer close(chout)
		defer func() {
			cw.chsInputMux.Lock()
			for i := range cw.chsInput {
				if cw.chsInput[i] == chProxy {
					cw.chsInput = append(cw.chsInput[:i], cw.chsInput[i+1:]...)
					close(chProxy)
					break
				}
			}
			cw.chsInputMux.Unlock()
		}()

		for {
			select {
			case <-stop:
				return
			case <-cw.chStop:
				return
			case inputMessage := <-chProxy:
				cw.sendInputMessage(chout, inputMessage, stop, sendInputMessageTimeout)
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) sendInputMessage(ch chan InputMessage, inputMessage InputMessage, stop <-chan struct{}, timeout time.Duration) {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()
	if cap(ch) == 0 {
		select {
		case ch <- inputMessage:
		case <-cw.chStop:
		case <-stop:
		case <-timer.C:
		}
	} else {
		for {
			select {
			case ch <- inputMessage:
				return
			case <-cw.chStop:
				return
			case <-stop:
				return
			case <-timer.C:
				return
			default:
				if len(ch) == cap(ch) {
					<-ch
				}
			}
		}
	}
}

func (cw *ConnectionWorker) write(chin <-chan []byte, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case data := <-chin:
				if err := cw.conn.WriteMessage(websocket.TextMessage, data); err != nil {
					// TODO: Handler error.
				}
			case <-stop:
				return
			}
		}
	}()
}

func (cw *ConnectionWorker) encode(stop <-chan struct{}, chins ...<-chan OutputMessage) <-chan []byte {
	chout := make(chan []byte, chanEncodeMessageBuffer)

	wg := sync.WaitGroup{}
	wg.Add(len(chins))

	for _, chin := range chins {
		go func(chin <-chan OutputMessage) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				case message, ok := <-chin:
					if !ok {
						return
					}
					if data, err := ffjson.Marshal(message); err != nil {
						// TODO: Handler error.
					} else {
						chout <- data
					}
				}
			}
		}(chin)
	}

	go func() {
		wg.Wait()
		close(chout)
	}()

	return chout
}

func (cw *ConnectionWorker) listenGameEvents(chin <-chan game.Event, stop <-chan struct{}) <-chan OutputMessage {
	chout := make(chan OutputMessage, chanOutputMessageBuffer)

	go func() {
		defer close(chout)

		for {
			select {
			// TODO: Fix case of closing channel chin.
			case event := <-chin:
				// TODO: Do stuff.

				outputMessage := OutputMessage{
					Type:    OutputMessageTypeGameEvent,
					Payload: event,
				}

				cw.logger.Info(event)

				select {
				case chout <- outputMessage:
				case <-stop:
					return
				}
			case <-stop:
				return
			}
		}
	}()

	return chout
}
