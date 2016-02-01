// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objects

import (
	"errors"
	"math/rand"
	"time"
)

type errCreateObject struct {
	err error
}

func (e *errCreateObject) Error() string {
	return "cannot create object: " + e.err.Error()
}

type errStartingObject struct {
	err error
}

func (e *errStartingObject) Error() string {
	return "cannot start object: " + e.err.Error()
}

var (
	errNilPlayground    = errors.New("playground is nil")
	errNilGameProcessor = errors.New("game processor is nil")
	errEmptyDotList     = errors.New("empty dot list")
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
