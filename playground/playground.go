package playground

import (
	"errors"
	"sync"

	"github.com/ivan1993spb/snake-server/concurrent-map"
	"github.com/ivan1993spb/snake-server/engine"
)

var FindRetriesNumber = 64

var ErrRetriesLimit = errors.New("retries limit was reached")

func prepareMap(object interface{}, location engine.Location) map[uint16]interface{} {
	m := make(map[uint16]interface{})
	for _, dot := range location {
		m[dot.Hash()] = object
	}
	return m
}

type Playground struct {
	cMap *cmap.ConcurrentMap

	objects    []interface{}
	objectsMux *sync.RWMutex

	area engine.Area
}

type ErrCreatePlayground struct {
	Err error
}

func (e ErrCreatePlayground) Error() string {
	return "cannot create playground: " + e.Err.Error()
}

func NewPlayground(width, height uint8) (*Playground, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	cMap, err := cmap.New(calcShardCount(area.Size()))
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	return &Playground{
		cMap:       cMap,
		objects:    make([]interface{}, 0),
		objectsMux: &sync.RWMutex{},
		area:       area,
	}, nil
}

func (pg *Playground) unsafeObjectExists(object interface{}) bool {
	for i := range pg.objects {
		if pg.objects[i] == object {
			return true
		}
	}
	return false
}

func (pg *Playground) unsafeAddObject(object interface{}) error {
	if pg.unsafeObjectExists(object) {
		return errors.New("cannot add object: object already exists")
	}

	pg.objects = append(pg.objects, object)
	return nil
}

func (pg *Playground) addObject(object interface{}) error {
	pg.objectsMux.Lock()
	defer pg.objectsMux.Unlock()
	return pg.unsafeAddObject(object)
}

func (pg *Playground) unsafeDeleteObject(object interface{}) error {
	for i := range pg.objects {
		if pg.objects[i] == object {
			pg.objects = append(pg.objects[:i], pg.objects[i+1:]...)
			return nil
		}
	}
	return errors.New("delete object error: object to delete not found")
}

func (pg *Playground) deleteObject(object interface{}) error {
	pg.objectsMux.Lock()
	defer pg.objectsMux.Unlock()
	return pg.unsafeDeleteObject(object)
}

func (pg *Playground) GetObjectByDot(dot engine.Dot) interface{} {
	if object, ok := pg.cMap.Get(dot.Hash()); ok {
		return object
	}
	return nil
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

	for _, object := range pg.cMap.MGet(keys) {
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

	return objects
}

type errCreateObject string

func (e errCreateObject) Error() string {
	return "error create object: " + string(e)
}

func (pg *Playground) CreateObject(object interface{}, location engine.Location) error {
	if location.Empty() {
		return errCreateObject("passed empty location")
	}

	if !pg.area.ContainsLocation(location) {
		return errCreateObject("area not contains location")
	}

	if !pg.cMap.MSetIfAllAbsent(prepareMap(object, location)) {
		return errCreateObject("location is occupied")
	}

	if err := pg.addObject(object); err != nil {
		// Rollback map if cannot add object.
		pg.cMap.MRemove(location.Hash())

		return errCreateObject(err.Error())
	}

	return nil
}

type errCreateObjectAvailableDots string

func (e errCreateObjectAvailableDots) Error() string {
	return "error create object available dots: " + string(e)
}

func (pg *Playground) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error) {
	if location.Empty() {
		return nil, errCreateObjectAvailableDots("passed empty location")
	}

	if !pg.area.ContainsLocation(location) {
		return nil, errCreateObjectAvailableDots("area not contains location")
	}

	hashes := pg.cMap.MSetIfAbsent(prepareMap(object, location))

	if len(hashes) == 0 {
		return nil, errCreateObjectAvailableDots("all dots in location are occupied")
	}

	resultLocation := engine.HashToLocation(hashes)

	if err := pg.addObject(object); err != nil {
		// Rollback map if cannot add object.
		pg.cMap.MRemove(resultLocation.Hash())

		return nil, errCreateObjectAvailableDots(err.Error())
	}

	return resultLocation, nil
}

type errDeleteObject string

func (e errDeleteObject) Error() string {
	return "error delete object: " + string(e)
}

func (pg *Playground) DeleteObject(object interface{}, location engine.Location) error {
	if !location.Empty() {
		pg.cMap.MRemoveCb(location.Hash(), func(key uint16, v interface{}, exists bool) bool {
			return exists && v == object
		})
	}

	if err := pg.deleteObject(object); err != nil {
		return errDeleteObject(err.Error())
	}

	return nil
}

type errUpdateObject string

func (e errUpdateObject) Error() string {
	return "error update object: " + string(e)
}

func (pg *Playground) UpdateObject(object interface{}, old, new engine.Location) error {
	diff := old.Difference(new)

	keysToRemove := make([]uint16, len(diff))
	dotsToSet := make(map[uint16]interface{})

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet[dot.Hash()] = object
		} else {
			keysToRemove = append(keysToRemove, dot.Hash())
		}
	}

	if !pg.cMap.MSetIfAllAbsent(dotsToSet) {
		return errUpdateObject("cannot occupy new location")
	}

	pg.cMap.MRemoveCb(keysToRemove, func(key uint16, v interface{}, exists bool) bool {
		return exists && v == object
	})

	return nil
}

type errUpdateObjectAvailableDots string

func (e errUpdateObjectAvailableDots) Error() string {
	return "error update object available dots: " + string(e)
}

func (pg *Playground) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error) {
	actualLocation := old.Copy()
	diff := old.Difference(new)

	keysToRemove := make([]uint16, len(diff))
	dotsToSet := make(map[uint16]interface{})

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet[dot.Hash()] = object
		} else {
			keysToRemove = append(keysToRemove, dot.Hash())
		}
	}

	if len(dotsToSet) > 0 {
		hashes := pg.cMap.MSetIfAbsent(dotsToSet)
		if len(hashes) > 0 {
			for _, hash := range hashes {
				actualLocation = actualLocation.Add(engine.HashToDot(hash))
			}
		}
	}

	if len(keysToRemove) > 0 {
		pg.cMap.MRemoveCb(keysToRemove, func(key uint16, v interface{}, exists bool) bool {
			return exists && v == object
		})
		for _, key := range keysToRemove {
			actualLocation = actualLocation.Delete(engine.HashToDot(key))
		}
	}

	if len(actualLocation) == 0 {
		return nil, errUpdateObjectAvailableDots("all dots to set are occupied")
	}

	return actualLocation, nil
}

type errCreateObjectRandomDot string

func (e errCreateObjectRandomDot) Error() string {
	return "error create object random dot: " + string(e)
}

func (pg *Playground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	for i := 0; i < FindRetriesNumber; i++ {
		dot := pg.area.NewRandomDot(0, 0)

		if pg.cMap.SetIfAbsent(dot.Hash(), object) {
			if err := pg.addObject(object); err != nil {
				// Rollback map if cannot add object.
				pg.cMap.Remove(dot.Hash())

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

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(rw, rh, 0, 0)
		if err != nil {
			continue
		}
		location := rect.Location()

		if pg.cMap.MSetIfAllAbsent(prepareMap(object, location)) {
			if err := pg.addObject(object); err != nil {
				// Rollback map if cannot add object.
				pg.cMap.MRemove(location.Hash())

				return nil, errCreateObjectRandomRect(err.Error())
			}

			return location, nil
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

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(rw+margin*2, rh+margin*2, 0, 0)
		if err != nil {
			continue
		}

		if pg.cMap.HasAny(rect.Location().Hash()) {
			continue
		}

		location := engine.NewRect(rect.X()+margin, rect.Y()+margin, rw, rh).Location()

		if pg.cMap.MSetIfAllAbsent(prepareMap(object, location)) {
			if err := pg.addObject(object); err != nil {
				// Rollback map if cannot add object.
				pg.cMap.MRemove(location.Hash())

				return nil, errCreateObjectRandomRectMargin(err.Error())
			}

			return location, nil
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

	for i := 0; i < FindRetriesNumber; i++ {
		rect, err := pg.area.NewRandomRect(dm.Width(), dm.Height(), 0, 0)
		if err != nil {
			continue
		}

		location := dm.Location(rect.X(), rect.Y())

		if pg.cMap.HasAny(location.Hash()) {
			continue
		}

		if pg.cMap.MSetIfAllAbsent(prepareMap(object, location)) {
			if err := pg.addObject(object); err != nil {
				// Rollback map if cannot add object.
				pg.cMap.MRemove(location.Hash())

				return nil, errCreateObjectRandomByDotsMask(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomByDotsMask(ErrRetriesLimit.Error())
}

func (pg *Playground) LocationOccupied(location engine.Location) bool {
	return pg.cMap.HasAll(location.Hash())
}

func (pg *Playground) Width() uint8 {
	return pg.area.Width()
}

func (pg *Playground) Height() uint8 {
	return pg.area.Height()
}

func (pg *Playground) Area() engine.Area {
	return pg.area
}

func (pg *Playground) unsafeGetObjects() []interface{} {
	objects := make([]interface{}, len(pg.objects))
	copy(objects, pg.objects)
	return objects
}

func (pg *Playground) GetObjects() []interface{} {
	pg.objectsMux.RLock()
	defer pg.objectsMux.RUnlock()
	return pg.unsafeGetObjects()
}
