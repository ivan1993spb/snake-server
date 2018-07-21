package connections

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/broadcast"
	"github.com/ivan1993spb/snake-server/game"
)

const (
	chanBroadcastBuffer  = 128
	chanGameEventsBuffer = 8192
	chanBytesProxyBuffer = 8192
	chanBytesOutBuffer   = 8192

	chanEncodeGroupMessageBuffer = 8192

	sendPreparedMessageTimeout = time.Millisecond * 50

	broadcastOutputMessageBufferMonitoringDelay = time.Second * 10
	gameOutputMessageBufferMonitoringDelay      = time.Second * 10
	encodeGroupMessageBufferMonitoringDelay     = time.Second * 10
	preparedMessageBufferMonitoringDelay        = time.Second * 10
)

type ConnectionGroup struct {
	limit      int
	counter    int
	counterMux *sync.RWMutex

	logger logrus.FieldLogger

	game      *game.Game
	broadcast *broadcast.GroupBroadcast

	chs    []chan *websocket.PreparedMessage
	chsMux *sync.RWMutex

	stop    chan struct{}
	stopper *sync.Once
}

type errCreateConnectionGroup string

func (e errCreateConnectionGroup) Error() string {
	return "cannot create connection group: " + string(e)
}

func NewConnectionGroup(logger logrus.FieldLogger, connectionLimit int, width, height uint8) (*ConnectionGroup, error) {
	g, err := game.NewGame(logger, width, height)
	if err != nil {
		return nil, errCreateConnectionGroup(err.Error())
	}

	if connectionLimit > 0 {
		return &ConnectionGroup{
			limit:      connectionLimit,
			counterMux: &sync.RWMutex{},
			game:       g,
			broadcast:  broadcast.NewGroupBroadcast(),
			logger:     logger,
			chs:        make([]chan *websocket.PreparedMessage, 0),
			chsMux:     &sync.RWMutex{},
			stop:       make(chan struct{}),
			stopper:    &sync.Once{},
		}, nil
	}

	return nil, errCreateConnectionGroup("invalid connection limit")
}

func (cg *ConnectionGroup) GetLimit() int {
	cg.counterMux.RLock()
	defer cg.counterMux.RUnlock()
	return cg.limit
}

func (cg *ConnectionGroup) SetLimit(limit int) {
	cg.counterMux.Lock()
	cg.limit = limit
	cg.counterMux.Unlock()
}

func (cg *ConnectionGroup) GetCount() int {
	cg.counterMux.RLock()
	defer cg.counterMux.RUnlock()
	return cg.counter
}

// unsafeIsFull returns true if group is full
func (cg *ConnectionGroup) unsafeIsFull() bool {
	return cg.counter == cg.limit
}

func (cg *ConnectionGroup) IsFull() bool {
	cg.counterMux.RLock()
	defer cg.counterMux.RUnlock()
	return cg.unsafeIsFull()
}

// unsafeIsEmpty returns true if group is empty
func (cg *ConnectionGroup) unsafeIsEmpty() bool {
	return cg.counter == 0
}

func (cg *ConnectionGroup) IsEmpty() bool {
	cg.counterMux.RLock()
	defer cg.counterMux.RUnlock()
	return cg.unsafeIsEmpty()
}

type ErrHandleConnection struct {
	Err error
}

func (e *ErrHandleConnection) Error() string {
	return "handle connection error: " + e.Err.Error()
}

var ErrGroupIsFull = errors.New("group is full")

func (cg *ConnectionGroup) Handle(connectionWorker *ConnectionWorker) error {
	cg.counterMux.Lock()
	if cg.unsafeIsFull() {
		cg.counterMux.Unlock()
		return &ErrHandleConnection{
			Err: ErrGroupIsFull,
		}
	}
	cg.counter += 1
	cg.counterMux.Unlock()

	defer func() {
		cg.counterMux.Lock()
		cg.counter -= 1
		cg.counterMux.Unlock()
	}()

	chStopHandle := make(chan struct{})
	defer close(chStopHandle)

	if err := connectionWorker.Start(cg.stop, cg.game, cg.broadcast, cg.proxyCh(chStopHandle, chanBytesOutBuffer)); err != nil {
		return &ErrHandleConnection{
			Err: err,
		}
	}

	return nil
}

func (cg *ConnectionGroup) Start() {
	cg.broadcast.Start(cg.stop)
	cg.game.Start(cg.stop)

	chMessagesGame := cg.listenGame(cg.stop, cg.game.ListenEvents(cg.stop, chanGameEventsBuffer))
	chMessagesBroadcast := cg.listenBroadcast(cg.stop, cg.broadcast.ListenMessages(cg.stop, chanBroadcastBuffer))
	chBytes := cg.encode(cg.stop, chMessagesGame, chMessagesBroadcast)
	chPreparedMessages := cg.prepare(cg.stop, chBytes)
	cg.broadcastPreparedMessages(chPreparedMessages)
}

func (cg *ConnectionGroup) broadcastPreparedMessages(chin <-chan *websocket.PreparedMessage) {
	go func() {
		for {
			select {
			case pm, ok := <-chin:
				if !ok {
					return
				}

				cg.doBroadcast(pm)
			case <-cg.stop:
				return
			}
		}
	}()
}

func (cg *ConnectionGroup) doBroadcast(pm *websocket.PreparedMessage) {
	cg.chsMux.RLock()
	defer cg.chsMux.RUnlock()

	for _, ch := range cg.chs {
		select {
		case ch <- pm:
		case <-cg.stop:
			return
		}
	}
}

func (cg *ConnectionGroup) Stop() {
	cg.stopper.Do(func() {
		close(cg.stop)
	})
}

func (cg *ConnectionGroup) GetWorldWidth() uint8 {
	return cg.game.World().Width()
}

func (cg *ConnectionGroup) GetWorldHeight() uint8 {
	return cg.game.World().Height()
}

func (cg *ConnectionGroup) GetObjects() []interface{} {
	return cg.game.World().GetObjects()
}

func (cg *ConnectionGroup) createChan() chan *websocket.PreparedMessage {
	ch := make(chan *websocket.PreparedMessage, chanBytesProxyBuffer)

	cg.chsMux.Lock()
	cg.chs = append(cg.chs, ch)
	cg.chsMux.Unlock()

	return ch
}

func (cg *ConnectionGroup) deleteChan(ch chan *websocket.PreparedMessage) {
	go func() {
		for range ch {
		}
	}()

	cg.chsMux.Lock()
	for i := range cg.chs {
		if cg.chs[i] == ch {
			cg.chs = append(cg.chs[:i], cg.chs[i+1:]...)
			close(ch)
			break
		}
	}
	cg.chsMux.Unlock()
}

func (cg *ConnectionGroup) proxyCh(stop <-chan struct{}, buffer uint) <-chan *websocket.PreparedMessage {
	ch := cg.createChan()
	chOut := make(chan *websocket.PreparedMessage, buffer)

	go func() {
		defer close(chOut)
		defer cg.deleteChan(ch)

		for {
			select {
			case <-stop:
				return
			case <-cg.stop:
				return
			case message, ok := <-ch:
				if !ok {
					return
				}
				cg.sendTimeout(chOut, message, stop, sendPreparedMessageTimeout)
			}
		}
	}()

	return chOut
}

func (cg *ConnectionGroup) sendTimeout(ch chan *websocket.PreparedMessage, pm *websocket.PreparedMessage, stop <-chan struct{}, timeout time.Duration) {
	const warnFormat = "game group message was not send to connection: %s"
	var timer = time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- pm:
	case <-cg.stop:
		cg.logger.Warnf(warnFormat, "game group stopped")
	case <-stop:
		cg.logger.Warnf(warnFormat, "connection handler stopped")
	case <-timer.C:
		cg.logger.Warnf(warnFormat, "time is out")
		if len(ch) == cap(ch) {
			cg.logger.Warn("connection group output channel buffer is overflow for connection")
		}
	}
}

func (cg *ConnectionGroup) listenGame(stop <-chan struct{}, chin <-chan game.Event) <-chan OutputMessage {
	chout := make(chan OutputMessage, cap(chin))

	go func() {
		defer close(chout)

		ticker := time.NewTicker(gameOutputMessageBufferMonitoringDelay)
		defer ticker.Stop()

		var count = 0

		for {
			select {
			case event, ok := <-chin:
				if !ok {
					return
				}

				// Do not send internal game errors to clients
				if event.Type == game.EventTypeError {
					continue
				}

				// Do not send checked events to clients
				if event.Type == game.EventTypeObjectChecked {
					continue
				}

				outputMessage := OutputMessage{
					Type:    OutputMessageTypeGame,
					Payload: event,
				}

				select {
				case chout <- outputMessage:
					count++
				case <-stop:
					return
				}
			case <-stop:
				return
			case <-ticker.C:
				cg.logger.WithFields(logrus.Fields{
					"buffered_messages": len(chout),
					"buffer_size":       cap(chout),
					"time_frame":        gameOutputMessageBufferMonitoringDelay,
					"count":             count,
				}).Debug("game output messages buffer monitoring")

				count = 0
			}
		}
	}()

	return chout
}

func (cg *ConnectionGroup) listenBroadcast(stop <-chan struct{}, chin <-chan broadcast.Message) <-chan OutputMessage {
	chout := make(chan OutputMessage, cap(chin))

	go func() {
		defer close(chout)

		ticker := time.NewTicker(broadcastOutputMessageBufferMonitoringDelay)
		defer ticker.Stop()

		var count = 0

		for {
			select {
			case message, ok := <-chin:
				if !ok {
					return
				}

				outputMessage := OutputMessage{
					Type:    OutputMessageTypeBroadcast,
					Payload: message,
				}

				select {
				case chout <- outputMessage:
					count++
				case <-stop:
					return
				}
			case <-stop:
				return
			case <-ticker.C:
				cg.logger.WithFields(logrus.Fields{
					"buffered_messages": len(chout),
					"buffer_size":       cap(chout),
					"time_frame":        broadcastOutputMessageBufferMonitoringDelay,
					"count":             count,
				}).Debug("broadcast output messages buffer monitoring")

				count = 0
			}
		}
	}()

	return chout
}

func (cg *ConnectionGroup) encode(stop <-chan struct{}, chins ...<-chan OutputMessage) <-chan []byte {
	chout := make(chan []byte, chanEncodeGroupMessageBuffer)

	wg := sync.WaitGroup{}
	wg.Add(len(chins))

	for i, chin := range chins {
		go func(i int, chin <-chan OutputMessage) {
			defer wg.Done()

			ticker := time.NewTicker(encodeGroupMessageBufferMonitoringDelay)
			defer ticker.Stop()

			var count = 0

			for {
				select {
				case <-stop:
					return
				case message, ok := <-chin:
					if !ok {
						return
					}

					if data, err := ffjson.Marshal(message); err != nil {
						cg.logger.Errorln("encode output message error:", err)
					} else {
						select {
						case chout <- data:
							count++
						case <-stop:
							return
						}
					}
				case <-ticker.C:
					cg.logger.WithFields(logrus.Fields{
						"buffered_messages": len(chout),
						"buffer_size":       chanEncodeGroupMessageBuffer,
						"time_frame":        encodeGroupMessageBufferMonitoringDelay,
						"count":             count,
						"channel":           i,
					}).Debug("encoded group messages buffer monitoring")

					count = 0
				}
			}
		}(i, chin)
	}

	go func() {
		wg.Wait()
		close(chout)
	}()

	return chout
}

func (cg *ConnectionGroup) prepare(stop <-chan struct{}, chin <-chan []byte) <-chan *websocket.PreparedMessage {
	chout := make(chan *websocket.PreparedMessage, cap(chin))

	go func() {
		defer close(chout)

		ticker := time.NewTicker(preparedMessageBufferMonitoringDelay)
		defer ticker.Stop()

		var count = 0

		for {
			select {
			case data, ok := <-chin:
				if !ok {
					return
				}

				if pm, err := websocket.NewPreparedMessage(websocket.TextMessage, data); err != nil {
					cg.logger.Errorln("prepare group output message error:", err)
				} else {
					select {
					case chout <- pm:
						count++
					case <-stop:
						return
					}
				}
			case <-stop:
				return
			case <-ticker.C:
				cg.logger.WithFields(logrus.Fields{
					"buffered_messages": len(chout),
					"buffer_size":       cap(chout),
					"time_frame":        preparedMessageBufferMonitoringDelay,
					"count":             count,
				}).Debug("prepared messages buffer monitoring")

				count = 0
			}
		}
	}()

	return chout
}

func (cg *ConnectionGroup) BroadcastMessageTimeout(message string, timeout time.Duration) bool {
	return cg.broadcast.BroadcastMessageTimeout(broadcast.Message(message), timeout)
}
