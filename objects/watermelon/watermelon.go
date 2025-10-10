package watermelon

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

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
		id:  world.IdentifierRegistry().Obtain(),
		mux: &sync.RWMutex{},
	}

	watermelon.mux.Lock()
	defer watermelon.mux.Unlock()

	location, err := world.CreateObjectRandomRect(watermelon, watermelonWidth, watermelonHeight)
	if err != nil {
		world.IdentifierRegistry().Release(watermelon.id)
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
			newLocation, err := w.world.UpdateObjectAvailableDotsV2(world.UpdateV2Params{
				Object: w,
				Payload: updateV2{
					ID:  w.id,
					Del: dot,
				},
				Old: w.location,
				New: newDots,
			})

			if err != nil {
				return 0, false, errWatermelonBite(err.Error())
			}

			if len(newLocation) > 0 {
				w.location = newLocation
				return watermelonNutritionalValue, true, nil
			}
		}

		w.world.IdentifierRegistry().Release(w.id)

		err := w.world.DeleteObjectV2(world.DeleteV2Params{
			Object:   w,
			Payload:  DeleteV2Payload{ID: w.id},
			Location: w.location,
		})

		if err != nil {
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
	return json.Marshal(&watermelon{
		ID:   w.id,
		Dots: w.location,
		Type: watermelonTypeLabel,
	})
}

type watermelon struct {
	ID   world.Identifier `json:"id"`
	Dots []engine.Dot     `json:"dots,omitempty"`
	Type string           `json:"type"`
}

type updateV2 struct {
	ID  world.Identifier `json:"id"`
	Del engine.Dot       `json:"del"`
}

type DeleteV2Payload struct {
	ID world.Identifier
}

func (p DeleteV2Payload) MarshalJSON() ([]byte, error) {
	return p.ID.MarshalJSON()
}
