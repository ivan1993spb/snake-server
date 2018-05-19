package engine

import (
	"errors"
	"fmt"
	"sync"
)

type ErrNotLocated struct {
	Location Location
}

func (e *ErrNotLocated) Error() string {
	return "not located"
}

type ErrLocated struct {
	Location Location
}

func (e *ErrLocated) Error() string {
	return "located"
}

type ErrDotsOccupied struct {
	Dots []Dot // List of occupied dots
}

func (e *ErrDotsOccupied) Error() string {
	return "dots is occupied"
}

// Scene contains locations
type Scene struct {
	area           Area
	locations      []Location
	locationsMutex *sync.RWMutex
}

// NewScene returns new empty scene
func NewScene(width, height uint8) (*Scene, error) {
	area, err := NewUsefulArea(width, height)
	if err != nil {
		return nil, fmt.Errorf("cannot create scene: %s", err)
	}

	return &Scene{
		area:           area,
		locations:      make([]Location, 0),
		locationsMutex: &sync.RWMutex{},
	}, nil
}

// unsafeLocated returns true if passed location is located on scene
func (s *Scene) unsafeLocated(location Location) bool {
	for i := range s.locations {
		if s.locations[i].Equals(location) {
			return true
		}
	}
	return false
}

func (s *Scene) Located(location Location) bool {
	s.locationsMutex.RLock()
	defer s.locationsMutex.RUnlock()
	return s.unsafeLocated(location)
}

// unsafeDotOccupied returns true if passed dot already used by a location on scene
func (s *Scene) unsafeDotOccupied(dot Dot) bool {
	if s.area.Contains(dot) {
		for _, location := range s.locations {
			if location.Contains(dot) {
				return true
			}
		}
	}
	return false
}

func (s *Scene) DotOccupied(dot Dot) bool {
	s.locationsMutex.RLock()
	defer s.locationsMutex.RUnlock()
	return s.unsafeDotOccupied(dot)
}

// unsafeGetLocationByDot returns location which contains passed dot
func (s *Scene) unsafeGetLocationByDot(dot Dot) Location {
	if s.area.Contains(dot) {
		for _, location := range s.locations {
			if location.Contains(dot) {
				return location.Copy()
			}
		}
	}
	return nil
}

func (s *Scene) GetLocationByDot(dot Dot) Location {
	s.locationsMutex.RLock()
	defer s.locationsMutex.RUnlock()
	return s.unsafeGetLocationByDot(dot)
}

type ErrLocate struct {
	Err error
}

func (e *ErrLocate) Error() string {
	return "cannot locate: " + e.Err.Error()
}

func (s *Scene) unsafeLocate(location Location) *ErrLocate {
	if s.unsafeLocated(location) {
		return &ErrLocate{
			Err: &ErrLocated{
				Location: location.Copy(),
			},
		}
	}

	location = location.Copy()
	occupiedDots := make([]Dot, 0)

	// Check each dot of passed location
	for i := uint16(0); i < location.DotCount(); i++ {
		var dot = location.Dot(i)

		if !s.area.Contains(dot) {
			return &ErrLocate{
				Err: &ErrAreaNotContainsDot{
					Dot: dot,
				},
			}
		}

		for i := range s.locations {
			if s.locations[i].Contains(dot) {
				occupiedDots = append(occupiedDots, dot)
				break
			}
		}
	}

	if len(occupiedDots) > 0 {
		return &ErrLocate{
			Err: &ErrDotsOccupied{
				Dots: occupiedDots,
			},
		}
	}

	// Add to location list of scene
	s.locations = append(s.locations, location)

	return nil
}

// Locate tries to create location to scene
func (s *Scene) Locate(location Location) *ErrLocate {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeLocate(location)
}

func (s *Scene) unsafeLocateAvailableDots(location Location) Location {
	if s.unsafeLocated(location) {
		return Location{}
	}

	location = location.Copy()
	locationMirror := location.Copy()

	// Check each dot of passed location
	for i := uint16(0); i < locationMirror.DotCount(); i++ {
		var dot = locationMirror.Dot(i)

		if !s.area.Contains(dot) {
			location = location.Delete(dot)
			continue
		}

		for i := range s.locations {
			if s.locations[i].Contains(dot) {
				location = location.Delete(dot)
				break
			}
		}
	}

	if len(location) > 0 {
		s.locations = append(s.locations, location)
	}

	return location.Copy()
}

func (s *Scene) LocateAvailableDots(location Location) Location {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeLocateAvailableDots(location)
}

type ErrDelete struct {
	Err error
}

func (e *ErrDelete) Error() string {
	return "cannot delete"
}

// unsafeDelete deletes passed location from scene and returns error if there is a problem
func (s *Scene) unsafeDelete(location Location) *ErrDelete {
	for i := range s.locations {
		if s.locations[i].Equals(location) {
			s.locations = append(s.locations[:i], s.locations[i+1:]...)
			return nil
		}
	}

	return &ErrDelete{
		Err: &ErrNotLocated{
			Location: location.Copy(),
		},
	}
}

// Delete deletes passed location from scene and returns error if there is a problem
func (s *Scene) Delete(location Location) *ErrDelete {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeDelete(location)
}

type ErrRelocate struct {
	Err error
}

func (e *ErrRelocate) Error() string {
	return "cannot relocate"
}

func (s *Scene) unsafeRelocate(old, new Location) *ErrRelocate {
	if !s.unsafeLocated(old) {
		return &ErrRelocate{
			Err: &ErrNotLocated{
				Location: old.Copy(),
			},
		}
	}

	if s.unsafeLocated(new) {
		return &ErrRelocate{
			Err: &ErrLocated{
				Location: new.Copy(),
			},
		}
	}

	if err := s.unsafeDelete(old); err != nil {
		return &ErrRelocate{
			Err: err,
		}
	}

	if err := s.unsafeLocate(new); err != nil {
		return &ErrRelocate{
			Err: err,
		}
	}

	return nil
}

func (s *Scene) Relocate(old, new Location) *ErrRelocate {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeRelocate(old, new)
}

type ErrRelocateAvailableDots struct {
	Err error
}

func (e *ErrRelocateAvailableDots) Error() string {
	return "cannot relocate with available dots"
}

func (s *Scene) unsafeRelocateAvailableDots(old, new Location) (Location, *ErrRelocateAvailableDots) {
	if !s.unsafeLocated(old) {
		return nil, &ErrRelocateAvailableDots{
			Err: &ErrNotLocated{
				Location: old.Copy(),
			},
		}
	}

	if err := s.unsafeDelete(old); err != nil {
		return nil, &ErrRelocateAvailableDots{
			Err: err,
		}
	}

	dots := s.unsafeLocateAvailableDots(new)
	if len(dots) == 0 {
		return nil, &ErrRelocateAvailableDots{
			Err: &ErrDotsOccupied{
				Dots: new,
			},
		}
	}

	return dots.Copy(), nil
}

func (s *Scene) RelocateAvailableDots(old, new Location) (Location, *ErrRelocateAvailableDots) {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeRelocateAvailableDots(old, new)
}

var FindRetriesNumber = 32

var ErrRetriesLimit = errors.New("retries limit was reached")

func (s *Scene) unsafeLocateRandomDot() (Location, error) {
	for count := 0; count < FindRetriesNumber; count++ {
		if dot := s.area.NewRandomDot(0, 0); !s.unsafeDotOccupied(dot) {
			if err := s.unsafeLocate(Location{dot}); err != nil {
				return nil, err
			}
			return Location{dot}, nil
		}
	}

	return nil, ErrRetriesLimit
}

func (s *Scene) LocateRandomDot() (Location, error) {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeLocateRandomDot()
}

func (s *Scene) unsafeLocateRandomRectTryOnce(rw, rh uint8) (Location, error) {
	if rect, err := s.area.NewRandomRect(rw, rh, 0, 0); err == nil {
		if err := s.unsafeLocate(rect.Location()); err != nil {
			return nil, err
		}
		return rect.Location(), nil
	} else {
		return nil, err
	}
}

func (s *Scene) unsafeLocateRandomRect(rw, rh uint8) (Location, error) {
	for count := 0; count < FindRetriesNumber; count++ {
		if rect, err := s.unsafeLocateRandomRectTryOnce(rw, rh); err == nil {
			return rect, nil
		}
	}

	return nil, ErrRetriesLimit
}

func (s *Scene) LocateRandomRect(rw, rh uint8) (Location, error) {
	s.locationsMutex.Lock()
	defer s.locationsMutex.Unlock()
	return s.unsafeLocateRandomRect(rw, rh)
}

func (s *Scene) Navigate(dot Dot, dir Direction, dis uint8) (Dot, error) {
	return s.area.Navigate(dot, dir, dis)
}

func (s *Scene) Size() uint16 {
	return s.area.Size()
}

func (s *Scene) Width() uint8 {
	return s.area.Width()
}

func (s *Scene) Height() uint8 {
	return s.area.Height()
}
