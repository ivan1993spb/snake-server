package engine

import (
	"errors"
	"fmt"
	"sync"
)

type ErrObjectNotLocated struct {
	Object DotList
}

func (e *ErrObjectNotLocated) Error() string {
	return "object is not located"
}

type ErrObjectLocated struct {
	Object DotList
}

func (e *ErrObjectLocated) Error() string {
	return "object is located"
}

type ErrDotsOccupied struct {
	Dots []*Dot // Occupied dots
}

func (e *ErrDotsOccupied) Error() string {
	return fmt.Sprintf("dots is occupied by an objects: %s", e.Dots)
}

// Scene object contains all objects on map
type Scene struct {
	area         *Area
	objects      []DotList
	objectsMutex *sync.RWMutex
}

// NewScene returns new empty scene
func NewScene(area *Area) (*Scene, error) {
	return &Scene{
		area:         area,
		objects:      make([]DotList, 0),
		objectsMutex: &sync.RWMutex{},
	}, nil
}

// GetAreaSize returns scene area size
func (s *Scene) GetArea() *Area {
	return s.area
}

// unsafeObjectLocated returns true if passed object is located on scene
func (s *Scene) unsafeObjectLocated(object DotList) bool {
	for i := range s.objects {
		if s.objects[i].Equals(object) {
			return true
		}
	}
	return false
}

func (s *Scene) ObjectLocated(object DotList) bool {
	s.objectsMutex.RLock()
	defer s.objectsMutex.RUnlock()
	return s.unsafeObjectLocated(object)
}

// unsafeDotOccupied returns true if passed dot already used by an object located on scene
func (s *Scene) unsafeDotOccupied(dot *Dot) bool {
	if s.area.Contains(dot) {
		for _, object := range s.objects {
			if object.Contains(dot) {
				return true
			}
		}
	}
	return false
}

func (s *Scene) DotOccupied(dot *Dot) bool {
	s.objectsMutex.RLock()
	defer s.objectsMutex.RUnlock()
	return s.unsafeDotOccupied(dot)
}

// unsafeGetObjectByDot returns object which contains passed dot
func (s *Scene) unsafeGetObjectByDot(dot *Dot) DotList {
	if s.area.Contains(dot) {
		for _, object := range s.objects {
			if object.Contains(dot) {
				return object.Copy()
			}
		}
	}
	return nil
}

func (s *Scene) GetObjectByDot(dot *Dot) DotList {
	s.objectsMutex.RLock()
	defer s.objectsMutex.RUnlock()
	return s.unsafeGetObjectByDot(dot)
}

type ErrLocateObject struct {
	Err error
}

func (e *ErrLocateObject) Error() string {
	return "cannot locate object: " + e.Err.Error()
}

func (s *Scene) unsafeLocateObject(object DotList) *ErrLocateObject {
	if s.unsafeObjectLocated(object) {
		return &ErrLocateObject{
			Err: &ErrObjectLocated{
				Object: object.Copy(),
			},
		}
	}

	object = object.Copy()
	occupiedDots := make([]*Dot, 0)

	// Check each dot of passed object
	for i := uint16(0); i < object.DotCount(); i++ {
		var dot = object.Dot(i)

		if s.area.Contains(dot) {
			return &ErrLocateObject{
				Err: &ErrAreaNotContainsDot{
					Dot: dot,
				},
			}
		}

		for i := range s.objects {
			if s.objects[i].Contains(dot) {
				occupiedDots = append(occupiedDots, dot)
				break
			}
		}
	}

	if len(occupiedDots) > 0 {
		return &ErrLocateObject{
			Err: &ErrDotsOccupied{
				Dots: occupiedDots,
			},
		}
	}

	// Add to object list of scene
	s.objects = append(s.objects, object)

	return nil
}

// LocateObject tries to create object on scene
func (s *Scene) LocateObject(object DotList) *ErrLocateObject {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeLocateObject(object)
}

func (s *Scene) unsafeLocateAvailableObjectDots(object DotList) DotList {
	if s.unsafeObjectLocated(object) {
		return DotList{}
	}

	object = object.Copy()
	objectMirror := object.Copy()

	// Check each dot of passed object
	for i := uint16(0); i < objectMirror.DotCount(); i++ {
		var dot = objectMirror.Dot(i)

		if !s.area.Contains(dot) {
			object = object.Delete(dot)
			continue
		}

		for i := range s.objects {
			if s.objects[i].Contains(dot) {
				object = object.Delete(dot)
				break
			}
		}
	}

	if len(object) > 0 {
		s.objects = append(s.objects, object)
	}

	return object.Copy()
}

func (s *Scene) LocateAvailableObjectDots(object DotList) DotList {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeLocateAvailableObjectDots(object)
}

type ErrDeleteObject struct {
	Err error
}

func (e *ErrDeleteObject) Error() string {
	return "cannot delete object"
}

// unsafeDeleteObject deletes passed object from scene and returns error if there is a problem
func (s *Scene) unsafeDeleteObject(object DotList) *ErrDeleteObject {
	for i := range s.objects {
		if s.objects[i].Equals(object) {
			s.objects = append(s.objects[:i], s.objects[i+1:]...)
			return nil
		}
	}

	return &ErrDeleteObject{
		Err: &ErrObjectNotLocated{
			Object: object.Copy(),
		},
	}
}

// DeleteObject deletes passed object from scene and returns error if there is a problem
func (s *Scene) DeleteObject(object DotList) *ErrDeleteObject {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeDeleteObject(object)
}

type ErrRelocateObject struct {
	Err error
}

func (e *ErrRelocateObject) Error() string {
	return "cannot relocate object: " + e.Err.Error()
}

func (s *Scene) unsafeRelocateObject(old, new DotList) *ErrRelocateObject {
	if !s.unsafeObjectLocated(old) {
		return &ErrRelocateObject{
			Err: &ErrObjectNotLocated{
				Object: old.Copy(),
			},
		}
	}

	if s.unsafeObjectLocated(new) {
		return &ErrRelocateObject{
			Err: &ErrObjectLocated{
				Object: new.Copy(),
			},
		}
	}

	if err := s.unsafeDeleteObject(old); err != nil {
		return &ErrRelocateObject{
			Err: err,
		}
	}

	if err := s.unsafeLocateObject(new); err != nil {
		return &ErrRelocateObject{
			Err: err,
		}
	}

	return nil
}

func (s *Scene) RelocateObject(old, new DotList) *ErrRelocateObject {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeRelocateObject(old, new)
}

type ErrRelocateObjectAvailableDots struct {
	Err error
}

func (e *ErrRelocateObjectAvailableDots) Error() string {
	return "cannot relocate object with available dots"
}

func (s *Scene) unsafeRelocateObjectAvailableDots(old, new DotList) (*ErrRelocateObjectAvailableDots, DotList) {
	if !s.unsafeObjectLocated(old) {
		return &ErrRelocateObjectAvailableDots{
			Err: &ErrObjectNotLocated{
				Object: old.Copy(),
			},
		}, nil
	}

	if err := s.unsafeDeleteObject(old); err != nil {
		return &ErrRelocateObjectAvailableDots{
			Err: err,
		}, nil
	}

	dots := s.unsafeLocateAvailableObjectDots(new)
	if len(dots) == 0 {
		return &ErrRelocateObjectAvailableDots{
			Err: fmt.Errorf("all dots are not available"),
		}, nil
	}

	return nil, dots.Copy()
}

func (s *Scene) RelocateObjectAvailableDots(old, new DotList) (error, DotList) {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeRelocateObjectAvailableDots(old, new)
}

var FindRetriesNumber = 32

var ErrRetriesLimit = errors.New("retries limit was reached")

func (s *Scene) unsafeLocateRandomDot() (*Dot, error) {
	for count := 0; count < FindRetriesNumber; count++ {
		if dot := s.area.NewRandomDot(0, 0); !s.unsafeDotOccupied(dot) {
			if err := s.unsafeLocateObject(DotList{dot}); err != nil {
				return nil, err
			}
			return dot, nil
		}
	}

	return nil, ErrRetriesLimit
}

func (s *Scene) LocateRandomDot() (*Dot, error) {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeLocateRandomDot()
}

func (s *Scene) unsafeLocateRandomRectTryOnce(rw, rh uint8) (*Rect, error) {
	if rect, err := s.area.NewRandomRect(rw, rh, 0, 0); err == nil {
		if err := s.unsafeLocateObject(rect.DotList()); err != nil {
			return nil, err
		}
		return rect, nil
	} else {
		return nil, err
	}
}

func (s *Scene) unsafeLocateRandomRect(rw, rh uint8) (*Rect, error) {
	for count := 0; count < FindRetriesNumber; count++ {
		if rect, err := s.unsafeLocateRandomRectTryOnce(rw, rh); err == nil {
			return rect, nil
		}
	}

	return nil, ErrRetriesLimit
}

func (s *Scene) LocateRandomRect(rw, rh uint8) (*Rect, error) {
	s.objectsMutex.Lock()
	defer s.objectsMutex.Unlock()
	return s.unsafeLocateRandomRect(rw, rh)
}
