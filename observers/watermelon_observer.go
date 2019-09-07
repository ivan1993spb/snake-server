package observers

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/watermelon"
	"github.com/ivan1993spb/snake-server/world"
)

const chanWatermelonObserverEventsBuffer = 64

const addWatermelonDelay = time.Second * 15

const addWatermelonsDuringTickLimit = 2

const oneWatermelonArea = 200

type WatermelonObserver struct {
	world  world.Interface
	logger logrus.FieldLogger

	watermelonCount    int32
	maxWatermelonCount int32
}

func NewWatermelonObserver(w world.Interface, logger logrus.FieldLogger) Observer {
	return &WatermelonObserver{
		world:  w,
		logger: logger,
	}
}

func (wo *WatermelonObserver) Observe(stop <-chan struct{}) {
	go wo.run(stop)
}

func (wo *WatermelonObserver) run(stop <-chan struct{}) {
	wo.init()

	if wo.maxWatermelonCount > 0 {
		go wo.schedule(stop)

		wo.listen(stop)
	}
}

func (wo *WatermelonObserver) init() {
	maxWatermelonCount := wo.calcMaxWatermelonCount()

	wo.logger.WithFields(logrus.Fields{
		"watermelon_count": maxWatermelonCount,
	}).Debug("watermelon observer")

	wo.maxWatermelonCount = maxWatermelonCount
}

// calcMaxWatermelonCount returns max possible watermelon count
func (wo *WatermelonObserver) calcMaxWatermelonCount() int32 {
	var size = int32(wo.world.Area().Size())
	var maxWatermelonCount = size / oneWatermelonArea
	return maxWatermelonCount
}

func (wo *WatermelonObserver) listen(stop <-chan struct{}) {
	for event := range wo.world.Events(stop, chanWatermelonObserverEventsBuffer) {
		wo.handleEvent(event)
	}
}

func (wo *WatermelonObserver) handleEvent(event world.Event) {
	if event.Type != world.EventTypeObjectDelete {
		return
	}

	if _, ok := event.Payload.(*watermelon.Watermelon); ok {
		atomic.AddInt32(&wo.watermelonCount, -1)
	}
}

func (wo *WatermelonObserver) schedule(stop <-chan struct{}) {
	ticker := time.NewTicker(addWatermelonDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wo.addWatermelons()
		case <-stop:
			return
		}
	}
}

func (wo *WatermelonObserver) addWatermelons() {
	var watermelonsAdded = 0

	for {
		if atomic.LoadInt32(&wo.watermelonCount) >= wo.maxWatermelonCount {
			return
		}

		if watermelonsAdded >= addWatermelonsDuringTickLimit {
			return
		}

		// TODO: Create abstraction layer for adding of objects.
		if _, err := watermelon.NewWatermelon(wo.world); err != nil {
			wo.logger.WithError(err).Error("cannot create watermelon")
			return
		}

		atomic.AddInt32(&wo.watermelonCount, 1)
		watermelonsAdded++
	}
}
