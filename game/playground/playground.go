package playground

import (
	"errors"
	"fmt"
)

var (
	ErrPGNotContainsDot = errors.New("playground doesn't contain dot")
	ErrInvalid_W_or_H   = errors.New("invalid width or height")
	ErrObjectNotLocated = errors.New("passed object is not located")
)

// Playground object contains all objects on map
type Playground struct {
	width, height uint8    // Width and height of map
	objects       []Object // All objects on map
}

// NewPlayground returns new empty playground
func NewPlayground(width, height uint8) (*Playground, error) {
	if width*height == 0 {
		return nil, fmt.Errorf("cannot create playground: %s",
			ErrInvalid_W_or_H)
	}

	return &Playground{width, height, make([]Object, 0)}, nil
}

// GetArea returns playground area
func (pg *Playground) GetArea() uint16 {
	return uint16(pg.width) * uint16(pg.height)
}

func (pg *Playground) GetSize() (uint8, uint8) {
	return pg.width, pg.height
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
	return "cannot locate object: " + e.err.Error()
}

// Locate tries to create object on playground. If occupy=true object
// may be located only if each object dot is not occupied
func (pg *Playground) Locate(object Object, occupy bool) error {
	// Return error if object is already located on playground
	if pg.Located(object) {
		return &errCannotLocate{
			errors.New("object is already located"),
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
		if occupy && pg.Occupied(dot) {
			return &errCannotLocate{errors.New("dot is occupied")}
		}
	}

	// Add to object list of playground
	pg.objects = append(pg.objects, object)

	return nil
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
		for i := range pg.objects {
			if pg.objects[i] == object {
				pg.objects = append(pg.objects[:i],
					pg.objects[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("cannot delocate: %s", ErrObjectNotLocated)
}

// RandomDot generates random dot located on playground
func (pg *Playground) RandomDot() *Dot {
	return NewRandomDotOnSquare(0, 0, pg.width, pg.height)
}

func (pg *Playground) RandomRect(rw, rh uint8) (*Rect, error) {
	return NewRandomRectOnSquare(rw, rh, 0, 0, pg.width, pg.height)
}
