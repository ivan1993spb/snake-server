// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// NewGamePoolFactory creates new game pool with passed params
func NewGamePoolFactory(cxt context.Context, connLimit uint16,
	pgW, pgH uint8) (PoolFactory, error) {
	if err := cxt.Err(); err != nil {
		return nil, fmt.Errorf("cannot create pool factory: %s", err)
	}

	return func() (*GamePool, error) {
		pool, err := NewGamePool(cxt, connLimit, pgW, pgH)
		if err != nil {
			return nil, err
		}

		return pool, nil
	}, nil
}

type GamePool struct {
	// conns is connections in the pool
	conns map[uint16]*WebsocketWrapper
	// Max connection count per pool
	connLimit uint16

	// Pool context
	cxt context.Context
	// cancel stops all pool goroutines
	cancel context.CancelFunc

	// stream is pool stream
	stream *Stream
	game   *game.Game
	// noticeC channel for pool notifications
	noticeC chan *OutputMessage
}

type errCannotCreatePool struct {
	err error
}

func (e *errCannotCreatePool) Error() string {
	return "cannot create pool: " + e.err.Error()
}

func NewGamePool(cxt context.Context, connLimit uint16,
	pgW, pgH uint8) (*GamePool, error) {
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

	// Pool game
	game, err := game.NewGame(pcxt, pgW, pgH)
	if err != nil {
		return nil, &errCannotCreatePool{err}
	}

	// Pool stream
	stream, err := NewStream(pcxt)
	if err != nil {
		return nil, &errCannotCreatePool{err}
	}

	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("stream was started")
	}

	// Create stream data sources:
	// Create channel for pool notifications
	noticeC := make(chan *OutputMessage)
	stream.AddSource(noticeC)
	// Create game data channel
	stream.AddSourceHeader(HEADER_GAME, game.StartGame())

	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("game was started")
	}

	return &GamePool{
		make(map[uint16]*WebsocketWrapper), connLimit,
		pcxt, cancel,
		stream, game, noticeC,
	}, nil
}

func (p *GamePool) Stop() {
	if p.cxt.Err() == nil {
		p.cancel()

		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("pool goroutines was canceled")
		}
	}

	// Close pool notification channel
	close(p.noticeC)

	// Stop pool data stream
	p.stream.Stop()

	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("pool data stream was stopped")
	}
}

// IsFull returns true if game pool is full
func (p *GamePool) IsFull() bool {
	return len(p.conns) == int(p.connLimit)
}

// IsEmpty returns true if game pool is empty
func (p *GamePool) IsEmpty() bool {
	return len(p.conns) == 0
}

type PoolFeatures struct {
	cxt         context.Context
	startPlayer game.StartPlayerFunc
	stream      *Stream
}

type errCannotAddConnToPool struct {
	err error
}

func (e *errCannotAddConnToPool) Error() string {
	return "cannot add connection to pool: " + e.err.Error()
}

// AddConn creates connection in game pool
func (p *GamePool) AddConn(ww *WebsocketWrapper,
) (*PoolFeatures, error) {
	if p.IsFull() {
		return nil, &errCannotAddConnToPool{
			errors.New("pool is full"),
		}
	}
	if p.HasConn(ww) {
		return nil, &errCannotAddConnToPool{
			errors.New("passed connection already added to pool"),
		}
	}

	for id := uint16(0); int(id) <= len(p.conns); id++ {
		if _, occupied := p.conns[id]; !occupied {
			p.conns[id] = ww

			if err := ww.Send(HEADER_CONN_ID, id); err != nil {
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

	p.Send(HEADER_INFO, "user created to pool")

	return &PoolFeatures{
		p.cxt,
		p.game.StartPlayer,
		p.stream,
	}, nil
}

// DelConn removes connection from game pool
func (p *GamePool) DelConn(ww *WebsocketWrapper) error {
	for id := range p.conns {
		// Find connection
		if p.conns[id] == ww {
			// Remove connection
			delete(p.conns, id)

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("connection was found and removed")
			}

			p.Send(HEADER_INFO, "user deleted from pool")

			if p.IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("pool is empty")
				}

				p.Stop()
			}

			return nil
		}
	}

	return errors.New("cannot delete connection from pool: connection was not found")
}

// HasConn returns true if passed connection belongs to current pool
func (p *GamePool) HasConn(ww *WebsocketWrapper) bool {
	for id := range p.conns {
		if p.conns[id] == ww {
			return true
		}
	}

	return false
}

// ConnCount returns connection count in game pool
func (p *GamePool) ConnCount() uint16 {
	return uint16(len(p.conns))
}

// ConnIds returns connection ids
func (p *GamePool) ConnIds() []uint16 {
	var ids = make([]uint16, 0, len(p.conns))

	for id := range p.conns {
		ids = append(ids, id)
	}

	return ids
}

// GetRequests returns requests
func (p *GamePool) GetRequests() []*http.Request {
	var requests = make([]*http.Request, 0, len(p.conns))

	for _, ww := range p.conns {
		requests = append(requests, ww.Request())
	}

	return requests
}

func (p *GamePool) Send(header string, data interface{}) {
	p.noticeC <- &OutputMessage{header, data}
}
