package observers

import (
	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

type SnakeObserver struct{}

func (SnakeObserver) Observe(stop <-chan struct{}, w *world.World) {
	go func() {
		// TODO: Create buffer const.
		for event := range w.Events(stop, 32) {
			if event.Type == world.EventTypeObjectDelete {
				switch event.Payload.(type) {
				case *snake.Snake:
					//corpse.NewCorpse(w)
				}
			}
		}
	}()
}
