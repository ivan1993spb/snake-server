package main

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func NewPGPoolFactory(cxt context.Context, connLimit uint16,
	pgW, pgH uint8) (PoolFactory, error) {
	if err := cxt.Err(); err != nil {
		return nil, fmt.Errorf("cannot create pool factory: %s", err)
	}

	return func() (Pool, error) {
		pool, err := NewPGPool(cxt, connLimit, pgW, pgH)
		if err != nil {
			return nil, err
		}

		return pool, nil
	}, nil
}

type PGPool struct {
	// conns is connections in the pool
	conns map[uint16]*websocket.Conn
	// Max connection count per pool
	connLimit uint16
	// Pool context
	cxt context.Context
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

func NewPGPool(cxt context.Context, connLimit uint16, pgW, pgH uint8,
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

	startStreamConn, stopStreamConn := StartGameStream(pcxt, chStream)
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("stream was started")
	}

	return &PGPool{
		make(map[uint16]*websocket.Conn),
		connLimit,
		pcxt,
		cancel,
		startStreamConn,
		stopStreamConn,
		startPlayer,
	}, nil
}

// Implementing Pool interface
func (p *PGPool) IsFull() bool {
	return len(p.conns) == int(p.connLimit)
}

// Implementing Pool interface
func (p *PGPool) IsEmpty() bool {
	return len(p.conns) == 0
}

type PoolFeatures struct {
	startStreamConn StartStreamConnFunc
	stopStreamConn  StopStreamConnFunc
	// startPlayer starts player
	startPlayer game.StartPlayerFunc
	cxt         context.Context
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

	for id := uint16(0); int(id) <= len(p.conns); id++ {
		if _, occupied := p.conns[id]; !occupied {
			p.conns[id] = ws

			err := SendMessage(ws, HEADER_CONN_ID, id)
			if err != nil {
				return nil, &errCannotAddConn{
					fmt.Errorf("cannot send connection id: %s", err),
				}
			}

			break
		}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("connection was created to pool")
	}

	return &PoolFeatures{
		p.startStreamConn,
		p.stopStreamConn,
		p.startPlayer,
		p.cxt,
	}, nil
}

// Implementing Pool interface
func (p *PGPool) DelConn(ws *websocket.Conn) error {
	for id := range p.conns {
		// Find connection
		if p.conns[id] == ws {
			// Remove connection
			delete(p.conns, id)

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("connection was found and removed")
			}

			if p.IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("pool is empty")
				}

				if p.cxt.Err() == nil {
					p.stopPool()

					if glog.V(INFOLOG_LEVEL_POOLS) {
						glog.Infoln("pool goroutines was canceled")
					}
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
	for id := range p.conns {
		if p.conns[id] == ws {
			return true
		}
	}

	return false
}

// Implementing Pool interface
func (p *PGPool) ConnCount() uint16 {
	return uint16(len(p.conns))
}

// Implementing Pool interface
func (p *PGPool) ConnIds() []uint16 {
	var ids = make([]uint16, 0, len(p.conns))

	for id := range p.conns {
		ids = append(ids, id)
	}

	return ids
}

// Implementing Pool interface
func (p *PGPool) GetRequests() []*http.Request {
	var requests = make([]*http.Request, 0, len(p.conns))

	for _, ws := range p.conns {
		requests = append(requests, ws.Request())
	}

	return requests
}
