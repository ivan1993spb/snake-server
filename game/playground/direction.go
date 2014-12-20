package playground

import "errors"

var ErrInvalidDirection = errors.New("Invalid direction")

// Direction indicates movement direction
type Direction uint8

const (
	DIR_NORTH = iota
	DIR_EAST
	DIR_SOUTH
	DIR_WEST
	_DIR_COUNT
)

// RandomDirection returns random direction
func RandomDirection() Direction {
	return Direction(random.Intn(_DIR_COUNT))
}

// ValidDirection returns true if passed direction is valid
func ValidDirection(dir Direction) bool {
	switch dir {
	case DIR_NORTH, DIR_EAST, DIR_SOUTH, DIR_WEST:
		return true
	}
	return false
}

// Implementing json.Marshaler interface
func (dir Direction) MarshalJSON() ([]byte, error) {
	switch dir {
	case DIR_NORTH:
		return []byte(`"n"`), nil
	case DIR_SOUTH:
		return []byte(`"s"`), nil
	case DIR_EAST:
		return []byte(`"e"`), nil
	case DIR_WEST:
		return []byte(`"w"`), nil
	}
	return nil, ErrInvalidDirection
}
