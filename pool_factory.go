package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/objects"
	"bitbucket.org/pushkin_ivan/clever-snake/playground"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func NewPGPoolFactory(rootCxt context.Context, connLimit,
	pgW, pgH uint8) (PoolFactory, error) {
	if err := rootCxt.Err(); err != nil {
		return nil, err
	}
	if connLimit == 0 {
		return nil, errors.New("Connection limit cannot be zero")
	}
	if pgW*pgH == 0 {
		return nil, playground.ErrInvalid_W_or_H
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
	conns []*websocket.Conn // Connection list

	// Goroutine management
	cxt    context.Context
	cancel context.CancelFunc

	pg *playground.Playground
}

func NewPGPool(cxt context.Context, connLimit uint8, pgW, pgH uint8,
) (*PGPool, error) {
	if err := cxt.Err(); err != nil {
		return nil, err
	}
	if connLimit == 0 {
		return nil, errors.New("Invalid connection limit")
	}
	if pgW*pgH == 0 {
		return nil, errors.New("Invalid playground size")
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                BEGIN INIT PLAYGROUND                *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	pg, err := playground.NewPlayground(pgW, pgH)
	if err != nil {
		return nil, err
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_POOLS) {
		glog.Infoln("Starting playground init")
	}

	// Create long wall to the playground
	if _, err := objects.CreateLongWall(pg); err != nil {
		return nil, err
	}

	// Create apple to the playground
	if _, err := objects.CreateApple(pg); err != nil {
		return nil, err
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                 END INIT PLAYGROUND                 *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	pcxt, cancel := context.WithCancel(cxt)

	return &PGPool{make([]*websocket.Conn, 0, connLimit), pcxt,
		cancel, pg}, nil
}

// Implementing Pool interface
func (p *PGPool) IsFull() bool {
	return len(p.conns) == cap(p.conns)
}

// Implementing Pool interface
func (p *PGPool) IsEmpty() bool {
	return len(p.conns) == 0
}

// Implementing Pool interface
func (p *PGPool) AddConn(ws *websocket.Conn) (
	pwshandler.Environment, error) {
	if p.IsFull() {
		return nil, errors.New("Pool is full")
	}

	p.conns = append(p.conns, ws)

	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Connection was created to pool")
	}

	return &GameData{p.cxt, p.pg}, nil
}

// Implementing Pool interface
func (p *PGPool) DelConn(ws *websocket.Conn) error {
	for i := range p.conns {
		// Find connection
		if p.conns[i] == ws {
			// Remove connection
			p.conns = append(p.conns[:i], p.conns[i+1:]...)

			if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
				glog.Infoln("Connection was found and removed")
			}

			// Stop all child goroutines if empty pool
			if p.IsEmpty() {
				if glog.V(INFOLOG_LEVEL_ABOUT_POOLS) {
					glog.Infoln("Pool is empty")
				}

				if p.cancel != nil {
					p.cancel()
					if glog.V(INFOLOG_LEVEL_ABOUT_POOLS) {
						glog.Infoln(
							"Pool goroutines was canceled")
					}
				} else {
					return errors.New(
						"CancelFunc of pool was not found")
				}
			}

			return nil
		}
	}

	return errors.New("Connection was not found in pool")
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
