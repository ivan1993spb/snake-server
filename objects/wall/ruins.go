package wall

import (
	"fmt"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const ruinsFactor = 0.15

var dotsMaskOne = engine.NewDotsMask([][]uint8{{1}})

var masks = []*engine.DotsMask{
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

func calcRuinsAreaLimit(size uint16) uint16 {
	return uint16(float32(size) * ruinsFactor)
}

type ErrGenerateRuins string

func (e ErrGenerateRuins) Error() string {
	return "generate ruins error: " + string(e)
}

const locationFindingSuccessiveErrorLimit = 16
const locationOccupiedSuccessiveLimit = 16
const newWallLocationSuccessiveErrorLimit = 8

func GenerateRuins(w world.Interface) ([]*Wall, error) {
	area, err := engine.NewArea(w.Width(), w.Height())
	if err != nil {
		return nil, ErrGenerateRuins("cannot create area: " + err.Error())
	}

	ruinsAreaLimit := calcRuinsAreaLimit(area.Size())

	walls := make([]*Wall, 0)

	var areaOccupiedSum uint16

	errLocationFindCounter := 0
	locationOccupiedCounter := 0
	errNewWallLocationCounter := 0

	for areaOccupiedSum < ruinsAreaLimit {
		for _, mask := range masks {
			location, err := findLocationLimit(area, mask, ruinsAreaLimit-areaOccupiedSum)
			if errorLimitReached(err != nil, &errLocationFindCounter, locationFindingSuccessiveErrorLimit) {
				return walls, ErrGenerateRuins("too many successive errors in location finding")
			}

			isOccupied := w.LocationOccupied(location)
			if errorLimitReached(isOccupied, &locationOccupiedCounter, locationOccupiedSuccessiveLimit) {
				return walls, ErrGenerateRuins("too many successive errors in location finding: occupied")
			}

			wall, err := NewWallLocation(w, location)
			if errorLimitReached(err != nil, &errNewWallLocationCounter, newWallLocationSuccessiveErrorLimit) {
				return walls, ErrGenerateRuins("to many successive errors on creating wall: " + err.Error())
			}

			areaOccupiedSum += wall.location.DotCount()
			walls = append(walls, wall)
		}
	}

	return walls, nil
}

func errorLimitReached(errorOccurred bool, counter *int, limit int) bool {
	if errorOccurred {
		*counter++
		return *counter >= limit
	}
	*counter = 0
	return false
}

func findLocationLimit(area engine.Area, mask *engine.DotsMask, limit uint16) (engine.Location, error) {
	mask = mask.TurnRandom()

	if area.Width() < mask.Width() || area.Height() < mask.Height() {
		return nil, fmt.Errorf("mask doesn't fit the area")
	}

	rect, err := area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get random rect: %s", err)
	}

	location := mask.Location(rect.X(), rect.Y())
	if location.DotCount() > limit {
		location = location[:limit]
	}

	return location, nil
}
