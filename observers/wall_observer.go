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
	wo.generateRuins()
}

func (wo *WallObserver) generateRuins() {
	ruinsGenerator, err := wall.NewRuinsGenerator(wo.world)
	if err != nil {
		wo.logger.WithError(err).Error("cannot create ruins generator")
		return
	}

	for !ruinsGenerator.Done() {
		if err := ruinsGenerator.Err(); err != nil {
			wo.logger.WithError(err).Error("error on ruins generation: interrupted")
			break
		}

		if _, err := ruinsGenerator.GenerateWall(); err != nil {
			wo.logger.WithError(err).Error("error on ruins generation")
		}
	}
}
