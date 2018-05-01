package apple

import (
	"github.com/olebedev/emitter"

	"github.com/ivan1993spb/clever-snake/engine"
	"github.com/ivan1993spb/clever-snake/playground"
)

type Apple struct {
	pg       *playground.Playground
	location engine.Location
}

// CreateApple creates and locates new apple
func CreateApple(pg *playground.Playground) (*Apple, error) {
	apple := &Apple{}

	location, err := pg.CreateObjectRandomDot(apple)
	if err != nil {
		return nil, err
	}

	apple.location = location
	apple.pg = pg

	return apple, nil
}

// Implementing playground.Location interface
func (*Apple) DotCount() uint16 {
	return 1
}

// Implementing playground.Location interface
func (a *Apple) Dot(i uint16) *engine.Dot {
	if i == 0 {
		return a.location.Dot(0)
	}

	return nil
}

// Implementing logic.Food interface
func (a *Apple) NutritionalValue(dot *engine.Dot) int8 {
	if a.location.Equals(engine.Location{dot}) {
		a.pg.DeleteObject(a, a.location)
		return 1
	}

	return 0
}

func (a *Apple) Run(emitter *emitter.Emitter) {
	go func() {

	}()
}
