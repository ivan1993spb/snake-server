package engine

import (
	"encoding/json"
	"fmt"
)

type Dot struct {
	x uint8
	y uint8
}

// NewDot creates dot object
func NewDot(x, y uint8) *Dot {
	return &Dot{
		x: x,
		y: y,
	}
}

// Equals compares two dots
func (d1 *Dot) Equals(d2 *Dot) bool {
	return d1 == d2 || (d1.x == d2.x && d1.y == d2.y)
}

// Implementing json.Marshaler interface
func (d *Dot) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint16{uint16(d.x), uint16(d.y)})
}

func (d *Dot) String() string {
	return fmt.Sprintf("[%d, %d]", d.x, d.y)
}

// DistanceTo calculates distance between two dots
func (from *Dot) DistanceTo(to *Dot) (res uint16) {
	if !from.Equals(to) {
		if from.x > to.x {
			res = uint16(from.x - to.x)
		} else {
			res = uint16(to.x - from.x)
		}

		if from.y > to.y {
			res += uint16(from.y - to.y)
		} else {
			res += uint16(to.y - from.y)
		}
	}

	return
}

func (d *Dot) X() uint8 {
	return d.x
}

func (d *Dot) Y() uint8 {
	return d.y
}
