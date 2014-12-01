package objects

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

const _APPLE_LOC_RETRIES_NUMBER = 10

type Apple struct {
	pg  *playground.Playground
	dot *playground.Dot
}

// Create and locate new apple
func NewApple(pg *playground.Playground) (*Apple, error) {
	if pg != nil {
		// Try to locate apple for X times
		for i := 0; i < _APPLE_LOC_RETRIES_NUMBER; i++ {
			apple := &Apple{pg, pg.RandomDot()}
			if err := pg.Locate(apple); err == nil {
				return apple, nil
			} else {
				return nil, err
			}
		}
	}
	return nil, errors.New("Cannot create apple")
}

// Implementing playground.Object interface
func (*Apple) DotCount() uint16 {
	return 1
}

// Implementing playground.Object interface
func (a *Apple) Dot(uint16) *playground.Dot {
	return a.dot
}

// Implementing playground.Object interface
func (a *Apple) Pack() string {
	return a.dot.Pack()
}

// Implementing logic.Food interface
func (a *Apple) NutritionalValue(dot *playground.Dot) int8 {
	if a.dot.Equals(dot) {
		a.pg.Delete(a)
		return 1
	}
	return 0
}
