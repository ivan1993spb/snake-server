package game

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type Game struct {
	world World
}

type ErrCreateGame struct {
	Err error
}

func (e *ErrCreateGame) Error() string {
	return "cannot create game: " + e.Err.Error()
}

func NewGame(width, height uint8) (*Game, error) {
	area, err := engine.NewArea(width, height)
	if err != nil {
		return nil, &ErrCreateGame{Err: err}
	}

	scene, err := engine.NewScene(area)
	if err != nil {
		return nil, &ErrCreateGame{Err: err}
	}

	pg := playground.NewPlayground(scene)
	world := newWorld(pg)
	world.run()

	return &Game{
		world: world,
	}, nil
}
