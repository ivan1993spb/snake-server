package engine

import (
	"bytes"
	"strconv"
)

// Rect struct represents a rectangle object
type Rect struct {
	x uint8
	y uint8
	w uint8
	h uint8
}

// NewRect returns a new rectangle
func NewRect(x, y, w, h uint8) Rect {
	return Rect{
		x: x,
		y: y,
		w: w,
		h: h,
	}
}

// Width returns rectangle's width
func (r Rect) Width() uint8 {
	return r.w
}

// Height returns rectangle's height
func (r Rect) Height() uint8 {
	return r.h
}

// X returns the X-coordinate of a rectangle
func (r Rect) X() uint8 {
	return r.x
}

// Y returns the Y-coordinate of a rectangle
func (r Rect) Y() uint8 {
	return r.y
}

// ContainsDot returns true if a rectangle contains a given dot
func (r Rect) ContainsDot(d Dot) bool {
	return r.x <= d.X && r.y <= d.Y && r.x+r.w > d.X && r.y+r.h > d.Y
}

// ContainsRect returns true if a rectangle contains another rectangle
func (r Rect) ContainsRect(rect Rect) bool {
	return r.x <= rect.x && r.y <= rect.y && r.x+r.w >= rect.x+rect.w && r.y+r.h >= rect.y+rect.h
}

// Equals returns true if rectangles are equal
func (r Rect) Equals(rect Rect) bool {
	return r == rect
}

// DotCount returns a dot number of a rectangle
func (r Rect) DotCount() uint16 {
	return uint16(r.w) * uint16(r.h)
}

// Dot returns a dot on a rectangle by its index
func (r Rect) Dot(i uint16) Dot {
	return Dot{uint8(i%uint16(r.w)) + r.x, uint8(i/uint16(r.w)) + r.y}
}

const rectExpectedSerializedSize = 20

// Implementing json.Marshaler interface
func (r Rect) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBuffer(make([]byte, 0, rectExpectedSerializedSize))
	buff.WriteByte('[')
	buff.WriteString(strconv.Itoa(int(r.x)))
	buff.WriteByte(',')
	buff.WriteString(strconv.Itoa(int(r.y)))
	buff.WriteByte(',')
	buff.WriteString(strconv.Itoa(int(r.w)))
	buff.WriteByte(',')
	buff.WriteString(strconv.Itoa(int(r.h)))
	buff.WriteByte(']')
	return buff.Bytes(), nil
}

// Dots returns a slice of all dots in a rectangle
func (r Rect) Dots() []Dot {
	dots := make([]Dot, 0, r.DotCount())

	for i := uint16(0); i < r.DotCount(); i++ {
		dots = append(dots, r.Dot(i))
	}

	return dots
}

// Location returns the location of a rectangle
func (r Rect) Location() Location {
	object := make(Location, 0, r.DotCount())

	for i := uint16(0); i < r.DotCount(); i++ {
		object = append(object, r.Dot(i))
	}

	return object
}
