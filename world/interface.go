package world

import "github.com/ivan1993spb/snake-server/engine"

type Interface interface {
	Start(stop <-chan struct{})
	Events(stop <-chan struct{}, buffer uint) <-chan Event
	GetObjectByDot(dot engine.Dot) interface{}
	GetObjectsByDots(dots []engine.Dot) []interface{}
	CreateObject(object interface{}, location engine.Location) error
	CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error)
	DeleteObject(object interface{}, location engine.Location) error
	UpdateObject(object interface{}, old, new engine.Location) error
	UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error)
	CreateObjectRandomDot(object interface{}) (engine.Location, error)
	CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error)
	CreateObjectRandomRectMargin(object interface{}, rw, rh, margin uint8) (engine.Location, error)
	CreateObjectRandomByDotsMask(object interface{}, dm *engine.DotsMask) (engine.Location, error)
	LocationOccupied(location engine.Location) bool
	Area() engine.Area
	GetObjects() []interface{}
	ObtainIdentifier() Identifier
	ReleaseIdentifier(id Identifier)
}
