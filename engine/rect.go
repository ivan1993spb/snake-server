package engine

import "encoding/json"

type Rect struct {
	x uint8
	y uint8
	w uint8
	h uint8
}

// NewRect creates rect
func NewRect(x, y, w, h uint8) Rect {
	return Rect{
		x: x,
		y: y,
		w: w,
		h: h,
	}
}

func (r Rect) Width() uint8 {
	return r.w
}

func (r Rect) Height() uint8 {
	return r.h
}

func (r Rect) X() uint8 {
	return r.x
}

func (r Rect) Y() uint8 {
	return r.y
}

func (r Rect) ContainsDot(d Dot) bool {
	return r.x <= d.X && r.y <= d.Y && r.x+r.w > d.X && r.y+r.h > d.Y
}

func (r1 Rect) ContainsRect(r2 Rect) bool {
	return r1.x <= r2.x && r1.y <= r2.y && r1.w >= r2.w && r1.h >= r2.h
}

func (r1 Rect) Equals(r2 Rect) bool {
	return r1 == r2 || (r1.x == r2.x && r1.y == r2.y && r1.w == r2.w && r1.h == r2.h)
}

func (r Rect) DotCount() uint16 {
	return uint16(r.w) * uint16(r.h)
}

func (r Rect) Dot(i uint16) Dot {
	return Dot{uint8(i%uint16(r.w)) + r.x, uint8(i/uint16(r.w)) + r.y}
}

// Implementing json.Marshaler interface
func (r Rect) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint16{
		uint16(r.x),
		uint16(r.y),
		uint16(r.w),
		uint16(r.h),
	})
}

func (r Rect) Dots() []Dot {
	dots := make([]Dot, 0, r.DotCount())

	for i := uint16(0); i < r.DotCount(); i++ {
		dots = append(dots, r.Dot(i))
	}

	return dots
}

func (r Rect) Location() Location {
	object := make(Location, 0, r.DotCount())

	for i := uint16(0); i < r.DotCount(); i++ {
		object = append(object, r.Dot(i))
	}

	return object
}
