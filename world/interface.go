package world

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type Interface interface {
	Start(stop <-chan struct{})
	Events(stop <-chan struct{}, buffer uint) <-chan Event

	IdentifierRegistry() *IdentifierRegistry

	playground.Playground

	UpdateV2
	DeleteV2
}

type UpdateV2Params struct {
	Object  engine.Object
	Payload interface{}
	Old     engine.Location
	New     engine.Location
}

type UpdateV2 interface {
	UpdateObjectV2(params UpdateV2Params) error
	UpdateObjectAvailableDotsV2(params UpdateV2Params) (engine.Location, error)
}

type DeleteV2Params struct {
	Object   engine.Object
	Payload  interface{}
	Location engine.Location
}

type DeleteV2 interface {
	DeleteObjectV2(params DeleteV2Params) error
}
