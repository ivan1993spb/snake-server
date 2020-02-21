package snake_observer

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/corpse"
	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/observers"
	"github.com/ivan1993spb/snake-server/world"
)

const chanSnakeObserverEventsBuffer = 64

type SnakeObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewSnakeObserver(w world.Interface, logger logrus.FieldLogger) observers.Observer {
	return &SnakeObserver{
		world:  w,
		logger: logger,
	}
}

func (so *SnakeObserver) Observe(stop <-chan struct{}) {
	go so.run(stop)
}

func (so *SnakeObserver) run(stop <-chan struct{}) {
	so.listen(stop)
}

func (so *SnakeObserver) listen(stop <-chan struct{}) {
	for event := range so.world.Events(stop, chanSnakeObserverEventsBuffer) {
		so.handleEvent(event, stop)
	}
}

func (so *SnakeObserver) handleEvent(event world.Event, stop <-chan struct{}) {
	if event.Type != world.EventTypeObjectDelete {
		return
	}

	if s, ok := event.Payload.(*snake.Snake); ok {
		location := s.GetLocation().Copy()
		if location.Empty() {
			so.logger.Warn("snake dies and returns empty location")
			return
		}

		// TODO: Create abstraction layer for adding of objects.
		if c, err := corpse.NewCorpse(so.world, location); err != nil {
			so.logger.WithError(err).Error("cannot create corpse")
		} else {
			c.Run(stop, so.logger)
		}
	}
}
