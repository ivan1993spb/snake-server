package cplayground

import (
	"sync"

	"github.com/ivan1993spb/snake-server/concurrent-map"
	"github.com/ivan1993spb/snake-server/engine"
)

type Playground struct {
	cMap *cmap.ConcurrentMap

	entities    []*entity
	entitiesMux *sync.RWMutex

	area engine.Area
}

func NewPlayground(width, height uint8) (*Playground, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) ObjectExists(object interface{}) bool {
	// TODO: Implement method.
	return false
}

func (pg *Playground) LocationExists(location engine.Location) bool {
	// TODO: Implement method.
	return false
}

func (pg *Playground) EntityExists(object interface{}, location engine.Location) bool {
	// TODO: Implement method.
	return false
}

func (pg *Playground) GetObjectByLocation(location engine.Location) interface{} {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) GetObjectByDot(dot engine.Dot) interface{} {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) GetEntityByDot(dot engine.Dot) (interface{}, engine.Location) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) GetObjectsByDots(dots []engine.Dot) []interface{} {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) CreateObject(object interface{}, location engine.Location) error {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) DeleteObject(object interface{}, location engine.Location) error {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) UpdateObject(object interface{}, old, new engine.Location) error {
	// TODO: Implement method.
	return nil
}

func (pg *Playground) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) CreateObjectRandomRectMargin(object interface{}, rw, rh, margin uint8) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) CreateObjectRandomByDotsMask(object interface{}, dm *engine.DotsMask) (engine.Location, error) {
	// TODO: Implement method.
	return nil, nil
}

func (pg *Playground) Navigate(dot engine.Dot, dir engine.Direction, dis uint8) (engine.Dot, error) {
	return pg.area.Navigate(dot, dir, dis)
}

func (pg *Playground) Size() uint16 {
	return pg.area.Size()
}

func (pg *Playground) Width() uint8 {
	return pg.area.Width()
}

func (pg *Playground) Height() uint8 {
	return pg.area.Height()
}

func (pg *Playground) GetObjects() []interface{} {
	// TODO: Implement method.
	return nil
}
