package cplayground

import (
	"errors"
	"sync"

	"github.com/ivan1993spb/snake-server/concurrent-map"
	"github.com/ivan1993spb/snake-server/engine"
)

const defaultShardCount = 32

func calcShardCount(size uint16) int {
	// TODO: Implement function.
	return defaultShardCount
}

var FindRetriesNumber = 32

var ErrRetriesLimit = errors.New("retries limit was reached")

type Playground struct {
	cMap *cmap.ConcurrentMap

	entities    []*entity
	entitiesMux *sync.RWMutex

	area engine.Area
}

type ErrCreatePlayground struct {
	Err error
}

func (e ErrCreatePlayground) Error() string {
	return "cannot create playground: " + e.Error()
}

func NewPlayground(width, height uint8) (*Playground, error) {
	area, err := engine.NewArea(height, width)
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	cMap, err := cmap.New(calcShardCount(area.Size()))
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	return &Playground{
		cMap:        cMap,
		entities:    make([]*entity, 0),
		entitiesMux: &sync.RWMutex{},
		area:        area,
	}, nil
}

func (pg *Playground) unsafeObjectExists(object interface{}) bool {
	for i := range pg.entities {
		if pg.entities[i].GetObject() == object {
			return true
		}
	}
	return false
}

func (pg *Playground) ObjectExists(object interface{}) bool {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeObjectExists(object)
}

func (pg *Playground) unsafeLocationExists(location engine.Location) bool {
	for i := range pg.entities {
		if pg.entities[i].GetLocation().Equals(location) {
			return true
		}
	}
	return false
}

func (pg *Playground) LocationExists(location engine.Location) bool {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeLocationExists(location)
}

func (pg *Playground) unsafeEntityExists(object interface{}, location engine.Location) bool {
	for i := range pg.entities {
		if pg.entities[i].GetObject() == object && pg.entities[i].GetLocation().Equals(location) {
			return true
		}
	}
	return false
}

func (pg *Playground) EntityExists(object interface{}, location engine.Location) bool {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeEntityExists(object, location)
}

func (pg *Playground) unsafeGetObjectByLocation(location engine.Location) interface{} {
	for i := range pg.entities {
		if pg.entities[i].GetLocation().Equals(location) {
			return pg.entities[i].GetObject()
		}
	}
	return nil
}

func (pg *Playground) GetObjectByLocation(location engine.Location) interface{} {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeGetObjectByLocation(location)
}

func (pg *Playground) GetObjectByDot(dot engine.Dot) interface{} {
	if v, ok := pg.cMap.Get(dot.Hash()); ok {
		if e, ok := v.(*entity); ok {
			return e.GetObject()
		}
	}
	return nil
}

func (pg *Playground) GetEntityByDot(dot engine.Dot) (interface{}, engine.Location) {
	if v, ok := pg.cMap.Get(dot.Hash()); ok {
		if e, ok := v.(*entity); ok {
			return e.GetObject(), e.GetLocation()
		}
	}
	return nil, nil
}

func (pg *Playground) GetObjectsByDots(dots []engine.Dot) []interface{} {
	if len(dots) == 0 {
		return nil
	}

	objects := make([]interface{}, 0)
	checked := make([]engine.Dot, 0, len(dots))

	for _, dot := range dots {
		flagDotChecked := false

		for i := range checked {
			if checked[i].Equals(dot) {
				flagDotChecked = true
				break
			}
		}

		if flagDotChecked {
			continue
		}

		if v, ok := pg.cMap.Get(dot.Hash()); ok {
			if e, ok := v.(*entity); ok {
				object := e.GetObject()
				flagObjectCreated := false

				for i := range objects {
					if objects[i] == object {
						flagObjectCreated = true
						break
					}
				}

				if !flagObjectCreated {
					objects = append(objects, object)
				}
			}
		}

		checked = append(checked, dot)
	}

	return objects
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
	// TODO: Create specific error
	if pg.ObjectExists(object) {
		return nil, errors.New("object already exists")
	}

	for i := 0; i < FindRetriesNumber; i++ {
		dot := pg.area.NewRandomDot(0, 0)
		e := newEntity(object, engine.Location{dot})

		if pg.cMap.SetIfAbsent(dot.Hash(), e) {

			pg.entitiesMux.Lock()
			pg.entities = append(pg.entities, e)
			pg.entitiesMux.Unlock()

			return engine.Location{dot}, nil
		}
	}
	return nil, ErrRetriesLimit
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

func (pg *Playground) unsafeGetObjects() []interface{} {
	objects := make([]interface{}, len(pg.entities))
	for i, entity := range pg.entities {
		objects[i] = entity.GetObject()
	}
	return objects
}

func (pg *Playground) GetObjects() []interface{} {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeGetObjects()
}
