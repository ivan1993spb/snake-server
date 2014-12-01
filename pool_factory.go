package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"bitbucket.org/pushkin_ivan/simple-2d-playground"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type PGPoolFactory struct {
	rootCxt   context.Context // Root context
	connLimit uint8           // Max connection number in pool
	pgW, pgH  uint8           // Playground size
}

func NewPGPoolFactory(rootCxt context.Context, connLimit,
	pgW, pgH uint8) PoolFactory {
	return &PGPoolFactory{rootCxt, connLimit, pgW, pgH}
}

// Implementing PoolFactory interface
func (f *PGPoolFactory) NewPool() (Pool, error) {
	var (
		pg  *playground.Playground
		err error
	)

	if pg, err = playground.NewPlayground(f.pgW, f.pgH); err != nil {
		return nil, err
	}

	// Create context for each pool
	cxt, cancel := context.WithCancel(f.rootCxt)

	return NewGamePool(cxt, cancel, f.connLimit, pg), nil
}

type GamePool struct {
	conns []*websocket.Conn // Connection list

	// Goroutine management
	cxt    context.Context
	cancel context.CancelFunc

	pg *playground.Playground
}

func NewGamePool(cxt context.Context, cancel context.CancelFunc,
	connLimit uint8, pg *playground.Playground) *GamePool {

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *       GAME LOGIC STARTS HERE. INIT PLAYGROUND       *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	return &GamePool{make([]*websocket.Conn, 0, connLimit), cxt,
		cancel, pg}
}

// Implementing Pool interface
func (p *GamePool) IsFull() bool {
	return len(p.conns) == cap(p.conns)
}

// Implementing Pool interface
func (p *GamePool) IsEmpty() bool {
	return len(p.conns) == 0
}

// Implementing Pool interface
func (p *GamePool) AddConn(conn *websocket.Conn) (
	pwshandler.Environment, error) {

	if p.IsFull() {
		return nil, errors.New("Pool is full")
	}
	if p.HasConn(conn) {
		return nil, errors.New("Pool already has passed connection")
	}

	p.conns = append(p.conns, conn)

	// Create context for each connection in pool
	cxt, _ := context.WithCancel(p.cxt)

	return &GameData{cxt, p.pg}, nil
}

// Implementing Pool interface
func (p *GamePool) DelConn(conn *websocket.Conn) {
	if p.HasConn(conn) {
		for i := range p.conns {
			// Find connection
			if p.conns[i] == conn {
				// Delete connection
				p.conns = append(p.conns[:i], p.conns[i+1:]...)
				// Stop all child goroutines if empty pool
				if p.IsEmpty() && p.cancel != nil {
					p.cancel()
				}

				return
			}
		}
	}
}

// Implementing Pool interface
func (p *GamePool) HasConn(conn *websocket.Conn) bool {
	for i := range p.conns {
		if p.conns[i] == conn {
			return true
		}
	}

	return false
}
