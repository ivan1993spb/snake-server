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
	chanBroadcastBuffer  = 32
	chanGameEventsBuffer = 32
)

type ConnectionGroup struct {
	limit   int
	counter int
	mutex   *sync.RWMutex

	logger logrus.FieldLogger

	game      *game.Game
	broadcast *broadcast.GroupBroadcast

	chs    []chan []byte
	chsMux *sync.RWMutex

	stop chan struct{}
}

func NewConnectionGroup(logger logrus.FieldLogger, connectionLimit int, width, height uint8) (*ConnectionGroup, error) {
	g, err := game.NewGame(logger, width, height)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection group: %s", err)
	}

	if connectionLimit > 0 {
		return &ConnectionGroup{
			limit:     connectionLimit,
			mutex:     &sync.RWMutex{},
			game:      g,
			broadcast: broadcast.NewGroupBroadcast(),
			logger:    logger,
			chs:       make([]chan []byte, 0),
			chsMux:    &sync.RWMutex{},
			stop:      make(chan struct{}),
		}, nil
	}

	return nil, errors.New("cannot create connection group: invalid connection limit")
}

func (cg *ConnectionGroup) GetLimit() int {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()
	return cg.limit
}

func (cg *ConnectionGroup) SetLimit(limit int) {
	cg.mutex.Lock()
	cg.limit = limit
	cg.mutex.Unlock()
}

func (cg *ConnectionGroup) GetCount() int {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()
	return cg.counter
}

// unsafeIsFull returns true if group is full
func (cg *ConnectionGroup) unsafeIsFull() bool {
	return cg.counter == cg.limit
}

func (cg *ConnectionGroup) IsFull() bool {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()
	return cg.unsafeIsFull()
}

// unsafeIsEmpty returns true if group is empty
func (cg *ConnectionGroup) unsafeIsEmpty() bool {
	return cg.counter == 0
}

func (cg *ConnectionGroup) IsEmpty() bool {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()
	return cg.unsafeIsEmpty()
}

type ErrRunConnection struct {
	Err error
}

func (e *ErrRunConnection) Error() string {
	return "run connection error: " + e.Err.Error()
}

var ErrGroupIsFull = errors.New("group is full")

func (cg *ConnectionGroup) Handle(connectionWorker *ConnectionWorker) *ErrRunConnection {
	cg.mutex.Lock()
	if cg.unsafeIsFull() {
		cg.mutex.Unlock()
		return &ErrRunConnection{
			Err: ErrGroupIsFull,
		}
	}
	cg.counter += 1
	cg.mutex.Unlock()

	defer func() {
		cg.mutex.Lock()
		cg.counter -= 1
		cg.mutex.Unlock()
	}()

	// TODO: Create proxy chan buffer.
	if err := connectionWorker.Start(cg.stop, cg.game, cg.broadcast, cg.proxyCh(cg.stop, 64)); err != nil {
		return &ErrRunConnection{
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

				cg.chsMux.RLock()
				for _, ch := range cg.chs {
					select {
					case ch <- data:
					case <-cg.stop:
					}
				}
				cg.chsMux.RUnlock()
			case <-cg.stop:
				return
			}
		}
	}()
}

func (cg *ConnectionGroup) Stop() {
	close(cg.stop)
}

func (cg *ConnectionGroup) GetWorldWidth() uint8 {
	return cg.game.World().Width()
}

func (cg *ConnectionGroup) GetWorldHeight() uint8 {
	return cg.game.World().Height()
}

func (cg *ConnectionGroup) createChan() chan []byte {
	// TODO: Create buffer const.
	ch := make(chan []byte, 64)

	cg.chsMux.Lock()
	cg.chs = append(cg.chs, ch)
	cg.chsMux.Unlock()

	return ch
}

func (cg *ConnectionGroup) deleteChan(ch chan []byte) {
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
				// TODO: Create timeout const.
				cg.send(chOut, message, stop, time.Millisecond*50)
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
					<-ch
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
			case event := <-chin:
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
			case message := <-chin:
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
