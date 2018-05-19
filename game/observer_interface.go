package game

import "github.com/ivan1993spb/snake-server/world"

type ObserverInterface interface {
	Observe(stop <-chan struct{}, world *world.World)
}
