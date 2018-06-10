package game

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/observers"
	"github.com/ivan1993spb/snake-server/world"
)

type Game struct {
	world  *world.World
	logger logrus.FieldLogger
}

type ErrCreateGame struct {
	Err error
}

func (e *ErrCreateGame) Error() string {
	return "cannot create game: " + e.Err.Error()
}

func NewGame(logger logrus.FieldLogger, width, height uint8) (*Game, error) {
	w, err := world.NewWorld(width, height)
	if err != nil {
		return nil, fmt.Errorf("cannot create game: %s", err)
	}

	return &Game{
		world:  w,
		logger: logger,
	}, nil
}

func (g *Game) Start(stop <-chan struct{}) {
	g.world.Start(stop)

	observers.LoggerObserver{}.Observe(stop, g.world, g.logger)
	observers.AppleObserver{}.Observe(stop, g.world, g.logger)
	observers.SnakeObserver{}.Observe(stop, g.world, g.logger)
}

func (g *Game) World() *world.World {
	return g.world
}

func (g *Game) ListenEvents(stop <-chan struct{}, buffer uint) <-chan Event {
	chout := make(chan Event, buffer)
	go func() {
		defer close(chout)
		for worldEvent := range g.world.Events(stop, buffer) {
			chout <- Event{
				Type:    worldEventTypeToGameEventType(worldEvent.Type),
				Payload: worldEvent.Payload,
			}
		}
	}()
	return chout
}
