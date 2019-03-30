package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/world"
)

const chanLoggerObserverEventsBuffer = 64

type LoggerObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewLoggerObserver(w world.Interface, logger logrus.FieldLogger) Observer {
	return &LoggerObserver{
		world:  w,
		logger: logger,
	}
}

func (lo *LoggerObserver) Observe(stop <-chan struct{}) {
	go lo.run(stop)
}

func (lo *LoggerObserver) run(stop <-chan struct{}) {
	lo.listen(stop)
}

func (lo *LoggerObserver) listen(stop <-chan struct{}) {
	for event := range lo.world.Events(stop, chanLoggerObserverEventsBuffer) {
		lo.handleEvent(event)
	}
}

func (lo *LoggerObserver) handleEvent(event world.Event) {
	switch event.Type {
	case world.EventTypeError:
		if err, ok := event.Payload.(error); ok {
			lo.logger.WithError(err).Error("world error")
		}
	case world.EventTypeObjectCreate, world.EventTypeObjectDelete, world.EventTypeObjectUpdate, world.EventTypeObjectChecked:
		lo.logger.WithFields(logrus.Fields{
			"payload": event.Payload,
			"type":    event.Type,
		}).Debug("world event")
	}
}
