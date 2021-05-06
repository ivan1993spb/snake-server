package playground

import "github.com/ivan1993spb/snake-server/engine"

// Playground interface declares a set of methods to be implemented by
// a playground structure operating with objects on a map.
type Playground interface {
	// CreateObject should create an object on a predefined location
	CreateObject(object engine.Object, location engine.Location) error
	// CreateObjectAvailableDots should create an object on a predefined location. The method
	// doesn't guarantee the fact that all the dots will be occupied by the object.
	CreateObjectAvailableDots(object engine.Object, location engine.Location) (engine.Location, error)

	// CreateObjectRandomDot should create an object at a random dot
	CreateObjectRandomDot(object engine.Object) (engine.Location, error)
	// CreateObjectRandomRect should create an object on a random rectangle location with
	// predefined width and height
	CreateObjectRandomRect(object engine.Object, rw, rh uint8) (engine.Location, error)
	// CreateObjectRandomRectMargin should create an object on a random rectangle location
	// with predefined width and height and a margin between the object and others on the map
	CreateObjectRandomRectMargin(object engine.Object, rw, rh, margin uint8) (engine.Location, error)
	// CreateObjectRandomByDotsMask should create an object on a random location with the shape
	// of the passed dot mask
	CreateObjectRandomByDotsMask(object engine.Object, dm *engine.DotsMask) (engine.Location, error)

	// UpdateObject should move the object from its old location to the new one
	UpdateObject(object engine.Object, old, new engine.Location) error
	// UpdateObjectAvailableDots should move the objects from its old location to the new one.
	// The method doesn't guarantee that all the dots of new location will be occupied by the
	// object.
	UpdateObjectAvailableDots(object engine.Object, old, new engine.Location) (engine.Location, error)

	// DeleteObject should delete the object from map
	DeleteObject(object engine.Object, location engine.Location) error

	// GetObjectsByDots should return a slice of objects which occupy the dots
	GetObjectsByDots(dots []engine.Dot) []engine.Object
	// GetObjectByDot should return an object occupying the passed dot
	GetObjectByDot(dot engine.Dot) engine.Object
	// LocationOccupied should return true if the location is occupied and false if it's vacant
	LocationOccupied(location engine.Location) bool
	// Area should return Area object
	Area() engine.Area
	// GetObjects should return a slice of all objects on the map
	GetObjects() []engine.Object
}
