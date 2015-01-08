// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logic

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
)

// Object game characteristics
type (
	Object interface{}

	Resistant interface {
		Object
		Strength() float32
	}

	Living interface {
		Object
		Resistant
		Die()      // Every living thing ever dies
		Feed(int8) // Living things are hungry
	}

	Notalive interface {
		Object
		Break(*playground.Dot)
	}

	Food interface {
		Object
		NutritionalValue(*playground.Dot) int8 // Nutritional value
	}
)

var ErrRecognizingObject = errors.New("cannot recognize object")

const _STRENGTH_FACTOR float32 = 1.3

// Clash implements logic of clash of two objects (first and second)
// in dot dot
func Clash(first Living, second Object, dot *playground.Dot) error {
	switch second := second.(type) {
	case Food:
		// Feed if second object is food
		first.Feed(second.NutritionalValue(dot))

	case Resistant:
		if first.Strength() > second.Strength()*_STRENGTH_FACTOR {
			switch second := second.(type) {
			case Living:
				// Living dies
				second.Die()

			case Notalive:
				// Not living breaks
				second.Break(dot)

			default:
				return ErrRecognizingObject

			}
		} else {
			// Clash with any hard object will result to dying
			first.Die()
		}

	default:
		return ErrRecognizingObject
	}

	return nil
}
