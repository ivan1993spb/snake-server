package observers

import (
	"github.com/ivan1993spb/snake-server/objects/apple"
	"github.com/ivan1993spb/snake-server/world"
)

type AppleObserver struct{}

func (AppleObserver) Observe(stop <-chan struct{}, w *world.World) {
	go func() {
		// TODO: Create apple counter.
		apple.NewApple(w)
		apple.NewApple(w)
		apple.NewApple(w)
		apple.NewApple(w)
		apple.NewApple(w)
		apple.NewApple(w)
		// TODO: Handler errors.
		// TODO: Create buffer const.
		for event := range w.Events(stop, 32) {
			if event.Type == world.EventTypeObjectDelete {
				if _, ok := event.Payload.(*apple.Apple); ok {
					apple.NewApple(w)
				}
			}
		}
	}()
}
