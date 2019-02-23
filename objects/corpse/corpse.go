package corpse

import (
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

// Time for which corpse will be lie on playground
const corpseMaxExperience = time.Second * 15

const corpseNutritionalValue uint16 = 2

const corpseTypeLabel = "corpse"

// Snakes can eat corpses
// ffjson: skip
type Corpse struct {
	id       world.Identifier
	world    *world.World
	location engine.Location
	mux      *sync.RWMutex
	stop     chan struct{}
	stopper  *sync.Once
}

type errCreateCorpse string

func (e errCreateCorpse) Error() string {
	return "error on corpse creation: " + string(e)
}

// Corpse are created when a snake dies
func NewCorpse(world *world.World, location engine.Location) (*Corpse, error) {
	if location.Empty() {
		return nil, errCreateCorpse("location is empty")
	}

	corpse := &Corpse{
		id:      world.ObtainIdentifier(),
		world:   world,
		mux:     &sync.RWMutex{},
		stop:    make(chan struct{}),
		stopper: &sync.Once{},
	}

	corpse.mux.Lock()
	defer corpse.mux.Unlock()

	location, err := world.CreateObjectAvailableDots(corpse, location)
	if err != nil {
		world.ReleaseIdentifier(corpse.id)
		return nil, errCreateCorpse(err.Error())
	}

	if location.Empty() {
		world.ReleaseIdentifier(corpse.id)
		if err := world.DeleteObject(corpse, location); err != nil {
			return nil, errCreateCorpse("no location located and cannot delete corpse")
		}
		return nil, errCreateCorpse("no location located")
	}

	corpse.location = location

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

		var err error

		c.stopper.Do(func() {
			close(c.stop)
			c.world.ReleaseIdentifier(c.id)
			err = c.world.DeleteObject(c, c.location)
		})

		if err != nil {
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

			var err error

			c.stopper.Do(func() {
				close(c.stop)
				c.world.ReleaseIdentifier(c.id)
				err = c.world.DeleteObject(c, c.location)
			})

			if err != nil {
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
		ID:   c.id,
		Dots: c.location,
		Type: corpseTypeLabel,
	})
}

//go:generate ffjson -force-regenerate $GOFILE

// ffjson: nodecoder
type corpse struct {
	ID   world.Identifier `json:"id"`
	Dots engine.Location  `json:"dots,omitempty"`
	Type string           `json:"type"`
}
