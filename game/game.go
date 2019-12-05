package game

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/observers/apple"
	"github.com/ivan1993spb/snake-server/observers/logger"
	"github.com/ivan1993spb/snake-server/observers/mouse"
	"github.com/ivan1993spb/snake-server/observers/snake"
	"github.com/ivan1993spb/snake-server/observers/wall"
	"github.com/ivan1993spb/snake-server/observers/watermelon"
	"github.com/ivan1993spb/snake-server/world"
)

type Game struct {
	world  world.Interface
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

	logger_observer.NewLoggerObserver(g.world, g.logger).Observe(stop)
	wall_observer.NewWallObserver(g.world, g.logger).Observe(stop)
	apple_observer.NewAppleObserver(g.world, g.logger).Observe(stop)
	snake_observer.NewSnakeObserver(g.world, g.logger).Observe(stop)
	watermelon_observer.NewWatermelonObserver(g.world, g.logger).Observe(stop)
	mouse_observer.NewMouseObserver(g.world, g.logger).Observe(stop)
}

func (g *Game) World() world.Interface {
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
