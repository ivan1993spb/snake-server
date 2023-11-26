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
)

// ExperimentalPlayground is a framework which allows locating for game objects
type ExperimentalPlayground struct {
	gameMap *engine.Map

	// objectsContainers is a mapping between game objects and their containers.
	// If an object and its container are registered in this map, the object must
	// be presented on the playground.
	objectsContainers    map[engine.Object]*engine.Container
	objectsContainersMux *sync.RWMutex
}

type errCreatePlayground string

func (e errCreatePlayground) Error() string {
	return "failed to create playground: " + string(e)
}

// NewExperimentalPlayground creates a new empty playground of the specified area
func NewExperimentalPlayground(width, height uint8) (*ExperimentalPlayground, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, errCreatePlayground(err.Error())
	}

	gameMap := engine.NewMap(area)

	return &ExperimentalPlayground{
		gameMap: gameMap,

		objectsContainers:    make(map[engine.Object]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}, nil
}

// unsafeObjectExists returns true if a game object has been registered in the playground
func (p *ExperimentalPlayground) unsafeObjectExists(object engine.Object) bool {
	_, ok := p.objectsContainers[object]
	return ok
}

// unsafeAddObject registers unsafely an object and its container in the mapping
func (p *ExperimentalPlayground) unsafeAddObject(object engine.Object, container *engine.Container) error {
	if p.unsafeObjectExists(object) {
		return errors.New("cannot add object: object already exists")
	}

	p.objectsContainers[object] = container
	return nil
}

// addObject registers safely an object and its container in the mapping
func (p *ExperimentalPlayground) addObject(object engine.Object, container *engine.Container) error {
	p.objectsContainersMux.Lock()
	defer p.objectsContainersMux.Unlock()
	return p.unsafeAddObject(object, container)
}

// unsafeDeleteObject discards unsafely an object from the mapping and so from the playground
func (p *ExperimentalPlayground) unsafeDeleteObject(object engine.Object) error {
	if !p.unsafeObjectExists(object) {
		return errors.New("delete object error: object to delete not found")
	}

	delete(p.objectsContainers, object)
	return nil
}

// deleteObject discards safely an object from the mapping and so from the playground
func (p *ExperimentalPlayground) deleteObject(object engine.Object) error {
	p.objectsContainersMux.Lock()
	defer p.objectsContainersMux.Unlock()
	return p.unsafeDeleteObject(object)
}

// unsafeGetContainerByObject looks unsafely for the container of a specified object in the mapping of objects
func (p *ExperimentalPlayground) unsafeGetContainerByObject(object engine.Object) (*engine.Container, error) {
	container, ok := p.objectsContainers[object]
	if !ok {
		return nil, errors.New("get container: object was not found")
	}
	return container, nil
}

// getContainerByObject looks safely for the container of a specified object in the mapping of objects
func (p *ExperimentalPlayground) getContainerByObject(object engine.Object) (*engine.Container, error) {
	p.objectsContainersMux.RLock()
	defer p.objectsContainersMux.RUnlock()
	return p.unsafeGetContainerByObject(object)
}

// GetObjectByDot returns an object by the given dot
func (p *ExperimentalPlayground) GetObjectByDot(dot engine.Dot) engine.Object {
	if object, ok := p.gameMap.Get(dot); ok {
		return object.GetObject()
	}
	return nil
}

// GetObjectByDot returns a list of objects located at the given dots
func (p *ExperimentalPlayground) GetObjectsByDots(dots []engine.Dot) []engine.Object {
	if len(dots) == 0 {
		return nil
	}

	objects := make([]engine.Object, 0)

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

type errCreateObject string

func (e errCreateObject) Error() string {
	return "failed to create object: " + string(e)
}

// CreateObject creates and registers an object at the given location on the playground.
// If some dots are occupied by other objects, the operation will be turn down with an error.
// Initial location could be empty.
func (p *ExperimentalPlayground) CreateObject(object engine.Object, location engine.Location) error {
	if !p.gameMap.Area().ContainsLocation(location) {
		return errCreateObject(errAreaDoesNotContainLocationMessage)
	}

	container := engine.NewContainer(object)

	if !location.Empty() && !p.gameMap.MSetIfAllVacant(location, container) {
		return errCreateObject("location is occupied")
	}

	if err := p.addObject(object, container); err != nil {
		// An object has been successfully placed but hasn't been registered in the playground.
		// Hence, roll the map back and return the error.
		p.gameMap.MRemove(location)

		return errCreateObject(err.Error())
	}

	return nil
}

type errCreateObjectAvailableDots string

func (e errCreateObjectAvailableDots) Error() string {
	return "failed to create object available dots: " + string(e)
}

// CreateObjectAvailableDots creates and registers an object at the given location on the playground.
// If some dots are occupied by other objects, the dots will be ignored. If all dots are occupied
// the object will be registered without location and no error will be returned.
func (p *ExperimentalPlayground) CreateObjectAvailableDots(object engine.Object, location engine.Location) (engine.Location, error) {
	if !p.gameMap.Area().ContainsLocation(location) {
		return nil, errCreateObjectAvailableDots(errAreaDoesNotContainLocationMessage)
	}

	container := engine.NewContainer(object)
	resultLocation := p.gameMap.MSetIfVacant(location, container)

	if err := p.addObject(object, container); err != nil {
		// An object has been successfully placed but hasn't been registered in the playground.
		// Hence, roll the map back and return the error.
		p.gameMap.MRemove(resultLocation)

		return nil, errCreateObjectAvailableDots(err.Error())
	}

	return resultLocation, nil
}

type errDeleteObject string

func (e errDeleteObject) Error() string {
	return "failed to delete object: " + string(e)
}

// DeleteObject deletes an object with the given location from the playground
func (p *ExperimentalPlayground) DeleteObject(object engine.Object, location engine.Location) error {
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

type errUpdateObject string

func (e errUpdateObject) Error() string {
	return "failed to update object: " + string(e)
}

// UpdateObject updates the object's location. All dots of the new location must be vacant.
func (p *ExperimentalPlayground) UpdateObject(object engine.Object, old, new engine.Location) error {
	diff := old.Difference(new)

	// Nothing has changed
	if len(diff) == 0 {
		return nil
	}

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

	if !p.gameMap.MSetIfAllVacant(dotsToSet, container) {
		return errUpdateObject("cannot occupy new location")
	}

	p.gameMap.MRemoveContainer(dotsToRemove, container)

	return nil
}

type errUpdateObjectAvailableDots string

func (e errUpdateObjectAvailableDots) Error() string {
	return "failed to update object available dots: " + string(e)
}

// UpdateObjectAvailableDots updates the object's location. If some dots of the new location are
// occupied by other objects, the dots will be skipped.
func (p *ExperimentalPlayground) UpdateObjectAvailableDots(object engine.Object, old, new engine.Location) (engine.Location, error) {
	actualLocation := old.Copy()
	diff := old.Difference(new)

	// Nothing has changed
	if len(diff) == 0 {
		return actualLocation, nil
	}

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
		resultDots := p.gameMap.MSetIfVacant(dotsToSet, container)
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

type errCreateObjectRandomDot string

func (e errCreateObjectRandomDot) Error() string {
	return "failed to create object random dot: " + string(e)
}

// CreateObjectRandomDot creates and registers an object to the playground at a random dot
// which will be returned.
func (p *ExperimentalPlayground) CreateObjectRandomDot(object engine.Object) (engine.Location, error) {
	container := engine.NewContainer(object)

	for i := 0; i < findRetriesNumber; i++ {
		dot := p.gameMap.Area().NewRandomDot(0, 0)

		if p.gameMap.SetIfVacant(dot, container) {
			if err := p.addObject(object, container); err != nil {
				// An object has been successfully placed but hasn't been registered in the playground.
				// Hence, roll the map back and return the error.
				p.gameMap.Remove(dot)

				return nil, errCreateObjectRandomDot(err.Error())
			}

			return engine.Location{dot}, nil
		}
	}

	return nil, errCreateObjectRandomDot(errRetriesLimitMessage)
}

type errCreateObjectRandomRect string

func (e errCreateObjectRandomRect) Error() string {
	return "failed create object random rect: " + string(e)
}

// CreateObjectRandomRect creates and registers an object to the playground at a random location
// of rectangle shape with the given size.
func (p *ExperimentalPlayground) CreateObjectRandomRect(object engine.Object, rw, rh uint8) (engine.Location, error) {
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

		if p.gameMap.MSetIfAllVacant(location, container) {
			if err := p.addObject(object, container); err != nil {
				// An object has been successfully placed but hasn't been registered in the playground.
				// Hence, roll the map back and return the error.
				p.gameMap.MRemove(location)

				return nil, errCreateObjectRandomRect(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRect(errRetriesLimitMessage)
}

type errCreateObjectRandomRectMargin string

func (e errCreateObjectRandomRectMargin) Error() string {
	return "failed to create object random rect with margin: " + string(e)
}

// CreateObjectRandomRectMargin creates and registers an object to the playground at a random
// rectangle location with the given size in at least X (=margin) dots apart from the other
// objects on the playground.
func (p *ExperimentalPlayground) CreateObjectRandomRectMargin(object engine.Object, rw, rh, margin uint8) (engine.Location, error) {
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

		if p.gameMap.MSetIfAllVacant(location, container) {
			if err := p.addObject(object, container); err != nil {
				// An object has been successfully placed but hasn't been registered in the playground.
				// Hence, roll the map back and return the error.
				p.gameMap.MRemoveContainer(location, container)

				return nil, errCreateObjectRandomRectMargin(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomRectMargin(errRetriesLimitMessage)
}

type errCreateObjectRandomByDotsMask string

func (e errCreateObjectRandomByDotsMask) Error() string {
	return "failed to create object random by dots mask: " + string(e)
}

// CreateObjectRandomByDotsMask creates and registers an object to the playground at a random
// location shaped in form of the given mask dm.
func (p *ExperimentalPlayground) CreateObjectRandomByDotsMask(object engine.Object, dm *engine.DotsMask) (engine.Location, error) {
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

		if p.gameMap.MSetIfAllVacant(location, container) {
			if err := p.addObject(object, container); err != nil {
				// An object has been successfully placed but hasn't been registered in the playground.
				// Hence, roll the map back and return the error.
				p.gameMap.MRemoveContainer(location, container)

				return nil, errCreateObjectRandomByDotsMask(err.Error())
			}

			return location, nil
		}
	}

	return nil, errCreateObjectRandomByDotsMask(errRetriesLimitMessage)
}

// LocationOccupied returns true if the location is fully occupied by objects on the
// playground
func (p *ExperimentalPlayground) LocationOccupied(location engine.Location) bool {
	return p.gameMap.HasAll(location)
}

// Area returns the playground's area object
func (p *ExperimentalPlayground) Area() engine.Area {
	return p.gameMap.Area()
}

// unsafeGetObjects collects unsafely and returns all the objects registered at the
// playground
func (p *ExperimentalPlayground) unsafeGetObjects() []engine.Object {
	objects := make([]engine.Object, 0, len(p.objectsContainers))
	for object := range p.objectsContainers {
		objects = append(objects, object)
	}
	return objects
}

// GetObjects collects and returns all the objects registered at the playground
func (p *ExperimentalPlayground) GetObjects() []engine.Object {
	p.objectsContainersMux.RLock()
	defer p.objectsContainersMux.RUnlock()
	return p.unsafeGetObjects()
}
