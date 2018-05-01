// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wall

import (
	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type Wall struct {
	dots engine.Location
}

const wallStrengthFactor = 1000

func CreateWall(dots engine.Location) (*Wall, error) {

	wall := &Wall{dots}

	return wall, nil
}

func CreateLongWall(pg *playground.Playground) (*Wall, error) {
	var (
		pgW = pg.Width()
		pgH = pg.Height()
		err error
		e   *engine.Rect
	)

	wall := &Wall{}

	switch engine.RandomDirection() {
	case engine.DirectionNorth, engine.DirectionSouth:
		e, err = pg.CreateObjectRandomRect(wall, 1, pgH)
	case engine.DirectionEast, engine.DirectionWest:
		e, err = pg.CreateObjectRandomRect(wall, pgW, 1)
	default:
		err = &engine.ErrInvalidDirection{}
	}
	if err != nil {
		return nil, err
	}

	return CreateWall(p, pg, e.Location())
}

// Implementing playground.Location interface
func (w *Wall) DotCount() uint16 {
	return w.dots.DotCount()
}

// Implementing playground.Location interface
func (w *Wall) Dot(i uint16) *engine.Dot {
	if w.dots.DotCount() > i {
		return w.dots[i]
	}

	return nil
}

// Implementing logic.Notalive interface
func (w *Wall) Break(dot *engine.Dot) {
	if w.dots.Contains(dot) {
		w.dots = w.dots.Delete(dot)

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
	return wallStrengthFactor
}
