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

	objects    []*engine.Object
	objectsMux *sync.RWMutex
}

func NewExperimentalPlayground(width, height uint8) (*ExperimentalPlayground, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, ErrCreatePlayground{err}
	}

	gameMap := engine.NewMap(area)

	return &ExperimentalPlayground{
		gameMap:    gameMap,
		objects:    make([]*engine.Object, 0),
		objectsMux: &sync.RWMutex{},
	}, nil
}

func (p *ExperimentalPlayground) unsafeObjectExists(object *engine.Object) bool {
	for i := range p.objects {
		if p.objects[i] == object {
			return true
		}
	}
	return false
}

func (p *ExperimentalPlayground) unsafeAddObject(object *engine.Object) error {
	if p.unsafeObjectExists(object) {
		return errors.New("cannot add object: object already exists")
	}

	p.objects = append(p.objects, object)
	return nil
}

func (p *ExperimentalPlayground) addObject(object *engine.Object) error {
	p.objectsMux.Lock()
	defer p.objectsMux.Unlock()
	return p.unsafeAddObject(object)
}

func (p *ExperimentalPlayground) unsafeDeleteObject(object *engine.Object) error {
	for i := range p.objects {
		if p.objects[i] == object {
			p.objects = append(p.objects[:i], p.objects[i+1:]...)
			return nil
		}
	}
	return errors.New("delete object error: object to delete not found")
}

func (p *ExperimentalPlayground) deleteObject(object *engine.Object) error {
	p.objectsMux.Lock()
	defer p.objectsMux.Unlock()
	return p.unsafeDeleteObject(object)
}

func (p *ExperimentalPlayground) GetObjectByDot(dot engine.Dot) *engine.Object {
	if object, ok := p.gameMap.Get(dot); ok {
		return object
	}
	return nil
}

func (p *ExperimentalPlayground) GetObjectsByDots(dots []engine.Dot) []*engine.Object {
	if len(dots) == 0 {
		return nil
	}

	objects := make([]*engine.Object, 0)

	for _, object := range p.gameMap.MGet(dots) {
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

func (p *ExperimentalPlayground) CreateObject(object *engine.Object, location engine.Location) error {
	if location.Empty() {
		return errCreateObject(errEmptyLocationMessage)
	}

	if !p.gameMap.Area().ContainsLocation(location) {
		return errCreateObject(errAreaDoesNotContainLocationMessage)
	}

	if !p.gameMap.MSetIfAllAbsent(location, object) {
		return errCreateObject("location is occupied")
	}

	if err := p.addObject(object); err != nil {
		// Rollback map if cannot add object.
		p.gameMap.MRemove(location)

		return errCreateObject(err.Error())
	}

	return nil
}

func (p *ExperimentalPlayground) CreateObjectAvailableDots(object *engine.Object, location engine.Location) (engine.Location, error) {
	if location.Empty() {
		return nil, errCreateObjectAvailableDots(errEmptyLocationMessage)
	}

	if !p.gameMap.Area().ContainsLocation(location) {
		return nil, errCreateObjectAvailableDots(errAreaDoesNotContainLocationMessage)
	}

	resultLocation := p.gameMap.MSetIfAbsent(location, object)

	if len(resultLocation) == 0 {
		return nil, errCreateObjectAvailableDots("all dots in location are occupied")
	}

	if err := p.addObject(object); err != nil {
		// Rollback map if cannot add object.
		p.gameMap.MRemove(resultLocation)

		return nil, errCreateObjectAvailableDots(err.Error())
	}

	return resultLocation, nil
}

func (p *ExperimentalPlayground) DeleteObject(object *engine.Object, location engine.Location) error {
	if !location.Empty() {
		p.gameMap.MRemoveObject(location, object)
	}

	if err := p.deleteObject(object); err != nil {
		return errDeleteObject(err.Error())
	}

	return nil
}

func (p *ExperimentalPlayground) UpdateObject(object *engine.Object, old, new engine.Location) error {
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

	if !p.gameMap.MSetIfAllAbsent(dotsToSet, object) {
		return errUpdateObject("cannot occupy new location")
	}

	p.gameMap.MRemoveObject(dotsToRemove, object)

	return nil
}

func (p *ExperimentalPlayground) UpdateObjectAvailableDots(object *engine.Object, old, new engine.Location) (engine.Location, error) {
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

	if len(dotsToSet) > 0 {
		resultDots := p.gameMap.MSetIfAbsent(dotsToSet, object)
		if len(resultDots) > 0 {
			for _, dot := range resultDots {
				actualLocation = actualLocation.Add(dot)
			}
		}
	}

	if len(dotsToRemove) > 0 {
		p.gameMap.MRemoveObject(dotsToRemove, object)
		for _, dot := range dotsToRemove {
			actualLocation = actualLocation.Delete(dot)
		}
	}

	if len(actualLocation) == 0 {
		return nil, errUpdateObjectAvailableDots("all dots to set are occupied")
	}

	return actualLocation, nil
}

func (p *ExperimentalPlayground) CreateObjectRandomDot(object *engine.Object) (engine.Location, error) {
	for i := 0; i < findRetriesNumber; i++ {
		dot := p.gameMap.Area().NewRandomDot(0, 0)

		if p.gameMap.SetIfAbsent(dot, object) {
			if err := p.addObject(object); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.Remove(dot)

				return nil, errCreateObjectRandomDot(err.Error())
			}

			return engine.Location{dot}, nil
		}
	}

	return nil, errCreateObjectRandomDot(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomRect(object *engine.Object, rw, rh uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRect("invalid rect size: 0")
	}

	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, rw, rh)) {
		return nil, errCreateObjectRandomRect("area cannot contain located rect")
	}

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(rw, rh, 0, 0)
		if err != nil {
			continue
		}
		location := rect.Location()

		if p.gameMap.MSetIfAllAbsent(location, object) {
			if err := p.addObject(object); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemove(location)

				return nil, errCreateObjectRandomRect(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRect(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomRectMargin(object *engine.Object, rw, rh, margin uint8) (engine.Location, error) {
	if rw*rh == 0 {
		return nil, errCreateObjectRandomRectMargin("invalid rect size: 0")
	}

	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, rw+margin*2, rh+margin*2)) {
		return nil, errCreateObjectRandomRectMargin("area cannot contain located rect with margin")
	}

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(rw+margin*2, rh+margin*2, 0, 0)
		if err != nil {
			continue
		}

		if p.gameMap.HasAny(rect.Location()) {
			continue
		}

		location := engine.NewRect(rect.X()+margin, rect.Y()+margin, rw, rh).Location()

		if p.gameMap.MSetIfAllAbsent(location, object) {
			if err := p.addObject(object); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemoveObject(location, object)

				return nil, errCreateObjectRandomRectMargin(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRectMargin(errRetriesLimitMessage)
}

func (p *ExperimentalPlayground) CreateObjectRandomByDotsMask(object *engine.Object, dm *engine.DotsMask) (engine.Location, error) {
	if !p.gameMap.Area().ContainsRect(engine.NewRect(0, 0, dm.Width(), dm.Height())) {
		return nil, errCreateObjectRandomByDotsMask("area cannot contain located by dots mask object")
	}

	for i := 0; i < findRetriesNumber; i++ {
		rect, err := p.gameMap.Area().NewRandomRect(dm.Width(), dm.Height(), 0, 0)
		if err != nil {
			continue
		}

		location := dm.Location(rect.X(), rect.Y())

		if p.gameMap.HasAny(location) {
			continue
		}

		if p.gameMap.MSetIfAllAbsent(location, object) {
			if err := p.addObject(object); err != nil {
				// Rollback map if cannot add object.
				p.gameMap.MRemoveObject(location, object)

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

func (p *ExperimentalPlayground) unsafeGetObjects() []*engine.Object {
	objects := make([]*engine.Object, len(p.objects))
	copy(objects, p.objects)
	return objects
}

func (p *ExperimentalPlayground) GetObjects() []*engine.Object {
	p.objectsMux.RLock()
	defer p.objectsMux.RUnlock()
	return p.unsafeGetObjects()
}
