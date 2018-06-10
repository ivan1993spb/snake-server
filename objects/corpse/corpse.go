package corpse

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

const corpseNutritionalValue uint16 = 2

// Snakes can eat corpses
type Corpse struct {
	world       *world.World
	location    engine.Location
	nippedPiece engine.Dot // last nipped piece
	mux         *sync.RWMutex
	stop        chan struct{}
}

// Corpse are created when a snake dies
func NewCorpse(world *world.World, location engine.Location) (*Corpse, error) {
	if location.Empty() {
		return nil, errors.New("location is empty")
	}

	corpse := &Corpse{
		mux: &sync.RWMutex{},
	}

	location, err := world.CreateObjectAvailableDots(corpse, location)
	if err != nil {
		// TODO:
		fmt.Println("create", err)
		return nil, err
	}

	if len(location) == 0 {
		return nil, errors.New("no location available")
	}

	corpse.mux.Lock()
	corpse.world = world
	corpse.location = location
	corpse.stop = make(chan struct{})
	corpse.mux.Unlock()

	return corpse, nil
}

func (c *Corpse) String() string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return fmt.Sprint("corpse ", c.location)
}

func (c *Corpse) NutritionalValue(dot engine.Dot) uint16 {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.location.Contains(dot) {
		newDots := c.location.Delete(dot)

		if len(newDots) > 0 {
			// TODO: Handle errors.
			newLoc, _ := c.world.UpdateObjectAvailableDots(c, c.location, newDots)
			c.location = newLoc
			c.nippedPiece = dot
		} else {
			c.world.DeleteObject(c, c.location)
			close(c.stop)
		}

		return corpseNutritionalValue
	}

	return 0
}

func (c *Corpse) Run(stop <-chan struct{}) {
	go func() {
		var timer = time.NewTimer(corpseMaxExperience)
		defer timer.Stop()
		select {
		case <-stop:
			// global stop
		case <-timer.C:
			c.mux.RLock()
			c.world.DeleteObject(c, c.location)
			c.mux.RUnlock()
			close(c.stop)
		case <-c.stop:
			// Corpse was eaten.
		}
	}()
}

func (c *Corpse) MarshalJSON() ([]byte, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return ffjson.Marshal(c.location)
}
