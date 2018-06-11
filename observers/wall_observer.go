package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/wall"
	"github.com/ivan1993spb/snake-server/world"
)

type WallObserver struct{}

func (WallObserver) Observe(stop <-chan struct{}, w *world.World, logger logrus.FieldLogger) {
	go func() {
		if _, err := wall.NewWallRuins(w); err != nil {
			logger.WithError(err).Error("cannot create wall ruins")
		}
	}()
}
