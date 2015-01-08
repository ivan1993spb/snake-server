// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playground

import (
	"math/rand"
	"time"
)

type Entity interface {
	Dot(i uint16) *Dot // Dot returns dot by index
	DotCount() uint16  // DotCount must return dot count
}

func EntityToDotList(e Entity) DotList {
	dots := make(DotList, 0, e.DotCount())

	for i := uint16(0); i < e.DotCount(); i++ {
		dots = append(dots, e.Dot(0))
	}

	return dots
}

type Object interface {
	Entity
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
