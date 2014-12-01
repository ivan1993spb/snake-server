package objects

import (
	"errors"

	"bitbucket.org/pushkin_ivan/simple-2d-playground"
)

type Wall playground.DotList

const _WALL_STRENGTH_FACTOR = 1000

func NewWall(pg *playground.Playground, dots playground.DotList,
) (Wall, error) {

	if pg != nil && len(dots) > 0 {
		var w Wall = Wall(dots)

		if err := pg.Locate(w); err == nil {
			return w, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Cannot create wall")
}

// RandomWall generates and locates wall on passed playground pg using
// area factor af
func RandomWall(pg *playground.Playground, af float32) (Wall, error) {
	count := int(float32(pg.GetArea()) * af)
	if count > 0 {
		dots := make(playground.DotList, count)
		for i := range dots {
			dots[i] = pg.RandomDot()
		}
		return NewWall(pg, dots)
	}

	return nil, errors.New("Cannot generate random wall")
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
