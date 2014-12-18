package playground

import (
	"errors"
	"fmt"
)

var (
	ErrPGNotContainsDot = errors.New("Playground doesn't contain dot")
	ErrInvalid_W_or_H   = errors.New("Invalid width or height")
)

// Playground object contains all objects on map
type Playground struct {
	width, height uint8             // Width and height of map
	objects       map[uint16]Object // All objects on map
}

// NewPlayground returns new empty playground
func NewPlayground(width, height uint8) (*Playground, error) {
	if width*height == 0 {
		return nil, fmt.Errorf("Cannot create playground: %s",
			ErrInvalid_W_or_H)
	}

	return &Playground{width, height, make(map[uint16]Object)}, nil
}

// GetArea returns playground area
func (pg *Playground) GetArea() uint16 {
	return uint16(pg.width) * uint16(pg.height)
}

// Occupied returns true if passed dot already used by any object
// located on playground
func (pg *Playground) Occupied(dot *Dot) bool {
	return pg.GetObjectByDot(dot) != nil
}

// GetObjectByDot returns object which contains passed dot
func (pg *Playground) GetObjectByDot(dot *Dot) Object {
	if pg.Contains(dot) {
		for _, object := range pg.objects {
			for i := uint16(0); i < object.DotCount(); i++ {
				if object.Dot(i).Equals(dot) {
					return object
				}
			}
		}
	}
	return nil
}

type errCannotLocate struct {
	err error
}

func (e *errCannotLocate) Error() string {
	return "Cannot locate object: " + e.err.Error()
}

// Locate tries to create object to playground
func (pg *Playground) Locate(object Object) error {
	// Return error if object is already located on playground
	if pg.Located(object) {
		return &errCannotLocate{
			errors.New("Object is already located"),
		}
	}
	// Check each dot of passed object
	for i := uint16(0); i < object.DotCount(); i++ {
		var dot = object.Dot(i)
		// Return error if any dot is invalid...
		if !pg.Contains(dot) {
			return &errCannotLocate{ErrPGNotContainsDot}
		}
		// ...or occupied
		if pg.Occupied(dot) {
			return &errCannotLocate{errors.New("Dot is occupied")}
		}
	}

	// Object count can't be more than playground area
	var maxId = pg.GetArea()

	// Add to object list of playground
	for id := uint16(0); id < maxId; id++ {
		if _, ok := pg.objects[id]; !ok {
			pg.objects[id] = object
			return nil
		}
	}

	return &errCannotLocate{errors.New("Playground is full")}
}

// Located returns true if passed object is located on playground
func (pg *Playground) Located(object Object) bool {
	for i := range pg.objects {
		if pg.objects[i] == object {
			return true
		}
	}

	return false
}

// Contains return true if playground contains passed dot
func (pg *Playground) Contains(dot *Dot) bool {
	return pg.width > dot.x && pg.height > dot.y
}

// Delete deletes passed object from playground and returns error if
// there is a problem
func (pg *Playground) Delete(object Object) error {
	if pg.Located(object) {
		for id := range pg.objects {
			if pg.objects[id] == object {
				// Delete object from object storage
				delete(pg.objects, id)

				return nil
			}
		}
	}

	return errors.New("Cannot delocate: Passed object isn't located")
}

// RandomDot generates random dot located on playground
func (pg *Playground) RandomDot() *Dot {
	return NewRandomDotOnSquare(0, 0, pg.width, pg.height)
}

func (pg Playground) RandomRect(rw, rh uint8) (*Rect, error) {
	return NewRandomRectOnSquare(rw, rh, 0, 0, pg.width, pg.height)
}
