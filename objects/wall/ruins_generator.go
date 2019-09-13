package wall

import (
	"fmt"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const (
	ruinsFactorAreaTiny     = 0
	ruinsFactorAreaSmall    = 0.09
	ruinsFactorAreaMedium   = 0.12
	ruinsFactorAreaLarge    = 0.15
	ruinsFactorAreaEnormous = 0.16

	sizeAreaTiny   = 15 * 15
	sizeAreaSmall  = 50 * 50
	sizeAreaMedium = 100 * 100
	sizeAreaLarge  = 150 * 150
)

const (
	findLocationSuccessiveErrorLimit    = 16
	locationOccupiedSuccessiveLimit     = 16
	newWallLocationSuccessiveErrorLimit = 8
)

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

func getRuinsFactor(size uint16) float32 {
	if size <= sizeAreaTiny {
		return ruinsFactorAreaTiny
	}
	if size <= sizeAreaSmall {
		return ruinsFactorAreaSmall
	}
	if size <= sizeAreaMedium {
		return ruinsFactorAreaMedium
	}
	if size <= sizeAreaLarge {
		return ruinsFactorAreaLarge
	}
	return ruinsFactorAreaEnormous
}

func calcRuinsAreaLimit(size uint16) uint16 {
	return uint16(float32(size) * getRuinsFactor(size))
}

type RuinsGenerator struct {
	world world.Interface
	area  engine.Area

	ruinsAreaLimit  uint16
	areaOccupiedSum uint16

	errFindLocationCounter    int
	locationOccupiedCounter   int
	errNewWallLocationCounter int

	maskIndex int
}

type ErrCreateRuinsGenerator string

func (e ErrCreateRuinsGenerator) Error() string {
	return "cannot create ruins generator: " + string(e)
}

func NewRuinsGenerator(w world.Interface) *RuinsGenerator {
	area := w.Area()

	return &RuinsGenerator{
		world: w,
		area:  area,

		ruinsAreaLimit: calcRuinsAreaLimit(area.Size()),
	}
}

func (rg *RuinsGenerator) Done() bool {
	return rg.areaOccupiedSum == rg.ruinsAreaLimit
}

type ErrGenerateWall string

func (e ErrGenerateWall) Error() string {
	return "generate wall error: " + string(e)
}

func (rg *RuinsGenerator) Err() error {
	if rg.errFindLocationCounter >= findLocationSuccessiveErrorLimit {
		return ErrGenerateWall("too many successive errors in location finding")
	}
	if rg.locationOccupiedCounter >= locationOccupiedSuccessiveLimit {
		return ErrGenerateWall("too many successive occupied locations")
	}
	if rg.errNewWallLocationCounter >= newWallLocationSuccessiveErrorLimit {
		return ErrGenerateWall("to many successive errors on wall creating")
	}
	return nil
}

func (rg *RuinsGenerator) GenerateWall() (*Wall, error) {
	if rg.Err() != nil {
		return nil, rg.Err()
	}

	if rg.areaOccupiedSum == rg.ruinsAreaLimit {
		return nil, ErrGenerateWall("ruins generation has been done")
	}

	mask := rg.getMask()

	location, err := rg.findLocation(mask)
	if err != nil {
		rg.errFindLocationCounter++
		return nil, ErrGenerateWall("find location error: " + err.Error())
	}
	rg.errFindLocationCounter = 0

	if rg.world.LocationOccupied(location) {
		rg.locationOccupiedCounter++
		return nil, ErrGenerateWall("location occupied")
	}
	rg.locationOccupiedCounter = 0

	wall, err := NewWallLocation(rg.world, location)
	if err != nil {
		rg.errNewWallLocationCounter++
		return nil, ErrGenerateWall("new wall error: " + err.Error())
	}
	rg.errNewWallLocationCounter = 0

	rg.areaOccupiedSum += wall.location.DotCount()

	return nil, err
}

func (rg *RuinsGenerator) getMask() *engine.DotsMask {
	if rg.maskIndex >= len(masks) {
		rg.maskIndex = 0
	}

	if rg.maskIndex < len(masks)-1 {
		mask := masks[rg.maskIndex]
		rg.maskIndex++
		return mask
	}

	mask := masks[rg.maskIndex]
	rg.maskIndex = 0
	return mask
}

func (rg *RuinsGenerator) findLocation(mask *engine.DotsMask) (engine.Location, error) {
	if rg.areaOccupiedSum == rg.ruinsAreaLimit {
		return engine.Location{}, nil
	}

	mask = mask.TurnRandom()

	if rg.area.Width() < mask.Width() || rg.area.Height() < mask.Height() {
		return nil, fmt.Errorf("mask doesn't fit the area")
	}

	rect, err := rg.area.NewRandomRect(mask.Width(), mask.Height(), 0, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot get random rect: %s", err)
	}

	location := mask.Location(rect.X(), rect.Y())
	limit := rg.ruinsAreaLimit - rg.areaOccupiedSum
	if location.DotCount() > limit {
		location = location[:limit]
	}

	return location, nil
}
