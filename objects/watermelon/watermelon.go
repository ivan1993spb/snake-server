package watermelon

import (
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const (
	watermelonWidth  = 2
	watermelonHeight = 2

	watermelonArea = watermelonWidth * watermelonHeight

	watermelonMinNutrValue = 1
	watermelonMaxNutrValue = 10

	watermelonNutrVar = watermelonMaxNutrValue - watermelonMinNutrValue

	watermelonMaxExperience = time.Second * 30
)

const watermelonTypeLabel = "watermelon"

const watermelonNutritionalValue = 5

// ffjson: skip
type Watermelon struct {
	id       world.Identifier
	world    world.Interface
	location engine.Location
	mux      *sync.RWMutex
}

type ErrCreateWatermelon string

func (e ErrCreateWatermelon) Error() string {
	return "cannot create watermelon: " + string(e)
}

func NewWatermelon(world world.Interface) (*Watermelon, error) {
	watermelon := &Watermelon{
		id:  world.ObtainIdentifier(),
		mux: &sync.RWMutex{},
	}

	watermelon.mux.Lock()
	defer watermelon.mux.Unlock()

	location, err := world.CreateObjectRandomRect(watermelon, watermelonWidth, watermelonHeight)
	if err != nil {
		world.ReleaseIdentifier(watermelon.id)
		return nil, ErrCreateWatermelon(err.Error())
	}

	watermelon.world = world
	watermelon.location = location

	return watermelon, nil
}

func (w *Watermelon) String() string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return fmt.Sprintf("watermelon %s", w.location)
}

type errWatermelonBite string

func (e errWatermelonBite) Error() string {
	return "watermelon bite error: " + string(e)
}

func (w *Watermelon) Bite(dot engine.Dot) (nv uint16, success bool, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.location.Contains(dot) {
		newDots := w.location.Delete(dot)

		if len(newDots) > 0 {
			newLocation, err := w.world.UpdateObjectAvailableDots(w, w.location, newDots)
			if err != nil {
				return 0, false, errWatermelonBite(err.Error())
			}
			if len(newLocation) > 0 {
				w.location = newLocation
				return watermelonNutritionalValue, true, nil
			}
		}

		w.world.ReleaseIdentifier(w.id)

		if err := w.world.DeleteObject(w, w.location); err != nil {
			return 0, false, errWatermelonBite(err.Error())
		}

		w.location = w.location[:0]

		return watermelonNutritionalValue, true, nil
	}

	return 0, false, errWatermelonBite("watermelon does not contain dot")
}

func (w *Watermelon) MarshalJSON() ([]byte, error) {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return ffjson.Marshal(&watermelon{
		ID:   w.id,
		Dots: w.location,
		Type: watermelonTypeLabel,
	})
}

//go:generate ffjson -force-regenerate $GOFILE

// ffjson: nodecoder
type watermelon struct {
	ID   world.Identifier `json:"id"`
	Dots []engine.Dot     `json:"dots,omitempty"`
	Type string           `json:"type"`
}
