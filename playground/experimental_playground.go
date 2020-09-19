package playground

import (
	"errors"
	"sync"

	"github.com/ivan1993spb/snake-server/engine"
)

const findRetriesNumber = 64

const (
	errRetriesLimitMessage               = "retries limit was reached"
	errAreaDoesNotContainLocationMessage = "area does not contain location"
	errEmptyLocationMessage              = "empty location"
)

type ExperimentalPlayground struct {
	gameMap *engine.Map

	objectsContainers    map[interface{}]*engine.Container
	objectsContainersMux *sync.RWMutex
}

func NewExperimentalPlayground(width, height uint8) (*ExperimentalPlayground, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	gameMap := engine.NewMap(area)

	return &ExperimentalPlayground{
		gameMap: gameMap,

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}, nil
}

func (p *ExperimentalPlayground) unsafeObjectExists(object interface{}) bool {
	_, ok := p.objectsContainers[object]
	return ok
}

func (p *ExperimentalPlayground) unsafeAddObject(object interface{}, container *engine.Container) error {
	if p.unsafeObjectExists(object) {
		return errors.New("cannot add object: object already exists")
	}

	p.objectsContainers[object] = container
	return nil
}

func (p *ExperimentalPlayground) addObject(object interface{}, container *engine.Container) error {
	p.objectsContainersMux.Lock()
	defer p.objectsContainersMux.Unlock()
	return p.unsafeAddObject(object, container)
}

func (p *ExperimentalPlayground) unsafeDeleteObject(object interface{}) error {
	if !p.unsafeObjectExists(object) {
		return errors.New("delete object error: object to delete not found")
	}

	delete(p.objectsContainers, object)
	return nil
}

func (p *ExperimentalPlayground) deleteObject(object interface{}) error {
	p.objectsContainersMux.Lock()
	defer p.objectsContainersMux.Unlock()
	return p.unsafeDeleteObject(object)
}

func (p *ExperimentalPlayground) unsafeGetContainerByObject(object interface{}) (*engine.Container, error) {
	container, ok := p.objectsContainers[object]
	if !ok {
		return nil, errors.New("get container: object was not found")
	}
	return container, nil
}

func (p *ExperimentalPlayground) getContainerByObject(object interface{}) (*engine.Container, error) {
	p.objectsContainersMux.RLock()
	defer p.objectsContainersMux.RUnlock()
	return p.unsafeGetContainerByObject(object)
}

func (p *ExperimentalPlayground) GetObjectByDot(dot engine.Dot) interface{} {
	if object, ok := p.gameMap.Get(dot); ok {
		return object.GetObject()
	}
	return nil
}

func (p *ExperimentalPlayground) GetObjectsByDots(dots []engine.Dot) []interface{} {
	if len(dots) == 0 {
		return nil
	}

	objects := make([]interface{}, 0)

	for _, container := range p.gameMap.MGet(dots) {
		flagObjectCreated := false
		object := container.GetObject()

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

func (p *ExperimentalPlayground) CreateObject(object interface{}, location engine.Location) error {
	if location.Empty() {
		return errCreateObject(errEmptyLocationMessage)
	}

	if !p.gameMap.Area().ContainsLocation(location) {
		return errCreateObject(errAreaDoesNotContainLocationMessage)
	}

	container := engine.NewContainer(object)

	if !p.gameMap.MSetIfAllAbsent(location, container) {
		return errCreateObject("location is occupied")
	}

	if err := p.addObject(object, container); err != nil {
		// Rollback map if cannot add object.
		p.gameMap.MRemove(location)

		return errCreateObject(err.Error())
	}

	return nil
}

func (p *ExperimentalPlayground) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error) {
	if location.Empty() {
		return nil, errCreateObjectAvailableDots(errEmptyLocationMessage)
	}

	if !p.gameMap.Area().ContainsLocation(location) {
		return nil, errCreateObjectAvailableDots(errAreaDoesNotContainLocationMessage)
	}

	container := engine.NewContainer(object)
	resultLocation := p.gameMap.MSetIfAbsent(location, container)

	if len(resultLocation) == 0 {
		return nil, errCreateObjectAvailableDots("all dots in location are occupied")
	}

	if err := p.addObject(object, container); err != nil {
		// Rollback map if cannot add object.
		p.gameMap.MRemove(resultLocation)

		return nil, errCreateObjectAvailableDots(err.Error())
	}

	return resultLocation, nil
}

func (p *ExperimentalPlayground) DeleteObject(object interface{}, location engine.Location) error {
	if !location.Empty() {
		container, err := p.getContainerByObject(object)
		if err != nil {
			return errDeleteObject(err.Error())
		}
		p.gameMap.MRemoveContainer(location, container)
	}

	if err := p.deleteObject(object); err != nil {
		return errDeleteObject(err.Error())
	}

	return nil
}

func (p *ExperimentalPlayground) UpdateObject(object interface{}, old, new engine.Location) error {
	diff := old.Difference(new)

	dotsToRemove := make([]engine.Dot, 0, len(diff))
	dotsToSet := make([]engine.Dot, 0, len(diff))

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet = append(dotsToSet, dot)
		} else {
			dotsToRemove = append(dotsToRemove, dot)
		}
	}

	container, err := p.getContainerByObject(object)
	if err != nil {
		return errUpdateObject(err.Error())
	}

	if !p.gameMap.MSetIfAllAbsent(dotsToSet, container) {
		return errUpdateObject("cannot occupy new location")
	}

	p.gameMap.MRemoveContainer(dotsToRemove, container)

	return nil
}

func (p *ExperimentalPlayground) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error) {
	actualLocation := old.Copy()
	diff := old.Difference(new)

	dotsToRemove := make([]engine.Dot, 0, len(diff))
	dotsToSet := make([]engine.Dot, 0, len(diff))

	for _, dot := range diff {
		if new.Contains(dot) {
			dotsToSet = append(dotsToSet, dot)
		} else {
			dotsToRemove = append(dotsToRemove, dot)
		}
	}

	container, err := p.getContainerByObject(object)
	if err != nil {
		return nil, errUpdateObjectAvailableDots(err.Error())
	}

	if len(dotsToSet) > 0 {
		resultDots := p.gameMap.MSetIfAbsent(dotsToSet, container)
		if len(resultDots) > 0 {
			for _, dot := range resultDots {
				actualLocation = actualLocation.Add(dot)
			}
		}
	}

	if len(dotsToRemove) > 0 {
		p.gameMap.MRemoveContainer(dotsToRemove, container)
		for _, dot := range dotsToRemove {
			actualLocation = actualLocation.Delete(dot)
		}
	}

	if len(actualLocation) == 0 {
		return nil, errUpdateObjectAvailableDots("all dots to set are occupied")
	}

	return actualLocation, nil
}

func (p *ExperimentalPlayground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	container := engine.NewContainer(object)

	for i := 0; i < findRetriesNumber; i++ {
		dot := p.gameMap.Area().NewRandomDot(0, 0)

		if p.gameMap.SetIfAbsent(dot, container) {
			if err := p.addObject(object, container); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.Remove(dot)

				return nil, errCreateObjectRandomDot(err.Error())
			}

			return engine.Location{dot}, nil
		}
	}

	return nil, errCreateObjectRandomDot(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRect("invalid rect size: 0")
	}

	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, rw, rh)) {
		return nil, errCreateObjectRandomRect("area cannot contain located rect")
	}

	container := engine.NewContainer(object)

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(rw, rh, 0, 0)
		if err != nil {
			continue
		}
		location := rect.Location()

		if p.gameMap.MSetIfAllAbsent(location, container) {
			if err := p.addObject(object, container); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemove(location)

				return nil, errCreateObjectRandomRect(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRect(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomRectMargin(object interface{}, rw, rh, margin uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRectMargin("invalid rect size: 0")
	}

	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, rw+margin*2, rh+margin*2)) {
		return nil, errCreateObjectRandomRectMargin("area cannot contain located rect with margin")
	}

	container := engine.NewContainer(object)

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(rw+margin*2, rh+margin*2, 0, 0)
		if err != nil {
			continue
		}

		if p.gameMap.HasAny(rect.Location()) {
			continue
		}

		location := engine.NewRect(rect.X()+margin, rect.Y()+margin, rw, rh).Location()

		if p.gameMap.MSetIfAllAbsent(location, container) {
			if err := p.addObject(object, container); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemoveContainer(location, container)

				return nil, errCreateObjectRandomRectMargin(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRectMargin(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomByDotsMask(object interface{}, dm *engine.DotsMask) (engine.Location, error) {
	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, dm.Width(), dm.Height())) {
		return nil, errCreateObjectRandomByDotsMask("area cannot contain located by dots mask object")
	}

	container := engine.NewContainer(object)

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(dm.Width(), dm.Height(), 0, 0)
		if err != nil {
			continue
		}

		location := dm.Location(rect.X(), rect.Y())

		if p.gameMap.HasAny(location) {
			continue
		}

		if p.gameMap.MSetIfAllAbsent(location, container) {
			if err := p.addObject(object, container); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemoveContainer(location, container)

				return nil, errCreateObjectRandomByDotsMask(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomByDotsMask(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) LocationOccupied(location engine.Location) bool {
	return p.gameMap.HasAll(location)
}

func (p *ExperimentalPlayground) Area() engine.Area {
	return p.gameMap.Area()
}

func (p *ExperimentalPlayground) unsafeGetObjects() []interface{} {
	objects := make([]interface{}, 0, len(p.objectsContainers))
	for object := range p.objectsContainers {
		objects = append(objects, object)
	}
	return objects
}

func (p *ExperimentalPlayground) GetObjects() []interface{} {
	p.objectsContainersMux.RLock()
	defer p.objectsContainersMux.RUnlock()
	return p.unsafeGetObjects()
}
