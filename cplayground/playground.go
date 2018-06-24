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

func (pg *Playground) deleteEntity(e *entity) {
	pg.entitiesMux.Lock()
	pg.unsafeDeleteEntity(e)
	pg.entitiesMux.Unlock()
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

func (pg *Playground) unsafeGetEntityByObject(object interface{}) *entity {
	for i := range pg.entities {
		if pg.entities[i].GetObject() == object {
			return pg.entities[i]
		}
	}
	return nil
}

func (pg *Playground) getEntityByObject(object interface{}) *entity {
	pg.entitiesMux.RLock()
	defer pg.entitiesMux.RUnlock()
	return pg.unsafeGetEntityByObject(object)
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

	keys := make([]uint16, len(dots))
	for i, dot := range dots {
		keys[i] = dot.Hash()
	}

	objects := make([]interface{}, 0)

	for _, value := range pg.cMap.MGet(keys) {
		if e, ok := value.(*entity); ok {
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

	return objects
}

type errCreateObject string

func (e errCreateObject) Error() string {
	return "error create object: " + string(e)
}

func (pg *Playground) CreateObject(object interface{}, location engine.Location) error {
	if !pg.area.ContainsLocation(location) {
		return errCreateObject("area not contains location")
	}

	e := &entity{
		object:   object,
		location: location,
		mux:      &sync.RWMutex{},
	}

	if !pg.cMap.MSetIfAllAbsent(e.GetPreparedMap()) {
		return errCreateObject("location is occupied")
	}

	if err := pg.addEntity(e); err != nil {
		// Rollback map if cannot add entity.
		pg.cMap.MRemove(e.GetLocation().Hash())

		return errCreateObject(err.Error())
	}

	return nil
}

type errCreateObjectAvailableDots string

func (e errCreateObjectAvailableDots) Error() string {
	return "error create object available dots: " + string(e)
}

func (pg *Playground) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error) {
	if !pg.area.ContainsLocation(location) {
		return nil, errCreateObjectAvailableDots("area not contains location")
	}

	e := &entity{
		object:   object,
		location: location,
		mux:      &sync.RWMutex{},
	}

	hashes := pg.cMap.MSetIfAbsent(e.GetPreparedMap())

	if len(hashes) == 0 {
		return nil, errCreateObjectAvailableDots("dots in location are occupied")
	}

	e.SetLocation(engine.HashToLocation(hashes))

	if err := pg.addEntity(e); err != nil {
		// Rollback map if cannot add entity.
		pg.cMap.MRemove(e.GetLocation().Hash())

		return nil, errCreateObjectAvailableDots(err.Error())
	}

	return e.GetLocation(), nil
}

type errDeleteObject string

func (e errDeleteObject) Error() string {
	return "error delete object: " + string(e)
}

func (pg *Playground) DeleteObject(object interface{}, location engine.Location) error {
	e := pg.getEntityByObject(object)
	if e == nil {
		return errDeleteObject("cannot find entity by object")
	}

	pg.cMap.MRemove(e.GetLocation().Hash())

	pg.deleteEntity(e)

	return nil
}

type errUpdateObject string

func (e errUpdateObject) Error() string {
	return "error update object: " + string(e)
}

func (pg *Playground) UpdateObject(object interface{}, old, new engine.Location) error {
	e := pg.getEntityByObject(object)

	if e == nil {
		return errUpdateObject("cannot find entity by object")
	}

	actualLocation := e.GetLocation()
	diff := actualLocation.Difference(new)

	keysToRemove := make([]uint16, len(diff))
	dotsToSet := make(map[uint16]interface{})

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet[dot.Hash()] = e
		} else {
			keysToRemove = append(keysToRemove, dot.Hash())
		}
	}

	if !pg.cMap.MSetIfAllAbsent(dotsToSet) {
		return errUpdateObject("cannot occupy new location")
	}

	pg.cMap.MRemove(keysToRemove)

	e.SetLocation(new)

	return nil
}

type errUpdateObjectAvailableDots string

func (e errUpdateObjectAvailableDots) Error() string {
	return "error update object available dots: " + string(e)
}

func (pg *Playground) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error) {
	e := pg.getEntityByObject(object)

	if e == nil {
		return nil, errUpdateObjectAvailableDots("cannot find entity by object")
	}

	actualLocation := e.GetLocation()
	diff := actualLocation.Difference(new)

	keysToRemove := make([]uint16, len(diff))
	dotsToSet := make(map[uint16]interface{})

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet[dot.Hash()] = e
		} else {
			keysToRemove = append(keysToRemove, dot.Hash())
		}
	}

	hashes := pg.cMap.MSetIfAbsent(dotsToSet)
	if len(hashes) == 0 {
		return nil, errUpdateObjectAvailableDots("all dots to set are occupied")
	}

	pg.cMap.MRemove(keysToRemove)

	for _, key := range keysToRemove {
		actualLocation = actualLocation.Delete(engine.HashToDot(key))
	}

	for _, hash := range hashes {
		actualLocation = actualLocation.Add(engine.HashToDot(hash))
	}

	e.SetLocation(actualLocation)

	return e.GetLocation(), nil
}

type errCreateObjectRandomDot string

func (e errCreateObjectRandomDot) Error() string {
	return "error create object random dot: " + string(e)
}

func (pg *Playground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	e := &entity{
		object: object,
		mux:    &sync.RWMutex{},
	}

	for i := 0; i < FindRetriesNumber; i++ {
		dot := pg.area.NewRandomDot(0, 0)
		e.location = engine.Location{dot}

		if pg.cMap.SetIfAbsent(dot.Hash(), e) {
			if err := pg.addEntity(e); err != nil {
				// Rollback map if cannot add entity.
				pg.cMap.MRemove(e.GetLocation().Hash())

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

	e := &entity{
		object: object,
		mux:    &sync.RWMutex{},
	}

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(rw, rh, 0, 0)
		if err != nil {
			continue
		}
		e.location = rect.Location()

		if pg.cMap.MSetIfAllAbsent(e.GetPreparedMap()) {
			if err := pg.addEntity(e); err != nil {
				// Rollback map if cannot add entity.
				pg.cMap.MRemove(e.GetLocation().Hash())

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

	e := &entity{
		object: object,
		mux:    &sync.RWMutex{},
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
				// Rollback map if cannot add entity.
				pg.cMap.MRemove(e.GetLocation().Hash())

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

	e := &entity{
		object: object,
		mux:    &sync.RWMutex{},
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
				// Rollback map if cannot add entity.
				pg.cMap.MRemove(e.GetLocation().Hash())

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
