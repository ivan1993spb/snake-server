// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watermelon

import (
	"math/rand"
	"time"

	"github.com/olebedev/emitter"

	"github.com/ivan1993spb/clever-snake/engine"
	"github.com/ivan1993spb/clever-snake/playground"
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
	pg       *playground.Playground
	location engine.Location
}

func CreateWatermelon(pg *playground.Playground) (*Watermelon, error) {

	watermelon := &Watermelon{}

	location, err := pg.CreateObjectRandomRect(watermelon, watermelonWidth, watermelonHeight)
	if err != nil {
		// TODO: Create specific error.
		return nil, err
	}

	watermelon.pg = pg
	watermelon.location = location

	return watermelon, nil
}

// Implementing playground.Location interface
func (w *Watermelon) DotCount() (c uint16) {
	for _, dot := range w.location {
		if dot != nil {
			c++
		}
	}

	return
}

// Implementing playground.Location interface
func (w *Watermelon) Dot(i uint16) (dot *engine.Dot) {
	if i < w.DotCount() {
		var j uint16
		for _, dot = range w.location {
			if dot != nil {
				if i == j {
					break
				}
				j++
			}
		}
	}
	return
}

// Implementing logic.Food interface
func (w *Watermelon) NutritionalValue(dot *engine.Dot) int8 {
	for _, dot := range w.location {
		if dot.Equals(dot) {
			location := w.location.Delete(dot)
			w.pg.UpdateObject(w, w.location, location)
			w.location = location
			return watermelonMinNutrValue + int8(rand.Intn(watermelonNutrVar))
		}
	}

	return 0
}

func (w *Watermelon) Run(emitter *emitter.Emitter) error {
	// TODO: Implement method.
}
