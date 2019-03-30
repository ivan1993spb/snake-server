package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/apple"
	"github.com/ivan1993spb/snake-server/world"
)

const chanAppleObserverEventsBuffer = 64

const defaultAppleCount = 1

const oneAppleArea = 50

type AppleObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewAppleObserver(w world.Interface, logger logrus.FieldLogger) Observer {
	return &AppleObserver{
		world:  w,
		logger: logger,
	}
}

func (ao *AppleObserver) Observe(stop <-chan struct{}) {
	go func() {
		appleCount := defaultAppleCount
		size := ao.world.Size()

		if size > oneAppleArea {
			appleCount = int(size / oneAppleArea)
		}

		ao.logger.Debugf("apple count for size %d = %d", size, appleCount)

		for i := 0; i < appleCount; i++ {
			if _, err := apple.NewApple(ao.world); err != nil {
				ao.logger.WithError(err).Error("cannot create apple")
			}
		}

		for event := range ao.world.Events(stop, chanAppleObserverEventsBuffer) {
			if event.Type == world.EventTypeObjectDelete {
				if _, ok := event.Payload.(*apple.Apple); ok {
					if _, err := apple.NewApple(ao.world); err != nil {
						ao.logger.WithError(err).Error("cannot create apple")
					}
				}
			}
		}
	}()
}
