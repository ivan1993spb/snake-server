package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/wall"
	"github.com/ivan1993spb/snake-server/world"
)

const wallPerNDots = 100

type WallObserver struct{}

func (WallObserver) Observe(stop <-chan struct{}, w *world.World, logger logrus.FieldLogger) {
	go func() {
		for i := uint16(0); i < w.Size()/wallPerNDots; i++ {
			if _, err := wall.NewRandWall(w); err != nil {
				logger.WithError(err).Error("cannot create rand wall")
			}
		}
	}()
}
