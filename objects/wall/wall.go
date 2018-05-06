package wall

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/game"
)

type Wall struct {
	world    game.WorldInterface
	location engine.Location
}

const wallStrengthFactor = 1000

func CreateWall(world game.WorldInterface, location engine.Location) (*Wall, error) {

	wall := &Wall{world, location}

	return wall, nil
}

func CreateLongWall(world game.WorldInterface) (*Wall, error) {
	var (
		pgW      = world.Width()
		pgH      = world.Height()
		err      error
		location engine.Location
	)

	wall := &Wall{}

	switch engine.RandomDirection() {
	case engine.DirectionNorth, engine.DirectionSouth:
		location, err = world.CreateObjectRandomRect(wall, 1, pgH)
	case engine.DirectionEast, engine.DirectionWest:
		location, err = world.CreateObjectRandomRect(wall, pgW, 1)
	default:
		err = &engine.ErrInvalidDirection{}
	}
	if err != nil {
		return nil, err
	}

	wall.location = location
	wall.world = world

	return wall, nil
}

func (w *Wall) Break(dot *engine.Dot) {
	if w.location.Contains(dot) {
		location := w.location.Delete(dot)

		if w.location.DotCount() > 0 {
			w.world.UpdateObject(w, w.location, location)
			w.location = location
			return
		}
	}

	if w.location.DotCount() == 0 {
		w.world.DeleteObject(w, w.location)
	}
}
