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

	objectsBuffer    []interface{}
	objectsBufferMux *sync.RWMutex

	area engine.Area
}

type ErrCreatePlayground struct {
	Err error
}

func (e ErrCreatePlayground) Error() string {
	return "cannot create playground: " + e.Err.Error()
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
		cMap:             cMap,
		entities:         make([]*entity, 0),
		entitiesMux:      &sync.RWMutex{},
		objectsBuffer:    make([]interface{}, 0),
		objectsBufferMux: &sync.RWMutex{},
		area:             area,
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

func (pg *Playground) unsafeAddEntity(e *entity) error {
	if pg.unsafeObjectExists(e.GetObject()) {
		return errors.New("cannot add entity: object already exists")
	}

	pg.entities = append(pg.entities, e)
	return nil
}

func (pg *Playground) addEntity(e *entity) error {
	pg.entitiesMux.Lock()
	defer pg.entitiesMux.Unlock()
	return pg.unsafeAddEntity(e)
}

func (pg *Playground) unsafeDeleteEntity(e *entity) {
	for i := range pg.entities {
		if pg.entities[i] == e {
			pg.entities = append(pg.entities[:i], pg.entities[i+1:]...)
			break
		}
	}
}

func (pg *Playground) unsafeObjectBuffered(object interface{}) bool {
	for i := range pg.objectsBuffer {
		if pg.objectsBuffer[i] == object {
			return true
		}
	}
	return false
}

func (pg *Playground) bufferAddObject(object interface{}) error {
	pg.objectsBufferMux.Lock()
	defer pg.objectsBufferMux.Unlock()

	if pg.ObjectExists(object) {
		return errors.New("cannot buffer object: object already exists in entities")
	}

	if pg.unsafeObjectBuffered(object) {
		return errors.New("object already buffered")
	}

	pg.objectsBuffer = append(pg.objectsBuffer, object)

	return nil
}

func (pg *Playground) bufferDeleteObject(object interface{}) {
	pg.objectsBufferMux.Lock()
	defer pg.objectsBufferMux.Unlock()

	for i := range pg.objectsBuffer {
		if pg.objectsBuffer[i] == object {
			pg.objectsBuffer = append(pg.objectsBuffer[:i], pg.objectsBuffer[i+1:]...)
			break
		}
	}
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

type errCreateObjectRandomDot string

func (e errCreateObjectRandomDot) Error() string {
	return "error create object random dot: " + string(e)
}

func (pg *Playground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	if err := pg.bufferAddObject(object); err != nil {
		return nil, errCreateObjectRandomDot(err.Error())
	}
	defer pg.bufferDeleteObject(object)

	e := &entity{
		object: object,
	}

	for i := 0; i < FindRetriesNumber; i++ {
		dot := pg.area.NewRandomDot(0, 0)
		e.location = engine.Location{dot}

		if pg.cMap.SetIfAbsent(dot.Hash(), e) {
			if err := pg.addEntity(e); err != nil {
				return nil, errCreateObjectRandomDot(err.Error())
			}

			return engine.Location{dot}, nil
		}
	}

	return nil, errCreateObjectRandomDot(ErrRetriesLimit.Error())
}

type errCreateObjectRandomRect string

func (e errCreateObjectRandomRect) Error() string {
	return "error create object random rect: " + string(e)
}

func (pg *Playground) CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRect("invalid rect size: 0")
	}

	if !pg.area.ContainsRect(engine.NewRect(0, 0, rw, rh)) {
		return nil, errCreateObjectRandomRect("area cannot contain located rect")
	}

	if err := pg.bufferAddObject(object); err != nil {
		return nil, errCreateObjectRandomRect(err.Error())
	}
	defer pg.bufferDeleteObject(object)

	e := &entity{
		object: object,
	}

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(rw, rh, 0, 0)
		if err != nil {
			continue
		}
		e.location = rect.Location()

		if pg.cMap.MSetIfAllAbsent(e.GetPreparedMap()) {
			if err := pg.addEntity(e); err != nil {
				return nil, errCreateObjectRandomRect(err.Error())
			}

			return e.GetLocation(), nil
		}
	}

	return nil, errCreateObjectRandomRect(ErrRetriesLimit.Error())
}

type errCreateObjectRandomRectMargin string

func (e errCreateObjectRandomRectMargin) Error() string {
	return "error create object random rect with margin: " + string(e)
}

func (pg *Playground) CreateObjectRandomRectMargin(object interface{}, rw, rh, margin uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRectMargin("invalid rect size: 0")
	}

	if !pg.area.ContainsRect(engine.NewRect(0, 0, rw+margin*2, rh+margin*2)) {
		return nil, errCreateObjectRandomRectMargin("area cannot contain located rect with margin")
	}

	if err := pg.bufferAddObject(object); err != nil {
		return nil, errCreateObjectRandomRectMargin(err.Error())
	}
	defer pg.bufferDeleteObject(object)

	e := &entity{
		object: object,
	}

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(rw+margin*2, rh+margin*2, 0, 0)
		if err != nil {
			continue
		}

		if pg.cMap.HasAny(rect.Location().Hash()) {
			continue
		}

		e.location = engine.NewRect(rect.X()+margin, rect.Y()+margin, rw, rh).Location()

		if pg.cMap.MSetIfAllAbsent(e.GetPreparedMap()) {
			if err := pg.addEntity(e); err != nil {
				return nil, errCreateObjectRandomRectMargin(err.Error())
			}

			return e.GetLocation(), nil
		}
	}

	return nil, errCreateObjectRandomRectMargin(ErrRetriesLimit.Error())
}

type errCreateObjectRandomByDotsMask string

func (e errCreateObjectRandomByDotsMask) Error() string {
	return "error create object random by dots mask: " + string(e)
}

func (pg *Playground) CreateObjectRandomByDotsMask(object interface{}, dm *engine.DotsMask) (engine.Location, error) {
	if !pg.area.ContainsRect(engine.NewRect(0, 0, dm.Width(), dm.Height())) {
		return nil, errCreateObjectRandomByDotsMask("area cannot contain located by dots mask object")
	}

	if err := pg.bufferAddObject(object); err != nil {
		return nil, errCreateObjectRandomByDotsMask(err.Error())
	}
	defer pg.bufferDeleteObject(object)

	e := &entity{
		object: object,
	}

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(dm.Width(), dm.Height(), 0, 0)
		if err != nil {
			continue
		}

		location := dm.Location(rect.X(), rect.Y())

		if pg.cMap.HasAny(location.Hash()) {
			continue
		}

		e.location = location

		if pg.cMap.MSetIfAllAbsent(e.GetPreparedMap()) {
			if err := pg.addEntity(e); err != nil {
				return nil, errCreateObjectRandomByDotsMask(err.Error())
			}

			return e.GetLocation(), nil
		}
	}

	return nil, errCreateObjectRandomByDotsMask(ErrRetriesLimit.Error())
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
