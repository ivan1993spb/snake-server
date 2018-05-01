package mouse

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type Mouse struct {
	pg        *playground.Playground
	location  engine.Location
	direction engine.Direction
}

func NewMouse(pg *playground.Playground) *Mouse {
	mouse := &Mouse{}
	location, err := pg.CreateObjectRandomDot(mouse)
	if err != nil {
		// TODO: return error
		return nil
	}

	return &Mouse{
		pg:        pg,
		location:  location,
		direction: engine.RandomDirection(),
	}
}
