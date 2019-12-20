package mouse_observer

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/mouse"
	"github.com/ivan1993spb/snake-server/observers"
	"github.com/ivan1993spb/snake-server/world"
)

const mouseTickerDelay = time.Minute

const chanMouseObserverEventsBuffer = 64

const oneMouseArea = 400

type MouseObserver struct {
	world  world.Interface
	logger logrus.FieldLogger

	mouseNumber    int32
	maxMouseNumber int32
}

func NewMouseObserver(w world.Interface, logger logrus.FieldLogger) observers.Observer {
	return &MouseObserver{
		world:  w,
		logger: logger,
	}
}

func (mo *MouseObserver) Observe(stop <-chan struct{}) {
	go mo.run(stop)
}

func (mo *MouseObserver) run(stop <-chan struct{}) {
	mo.init()

	if mo.maxMouseNumber > 0 {
		go mo.schedule(stop)

		mo.listen(stop)
	}
}

func (mo *MouseObserver) init() {
	maxMouseNumber := mo.calcMaxMouseCount()

	mo.logger.WithFields(logrus.Fields{
		"mouse_count": maxMouseNumber,
	}).Debug("mouse observer")

	mo.maxMouseNumber = maxMouseNumber
}

func (mo *MouseObserver) calcMaxMouseCount() int32 {
	var size = int32(mo.world.Area().Size())
	var maxMouseNumber = size / oneMouseArea
	return maxMouseNumber
}

func (mo *MouseObserver) schedule(stop <-chan struct{}) {
	ticker := time.NewTicker(mouseTickerDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mo.addMouse(stop)
		case <-stop:
			return
		}
	}
}

const addMouseDuringTickLimit = 1

func (mo *MouseObserver) addMouse(stop <-chan struct{}) {
	var mouseAdded = 0

	for {
		if atomic.LoadInt32(&mo.mouseNumber) >= mo.maxMouseNumber {
			break
		}

		if mouseAdded >= addMouseDuringTickLimit {
			break
		}

		// TODO: Create abstraction layer for adding of objects.
		if m, err := mouse.NewMouse(mo.world); err != nil {
			mo.logger.WithError(err).Error("cannot create mouse")
			break
		} else {
			m.Run(stop)
		}

		atomic.AddInt32(&mo.mouseNumber, 1)
		mouseAdded++
	}
}

func (mo *MouseObserver) listen(stop <-chan struct{}) {
	for event := range mo.world.Events(stop, chanMouseObserverEventsBuffer) {
		mo.handleEvent(event)
	}
}

func (mo *MouseObserver) handleEvent(event world.Event) {
	if event.Type != world.EventTypeObjectDelete {
		return
	}

	if _, ok := event.Payload.(*mouse.Mouse); ok {
		atomic.AddInt32(&mo.mouseNumber, -1)
	}
}
