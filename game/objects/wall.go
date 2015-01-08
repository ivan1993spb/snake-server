package objects

import "bitbucket.org/pushkin_ivan/clever-snake/game/playground"

type Wall struct {
	p    GameProcessor
	pg   *playground.Playground
	dots playground.DotList
}

const _WALL_STRENGTH_FACTOR = 1000

func CreateWall(p GameProcessor, pg *playground.Playground,
	dots playground.DotList) (*Wall, error) {
	if p == nil {
		return nil, &errCreateObject{errNilGameProcessor}
	}
	if pg == nil {
		return nil, &errCreateObject{errNilPlayground}
	}
	if len(dots) == 0 {
		return nil, &errCreateObject{errEmptyDotList}
	}

	wall := &Wall{p, pg, dots}

	if err := pg.Locate(wall, true); err != nil {
		return nil, &errCreateObject{err}
	}

	p.OccurredCreating(wall)

	return wall, nil
}

func CreateLongWall(p GameProcessor, pg *playground.Playground,
) (*Wall, error) {
	var (
		pgW, pgH = pg.GetSize()
		err      error
		e        playground.Entity
	)

	switch playground.RandomDirection() {
	case playground.DIR_NORTH, playground.DIR_SOUTH:
		e, err = pg.GetRandomEmptyRect(1, pgH)
	case playground.DIR_EAST, playground.DIR_WEST:
		e, err = pg.GetRandomEmptyRect(pgW, 1)
	default:
		err = playground.ErrInvalidDirection
	}
	if err != nil {
		return nil, err
	}

	return CreateWall(p, pg, playground.EntityToDotList(e))
}

// Implementing playground.Object interface
func (w *Wall) DotCount() uint16 {
	return w.dots.DotCount()
}

// Implementing playground.Object interface
func (w *Wall) Dot(i uint16) *playground.Dot {
	if w.dots.DotCount() > i {
		return w.dots[i]
	}

	return nil
}

// Implementing logic.Notalive interface
func (w *Wall) Break(dot *playground.Dot) {
	if w.dots.Contains(dot) {
		w.dots.Delete(dot)

		if w.dots.DotCount() > 0 {
			w.p.OccurredUpdating(w)
			return
		}
	}

	if w.dots.DotCount() == 0 {
		w.pg.Delete(w)
		w.p.OccurredDeleting(w)
	}
}

// Implementing logic.Resistant interface
func (*Wall) Strength() float32 {
	return _WALL_STRENGTH_FACTOR
}
