// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import (
	"fmt"

	// "bitbucket.org/pushkin_ivan/clever-snake/game/objects"
)

type GameProcessor chan interface{}

func NewGameProcessor() GameProcessor {
	return make(chan interface{})
}

// Implementing game objects.GameProcessor
func (p GameProcessor) OccurredError(object interface{}, err error) {
	//

	//

	// работать нужно только тогда, когда контекст еще жив

	//

	//

	fmt.Printf("game processor: occurred error: %+v\n", object)
}

// Implementing game objects.GameProcessor
func (p GameProcessor) OccurredCreating(object interface{}) {
	fmt.Printf("game processor: occurred creating: %+v\n", object)
}

// Implementing game objects.GameProcessor
func (p GameProcessor) OccurredDeleting(object interface{}) {
	fmt.Printf("game processor: occurred deleting: %+v\n", object)
}

// Implementing game objects.GameProcessor
func (p GameProcessor) OccurredUpdating(object interface{}) {
	fmt.Printf("game processor: occurred updating: %+v\n", object)
}
