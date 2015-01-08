package objects

import (
	"errors"
	"math/rand"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
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

type GameProcessor interface {
	OccurredError(err error)
	OccurredCreating(object playground.Object)
	OccurredDeleting(object playground.Object)
	OccurredUpdating(object playground.Object)
}

var (
	errNilPlayground    = errors.New("playground is nil")
	errNilGameProcessor = errors.New("game processor is nil")
	errEmptyDotList     = errors.New("empty dot list")
)
