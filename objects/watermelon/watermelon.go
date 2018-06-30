package watermelon

import (
	"fmt"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"

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

type Watermelon struct {
	uuid     string
	world    *world.World
	location engine.Location
	mux      *sync.RWMutex
}

type ErrCreateWatermelon string

func (e ErrCreateWatermelon) Error() string {
	return "cannot create watermelon: " + string(e)
}

func NewWatermelon(world *world.World) (*Watermelon, error) {
	watermelon := &Watermelon{
		uuid: uuid.Must(uuid.NewV4()).String(),
		mux:  &sync.RWMutex{},
	}

	location, err := world.CreateObjectRandomRect(watermelon, watermelonWidth, watermelonHeight)
	if err != nil {
		return nil, ErrCreateWatermelon(err.Error())
	}

	watermelon.mux.Lock()
	watermelon.world = world
	watermelon.location = location
	watermelon.mux.Unlock()

	return watermelon, nil
}

func (w *Watermelon) String() string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return fmt.Sprintf("watermelon %s", w.location)
}

func (w *Watermelon) NutritionalValue(dot engine.Dot) uint16 {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.location.Contains(dot) {
		newDots := w.location.Delete(dot)

		if len(newDots) > 0 {
			// TODO: Handle errors?
			newLoc, _ := w.world.UpdateObjectAvailableDots(w, w.location, newDots)
			w.location = newLoc
		} else {
			w.world.DeleteObject(w, w.location)
		}

		return watermelonNutritionalValue
	}

	return 0
}

func (w *Watermelon) MarshalJSON() ([]byte, error) {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return ffjson.Marshal(&watermelon{
		UUID: w.uuid,
		Dots: w.location,
		Type: watermelonTypeLabel,
	})
}

type watermelon struct {
	UUID string       `json:"uuid"`
	Dots []engine.Dot `json:"dots"`
	Type string       `json:"type"`
}
