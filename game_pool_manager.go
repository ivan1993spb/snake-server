// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang/glog"
)

// Form key for passing pool id
const FORM_KEY_POOL_ID = "pool_id"

// PoolFactory must generate new pool
type PoolFactory func() (*GamePool, error)

type GamePoolManager struct {
	addPool   PoolFactory
	pools     map[uint16]*GamePool
	poolLimit uint16
}

// NewGamePoolManager creates new GamePoolManager with fixed max
// number of pools specified by poolLimit
func NewGamePoolManager(factory PoolFactory, poolLimit uint16,
) (*GamePoolManager, error) {
	if poolLimit == 0 {
		return nil, errors.New(
			"cannot create pool manager: invalid pool limit",
		)
	}

	return &GamePoolManager{
		factory,
		make(map[uint16]*GamePool),
		poolLimit,
	}, nil
}

type errCannotAddConn struct {
	err error
}

func (e *errCannotAddConn) Error() string {
	return "cannot add connection: " + e.err.Error()
}

type errCannotSendPoolId struct {
	err error
}

func (e *errCannotSendPoolId) Error() string {
	return "cannot send pool id: " + e.err.Error()
}

func (pm *GamePoolManager) AddConn(ww *WebsocketWrapper,
) (*PoolFeatures, error) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("try to add connection in a pool")
	}

	// Try to add connection in selected pool if passed pool id
	if id := ww.Request().FormValue(FORM_KEY_POOL_ID); len(id) > 0 {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received pool id")
		}

		if id, err := strconv.Atoi(id); err == nil && id > -1 {
			id := uint16(id)

			if pm.pools[id].IsFull() {
				return nil, &errCannotAddConn{
					errors.New("selected pool is full"),
				}
			}

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("creating connection in selected pool")
			}

			if err := ww.Send(HEADER_POOL_ID, id); err != nil {
				return nil, &errCannotAddConn{
					&errCannotSendPoolId{err},
				}
			}

			return pm.pools[id].AddConn(ww)
		}

		return nil, &errCannotAddConn{errors.New("invalid pool id")}
	}

	// Try to add connection in first not full pool

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("try to find not full pool")
	}
	// Try to find not full pool
	for id := range pm.pools {
		if !pm.pools[id].IsFull() {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("was found not full pool")
				glog.Infoln("creating connection to pool")
			}

			if err := ww.Send(HEADER_POOL_ID, id); err != nil {
				return nil, &errCannotAddConn{
					&errCannotSendPoolId{err},
				}
			}

			return pm.pools[id].AddConn(ww)
		}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("cannot find not full pool")
	}

	if len(pm.pools) == int(pm.poolLimit) {
		return nil, &errCannotAddConn{errors.New("server is full")}
	}

	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("server is not full")
		glog.Infoln("creating new pool")
	}

	pool, err := pm.addPool()
	if err != nil {
		return nil, &errCannotAddConn{err}
	}

	// Save the pool
	for id := uint16(0); int(id) <= len(pm.pools); id++ {
		if _, occupied := pm.pools[id]; !occupied {
			pm.pools[id] = pool

			if err := ww.Send(HEADER_POOL_ID, id); err != nil {
				return nil, &errCannotAddConn{
					&errCannotSendPoolId{err},
				}
			}

			break
		}
	}

	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("pool was created")
	}
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("creating connection to pool")
	}

	return pool.AddConn(ww)
}

type errCannotDelConn struct {
	err error
}

func (e *errCannotDelConn) Error() string {
	return "cannot delete connection: " + e.err.Error()
}

func (pm *GamePoolManager) DelConn(ww *WebsocketWrapper) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("try to remove information about connection")
		glog.Infoln("try to find pool of connection")
	}

	for id := range pm.pools {
		// If current pool has the connection...
		if pm.pools[id].HasConn(ww) {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("pool of connection was found")
				glog.Infoln("removing connection from pool")
			}

			if err := pm.pools[id].DelConn(ww); err != nil {
				return &errCannotDelConn{err}
			}

			if pm.pools[id].IsEmpty() {
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("removing empty pool")
				}

				// Delete pool
				delete(pm.pools, id)

				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("empty pool was removed")
				}
			}

			return nil
		}
	}

	return &errCannotDelConn{errors.New("connection was not found")}
}

type PoolInfo struct {
	PoolId    uint16 `json:"pool_id"`
	ConnCount uint16 `json:"conn_count"`
}

func (pm *GamePoolManager) PoolInfoList() []*PoolInfo {
	var info = make([]*PoolInfo, 0, len(pm.pools))

	for id, pool := range pm.pools {
		info = append(info, &PoolInfo{id, pool.ConnCount()})
	}

	return info
}

func (pm *GamePoolManager) PoolCount() uint16 {
	return uint16(len(pm.pools))
}

func (pm *GamePoolManager) ConnCount() (connCount uint32) {
	for i := range pm.pools {
		connCount += uint32(pm.pools[i].ConnCount())
	}

	return
}

func (pm *GamePoolManager) GetPool(id uint16) (*GamePool, error) {
	if pool, found := pm.pools[id]; found {
		return pool, nil
	}

	return nil, errors.New("cannot get pool: pool was not found")
}

func (pm *GamePoolManager) GetRequests() []*http.Request {
	var requests = make([]*http.Request, 0)

	for _, pool := range pm.pools {
		requests = append(requests, pool.GetRequests()...)
	}

	return requests
}
