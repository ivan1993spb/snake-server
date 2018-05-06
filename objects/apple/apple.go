package apple

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/game"
)

type Apple struct {
	world    game.WorldInterface
	location engine.Location
}

type ErrCreateApple string

func (e ErrCreateApple) Error() string {
	return "cannot create apple: " + string(e)
}

// CreateApple creates and locates new apple
func CreateApple(world game.WorldInterface) (*Apple, error) {
	apple := &Apple{}

	location, err := world.CreateObjectRandomDot(apple)
	if err != nil {
		return nil, ErrCreateApple(err.Error())
	}

	apple.location = location
	apple.world = world

	return apple, nil
}

func (a *Apple) NutritionalValue(dot *engine.Dot) int8 {
	if a.location.Equals(engine.Location{dot}) {
		a.world.DeleteObject(a, a.location)
	}

	return 0
}
