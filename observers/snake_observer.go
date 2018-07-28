package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/corpse"
	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

const chanSnakeObserverEventsBuffer = 64

type SnakeObserver struct{}

func (SnakeObserver) Observe(stop <-chan struct{}, w *world.World, logger logrus.FieldLogger) {
	go func() {
		for event := range w.Events(stop, chanSnakeObserverEventsBuffer) {
			if event.Type == world.EventTypeObjectDelete {
				if s, ok := event.Payload.(*snake.Snake); ok {
					location := s.GetLocation().Copy()
					if location.Empty() {
						logger.Warn("snake dies and returns empty location")
						continue
					}

					if c, err := corpse.NewCorpse(w, location); err != nil {
						logger.WithError(err).Error("cannot create corpse")
					} else {
						c.Run(stop, logger)
					}
				}
			}
		}
	}()
}
