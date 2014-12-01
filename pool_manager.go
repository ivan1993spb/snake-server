package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
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
) pwshandler.PoolManager {
	return &GamePoolManager{factory, make([]Pool, 0, poolLimit)}
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) AddConn(conn *websocket.Conn,
) (pwshandler.Environment, error) {
	// Try to find not full pool
	for i := range pm.pools {
		if !pm.pools[i].IsFull() {
			return pm.pools[i].AddConn(conn)
		}
	}

	// Try to create pool
	if !pm.isFull() {
		// Generate new pool
		if newPool, err := pm.factory.NewPool(); err == nil {
			// Save the pool
			pm.pools = append(pm.pools, newPool)
			// Create connection to the pool
			return newPool.AddConn(conn)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Pool manager refuses to add connection")
}

// Implementing pwshandler.ConnManager interface
func (pm *GamePoolManager) DelConn(conn *websocket.Conn) {
	for i := range pm.pools {
		// If current pool has the connection...
		if pm.pools[i].HasConn(conn) {
			// Remove it
			pm.pools[i].DelConn(conn)

			// And now if pool is empty
			if pm.pools[i].IsEmpty() {
				// Delete pool
				pm.pools = append(pm.pools[:i], pm.pools[i+1:]...)
			}

			return
		}
	}
}

// isFull returns true if pool storage is full
func (pm *GamePoolManager) isFull() bool {
	return len(pm.pools) == cap(pm.pools)
}
