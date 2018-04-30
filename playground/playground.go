package playground

import (
	"errors"
	"sync"

	"github.com/ivan1993spb/clever-snake/engine"
)

type entity struct {
	object   interface{}
	location engine.Location
}

var ErrEmptyLocation = errors.New("empty location")

type Playground struct {
	scene         *engine.Scene
	entities      []entity
	entitiesMutex *sync.RWMutex
}

func NewPlayground(scene *engine.Scene) *Playground {
	return &Playground{
		scene:         scene,
		entities:      []entity{},
		entitiesMutex: &sync.RWMutex{},
	}
}

func (pg *Playground) unsafeObjectExists(object interface{}) bool {
	for i := range pg.entities {
		if pg.entities[i].object == object {
			return true
		}
	}
	return false
}

func (pg *Playground) ObjectExists(object interface{}) bool {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeObjectExists(object)
}

func (pg *Playground) unsafeLocationExists(location engine.Location) bool {
	for i := range pg.entities {
		if pg.entities[i].location.Equals(location) {
			return true
		}
	}
	return false
}

func (pg *Playground) LocationExists(location engine.Location) bool {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeLocationExists(location)
}

func (pg *Playground) unsafeEntityExists(object interface{}, location engine.Location) bool {
	for i := range pg.entities {
		if pg.entities[i].object == object && pg.entities[i].location.Equals(location) {
			return true
		}
	}
	return false
}

func (pg *Playground) EntityExists(object interface{}, location engine.Location) bool {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeEntityExists(object, location)
}

func (pg *Playground) unsafeGetObjectByLocation(location engine.Location) interface{} {
	for i := range pg.entities {
		if pg.entities[i].location.Equals(location) {
			return pg.entities[i].object
		}
	}
	return nil
}

func (pg *Playground) GetObjectByLocation(location engine.Location) interface{} {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeGetObjectByLocation(location)
}

func (pg *Playground) unsafeGetObjectByDot(dot *engine.Dot) interface{} {
	for i := range pg.entities {
		if pg.entities[i].location.Contains(dot) {
			return pg.entities[i].object
		}
	}
	return nil
}

func (pg *Playground) GetObjectByDot(dot *engine.Dot) interface{} {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeGetObjectByDot(dot)
}

func (pg *Playground) unsafeGetEntityByDot(dot *engine.Dot) (interface{}, engine.Location) {
	for i := range pg.entities {
		if pg.entities[i].location.Contains(dot) {
			return pg.entities[i].object, pg.entities[i].location
		}
	}
	return nil, nil
}

func (pg *Playground) GetEntityByDot(dot *engine.Dot) (interface{}, engine.Location) {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeGetEntityByDot(dot)
}

func (pg *Playground) unsafeGetObjectsByDots(dots []*engine.Dot) []interface{} {
	objects := make([]interface{}, 0)
	for _, dot := range dots {
		objects = append(objects, pg.unsafeGetObjectByDot(dot))
	}
	return objects
}

func (pg *Playground) GetObjectsByDots(dots []*engine.Dot) []interface{} {
	pg.entitiesMutex.RLock()
	defer pg.entitiesMutex.RUnlock()
	return pg.unsafeGetObjectsByDots(dots)
}

type ErrCreateObject struct {
	Err error
}

func (e *ErrCreateObject) Error() string {
	return "error create object"
}

type ErrLocationUsedByObject struct {
	Location engine.Location
}

func (e *ErrLocationUsedByObject) Error() string {
	return "passed location used by an object"
}

type ErrLocationDotsOccupiedByObjects struct {
	Objects []interface{}
}

func (e *ErrLocationDotsOccupiedByObjects) Error() string {
	return "dots of location is occupied by objects"
}

func (pg *Playground) unsafeCreateEntity(object interface{}, location engine.Location) {
	pg.entities = append(pg.entities, entity{
		object:   object,
		location: location,
	})
}

func (pg *Playground) CreateObject(object interface{}, location engine.Location) *ErrCreateObject {
	// TODO: Check if location is empty

	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if pg.unsafeObjectExists(object) {
		return &ErrCreateObject{
			Err: errors.New("passed object exists on playground"),
		}
	}

	if pg.unsafeLocationExists(location) {
		return &ErrCreateObject{
			Err: &ErrLocationUsedByObject{
				Location: location,
			},
		}
	}

	if err := pg.scene.Locate(location); err != nil {
		switch errLocateReason := err.Err.(type) {
		case *engine.ErrLocated:
			// Location is occupied on scene
			return &ErrCreateObject{
				Err: errLocateReason,
			}
		case *engine.ErrAreaNotContainsDot:
			// A dot are not contained in area
			return &ErrCreateObject{
				Err: errLocateReason,
			}
		case *engine.ErrDotsOccupied:
			// Dots of location is occupied by objects
			return &ErrCreateObject{
				Err: &ErrLocationDotsOccupiedByObjects{
					Objects: pg.unsafeGetObjectsByDots(errLocateReason.Dots),
				},
			}
		default:
			// Unknown location error
			return &ErrCreateObject{
				Err: errLocateReason,
			}
		}
	}

	pg.unsafeCreateEntity(object, location.Copy())

	return nil
}

type ErrCreateObjectAvailableDots struct {
	Err error
}

func (e *ErrCreateObjectAvailableDots) Error() string {
	return "error on creating objects available dots"
}

func (pg *Playground) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, *ErrCreateObjectAvailableDots) {
	// TODO: Check if location is empty

	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if pg.unsafeObjectExists(object) {
		return nil, &ErrCreateObjectAvailableDots{
			Err: errors.New("passed object exists on playground"),
		}
	}

	if pg.unsafeLocationExists(location) {
		return nil, &ErrCreateObjectAvailableDots{
			Err: &ErrLocationUsedByObject{
				Location: location,
			},
		}
	}

	location = pg.scene.LocateAvailableDots(location)

	if len(location) == 0 {
		return nil, &ErrCreateObjectAvailableDots{
			Err: errors.New("location dots are occupied"),
		}
	}

	pg.unsafeCreateEntity(object, location.Copy())

	return location.Copy(), nil
}

func (pg *Playground) unsafeDeleteEntity(object interface{}, location engine.Location) error {
	for i := range pg.entities {
		if pg.entities[i].object == object && pg.entities[i].location.Equals(location) {
			pg.entities = append(pg.entities[:i], pg.entities[:i+1]...)
			return nil
		}
	}

	return errors.New("cannot delete entity: entity not found")
}

type ErrDeleteObject struct {
	Err error
}

func (e *ErrDeleteObject) Error() string {
	return "error on object deletion"
}

func (pg *Playground) DeleteObject(object interface{}, location engine.Location) *ErrDeleteObject {
	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if !pg.unsafeEntityExists(object, location) {
		return &ErrDeleteObject{
			Err: errors.New("passed object and location entity does not exists"),
		}
	}

	if err := pg.scene.Delete(location); err != nil {
		switch errDeleteReason := err.Err.(type) {
		case *engine.ErrNotLocated:
			if err := pg.unsafeDeleteEntity(object, location); err != nil {
				// Concurrent invocation of unsafe method of playground
				return &ErrDeleteObject{
					Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
				}
			}
			return &ErrDeleteObject{
				Err: errors.New("object is not located on scene"),
			}
		default:
			// Unknown deletion error
			return &ErrDeleteObject{
				Err: errDeleteReason,
			}
		}
	}

	if err := pg.unsafeDeleteEntity(object, location); err != nil {
		// Concurrent invocation of unsafe method of playground
		return &ErrDeleteObject{
			Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
		}
	}

	return nil
}

type ErrUpdateObject struct {
	Err error
}

func (e *ErrUpdateObject) Error() string {
	return "update object error"
}

func (pg *Playground) UpdateObject(object interface{}, old, new engine.Location) *ErrUpdateObject {
	// TODO: Check if old or new location is empty

	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if !pg.unsafeEntityExists(object, old) {
		return &ErrUpdateObject{
			Err: errors.New("passed object and location entity does not exists"),
		}
	}

	if pg.unsafeLocationExists(new) {
		return &ErrUpdateObject{
			Err: &ErrLocationUsedByObject{
				Location: new,
			},
		}
	}

	if err := pg.scene.Relocate(old, new); err != nil {
		switch errRelocateReason := err.Err.(type) {
		case *engine.ErrLocated:
			// New location is already occupied on scene but is not registered on playground
			return &ErrUpdateObject{
				Err: errRelocateReason,
			}
		case *engine.ErrNotLocated:
			// Старый объект не находится в том месте на сцене !
			if err := pg.unsafeDeleteEntity(object, old); err != nil {
				// Concurrent invocation of unsafe method of playground
				return &ErrUpdateObject{
					Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
				}
			}
			return &ErrUpdateObject{
				Err: errors.New("object is not located on scene"),
			}
		case *engine.ErrDelete:
			// Ошибка при удалении старого объекта
			switch errDeleteReason := errRelocateReason.Err.(type) {
			case *engine.ErrNotLocated:
				// Старый объект не находится на сцене
				if err := pg.unsafeDeleteEntity(object, old); err != nil {
					// Concurrent invocation of unsafe method of playground
					return &ErrUpdateObject{
						Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
					}
				}
				return &ErrUpdateObject{
					Err: errors.New("object is not located on scene"),
				}
			default:
				// Unknown deletion error
				return &ErrUpdateObject{
					Err: errDeleteReason,
				}
			}
		case *engine.ErrLocate:
			// Ошибка при размещении нового объекта
			switch errLocateReason := errRelocateReason.Err.(type) {
			case *engine.ErrLocated:
				// Новый объект уже на сцене аллоцирован
				return &ErrUpdateObject{
					Err: errLocateReason,
				}
			case *engine.ErrAreaNotContainsDot:
				// Точка выходит за рамки карты !
				return &ErrUpdateObject{
					Err: errLocateReason,
				}
			case *engine.ErrDotsOccupied:
				return &ErrUpdateObject{
					Err: &ErrLocationDotsOccupiedByObjects{
						Objects: pg.unsafeGetObjectsByDots(errLocateReason.Dots),
					},
				}
			default:
				// Unknown location error
				return &ErrUpdateObject{
					Err: errLocateReason,
				}
			}
		default:
			// Unknown relocation error
			// TODO: Create ErrUnknown{}
			return &ErrUpdateObject{
				Err: errRelocateReason,
			}
		}
	}

	if err := pg.unsafeDeleteEntity(object, old); err != nil {
		return &ErrUpdateObject{
			Err: errors.New("concurrent invocation of unsafe methods on playground"),
		}
	}

	pg.unsafeCreateEntity(object, new.Copy())

	return nil
}

type ErrUpdateObjectAvailableDots struct {
	Err error
}

func (err *ErrUpdateObjectAvailableDots) Error() string {
	return "error update object available dots"
}

func (pg *Playground) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, *ErrUpdateObjectAvailableDots) {
	// TODO: Check if old or new location is empty
	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if !pg.unsafeEntityExists(object, old) {
		return nil, &ErrUpdateObjectAvailableDots{
			Err: errors.New("passed object and location entity does not exists"),
		}
	}

	if pg.unsafeLocationExists(new) {
		return nil, &ErrUpdateObjectAvailableDots{
			Err: &ErrLocationUsedByObject{
				Location: new,
			},
		}
	}

	location, err := pg.scene.RelocateAvailableDots(old, new)
	switch errRelocateAvailableDotsReason := err.Err.(type) {
	case *engine.ErrNotLocated:
		// Old location is not actual for object
		if err := pg.unsafeDeleteEntity(object, location); err != nil {
			// Concurrent invocation of unsafe method of playground
			return nil, &ErrUpdateObjectAvailableDots{
				Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
			}
		}

		return nil, &ErrUpdateObjectAvailableDots{
			Err: errors.New("object is not located on scene"),
		}
	case *engine.ErrDelete:
		// Ошибка при удалении старого объекта
		switch errDeleteReason := errRelocateAvailableDotsReason.Err.(type) {
		case *engine.ErrNotLocated:
			// Старый объект не находится на сцене
			if err := pg.unsafeDeleteEntity(object, old); err != nil {
				// Concurrent invocation of unsafe method of playground
				return nil, &ErrUpdateObjectAvailableDots{
					Err: errors.New("cannot delete entity: concurrent invocation of unsafe methods of playground"),
				}
			}
			return nil, &ErrUpdateObjectAvailableDots{
				Err: errors.New("object is not located on scene"),
			}
		default:
			// Unknown deletion error
			return nil, &ErrUpdateObjectAvailableDots{
				Err: errDeleteReason,
			}
		}
	case *engine.ErrDotsOccupied:
		// Case of all dots of new location are occupied
		return nil, &ErrUpdateObjectAvailableDots{
			Err: &ErrLocationDotsOccupiedByObjects{
				Objects: pg.unsafeGetObjectsByDots(errRelocateAvailableDotsReason.Dots),
			},
		}
	default:
		// Unknown relocate available dots error
		// TODO: Create ErrUnknown{}
		return nil, &ErrUpdateObjectAvailableDots{
			Err: errRelocateAvailableDotsReason,
		}
	}

	if len(location) == 0 {
		// RelocateAvailableDots did not return error but return empty location

		// Concurrent invocation of unsafe method of playground
		if err := pg.unsafeDeleteEntity(object, old); err != nil {
			return nil, &ErrUpdateObjectAvailableDots{
				Err: err,
			}
		}
		return nil, &ErrUpdateObjectAvailableDots{
			Err: errors.New("all dots is occupied"),
		}
	}

	if err := pg.unsafeDeleteEntity(object, old); err != nil {
		// Concurrent invocation of unsafe method of playground
		return nil, &ErrUpdateObjectAvailableDots{
			Err: err,
		}
	}

	pg.unsafeCreateEntity(object, location.Copy())

	return location.Copy(), nil
}

type ErrCreateObjectRandomDot string

func (e ErrCreateObjectRandomDot) Error() string {
	return "cannot create object of random dot: " + string(e)
}

func (pg *Playground) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if !pg.unsafeObjectExists(object) {
		return nil, ErrCreateObjectRandomDot("object to create already created")
	}

	location, err := pg.scene.LocateRandomDot()
	if err != nil {
		return nil, ErrCreateObjectRandomDot(err.Error())
	}

	pg.unsafeCreateEntity(object, location.Copy())

	return location.Copy(), nil
}

type ErrCreateRandomRectObject string

func (e ErrCreateRandomRectObject) Error() string {
	return "cannot create random rect object: " + string(e)
}

func (pg *Playground) CreateRandomRectObject(object interface{}, rw, rh uint8) (engine.Location, error) {
	// TODO: Check rw*rh > 0 and return error if it is not true

	pg.entitiesMutex.Lock()
	defer pg.entitiesMutex.Unlock()

	if !pg.unsafeObjectExists(object) {
		return nil, ErrCreateRandomRectObject("object to create already created")
	}

	location, err := pg.scene.LocateRandomRect(rw, rh)
	if err != nil {
		return nil, ErrCreateRandomRectObject(err.Error())
	}

	pg.unsafeCreateEntity(object, location.Copy())

	return location.Copy(), nil
}

func (pg *Playground) Navigate(dot *engine.Dot, dir engine.Direction, dis uint8) (*engine.Dot, error) {
	return pg.scene.Navigate(dot, dir, dis)
}
