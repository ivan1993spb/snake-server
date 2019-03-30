package observers

import (
	"fmt"

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
	go ao.run(stop)
}

func (ao *AppleObserver) run(stop <-chan struct{}) {
	ao.init()
	ao.listen(stop)
}

func (ao *AppleObserver) init() {
	for i := 0; i < ao.calcAppleCount(); i++ {
		if _, err := apple.NewApple(ao.world); err != nil {
			ao.logger.WithError(err).Error("cannot create apple")
		}
	}
}

func (ao *AppleObserver) listen(stop <-chan struct{}) {
	for event := range ao.world.Events(stop, chanAppleObserverEventsBuffer) {
		if err := ao.handleEvent(event); err != nil {
			ao.logger.WithError(err).Error("handling event error")
		}
	}
}

func (ao *AppleObserver) calcAppleCount() int {
	appleCount := defaultAppleCount
	size := ao.world.Size()

	if size > oneAppleArea {
		appleCount = int(size / oneAppleArea)
	}

	return appleCount
}

func (ao *AppleObserver) handleEvent(event world.Event) error {
	// Event type is only delete
	if event.Type != world.EventTypeObjectDelete {
		return nil
	}

	// Object is only apple
	if _, ok := event.Payload.(*apple.Apple); !ok {
		return nil
	}

	if _, err := apple.NewApple(ao.world); err != nil {
		return fmt.Errorf("cannot create apple: %s", err)
	}

	return nil
}
