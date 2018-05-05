package connections

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionGroup struct {
	logger           *logrus.Logger
	connections      map[int]*websocket.Conn
	connectionsMutex *sync.RWMutex
	connectionLimit  int
}

// TODO: Is it necessary to pass logger in group manager ?
func NewConnectionGroup(logger *logrus.Logger, connectionLimit int, mapWidth, mapHeight uint8) (*ConnectionGroup, error) {
	// TODO: Check input params.

	// TODO: Create game. (?)

	return &ConnectionGroup{
		logger:           logger,
		connections:      make(map[int]*websocket.Conn),
		connectionsMutex: &sync.RWMutex{},
		connectionLimit:  connectionLimit,
	}, nil
}

// unsafeIsFull returns true if group is full
func (cg *ConnectionGroup) unsafeIsFull() bool {
	return len(cg.connections) == cg.connectionLimit
}

// unsafeIsEmpty returns true if group is empty
func (cg *ConnectionGroup) unsafeIsEmpty() bool {
	return len(cg.connections) == 0
}

func (cg *ConnectionGroup) IsEmpty() bool {
	cg.connectionsMutex.RLock()
	defer cg.connectionsMutex.RUnlock()
	return cg.unsafeIsEmpty()
}

type ErrAddConnectionToGroup string

func (e ErrAddConnectionToGroup) Error() string {
	return "cannot add connection to group: " + string(e)
}

// AddConn adds connection to group
func (cg *ConnectionGroup) Add(conn *websocket.Conn) error {
	cg.connectionsMutex.Lock()
	defer cg.connectionsMutex.Unlock()

	if cg.unsafeIsFull() {
		return ErrAddConnectionToGroup("group is full")
	}

	for id := 0; id <= len(cg.connections); id++ {
		if _, occupied := cg.connections[id]; !occupied {
			cg.connections[id] = conn
			return nil
		}
	}

	return ErrAddConnectionToGroup("cannot get free id")
}

// Delete removes connection from group
func (cg *ConnectionGroup) Delete(conn *websocket.Conn) error {
	cg.connectionsMutex.Lock()
	defer cg.connectionsMutex.Unlock()

	for id := range cg.connections {
		if cg.connections[id] == conn {
			delete(cg.connections, id)
			return nil
		}
	}

	return errors.New("cannot delete connection from group: connection was not found")
}
