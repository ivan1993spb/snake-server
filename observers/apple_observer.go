package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/apple"
	"github.com/ivan1993spb/snake-server/world"
)

const chanAppleObserverEventsBuffer = 32

const defaultAppleCount = 1

const oneAppleArea = 50

type AppleObserver struct{}

func (AppleObserver) Observe(stop <-chan struct{}, w *world.World, logger logrus.FieldLogger) {
	go func() {
		appleCount := defaultAppleCount
		size := w.Size()

		if size > oneAppleArea {
			appleCount = int(size / oneAppleArea)
		}

		logger.Debugf("apple count for size %d = %d", size, appleCount)

		for i := 0; i < appleCount; i++ {
			if _, err := apple.NewApple(w); err != nil {
				logger.WithError(err).Error("cannot create apple")
			}
		}

		for event := range w.Events(stop, chanAppleObserverEventsBuffer) {
			if event.Type == world.EventTypeObjectDelete {
				if _, ok := event.Payload.(*apple.Apple); ok {
					if _, err := apple.NewApple(w); err != nil {
						logger.WithError(err).Error("cannot create apple")
					}
				}
			}
		}
	}()
}
