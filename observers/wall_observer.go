package observers

import (
	"fmt"

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

type WallObserver struct {
	world  world.Interface
	logger logrus.FieldLogger
	area   engine.Area

	ruinsCount uint16
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
	if err := wo.init(); err != nil {
		wo.logger.WithError(err).Error("error in wall observer")
		return
	}

	wo.addWalls()
}

func (wo *WallObserver) init() error {
	area, err := engine.NewArea(wo.world.Width(), wo.world.Height())
	if err != nil {
		return fmt.Errorf("cannot create area: %s", err)
	}

	wo.area = area
	wo.ruinsCount = wo.calcRuinsCount(wo.area.Size())

	return nil
}

func (wo *WallObserver) calcRuinsCount(size uint16) uint16 {
	return uint16(float32(size) * ruinsFactor)
}

func (wo *WallObserver) addWalls() {
	var counter uint16

	for counter < wo.ruinsCount {
		for i := 0; i < len(ruins); i++ {
			// Pass one of the ruins and maximum dots count to occupy by new wall
			counter += wo.addWallFromMask(ruins[i], wo.ruinsCount-counter)
		}
	}
}

func (wo *WallObserver) addWallFromMask(mask *engine.DotsMask, dotsLimit uint16) (dotsResult uint16) {
	mask = mask.TurnRandom()

	if wo.area.Width() < mask.Width() || wo.area.Height() < mask.Height() {
		return
	}

	rect, err := wo.area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
	if err != nil {
		return
	}

	location := mask.Location(rect.X(), rect.Y())
	if location.DotCount() > dotsLimit {
		location = location[:dotsLimit]
	}

	if wo.world.LocationOccupied(location) {
		return
	}

	// TODO: Create abstraction layer for adding of objects.
	if _, err := wall.NewWallLocation(wo.world, location); err != nil {
		wo.logger.WithError(err).Error("error on wall creation")
		return
	}

	dotsResult = location.DotCount()

	return
}
