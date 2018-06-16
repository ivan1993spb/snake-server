package cplayground

import (
	"sync"

	"github.com/ivan1993spb/snake-server/engine"
)

type entity struct {
	object   interface{}
	location engine.Location
	mux      *sync.RWMutex
}

func newEntity(object interface{}, location engine.Location) *entity {
	return &entity{
		object:   object,
		location: location,
		mux:      &sync.RWMutex{},
	}
}

func (e *entity) GetLocation() engine.Location {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.location.Copy()
}

func (e *entity) GetObject() interface{} {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.object
}

func (e *entity) SetLocation(location engine.Location) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.location = location.Copy()
}

func (e *entity) GetPreparedMap() map[uint16]interface{} {
	e.mux.RLock()
	defer e.mux.RUnlock()

	m := make(map[uint16]interface{})
	for _, dot := range e.location {
		m[dot.Hash()] = e
	}

	return m
}
