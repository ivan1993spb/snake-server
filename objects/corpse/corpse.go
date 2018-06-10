package corpse

import (
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

const corpseNutritionalValue uint16 = 2

const corpseTypeLabel = "corpse"

// Snakes can eat corpses
type Corpse struct {
	uuid     string
	world    *world.World
	location engine.Location
	mux      *sync.RWMutex
	stop     chan struct{}
}

type ErrCreateCorpse string

func (e ErrCreateCorpse) Error() string {
	return "error on corpse creation: " + e.Error()
}

// Corpse are created when a snake dies
func NewCorpse(world *world.World, location engine.Location) (*Corpse, error) {
	if location.Empty() {
		return nil, ErrCreateCorpse("location is empty")
	}

	corpse := &Corpse{
		uuid: uuid.Must(uuid.NewV4()).String(),
		mux:  &sync.RWMutex{},
	}

	location, err := world.CreateObjectAvailableDots(corpse, location)
	if err != nil {
		return nil, ErrCreateCorpse(err.Error())
	}

	if location.Empty() {
		return nil, ErrCreateCorpse("no location available")
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
			// TODO: Handle errors?
			newLoc, _ := c.world.UpdateObjectAvailableDots(c, c.location, newDots)
			c.location = newLoc
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
			c.mux.Lock()
			c.world.DeleteObject(c, c.location)
			close(c.stop)
			c.mux.Unlock()
		case <-c.stop:
			// Corpse was eaten.
		}
	}()
}

func (c *Corpse) MarshalJSON() ([]byte, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return ffjson.Marshal(&corpse{
		UUID: c.uuid,
		Dots: c.location,
		Type: corpseTypeLabel,
	})
}

type corpse struct {
	UUID string          `json:"uuid"`
	Dots engine.Location `json:"dots"`
	Type string          `json:"type"`
}
