package corpse

import (
	"errors"
	"time"

	"github.com/olebedev/emitter"

	"github.com/ivan1993spb/clever-snake/engine"
	"github.com/ivan1993spb/clever-snake/playground"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

// Snakes can eat corpses
type Corpse struct {
	pg          *playground.Playground
	location    engine.Location
	nippedPiece *engine.Dot // last nipped piece
}

// Corpses are created when a snake dies
func CreateCorpse(pg *playground.Playground, location engine.Location) (*Corpse, error) {
	// TODO: Check location

	corpse := &Corpse{}

	if location, _ := pg.CreateObjectAvailableDots(corpse, location); len(location) > 0 {
		return &Corpse{
			pg:          pg,
			location:    location,
			nippedPiece: nil,
		}, nil
	}
	return nil, errors.New("")
}

// Implementing playground.Location interface
func (c *Corpse) DotCount() uint16 {
	return uint16(len(c.location))
}

// Implementing playground.Location interface
func (c *Corpse) Dot(i uint16) *engine.Dot {
	if c.DotCount() > i {
		return c.location[i]
	}
	return nil
}

// Implementing logic.Food interface
func (c *Corpse) NutritionalValue(dot *engine.Dot) int8 {
	if c.location.Contains(dot) {
		newDots := c.location.Delete(dot)

		if len(c.location) > 0 {
			c.pg.UpdateObjectAvailableDots(c, c.location, newDots)
			c.location = newDots
			c.nippedPiece = dot
		} else {
			c.pg.DeleteObject(c, c.location)
		}

		return 2
	}

	return 0
}

func (c *Corpse) Run(emitter *emitter.Emitter) {
	go func() {

	}()
}
