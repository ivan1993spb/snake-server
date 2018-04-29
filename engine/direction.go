package engine

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrInvalidDirection = errors.New("invalid direction")

// Direction indicates movement direction
type Direction uint8

const (
	DirNorth = iota
	DirEast
	DirSouth
	DirWest
	dirCount
)

var directionsJSON = map[Direction][]byte{
	DirNorth: []byte(`"n"`),
	DirEast:  []byte(`"e"`),
	DirSouth: []byte(`"s"`),
	DirWest:  []byte(`"w"`),
}

var unknownDirectionJSON = []byte(`"-"`)

// RandomDirection returns random direction
func RandomDirection() Direction {
	return Direction(rand.Intn(dirCount))
}

// CalculateDirection calculates direction by two passed dots
func CalculateDirection(from, to *Dot) Direction {
	if !from.Equals(to) {
		var diffX, diffY uint8

		if from.x > to.x {
			diffX = from.x - to.x
		} else {
			diffX = to.x - from.x
		}
		if from.y > to.y {
			diffY = from.y - to.y
		} else {
			diffY = to.y - from.y
		}

		if diffX > diffY {
			if to.x > from.x {
				return DirEast
			}
			return DirWest
		}

		if diffY > diffX {
			if to.y > from.y {
				return DirSouth
			}
			return DirNorth
		}
	}

	return RandomDirection()
}

// ValidDirection returns true if passed direction is valid
func ValidDirection(dir Direction) bool {
	return dirCount > dir
}

// Implementing json.Marshaler interface
func (dir Direction) MarshalJSON() ([]byte, error) {
	if dirJSON, ok := directionsJSON[dir]; ok {
		return dirJSON, nil
	}

	// Invalid direction
	return unknownDirectionJSON, fmt.Errorf("cannot marshal direction: %s", ErrInvalidDirection)
}

// Reverse reverses direction
func (dir Direction) Reverse() (Direction, error) {
	switch dir {
	case DirNorth:
		return DirSouth, nil
	case DirEast:
		return DirWest, nil
	case DirSouth:
		return DirNorth, nil
	case DirWest:
		return DirEast, nil
	}

	return 0, fmt.Errorf("cannot reverse direction: %s", ErrInvalidDirection)
}
