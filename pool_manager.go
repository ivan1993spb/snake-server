package main

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

// Pool interface represents pool with connections
type Pool interface {
	// IsFull returns true if pool is full
	IsFull() bool
	// IsEmpty returns true if pool is empty
	IsEmpty() bool
	// AddConn creates connection in the pool
	AddConn(ws *websocket.Conn) (pwshandler.Environment, error)
	// DelConn removes connection from pool
	DelConn(ws *websocket.Conn) error
	// HasConn returns true if passed connection belongs to the pool
	HasConn(ws *websocket.Conn) bool
}

type errCreatingPoolManager struct {
	err error
}

func (e *errCreatingPoolManager) Error() string {
	return "Cannot create pool manager: " + e.err.Error()
}

// PoolFactory must generate new pool
type PoolFactory func() (Pool, error)

type GamePoolManager struct {
	addPool PoolFactory
	pools   []Pool
}

// NewGamePoolManager creates new GamePoolManager with fixed max
// number of pools specified by poolLimit
func NewGamePoolManager(factory PoolFactory, poolLimit uint8,
) (pwshandler.PoolManager, error) {
	if factory == nil {
		return nil, &errCreatingPoolManager{
			errors.New("Passed nil pool factory"),
		}
	}
	if poolLimit == 0 {
		return nil, &errCreatingPoolManager{
			errors.New("Invalid pool limit"),
		}
	}

	return &GamePoolManager{factory, make([]Pool, 0, poolLimit)}, nil
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) AddConn(ws *websocket.Conn,
) (pwshandler.Environment, error) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Try to add new connection in a pool")
		glog.Infoln("Try to find not full pool")
	}
	// Try to find not full pool
	for i := range pm.pools {
		if !pm.pools[i].IsFull() {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Was found not full pool")
				glog.Infoln("Creating connection to pool")
			}
			return pm.pools[i].AddConn(ws)
		}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Cannot find not full pool")
	}

	// Try to create pool
	if len(pm.pools) != cap(pm.pools) {
		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("Server is not full so create new pool")
		}

		// Creating new pool
		pool, err := pm.addPool()

		if err == nil {
			// Save the pool
			pm.pools = append(pm.pools, pool)

			if glog.V(INFOLOG_LEVEL_POOLS) {
				glog.Infoln("New pool was created")
			}
			// Create connection to new pool
			return pool.AddConn(ws)
		}

		return nil, fmt.Errorf("Cannot create new pool: %s", err)
	}

	return nil, errors.New("Cannot create new pool: server is full")
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) DelConn(ws *websocket.Conn) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Try to remove information about connection")
		glog.Infoln("Try to find pool of closed connection")
	}

	for i := range pm.pools {
		// If current pool has the connection...
		if pm.pools[i].HasConn(ws) {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Pool of closed connection was found")
				glog.Infoln("Removing closed connection from pool")
			}

			// Remove it
			err := pm.pools[i].DelConn(ws)

			// And now if pool is empty
			if pm.pools[i].IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("Removing empty pool")
				}
				// Delete pool
				pm.pools = append(pm.pools[:i], pm.pools[i+1:]...)

				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("Empty pool was removed")
				}
			}

			return err
		}
	}

	return errors.New("Connection to removing was not found")
}
