package playground

import "strconv"

// Dot are used to understand object position
type Dot struct {
	x, y uint8
}

// NewDot creates new dot object
func NewDot(x, y uint8) *Dot {
	return &Dot{x, y}
}

func NewDefaultDot() *Dot {
	return &Dot{0, 0}
}

// NewRandomDotOnSquare generates random dot on square with
// coordinates x and y and width w and height h
func NewRandomDotOnSquare(x, y, w, h uint8) *Dot {
	return &Dot{
		x + uint8(random.Intn(int(w))),
		y + uint8(random.Intn(int(h))),
	}
}

// Equals compares two dots
func (d1 *Dot) Equals(d2 *Dot) bool {
	var (
		d1X, d1Y = d1.Position()
		d2X, d2Y = d2.Position()
	)
	return d1X == d2X && d1Y == d2Y
}

// Position returns coordinates of dot
func (d *Dot) Position() (uint8, uint8) {
	if d == nil {
		return 0, 0
	}
	return d.x, d.y
}

// Pack packs dot to string in accordance with standard ST_1
func (d *Dot) Pack() string {
	var x, y = d.Position()
	return strconv.Itoa(int(x)) + "&" + strconv.Itoa(int(y))
}

// CalculatePath calculates path between two dots
func CalculatePath(from, to *Dot) uint16 {
	if from.Equals(to) {
		return 0
	}

	var (
		fromX, fromY = from.Position()
		toX, toY     = to.Position()
		diffX, diffY uint8
	)

	if fromX > toX {
		diffX = fromX - toX
	} else {
		diffX = toX - fromX
	}
	if fromY > toY {
		diffY = fromY - toY
	} else {
		diffY = toY - fromY
	}

	return uint16(diffX) + uint16(diffY)
}
