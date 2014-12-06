package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// Pool represents pool with connections
type Pool interface {
	// IsFull returns true if pool is full
	IsFull() bool
	// IsEmpty returns true if pool is empty
	IsEmpty() bool
	// AddConn creates connection in the pool
	AddConn(*websocket.Conn) (pwshandler.Environment, error)
	// DelConn removes connection from pool and stops all pool
	// goroutines
	DelConn(*websocket.Conn)
	// HasConn returns true if passed connection belongs to the pool
	HasConn(*websocket.Conn) bool
}

// PoolFactory represents pool factory
type PoolFactory interface {
	// NewPool creates new Pool
	NewPool() (Pool, error)
}

type GamePoolManager struct {
	factory PoolFactory
	pools   []Pool // Pool storage
}

// NewGamePoolManager creates new GamePoolManager with fixed max
// number of pools specified by poolLimit
func NewGamePoolManager(factory PoolFactory, poolLimit uint8,
) (pwshandler.PoolManager, error) {
	if factory == nil {
		return nil, errors.New("Passed nil pool factory")
	}
	if poolLimit == 0 {
		return nil, errors.New("Invalid pool limit")
	}
	return &GamePoolManager{factory, make([]Pool, 0, poolLimit)}, nil
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) AddConn(conn *websocket.Conn,
) (pwshandler.Environment, error) {
	if glog.V(3) {
		glog.Infoln("Accepted new connection")
	}
	if glog.V(4) {
		glog.Infoln("Try to find not full pool")
	}
	// Try to find not full pool
	for i := range pm.pools {
		if !pm.pools[i].IsFull() {
			if glog.V(4) {
				glog.Infoln("Was found not full pool")
			}
			if glog.V(3) {
				glog.Infoln("Creating connection to pool")
			}
			return pm.pools[i].AddConn(conn)
		}
	}

	if glog.V(3) {
		glog.Infoln("Cannot find not full pool")
	}

	// Try to create pool
	if !pm.isFull() {
		if glog.V(3) {
			glog.Infoln("Server is not full: create new pool")
		}

		// Creating new pool
		if newPool, err := pm.factory.NewPool(); err == nil {
			if glog.V(3) {
				glog.Infoln("New pool was created")
			}
			// Save the pool
			pm.pools = append(pm.pools, newPool)
			// Create connection to the pool
			return newPool.AddConn(conn)
		} else {
			if glog.V(3) {
				glog.Infoln("Cannot create new pool")
			}
			return nil, err
		}
	} else {
		if glog.V(3) {
			glog.Infoln("Cannot create new pool: server is full")
		}
	}

	return nil, errors.New("Pool manager refuses to add connection")
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) DelConn(conn *websocket.Conn) {
	if glog.V(4) {
		glog.Infoln("Deleting information about closed connection")
	}
	if glog.V(3) {
		glog.Infoln("Try to find pool of closed connection")
	}
	for i := range pm.pools {
		// If current pool has the connection...
		if pm.pools[i].HasConn(conn) {
			if glog.V(3) {
				glog.Infoln("Pool of removing connection was found")
				glog.Infoln("Removing closed connection from pool")
			}
			// Remove it
			pm.pools[i].DelConn(conn)

			// And now if pool is empty
			if pm.pools[i].IsEmpty() {
				if glog.V(3) {
					glog.Infoln("Removing empty pool")
				}
				// Delete pool
				pm.pools = append(pm.pools[:i], pm.pools[i+1:]...)
				if glog.V(4) {
					glog.Infoln("Empty pool was removed")
				}
			}

			return
		}
	}
}

// isFull returns true if pool storage is full
func (pm *GamePoolManager) isFull() bool {
	return len(pm.pools) == cap(pm.pools)
}
