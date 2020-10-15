package world

import "github.com/ivan1993spb/snake-server/playground"

type Interface interface {
	Start(stop <-chan struct{})
	Events(stop <-chan struct{}, buffer uint) <-chan Event

	IdentifierRegistry() *IdentifierRegistry

	playground.Playground
}
