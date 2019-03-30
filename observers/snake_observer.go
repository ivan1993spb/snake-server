package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/corpse"
	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

const chanSnakeObserverEventsBuffer = 64

type SnakeObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewSnakeObserver(w world.Interface, logger logrus.FieldLogger) Observer {
	return &SnakeObserver{
		world:  w,
		logger: logger,
	}
}

func (so *SnakeObserver) Observe(stop <-chan struct{}) {
	go func() {
		for event := range so.world.Events(stop, chanSnakeObserverEventsBuffer) {
			if event.Type == world.EventTypeObjectDelete {
				if s, ok := event.Payload.(*snake.Snake); ok {
					location := s.GetLocation().Copy()
					if location.Empty() {
						so.logger.Warn("snake dies and returns empty location")
						continue
					}

					if c, err := corpse.NewCorpse(so.world, location); err != nil {
						so.logger.WithError(err).Error("cannot create corpse")
					} else {
						c.Run(stop, so.logger)
					}
				}
			}
		}
	}()
}
