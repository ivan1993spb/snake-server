package engine

import (
	"bytes"
	"errors"
	"math/rand"
	"strconv"
)

const (
	minAreaWidth  = 10
	minAreaHeight = 10
)

type Area struct {
	width  uint8
	height uint8
}

type ErrInvalidAreaSize struct {
	Width  uint8
	Height uint8
}

func (e *ErrInvalidAreaSize) Error() string {
	return "invalid area size"
}

func NewArea(width, height uint8) (Area, error) {
	if uint16(width)*uint16(height) == 0 {
		return Area{}, &ErrInvalidAreaSize{
			Width:  width,
			Height: height,
		}
	}

	return Area{
		width:  width,
		height: height,
	}, nil
}

func MustArea(width, height uint8) Area {
	area, err := NewArea(width, height)
	if err != nil {
		panic(err)
	}
	return area
}

func NewUsefulArea(width, height uint8) (Area, error) {
	if width < minAreaWidth || height < minAreaHeight {
		return Area{}, errors.New("cannot add useless area with extra small size")
	}

	return Area{
		width:  width,
		height: height,
	}, nil
}

// Size returns area size
func (a Area) Size() uint16 {
	return uint16(a.width) * uint16(a.height)
}

func (a Area) Width() uint8 {
	return a.width
}

func (a Area) Height() uint8 {
	return a.height
}

func (a Area) Dots() []Dot {
	var x, y uint8
	var dots = make([]Dot, 0, uint16(a.width)*uint16(a.height))

	for x = 0; x < a.width; x++ {
		for y = 0; y < a.height; y++ {
			dots = append(dots, Dot{
				X: x,
				Y: y,
			})
		}
	}

	return dots
}

func (a Area) ContainsDot(dot Dot) bool {
	return a.width > dot.X && a.height > dot.Y
}

func (a Area) ContainsRect(rect Rect) bool {
	return a.width > rect.w+rect.x && a.height > rect.h+rect.y
}

func (a Area) ContainsLocation(location Location) bool {
	for _, dot := range location {
		if a.width <= dot.X || a.height <= dot.Y {
			return false
		}
	}
	return true
}

// NewRandomDot generates random dot on area with starting coordinates X and Y
func (a Area) NewRandomDot(x, y uint8) Dot {
	return Dot{
		X: x + uint8(rand.Intn(int(a.width-x))),
		Y: y + uint8(rand.Intn(int(a.height-y))),
	}
}

func (a Area) NewRandomRect(rw, rh, sx, sy uint8) (*Rect, error) {
	if rw+sx > a.width || rh+sy > a.height {
		return nil, errors.New("cannot get random rect on square: invalid Width or Height")
	}

	var r = &Rect{
		x: sx,
		y: sy,
		w: rw,
		h: rh,
	}

	if a.width-r.w-r.x > 0 {
		r.x += uint8(rand.Intn(int(a.width - r.w - r.x)))
	}

	if a.height-r.h-r.y > 0 {
		r.y += uint8(rand.Intn(int(a.height - r.h - r.y)))
	}

	return r, nil
}

type ErrNavigation struct {
	Err error
}

func (e *ErrNavigation) Error() string {
	return "navigation error: " + e.Err.Error()
}

type ErrAreaNotContainsDot struct {
	Dot Dot
}

func (e *ErrAreaNotContainsDot) Error() string {
	return "area does not contain dot: " + e.Dot.String()
}

// Navigate calculates and returns dot placed on distance dis dots from passed dot in direction dir
func (a Area) Navigate(dot Dot, dir Direction, dis uint8) (Dot, error) {
	// If distance is zero return passed dot
	if dis == 0 {
		return dot, nil
	}

	// Area must contain passed dot
	if !a.ContainsDot(dot) {
		return Dot{}, &ErrNavigation{
			Err: &ErrAreaNotContainsDot{
				Dot: dot,
			},
		}
	}

	switch dir {
	case DirectionNorth, DirectionSouth:
		if dis > a.height {
			dis %= a.height
		}

		// North
		if dir == DirectionNorth {
			if dis > dot.Y {
				return Dot{
					X: dot.X,
					Y: a.height - dis + dot.Y,
				}, nil
			}
			return Dot{
				X: dot.X,
				Y: dot.Y - dis,
			}, nil
		}

		// South
		if a.height > dot.Y+dis {
			return Dot{
				X: dot.X,
				Y: dot.Y + dis,
			}, nil
		}
		return Dot{
			X: dot.X,
			Y: dis - a.height + dot.Y,
		}, nil

	case DirectionWest, DirectionEast:
		if dis > a.width {
			dis %= a.width
		}

		// East
		if dir == DirectionEast {
			if a.width > dot.X+dis {
				return Dot{
					X: dot.X + dis,
					Y: dot.Y,
				}, nil
			}
			return Dot{
				X: dis - a.width + dot.X,
				Y: dot.Y,
			}, nil
		}

		// West
		if dis > dot.X {
			return Dot{
				X: a.width - dis + dot.X,
				Y: dot.Y,
			}, nil
		}
		return Dot{
			X: dot.X - dis,
			Y: dot.Y,
		}, nil
	}

	return Dot{}, &ErrNavigation{
		Err: &ErrInvalidDirection{
			Direction: dir,
		},
	}
}

const areaExpectedSerializedSize = 10

// Implementing json.Marshaler interface
func (a Area) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBuffer(make([]byte, 0, areaExpectedSerializedSize))
	buff.WriteByte('[')
	buff.WriteString(strconv.Itoa(int(a.width)))
	buff.WriteByte(',')
	buff.WriteString(strconv.Itoa(int(a.height)))
	buff.WriteByte(']')
	return buff.Bytes(), nil
}
