package engine

import "math/rand"

type ErrInvalidDirection struct {
	Direction Direction
}

func (e *ErrInvalidDirection) Error() string {
	return "invalid direction"
}

// Direction indicates movement direction
type Direction uint8

const (
	DirectionNorth Direction = iota
	DirectionEast
	DirectionSouth
	DirectionWest
	directionCount
)

var directionsJSON = map[Direction][]byte{
	DirectionNorth: []byte(`"north"`),
	DirectionEast:  []byte(`"east"`),
	DirectionSouth: []byte(`"south"`),
	DirectionWest:  []byte(`"west"`),
}

var directionsLabels = map[Direction]string{
	DirectionNorth: "north",
	DirectionEast:  "east",
	DirectionSouth: "south",
	DirectionWest:  "west",
}

func (dir Direction) String() string {
	if label, ok := directionsLabels[dir]; ok {
		return label
	}
	return "unknown"
}

var unknownDirectionJSON = []byte(`"-"`)

// RandomDirection returns random direction
func RandomDirection() Direction {
	return Direction(rand.Intn(int(directionCount)))
}

// CalculateDirection calculates direction by two passed dots
func CalculateDirection(from, to Dot) Direction {
	if !from.Equals(to) {
		var diffX, diffY uint8

		if from.X > to.X {
			diffX = from.X - to.X
		} else {
			diffX = to.X - from.X
		}
		if from.Y > to.Y {
			diffY = from.Y - to.Y
		} else {
			diffY = to.Y - from.Y
		}

		if diffX > diffY {
			if to.X > from.X {
				return DirectionEast
			}
			return DirectionWest
		}

		if diffY > diffX {
			if to.Y > from.Y {
				return DirectionSouth
			}
			return DirectionNorth
		}
	}

	return RandomDirection()
}

// ValidDirection returns true if passed direction is valid
func ValidDirection(dir Direction) bool {
	return directionCount > dir
}

type ErrDirectionMarshal struct {
	Err error
}

func (e *ErrDirectionMarshal) Error() string {
	return "cannot marshal direction"
}

// Implementing json.Marshaler interface
func (dir Direction) MarshalJSON() ([]byte, error) {
	if dirJSON, ok := directionsJSON[dir]; ok {
		return dirJSON, nil
	}

	// Invalid direction
	return unknownDirectionJSON, &ErrDirectionMarshal{
		Err: &ErrInvalidDirection{
			Direction: dir,
		},
	}
}

type ErrReverseDirection struct {
	Err error
}

func (e ErrReverseDirection) Error() string {
	return "cannot reverse direction"
}

// Reverse reverses direction
func (dir Direction) Reverse() (Direction, error) {
	switch dir {
	case DirectionNorth:
		return DirectionSouth, nil
	case DirectionEast:
		return DirectionWest, nil
	case DirectionSouth:
		return DirectionNorth, nil
	case DirectionWest:
		return DirectionEast, nil
	}

	return 0, &ErrReverseDirection{
		Err: &ErrInvalidDirection{
			Direction: dir,
		},
	}
}
