package main

import (
	"errors"

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
	// DelConn removes connection from the pool
	DelConn(ws *websocket.Conn) error
	// HasConn returns true if passed connection belongs to the pool
	HasConn(ws *websocket.Conn) bool
}

type errCreatingPoolManager struct {
	err error
}

func (e *errCreatingPoolManager) Error() string {
	return "cannot create pool manager: " + e.err.Error()
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
			errors.New("passed nil pool factory"),
		}
	}
	if poolLimit == 0 {
		return nil, &errCreatingPoolManager{
			errors.New("invalid pool limit"),
		}
	}

	return &GamePoolManager{factory, make([]Pool, 0, poolLimit)}, nil
}

type errCannotAddConn struct {
	err error
}

func (e *errCannotAddConn) Error() string {
	return "cannot add connection: " + e.err.Error()
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) AddConn(ws *websocket.Conn,
) (pwshandler.Environment, error) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("try to add new connection in a pool")
		glog.Infoln("try to find not full pool")
	}
	// Try to find not full pool
	for i := range pm.pools {
		if !pm.pools[i].IsFull() {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("was found not full pool")
				glog.Infoln("creating connection to the pool")
			}
			return pm.pools[i].AddConn(ws)
		}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("cannot find not full pool")
	}

	// Try to create pool if server is not full
	if len(pm.pools) != cap(pm.pools) {
		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("server is not full")
			glog.Infoln("creating pool")
		}

		pool, err := pm.addPool()

		if err == nil {
			// Save the pool
			pm.pools = append(pm.pools, pool)

			if glog.V(INFOLOG_LEVEL_POOLS) {
				glog.Infoln("pool was created")
			}
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("creating connection to pool")
			}

			// Create connection to new pool
			return pool.AddConn(ws)
		}

		return nil, &errCannotAddConn{err}
	}

	return nil, &errCannotAddConn{errors.New("server is full")}
}

type errCannotDelConn struct {
	err error
}

func (e *errCannotDelConn) Error() string {
	return "cannot delete connection: " + e.err.Error()
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) DelConn(ws *websocket.Conn) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("try to remove information about connection")
		glog.Infoln("try to find pool of connection")
	}

	for i := range pm.pools {
		// If current pool has the connection...
		if pm.pools[i].HasConn(ws) {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("pool of connection was found")
				glog.Infoln("removing connection from pool")
			}

			if err := pm.pools[i].DelConn(ws); err != nil {
				return &errCannotDelConn{err}
			}

			if pm.pools[i].IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("removing empty pool")
				}
				// Delete pool
				pm.pools = append(pm.pools[:i], pm.pools[i+1:]...)

				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("empty pool was removed")
				}
			}

			return nil
		}
	}

	return &errCannotDelConn{errors.New("connection was not found")}
}
