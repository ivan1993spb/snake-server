package corpse

import (
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

const corpseNutritionalValue uint16 = 2

const corpseTypeLabel = "corpse"

// Snakes can eat corpses
type Corpse struct {
	uuid      string
	world     *world.World
	location  engine.Location
	mux       *sync.RWMutex
	stop      chan struct{}
	isStopped bool
}

type ErrCreateCorpse string

func (e ErrCreateCorpse) Error() string {
	return "error on corpse creation: " + string(e)
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

type errCorpseBite string

func (e errCorpseBite) Error() string {
	return "corpse bite error: " + string(e)
}

func (c *Corpse) Bite(dot engine.Dot) (nv uint16, success bool, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.location.Contains(dot) {
		newDots := c.location.Delete(dot)

		if len(newDots) > 0 {
			newLocation, err := c.world.UpdateObjectAvailableDots(c, c.location, newDots)
			if err != nil {
				return 0, false, errCorpseBite(err.Error())
			}
			if len(newLocation) > 0 {
				c.location = newLocation
				return corpseNutritionalValue, true, nil
			}
		}

		if !c.isStopped {
			close(c.stop)
			c.isStopped = true
		}

		if err := c.world.DeleteObject(c, c.location); err != nil {
			return 0, false, errCorpseBite(err.Error())
		}

		c.location = c.location[:0]

		return corpseNutritionalValue, true, nil
	}

	return 0, false, nil
}

func (c *Corpse) Run(stop <-chan struct{}, logger logrus.FieldLogger) {
	go func() {
		var timer = time.NewTimer(corpseMaxExperience)
		defer timer.Stop()
		select {
		case <-stop:
			// global stop
		case <-timer.C:
			c.mux.Lock()

			if !c.isStopped {
				close(c.stop)
				c.isStopped = true
			}

			if err := c.world.DeleteObject(c, c.location); err != nil {
				logger.WithError(err).Error("corpse stop error")
			}

			c.location = c.location[:0]

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
	Dots engine.Location `json:"dots,omitempty"`
	Type string          `json:"type"`
}
