package wall

import (
	"sync"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const ruinsFactor = 0.20

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
