package engine

import (
	"encoding/json"
	"errors"
	"math/rand"
)

// TODO: Try use Area instead *Area.

const (
	minAreaWidth  = 10
	minAreaHeight = 10
)

// TODO: Create constructor checking for minimal area width and height.

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

func NewArea(width, height uint8) (*Area, error) {
	if width*height == 0 {
		return nil, &ErrInvalidAreaSize{
			Width:  width,
			Height: height,
		}
	}

	return &Area{
		width:  width,
		height: height,
	}, nil
}

// Size returns area size
func (a *Area) Size() uint16 {
	return uint16(a.width) * uint16(a.height)
}

func (a *Area) Width() uint8 {
	return a.width
}

func (a *Area) Height() uint8 {
	return a.height
}

func (a *Area) Contains(dot *Dot) bool {
	return a.width > dot.x && a.height > dot.y
}

// NewRandomDot generates random dot on area with starting coordinates x and y
func (a *Area) NewRandomDot(x, y uint8) *Dot {
	return &Dot{
		x: x + uint8(rand.Intn(int(a.width))),
		y: y + uint8(rand.Intn(int(a.height))),
	}
}

func (a *Area) NewRandomRect(rw, rh, sx, sy uint8) (*Rect, error) {
	if rw > a.width || rh > a.height {
		return nil, errors.New("cannot get random rect on square: invalid Width or Height")
	}

	var r = &Rect{
		x: sx,
		y: sy,
		w: rw,
		h: rh,
	}

	if a.width-r.w > 0 {
		r.x = uint8(rand.Intn(int(a.width - r.w)))
	}

	if a.height-r.h > 0 {
		r.y = uint8(rand.Intn(int(a.height - r.h)))
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
	Dot *Dot
}

func (e *ErrAreaNotContainsDot) Error() string {
	return "area does not contain dot"
}

// Navigate calculates and returns dot placed on distance dis dots from passed dot in direction dir
func (a *Area) Navigate(dot *Dot, dir Direction, dis uint8) (*Dot, error) {
	// If distance is zero return passed dot
	if dis == 0 {
		return dot, nil
	}

	// Area must contain passed dot
	if !a.Contains(dot) {
		return nil, &ErrNavigation{
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
			if dis > dot.y {
				return &Dot{
					x: dot.x,
					y: a.height - dis + dot.y,
				}, nil
			}
			return &Dot{
				x: dot.x,
				y: dot.y - dis,
			}, nil
		}

		// South
		if dot.y+dis+1 > a.height {
			return &Dot{
				x: dot.x,
				y: dis - a.height + dot.y,
			}, nil
		}
		return &Dot{
			x: dot.x,
			y: dot.y + dis,
		}, nil

	case DirectionWest, DirectionEast:
		if dis > a.width {
			dis %= a.width
		}

		// East
		if dir == DirectionEast {
			if a.width > dot.x+dis {
				return &Dot{
					x: dot.x + dis,
					y: dot.y,
				}, nil
			}
			return &Dot{
				x: dis - a.width + dot.x,
				y: dot.y,
			}, nil
		}

		// West
		if dis > dot.x {
			return &Dot{
				x: a.width - dis + dot.x,
				y: dot.y,
			}, nil
		}
		return &Dot{
			x: dot.x - dis,
			y: dot.y,
		}, nil
	}

	return nil, &ErrNavigation{
		Err: &ErrInvalidDirection{
			Direction: dir,
		},
	}
}

// Implementing json.Marshaler interface
func (a *Area) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint8{
		a.width,
		a.height,
	})
}
