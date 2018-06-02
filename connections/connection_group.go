package connections

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/broadcast"
	"github.com/ivan1993spb/snake-server/game"
)

type ConnectionGroup struct {
	limit   int
	counter int
	mutex   *sync.RWMutex

	logger logrus.FieldLogger

	game      *game.Game
	broadcast *broadcast.GroupBroadcast

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

	if err := connectionWorker.Start(cg.stop, cg.game, cg.broadcast); err != nil {
		return &ErrRunConnection{
			Err: err,
		}
	}

	return nil
}

func (cg *ConnectionGroup) Start() {
	cg.broadcast.Start(cg.stop)
	cg.game.Start(cg.stop)
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
