package wall

import (
	"fmt"
	"sync"

	"github.com/ivan1993spb/snake-server/engine"
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

func calcRuinsCount(size uint16) int {
	return int(float32(size) * ruinsFactor)
}

type ErrCreateWallRuins string

func (e ErrCreateWallRuins) Error() string {
	return "cannot create wall ruins: " + string(e)
}

func NewWallRuins(world world.Interface) (*Wall, error) {
	area, err := engine.NewArea(world.Width(), world.Height())
	if err != nil {
		return nil, ErrCreateWallRuins(err.Error())
	}

	size := area.Size()
	ruinsCount := calcRuinsCount(size)
	wallLocation := make(engine.Location, 0, ruinsCount)

	for len(wallLocation) < ruinsCount {
		for i := 0; i < len(ruins); i++ {
			mask := ruins[i].TurnRandom()

			if area.Width() >= mask.Width() && area.Height() >= mask.Height() {
				rect, err := area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
				if err != nil {
					continue
				}

				location := mask.Location(rect.X(), rect.Y())
				if len(location) > ruinsCount-len(wallLocation) {
					location = location[:ruinsCount-len(wallLocation)]
				}

				wallLocation = append(wallLocation, location...)
			}
		}
	}

	wall := &Wall{
		id:    world.ObtainIdentifier(),
		world: world,
		mux:   &sync.RWMutex{},
	}

	wall.mux.Lock()
	defer wall.mux.Unlock()

	if resultLocation, err := world.CreateObjectAvailableDots(wall, wallLocation); err != nil {
		world.ReleaseIdentifier(wall.id)
		return nil, ErrCreateWallRuins(err.Error())
	} else {
		wall.location = resultLocation
	}

	return wall, nil
}

type ErrRuins string

func (e ErrRuins) Error() string {
	return "ruins error: " + string(e)
}

func Ruins(w world.Interface) ([]*Wall, error) {
	area, err := engine.NewArea(w.Width(), w.Height())
	if err != nil {
		return nil, ErrRuins("cannot create area: " + err.Error())
	}

	ruinsCount := uint16(float32(area.Size()) * ruinsFactor)

	walls := make([]*Wall, 0)

	var counter uint16

	for counter < ruinsCount {
		for i := 0; i < len(ruins); i++ {
			// Pass one of the ruins and maximum dots count to occupy by new wall
			wall, dotsResult, err := addWallFromMask(w, area, ruins[i], ruinsCount-counter)
			if err != nil {
				return walls, err
			}

			counter += dotsResult
			walls = append(walls, wall)
		}
	}

	return walls, nil
}

func addWallFromMask(w world.Interface, area engine.Area, mask *engine.DotsMask, dotsLimit uint16) (*Wall, uint16, error) {
	mask = mask.TurnRandom()

	if area.Width() < mask.Width() || area.Height() < mask.Height() {
		return nil, 0, nil
	}

	rect, err := area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get random rect: %s", err)
	}

	location := mask.Location(rect.X(), rect.Y())
	if location.DotCount() > dotsLimit {
		location = location[:dotsLimit]
	}

	if w.LocationOccupied(location) {
		return nil, 0, nil
	}

	wall, err := NewWallLocation(w, location)

	if err != nil {
		return nil, 0, fmt.Errorf("cannot create wall: %s", err)
	}

	dotsResult := location.DotCount()

	return wall, dotsResult, nil
}
