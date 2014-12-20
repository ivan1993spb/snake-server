package playground

import (
	"encoding/json"
	"fmt"
)

type Rect struct {
	x, y, w, h uint8
}

// NewRect creates rect
func NewRect(x, y, w, h uint8) *Rect {
	return &Rect{x, y, w, h}
}

func NewRandomRectOnSquare(rw, rh, sx, sy, sw, sh uint8,
) (*Rect, error) {
	if rw > sw || rh > sh {
		return nil, fmt.Errorf("Cannot get random rect on square: %s",
			ErrInvalid_W_or_H)
	}

	var r = &Rect{sx, sy, rw, rh}

	if sw-r.w > 0 {
		r.x = uint8(random.Intn(int(sw - r.w)))
	}

	if sh-r.h > 0 {
		r.y = uint8(random.Intn(int(sh - r.h)))
	}

	return r, nil
}

func (r *Rect) ContainsDot(d *Dot) bool {
	return r.x <= d.x && r.y <= d.y && r.x+r.w > d.x && r.y+r.h > d.y
}

func (r1 *Rect) ContainsRect(r2 *Rect) bool {
	return r1.x <= r2.x && r1.y <= r2.y &&
		r1.w >= r2.w && r1.h >= r2.h
}

func (r1 *Rect) Equals(r2 *Rect) bool {
	return r1 == r2 ||
		(r1.x == r2.x && r1.y == r2.y && r1.w == r2.w && r1.h == r2.h)
}

// Implementing Entity interface
func (r *Rect) DotCount() uint16 {
	return uint16(r.w) * uint16(r.h)
}

// Implementing Entity interface
func (r *Rect) Dot(i uint16) *Dot {
	return NewDot(uint8(i/uint16(r.w))+r.x, uint8(i%uint16(r.h))+r.y)
}

// RandomDotOnRect returns random dot on rect
func (r *Rect) RandomDotOnRect() *Dot {
	return NewRandomDotOnSquare(0, 0, r.w, r.h)
}

// Implementing json.Marshaler interface
func (r *Rect) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint16{
		uint16(r.x),
		uint16(r.y),
		uint16(r.w),
		uint16(r.h),
	})
}
