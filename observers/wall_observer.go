package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/wall"
	"github.com/ivan1993spb/snake-server/world"
)

type WallObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewWallObserver(w world.Interface, logger logrus.FieldLogger) Observer {
	return &WallObserver{
		world:  w,
		logger: logger,
	}
}

func (wo *WallObserver) Observe(stop <-chan struct{}) {
	go wo.run(stop)
}

func (wo *WallObserver) run(stop <-chan struct{}) {
	wo.addRuins()
}

func (wo *WallObserver) addRuins() {
	// TODO: Create abstraction layer for adding of objects.
	if _, err := wall.Ruins(wo.world); err != nil {
		wo.logger.WithError(err).Error("error on ruins creation")
	}
}
