package connections

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
)

type ConnectionWorker struct {
	conn        *websocket.Conn
	logger      *logrus.Logger
	chStop      chan struct{}
	chStopErr   chan error
	chOutput    chan OutputMessage
	chsInput    []chan InputMessage
	chsInputMux *sync.RWMutex

	flagStarted bool
	flagStopped bool
}

func NewConnectionWorker(conn *websocket.Conn, logger *logrus.Logger) *ConnectionWorker {
	return &ConnectionWorker{
		conn:      conn,
		logger:    logger,
		chStopErr: make(chan error, 0),
	}
}

func (cw *ConnectionWorker) Start(game *game.Game) error {
	//return cw.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"), time.Now().Add(time.Second))

	if cw.flagStarted {
		// Return error
		return nil
	}

	// Start connection read
	// Start connection write
	// Listen game events
	// Parse input messages (?) and send game commands

	//game.RunObserver(cw, 16)

	cw.flagStarted = true

	return <-cw.chStopErr
}

func (cw *ConnectionWorker) read() <-chan []byte {
	// TODO: Create buffer.
	chout := make(chan []byte, 0)

	go func() {
		defer close(chout)

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

	return chout
}

func (cw *ConnectionWorker) decode(chin <-chan []byte, stop <-chan struct{}) <-chan InputMessage {
	// TODO: Create buffer.
	chout := make(chan InputMessage, 0)

	go func() {
		defer close(chout)

		var decoder = ffjson.NewDecoder()

		for {
			select {
			case data := <-chin:
				var inputMessage *InputMessage
				if err := decoder.DecodeFast(data, &inputMessage); err != nil {
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
	// TODO: Create buffer.
	chProxy := make(chan InputMessage, 0)

	cw.chsInputMux.Lock()
	cw.chsInput = append(cw.chsInput, chProxy)
	cw.chsInputMux.Unlock()

	// TODO: Create buffer.
	chout := make(chan InputMessage, 0)

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
				// TODO: Create timeout.
				cw.sendInputMessage(chout, inputMessage, stop, time.Second)
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

func (cw *ConnectionWorker) outputMessage(outputMessage OutputMessage) {
	select {
	case cw.chOutput <- outputMessage:
	case <-cw.chStop:
	}
}

func (cw *ConnectionWorker) encode(chin <-chan OutputMessage, stop <-chan struct{}) <-chan []byte {
	// TODO: Create buffer.
	chout := make(chan []byte, 0)

	go func() {
		defer close(chout)

		for {
			select {
			case message := <-chin:
				if data, err := ffjson.MarshalFast(message); err != nil {
					// TODO: Handler error.
				} else {
					chout <- data
				}
			case <-stop:
				return
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) listen(chin <-chan game.Event, stop <-chan struct{}) {
	for {
		select {
		case event := <-chin:
			// TODO: Do stuff.
			cw.logger.Info(event)
		case <-stop:
			return
		}
	}
}
