package corpse

import (
	"errors"
	"fmt"
	"time"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

// Snakes can eat corpses
type Corpse struct {
	world       *world.World
	location    engine.Location
	nippedPiece *engine.Dot // last nipped piece
	stop        chan struct{}
}

// Corpse are created when a snake dies
func NewCorpse(world *world.World, location engine.Location) (*Corpse, error) {
	// TODO: Check location.

	// TODO: Pass stopper.

	corpse := &Corpse{}

	if location, _ := world.CreateObjectAvailableDots(corpse, location); len(location) > 0 {
		corpse.world = world
		corpse.location = location
		corpse.stop = make(chan struct{}, 0)
		go corpse.run()
		return corpse, nil
	}

	return nil, errors.New("")
}

func (c *Corpse) String() string {
	return fmt.Sprintln("corpse", c.location)
}

func (c *Corpse) NutritionalValue(dot *engine.Dot) uint16 {
	if c.location.Contains(dot) {
		newDots := c.location.Delete(dot)

		if len(c.location) > 0 {
			c.world.UpdateObjectAvailableDots(c, c.location, newDots)
			c.location = newDots
			c.nippedPiece = dot
		} else {
			c.world.DeleteObject(c, c.location)
			close(c.stop)
		}

		return 2
	}

	return 0
}

func (c *Corpse) run() {
	var timer = time.NewTimer(corpseMaxExperience)
	defer timer.Stop()
	select {
	case <-timer.C:
		c.world.DeleteObject(c, c.location)
		close(c.stop)
	case <-c.stop:
	}
}
