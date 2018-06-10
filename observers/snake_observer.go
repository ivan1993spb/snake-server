package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/corpse"
	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

type SnakeObserver struct{}

func (SnakeObserver) Observe(stop <-chan struct{}, w *world.World, logger logrus.FieldLogger) {
	go func() {
		// TODO: Create buffer const.
		for event := range w.Events(stop, 32) {
			if event.Type == world.EventTypeObjectDelete {
				if s, ok := event.Payload.(*snake.Snake); ok {
					// TODO: Handle error.
					c, err := corpse.NewCorpse(w, s.GetLocation())
					if err == nil {
						c.Run(stop)
					}
				}
			}
		}
	}()
}
