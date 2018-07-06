package connections

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/broadcast"
	"github.com/ivan1993spb/snake-server/game"
)

const (
	chanBroadcastBuffer  = 128
	chanGameEventsBuffer = 512
	chanBytesProxyBuffer = 512
	chanBytesOutBuffer   = 128
)

const sendBytesTimeout = time.Millisecond * 50

type ConnectionGroup struct {
	limit      int
	counter    int
	counterMux *sync.RWMutex

	logger logrus.FieldLogger

	game      *game.Game
	broadcast *broadcast.GroupBroadcast

	chs    []chan []byte
	chsMux *sync.RWMutex

	stop    chan struct{}
	stopper *sync.Once
}

func NewConnectionGroup(logger logrus.FieldLogger, connectionLimit int, width, height uint8) (*ConnectionGroup, error) {
	g, err := game.NewGame(logger, width, height)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection group: %s", err)
	}

	if connectionLimit > 0 {
		return &ConnectionGroup{
			limit:      connectionLimit,
			counterMux: &sync.RWMutex{},
			game:       g,
			broadcast:  broadcast.NewGroupBroadcast(),
			logger:     logger,
			chs:        make([]chan []byte, 0),
			chsMux:     &sync.RWMutex{},
			stop:       make(chan struct{}),
			stopper:    &sync.Once{},
		}, nil
	}

	return nil, errors.New("cannot create connection group: invalid connection limit")
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
	cg.broadcastBytes(cg.encode(cg.stop, chMessagesGame, chMessagesBroadcast))
}

func (cg *ConnectionGroup) broadcastBytes(chin <-chan []byte) {
	go func() {
		for {
			select {
			case data, ok := <-chin:
				if !ok {
					return
				}

				cg.doBroadcast(data)
			case <-cg.stop:
				return
			}
		}
	}()
}

func (cg *ConnectionGroup) doBroadcast(data []byte) {
	cg.chsMux.RLock()
	defer cg.chsMux.RUnlock()

	for _, ch := range cg.chs {
		select {
		case ch <- data:
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

func (cg *ConnectionGroup) createChan() chan []byte {
	ch := make(chan []byte, chanBytesProxyBuffer)

	cg.chsMux.Lock()
	cg.chs = append(cg.chs, ch)
	cg.chsMux.Unlock()

	return ch
}

func (cg *ConnectionGroup) deleteChan(ch chan []byte) {
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

func (cg *ConnectionGroup) proxyCh(stop <-chan struct{}, buffer uint) <-chan []byte {
	ch := cg.createChan()
	chOut := make(chan []byte, buffer)

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
				cg.send(chOut, message, stop, sendBytesTimeout)
			}
		}
	}()

	return chOut
}

func (cg *ConnectionGroup) send(ch chan []byte, data []byte, stop <-chan struct{}, timeout time.Duration) {
	const tickSize = 5

	var timer = time.NewTimer(timeout)
	defer timer.Stop()

	var ticker = time.NewTicker(timeout / tickSize)
	defer ticker.Stop()

	if cap(ch) == 0 {
		select {
		case ch <- data:
		case <-cg.stop:
		case <-stop:
		case <-timer.C:
		}
	} else {
		for {
			select {
			case ch <- data:
				return
			case <-cg.stop:
				return
			case <-stop:
				return
			case <-timer.C:
				return
			case <-ticker.C:
				if len(ch) == cap(ch) {
					select {
					case <-ch:
					case ch <- data:
						return
					case <-stop:
						return
					case <-cg.stop:
						return
					case <-timer.C:
						return
					}
				}
			}
		}
	}
}

func (cg *ConnectionGroup) listenGame(stop <-chan struct{}, chin <-chan game.Event) <-chan OutputMessage {
	chout := make(chan OutputMessage, chanOutputMessageBuffer)

	go func() {
		defer close(chout)

		for {
			select {
			case event, ok := <-chin:
				if !ok {
					return
				}

				outputMessage := OutputMessage{
					Type:    OutputMessageTypeGame,
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

func (cg *ConnectionGroup) listenBroadcast(stop <-chan struct{}, chin <-chan broadcast.BroadcastMessage) <-chan OutputMessage {
	chout := make(chan OutputMessage, chanOutputMessageBuffer)

	go func() {
		defer close(chout)

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

func (cg *ConnectionGroup) encode(stop <-chan struct{}, chins ...<-chan OutputMessage) <-chan []byte {
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
						cg.logger.Errorln("encode output message error:", err)
					} else {
						select {
						case chout <- data:
						case <-stop:
							return
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

func (cg *ConnectionGroup) BroadcastMessageTimeout(message string, timeout time.Duration) bool {
	return cg.broadcast.BroadcastMessageTimeout(broadcast.BroadcastMessage(message), timeout)
}
