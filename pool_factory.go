package main

import (
	"errors"
	"fmt"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func NewPGPoolFactory(rootCxt context.Context, connLimit,
	pgW, pgH uint8) (PoolFactory, error) {
	if err := rootCxt.Err(); err != nil {
		return nil, fmt.Errorf("cannot create pool factory: %s", err)
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
	// stopPool stops all pool goroutines
	stopPool context.CancelFunc
	// startStreamConn starts stream for passed websocket connection
	startStreamConn StartStreamConnFunc
	// stopStreamConn stops stream for passed websocket connection
	stopStreamConn StopStreamConnFunc
	// startPlayer starts new player
	startPlayer game.StartPlayerFunc
}

type errCannotCreatePool struct {
	err error
}

func (e *errCannotCreatePool) Error() string {
	return "cannot create pool: " + e.err.Error()
}

func NewPGPool(cxt context.Context, connLimit uint8, pgW, pgH uint8,
) (*PGPool, error) {
	if err := cxt.Err(); err != nil {
		return nil, &errCannotCreatePool{err}
	}
	if connLimit == 0 {
		return nil, &errCannotCreatePool{
			errors.New("invalid connection limit"),
		}
	}

	// Pool context
	pcxt, cancel := context.WithCancel(cxt)

	// chStream common game channel for data of all players in pool
	chStream, startPlayer, err := game.StartGame(pcxt, pgW, pgH)
	if err != nil {
		return nil, &errCannotCreatePool{err}
	}
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("game was started")
	}

	startStreamConn, stopStreamConn := StartGameStream(chStream)
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("stream was started")
	}

	return &PGPool{
		make([]*websocket.Conn, 0, connLimit),
		cancel,
		startStreamConn,
		stopStreamConn,
		startPlayer,
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
	return "cannot add connection to pool: " + e.err.Error()
}

// Implementing Pool interface
func (p *PGPool) AddConn(ws *websocket.Conn) (
	pwshandler.Environment, error) {
	if p.IsFull() {
		return nil, &errCannotAddConnToPool{
			errors.New("pool is full"),
		}
	}
	if p.HasConn(ws) {
		return nil, &errCannotAddConnToPool{
			errors.New("passed connection already added in pool"),
		}
	}

	p.conns = append(p.conns, ws)

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("connection was created to pool")
	}

	return &PoolFeatures{
		p.startStreamConn,
		p.stopStreamConn,
		p.startPlayer,
	}, nil
}

// Implementing Pool interface
func (p *PGPool) DelConn(ws *websocket.Conn) error {
	for i := range p.conns {
		// Find connection
		if p.conns[i] == ws {
			// Remove connection
			p.conns = append(p.conns[:i], p.conns[i+1:]...)

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("connection was found and removed")
			}

			if p.IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("pool is empty")
				}

				// Stop all pool goroutines
				p.stopPool()

				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("pool goroutines was canceled")
				}
			}

			return nil
		}
	}

	return errors.New("cannot delete connection from pool: " +
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
