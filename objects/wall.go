package objects

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

type Wall playground.DotList

const _WALL_STRENGTH_FACTOR = 1000

func CreateWall(pg *playground.Playground, dots playground.DotList,
) (Wall, error) {
	if pg == nil {
		return nil, errors.New("Passed nil playground")
	}
	if len(dots) == 0 {
		return nil, errors.New("Passed empty dot list")
	}

	var w Wall = Wall(dots)

	if err := pg.Locate(w); err != nil {
		return nil, err
	}

	return w, nil
}

func CreateLongWall(pg *playground.Playground) (Wall, error) {
	if pg == nil {
		return nil, errors.New("Passed nil playground")
	}

	var (
		pgW, pgH = pg.GetSize()
		err      error
		dots     playground.DotList
	)

	switch playground.RandomDirection() {
	case playground.DIR_NORTH, playground.DIR_SOUTH:
		dots, err = pg.GetEmptyField(1, pgH)
	case playground.DIR_EAST, playground.DIR_WEST:
		dots, err = pg.GetEmptyField(pgW, 1)
	default:
		err = errors.New("Invalid direction")
	}

	if err != nil {
		return nil, err
	}

	return CreateWall(pg, dots)
}

// Implementing playground.Object interface
func (w Wall) DotCount() uint16 {
	return uint16(len(w))
}

// Implementing playground.Object interface
func (w Wall) Dot(i uint16) *playground.Dot {
	if uint16(len(w)) > i {
		return w[i]
	}
	return nil
}

// Implementing playground.Object interface
func (w Wall) Pack() string {
	return playground.DotList(w).Pack()
}

// Implementing logic.Notalive interface
func (w Wall) Break(dot *playground.Dot) {
	if dl := playground.DotList(w); dl.Contains(dot) {
		dl.Delete(dot)
	}
}

// Implementing logic.Resistant interface
func (Wall) Strength() float32 {
	return _WALL_STRENGTH_FACTOR
}
