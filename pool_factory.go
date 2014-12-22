package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var ErrInvalidConnLimit = errors.New("Invalid connection limit")

type errCreatingPoolFactory struct {
	err error
}

func (e *errCreatingPoolFactory) Error() string {
	return "Cannot create pool factory: " + e.err.Error()
}

func NewPGPoolFactory(rootCxt context.Context, connLimit,
	pgW, pgH uint8) (PoolFactory, error) {
	if err := rootCxt.Err(); err != nil {
		return nil, &errCreatingPoolFactory{err}
	}
	if connLimit == 0 {
		return nil, &errCreatingPoolFactory{ErrInvalidConnLimit}
	}

	return func() (Pool, error) {
		pool, err := NewPGPool(rootCxt, connLimit, pgW, pgH)
		if err != nil {
			return nil, err
		}

		return pool, nil
	}, nil
}

type PGPool struct {
	// conns is connections in the pool
	conns []*websocket.Conn
	// Pool context
	cxt context.Context
	// cancel stops all pool goroutines
	cancel context.CancelFunc
	// startStream starts stream for passed websocket connection
	startStream StartStreamFunc
	// startPlayer starts new player
	startPlayer game.StartPlayerFunc
}

type errCannotCreatePool struct {
	err error
}

func (e *errCannotCreatePool) Error() string {
	return "Cannot create pool: " + e.err.Error()
}

func NewPGPool(cxt context.Context, connLimit uint8, pgW, pgH uint8,
) (*PGPool, error) {
	if err := cxt.Err(); err != nil {
		return nil, &errCannotCreatePool{err}
	}
	if connLimit == 0 {
		return nil, &errCannotCreatePool{ErrInvalidConnLimit}
	}

	// Pool context
	pcxt, cancel := context.WithCancel(cxt)

	// chStream common game data for all players in current pool
	chStream, startPlayerFunc, err := game.StartGame(pcxt, pgW, pgH)
	if err != nil {
		return nil, &errCannotCreatePool{err}
	}
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("Game was started")
	}

	startStreamFunc, err := StartStream(pcxt, chStream)
	if err != nil {
		return nil, &errCannotCreatePool{err}
	}
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("Stream was started")
	}

	return &PGPool{
		make([]*websocket.Conn, 0, connLimit),
		pcxt,
		cancel,
		startStreamFunc,
		startPlayerFunc,
	}, nil
}

// Implementing Pool interface
func (p *PGPool) IsFull() bool {
	return cap(p.conns) == len(p.conns)
}

// Implementing Pool interface
func (p *PGPool) IsEmpty() bool {
	return len(p.conns) == 0
}

type errCannotAddConnToPool struct {
	err error
}

func (e *errCannotAddConnToPool) Error() string {
	return "Cannot add connection to pool: " + e.err.Error()
}

// Implementing Pool interface
func (p *PGPool) AddConn(ws *websocket.Conn) (
	pwshandler.Environment, error) {
	if p.IsFull() {
		return nil, &errCannotAddConnToPool{
			errors.New("Pool is full"),
		}
	}
	if p.HasConn(ws) {
		return nil, &errCannotAddConnToPool{
			errors.New("Passed connection already added in pool"),
		}
	}

	p.conns = append(p.conns, ws)

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Connection was created to the pool")
	}

	return &PoolFeatures{p.startStream, p.startPlayer, p.cxt}, nil
}

// Implementing Pool interface
func (p *PGPool) DelConn(ws *websocket.Conn) error {
	for i := range p.conns {
		// Find connection
		if p.conns[i] == ws {
			// Remove connection
			p.conns = append(p.conns[:i], p.conns[i+1:]...)

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Connection was found and removed")
			}

			if p.IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("Pool is empty")
				}

				// Stop all pool goroutines
				p.cancel()

				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("Pool goroutines was canceled")
				}
			}

			return nil
		}
	}

	return errors.New("Cannot delete connection from pool: " +
		"connection was not found in pool")
}

// Implementing Pool interface
func (p *PGPool) HasConn(ws *websocket.Conn) bool {
	for i := range p.conns {
		if p.conns[i] == ws {
			return true
		}
	}

	return false
}
