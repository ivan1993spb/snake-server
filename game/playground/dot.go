package playground

import "encoding/json"

type Dot struct {
	x, y uint8
}

// NewDot creates dot object
func NewDot(x, y uint8) *Dot {
	return &Dot{x, y}
}

// NewRandomDotOnSquare generates random dot on square with
// coordinates x and y, width w and height h
func NewRandomDotOnSquare(x, y, w, h uint8) *Dot {
	return &Dot{
		x + uint8(random.Intn(int(w))),
		y + uint8(random.Intn(int(h))),
	}
}

// Equals compares two dots
func (d1 *Dot) Equals(d2 *Dot) bool {
	return d1 == d2 || (d1.x == d2.x && d1.y == d2.y)
}

// PackJson packs dot
func (d *Dot) PackJson() (json.RawMessage, error) {
	return json.Marshal([]uint16{uint16(d.x), uint16(d.y)})
}
