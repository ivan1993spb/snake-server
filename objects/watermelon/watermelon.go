package watermelon

import (
	"math/rand"
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

type Watermelon struct {
	world    *world.World
	location engine.Location
}

func CreateWatermelon(world *world.World) (*Watermelon, error) {

	watermelon := &Watermelon{}

	location, err := world.CreateObjectRandomRect(watermelon, watermelonWidth, watermelonHeight)
	if err != nil {
		// TODO: Create specific error.
		return nil, err
	}

	watermelon.world = world
	watermelon.location = location

	return watermelon, nil
}

func (w *Watermelon) NutritionalValue(dot engine.Dot) int8 {
	for _, dot := range w.location {
		if dot.Equals(dot) {
			location := w.location.Delete(dot)
			w.world.UpdateObject(w, w.location, location)
			w.location = location
			return watermelonMinNutrValue + int8(rand.Intn(watermelonNutrVar))
		}
	}

	return 0
}
