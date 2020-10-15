package world

import "github.com/ivan1993spb/snake-server/engine"

type Interface interface {
	Start(stop <-chan struct{})
	Events(stop <-chan struct{}, buffer uint) <-chan Event
	GetObjectByDot(dot engine.Dot) engine.Object
	GetObjectsByDots(dots []engine.Dot) []engine.Object
	CreateObject(object engine.Object, location engine.Location) error
	CreateObjectAvailableDots(object engine.Object, location engine.Location) (engine.Location, error)
	DeleteObject(object engine.Object, location engine.Location) error
	UpdateObject(object engine.Object, old, new engine.Location) error
	UpdateObjectAvailableDots(object engine.Object, old, new engine.Location) (engine.Location, error)
	CreateObjectRandomDot(object engine.Object) (engine.Location, error)
	CreateObjectRandomRect(object engine.Object, rw, rh uint8) (engine.Location, error)
	CreateObjectRandomRectMargin(object engine.Object, rw, rh, margin uint8) (engine.Location, error)
	CreateObjectRandomByDotsMask(object engine.Object, dm *engine.DotsMask) (engine.Location, error)
	LocationOccupied(location engine.Location) bool
	Area() engine.Area
	GetObjects() []engine.Object
	IdentifierRegistry() *IdentifierRegistry
}
