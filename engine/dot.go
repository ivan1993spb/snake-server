package engine

import "fmt"

type Dot struct {
	X uint8
	Y uint8
}

// Equals compares two dots
func (d1 Dot) Equals(d2 Dot) bool {
	return d1 == d2 || (d1.X == d2.X && d1.Y == d2.Y)
}

// Implementing json.Marshaler interface
func (d Dot) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%d]", d.X, d.Y)), nil
}

func (d Dot) String() string {
	return fmt.Sprintf("[%d, %d]", d.X, d.Y)
}

// DistanceTo calculates distance between two dots
func (from Dot) DistanceTo(to Dot) (res uint16) {
	if !from.Equals(to) {
		if from.X > to.X {
			res = uint16(from.X - to.X)
		} else {
			res = uint16(to.X - from.X)
		}

		if from.Y > to.Y {
			res += uint16(from.Y - to.Y)
		} else {
			res += uint16(to.Y - from.Y)
		}
	}

	return
}
