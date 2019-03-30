package observers

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/objects/wall"
	"github.com/ivan1993spb/snake-server/world"
)

const ruinsFactor = 0.15

var dotsMaskOne = engine.NewDotsMask([][]uint8{{1}})

var ruins = []*engine.DotsMask{
	dotsMaskOne,
	engine.DotsMaskSquare2x2,
	engine.DotsMaskTank,
	engine.DotsMaskHome1,
	engine.DotsMaskHome2,
	engine.DotsMaskCross,
	engine.DotsMaskDiagonal,
	engine.DotsMaskCrossSmall,
	engine.DotsMaskDiagonalSmall,
	engine.DotsMaskLabyrinth,
	engine.DotsMaskTunnel1,
	engine.DotsMaskTunnel2,
	engine.DotsMaskBigHome,
}

func calcRuinsCount(size uint16) uint16 {
	return uint16(float32(size) * ruinsFactor)
}

type WallObserver struct{}

func (WallObserver) Observe(stop <-chan struct{}, w world.Interface, logger logrus.FieldLogger) {
	go func() {
		area, err := engine.NewArea(w.Width(), w.Height())
		if err != nil {
			logger.WithError(err).Error("cannot create area in wall observer")
			return
		}

		size := area.Size()
		ruinsCount := calcRuinsCount(size)
		var counter uint16

		for counter < ruinsCount {
			for i := 0; i < len(ruins); i++ {
				mask := ruins[i].TurnRandom()

				if area.Width() >= mask.Width() && area.Height() >= mask.Height() {
					rect, err := area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
					if err != nil {
						continue
					}

					location := mask.Location(rect.X(), rect.Y())
					if location.DotCount() > ruinsCount-counter {
						location = location[:ruinsCount-counter]
					}

					if w.LocationOccupied(location) {
						continue
					}

					if _, err := wall.NewWallLocation(w, location); err != nil {
						logger.WithError(err).Error("error on wall creation")
					} else {
						counter += location.DotCount()
					}
				}
			}
		}
	}()
}
