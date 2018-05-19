package game

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type Game struct {
	world  *World
	logger logrus.FieldLogger
}

type ErrCreateGame struct {
	Err error
}

func (e *ErrCreateGame) Error() string {
	return "cannot create game: " + e.Err.Error()
}

func NewGame(logger logrus.FieldLogger, width, height uint8) (*Game, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, &ErrCreateGame{Err: err}
	}

	scene, err := engine.NewScene(area)
	if err != nil {
		return nil, &ErrCreateGame{Err: err}
	}

	pg := playground.NewPlayground(scene)
	world := NewWorld(pg)

	return &Game{
		world:  world,
		logger: logger,
	}, nil
}

func (g *Game) Start(stop <-chan struct{}) {
	g.world.start(stop)

	// TODO: Start observers.
}

func (g *Game) World() *World {
	return g.world
}

func (g *Game) Events(stop <-chan struct{}, buffer uint) <-chan Event {
	return g.world.Events(stop, buffer)
}
