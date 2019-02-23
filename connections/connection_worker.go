package connections

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/broadcast"
	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/player"
)

const (
	chanPlayerOutputMessageBuffer  = 256
	chanPlayerEncodedMessageBuffer = 1024

	chanMergePreparedMessageBuffer = 8192

	chanReadMessagesBuffer           = 64
	chanDecodeMessageBuffer          = 64
	chanProxyInputMessageBuffer      = 64
	chanInputMessagesSnakeBuffer     = 64
	chanInputMessagesBroadcastBuffer = 64
	chanSnakeCommandsBuffer          = 64

	sendInputMessageTimeout  = time.Millisecond * 5
	sendOutputMessageTimeout = time.Millisecond * 25

	broadcastDelay = time.Second * 30

	ignoredBroadcastsCountToDisconnect = 100
)

type ConnectionWorker struct {
	conn   *websocket.Conn
	logger logrus.FieldLogger

	chsInput    []chan InputMessage
	chsInputMux *sync.RWMutex

	flagStarted bool
	startedMux  *sync.Mutex
}

func NewConnectionWorker(conn *websocket.Conn, logger logrus.FieldLogger) *ConnectionWorker {
	return &ConnectionWorker{
		conn:        conn,
		logger:      logger,
		chsInput:    make([]chan InputMessage, 0),
		chsInputMux: &sync.RWMutex{},

		flagStarted: false,
		startedMux:  &sync.Mutex{},
	}
}

type ErrStartConnectionWorker string

func (e ErrStartConnectionWorker) Error() string {
	return "error start connection worker: " + string(e)
}

func (cw *ConnectionWorker) Start(stop <-chan struct{}, game *game.Game, broadcast *broadcast.GroupBroadcast, gamePreparedMessages <-chan *websocket.PreparedMessage) error {
	cw.startedMux.Lock()
	if cw.flagStarted {
		cw.startedMux.Unlock()
		return ErrStartConnectionWorker("connection worker already started")
	}
	cw.flagStarted = true
	cw.startedMux.Unlock()

	broadcast.BroadcastMessage("user joined your game group")

	// Input
	chInputBytes, chStop := cw.read()
	chInputMessages := cw.decode(chInputBytes, chStop)
	cw.broadcastInputMessage(chInputMessages, chStop)
	chCommands := cw.listenSnakeCommands(chStop, cw.input(chStop, chanInputMessagesSnakeBuffer))
	cw.listenPlayerBroadcasts(chStop, cw.input(chStop, chanInputMessagesBroadcastBuffer), broadcast, broadcastDelay)

	p := player.NewPlayer(cw.logger, game.World())

	// Output
	chPlayer := p.Start(chStop, chCommands)
	chOutputBytes := cw.encode(chStop, cw.listenPlayer(chStop, chPlayer))
	chPlayerPreparedMessages := cw.prepare(chStop, chOutputBytes)
	chPreparedMessages := cw.mergePreparedMessagesChs(chStop, chPlayerPreparedMessages, gamePreparedMessages)
	chPreparedMessagesTimeout := cw.chPreparedMessageTimeout(chPreparedMessages, chStop, sendOutputMessageTimeout)
	cw.write(chPreparedMessagesTimeout, chStop)

	select {
	case <-chStop:
		// On connection error
	case <-stop:
		// External stop
		cw.logger.Warn("stop connection worker from external stopper channel")
	}

	broadcast.BroadcastMessage("user left your game group")

	cw.stopInputs()

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
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					cw.logger.Errorln("read input message error:", err)
				}
				return
			}

			if websocket.TextMessage != messageType {
				cw.logger.Warning("unexpected input message type")
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
			case data, ok := <-chin:
				if !ok {
					return
				}

				var inputMessage InputMessage
				if err := decoder.Decode(data, &inputMessage); err != nil {
					cw.logger.Errorln("decode input message error:", err)
				} else {
					select {
					case <-stop:
						return
					case chout <- inputMessage:
					}
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
			case inputMessage, ok := <-chin:
				if !ok {
					return
				}

				cw.doBroadcast(inputMessage, stop)
			case <-stop:
				return
			}
		}
	}()
}

func (cw *ConnectionWorker) doBroadcast(message InputMessage, stop <-chan struct{}) {
	cw.chsInputMux.RLock()
	defer cw.chsInputMux.RUnlock()

	for _, ch := range cw.chsInput {
		select {
		case ch <- message:
		case <-stop:
			return
		}
	}
}

func (cw *ConnectionWorker) input(stop <-chan struct{}, buffer uint) <-chan InputMessage {
	chProxy := make(chan InputMessage, chanProxyInputMessageBuffer)

	cw.chsInputMux.Lock()
	cw.chsInput = append(cw.chsInput, chProxy)
	cw.chsInputMux.Unlock()

	chout := make(chan InputMessage, buffer)

	go func() {
		defer close(chout)
		defer func() {
			go func() {
				for range chProxy {
				}
			}()

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
			case inputMessage, ok := <-chProxy:
				if !ok {
					return
				}
				cw.sendInputMessage(chout, inputMessage, stop, sendInputMessageTimeout)
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) sendInputMessage(ch chan InputMessage, inputMessage InputMessage, stop <-chan struct{}, timeout time.Duration) {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case ch <- inputMessage:
	case <-stop:
	case <-timer.C:
		cw.logger.WithFields(logrus.Fields{
			"timeout":      timeout,
			"message_type": inputMessage.Type,
		}).Warn("send input message to processing: time is out")
	}
}

func (cw *ConnectionWorker) listenSnakeCommands(stop <-chan struct{}, chin <-chan InputMessage) <-chan string {
	chout := make(chan string, chanSnakeCommandsBuffer)

	go func() {
		defer close(chout)

		for {
			select {
			case message, ok := <-chin:
				if !ok {
					return
				}

				if message.Type == InputMessageTypeSnakeCommand {
					select {
					case chout <- message.Payload:
					case <-stop:
						return
					}
				}
			case <-stop:
				return
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) listenPlayerBroadcasts(stop <-chan struct{}, chin <-chan InputMessage, b *broadcast.GroupBroadcast, delay time.Duration) {
	go func() {
		var (
			lastBroadcastTime time.Time
			ignored           int
		)

		for {
			select {
			case message, ok := <-chin:
				if !ok {
					return
				}

				if message.Type == InputMessageTypeBroadcast {
					if time.Since(lastBroadcastTime) > delay {
						b.BroadcastMessage(broadcast.Message(message.Payload))
						lastBroadcastTime = time.Now()
						ignored = 0
					} else {
						ignored += 1
						cw.logger.Warn("ignore broadcast: delay")

						if ignored > ignoredBroadcastsCountToDisconnect {
							cw.logger.Warn("ignored broadcasts limit reached")
							if err := cw.conn.Close(); err != nil {
								cw.logger.WithError(err).Error("close connection error")
							}
						}
					}
				}
			case <-stop:
				return
			}
		}
	}()
}

func (cw *ConnectionWorker) listenPlayer(stop <-chan struct{}, chin <-chan player.Message) <-chan OutputMessage {
	chout := make(chan OutputMessage, chanPlayerOutputMessageBuffer)

	go func() {
		defer close(chout)

		for {
			select {
			case event, ok := <-chin:
				if !ok {
					return
				}

				outputMessage := OutputMessage{
					Type:    OutputMessageTypePlayer,
					Payload: event,
				}

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

func (cw *ConnectionWorker) encode(stop <-chan struct{}, chins ...<-chan OutputMessage) <-chan []byte {
	chout := make(chan []byte, chanPlayerEncodedMessageBuffer)

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
						cw.logger.Errorln("encode output message error:", err)
					} else {
						select {
						case <-stop:
							return
						case chout <- data:
						}
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

func (cw *ConnectionWorker) prepare(stop <-chan struct{}, chin <-chan []byte) <-chan *websocket.PreparedMessage {
	chout := make(chan *websocket.PreparedMessage, cap(chin))

	go func() {
		defer close(chout)

		for {
			select {
			case data, ok := <-chin:
				if !ok {
					return
				}

				if pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data); err != nil {
					cw.logger.Errorln("prepare player output message error:", err)
				} else {
					select {
					case chout <- pm:
					case <-stop:
						return
					}
				}

			case <-stop:
				return
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) mergePreparedMessagesChs(stop <-chan struct{}, chins ...<-chan *websocket.PreparedMessage) <-chan *websocket.PreparedMessage {
	chout := make(chan *websocket.PreparedMessage, chanMergePreparedMessageBuffer)

	wg := sync.WaitGroup{}
	wg.Add(len(chins))

	for _, chin := range chins {
		go func(chin <-chan *websocket.PreparedMessage) {
			defer wg.Done()

			for {
				select {
				case <-stop:
					return
				case data, ok := <-chin:
					if !ok {
						return
					}

					select {
					case chout <- data:
					case <-stop:
						return
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

func (cw *ConnectionWorker) chPreparedMessageTimeout(chin <-chan *websocket.PreparedMessage, stop <-chan struct{}, timeout time.Duration) <-chan *websocket.PreparedMessage {
	chout := make(chan *websocket.PreparedMessage, cap(chin))

	go func() {
		defer close(chout)

		for {
			select {
			case pm, ok := <-chin:
				if !ok {
					return
				}

				timer := time.NewTimer(timeout)

				select {
				case chout <- pm:
					timer.Stop()
				case <-timer.C:
					cw.logger.Warn("send message for writing time is out")
					timer.Stop()
					continue
				case <-stop:
					timer.Stop()
					return
				}
			case <-stop:
				return
			}
		}
	}()

	return chout
}

func (cw *ConnectionWorker) write(chin <-chan *websocket.PreparedMessage, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case pm, ok := <-chin:
				if !ok {
					return
				}

				if err := cw.conn.WritePreparedMessage(pm); err != nil {
					if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						cw.logger.Errorln("write output message error:", err)
					}
					return
				}
			case <-stop:
				return
			}
		}
	}()
}
