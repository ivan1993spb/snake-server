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

type WatermelonObserver struct{}

func (WatermelonObserver) Observe(stop <-chan struct{}, w world.Interface, logger logrus.FieldLogger) {
	var size = int32(w.Size())
	var maxWatermelonCount = size / oneWatermelonArea

	logger.WithFields(logrus.Fields{
		"map_size":         size,
		"watermelon_count": maxWatermelonCount,
	}).Debug("watermelon observer")

	if maxWatermelonCount == 0 {
		return
	}

	var watermelonCount int32 = 0

	go func() {
		for event := range w.Events(stop, chanWatermelonObserverEventsBuffer) {
			if event.Type == world.EventTypeObjectDelete {
				if _, ok := event.Payload.(*watermelon.Watermelon); ok {
					atomic.AddInt32(&watermelonCount, -1)
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(addWatermelonDelay)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				var watermelonsAddedDuringTick = 0

				for atomic.LoadInt32(&watermelonCount) < maxWatermelonCount && watermelonsAddedDuringTick < addWatermelonsDuringTickLimit {
					if _, err := watermelon.NewWatermelon(w); err != nil {
						logger.WithError(err).Error("cannot create watermelon")
					} else {
						atomic.AddInt32(&watermelonCount, 1)
						watermelonsAddedDuringTick++
					}
				}
			case <-stop:
				return
			}
		}
	}()
}
