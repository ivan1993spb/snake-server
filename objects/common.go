package objects

import (
	"math/rand"
	"time"
)

type errCreateObject struct {
	err error
}

func (e *errCreateObject) Error() string {
	return "Cannot create object: " + e.err.Error()
}

type errStartingObject struct {
	err error
}

func (e *errStartingObject) Error() string {
	return "Cannot start object: " + e.err.Error()
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
