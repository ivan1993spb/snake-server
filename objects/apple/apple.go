package apple

import (
	"fmt"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

type Apple struct {
	world    *world.World
	location engine.Location
}

type ErrCreateApple string

func (e ErrCreateApple) Error() string {
	return "cannot create apple: " + string(e)
}

// NewApple creates and locates new apple
func NewApple(world *world.World) (*Apple, error) {
	apple := &Apple{}

	location, err := world.CreateObjectRandomDot(apple)
	if err != nil {
		return nil, ErrCreateApple(err.Error())
	}

	apple.location = location
	apple.world = world

	return apple, nil
}

func (a *Apple) String() string {
	return fmt.Sprintf("apple %s", a.location)
}

func (a *Apple) NutritionalValue(dot *engine.Dot) uint16 {
	if a.location.Equals(engine.Location{dot}) {
		// TODO: Handle error.
		a.world.DeleteObject(a, a.location)
		return 1
	}

	return 0
}
