package engine

import (
	"math"
	"math/rand"
)

type DotsMask struct {
	mask [][]uint8
}

var DotsMaskSquare2x2 = NewDotsMask([][]uint8{
	{1, 1},
	{1, 1},
})

var DotsMaskTank = NewDotsMask([][]uint8{
	{0, 1, 0},
	{1, 1, 1},
	{1, 0, 1},
})

var DotsMaskHome1 = NewDotsMask([][]uint8{
	{1, 1, 0, 1, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 1, 1, 1, 1},
})

var DotsMaskHome2 = NewDotsMask([][]uint8{
	{1, 1, 1, 0, 0, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{1, 0, 1, 0, 0, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 0, 0, 1, 1, 1},
})

var DotsMaskCross = NewDotsMask([][]uint8{
	{0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0},
	{1, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0},
})

var DotsMaskDiagonal = NewDotsMask([][]uint8{
	{0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 1, 0},
	{0, 0, 0, 0, 1, 0, 0},
	{0, 0, 0, 1, 0, 0, 0},
	{0, 0, 1, 0, 0, 0, 0},
	{0, 1, 0, 0, 0, 0, 0},
	{1, 0, 0, 0, 0, 0, 0},
})

var DotsMaskCrossSmall = NewDotsMask([][]uint8{
	{0, 1, 0},
	{1, 1, 1},
	{0, 1, 0},
})

var DotsMaskDiagonalSmall = NewDotsMask([][]uint8{
	{1, 0, 1},
	{0, 1, 0},
	{1, 0, 1},
})

var DotsMaskLabyrinth = NewDotsMask([][]uint8{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 1, 0, 0, 0, 1, 1, 0, 1},
	{1, 0, 1, 1, 1, 0, 1, 0, 0, 1},
	{0, 0, 1, 0, 1, 0, 1, 1, 1, 1},
	{0, 0, 1, 0, 0, 0, 0, 0, 1, 0},
	{1, 0, 1, 0, 1, 0, 1, 0, 0, 0},
	{1, 0, 1, 0, 1, 0, 1, 1, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
})

var DotsMaskTunnel1 = NewDotsMask([][]uint8{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
})

var DotsMaskTunnel2 = NewDotsMask([][]uint8{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
})

var DotsMaskBigHome = NewDotsMask([][]uint8{
	{1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1},
})

func NewDotsMask(mask [][]uint8) *DotsMask {
	if len(mask) > math.MaxUint8 {
		mask = mask[:math.MaxUint8+1]
	}

	copyMask := make([][]uint8, len(mask))

	for i, row := range mask {
		if len(row) > math.MaxUint8 {
			copyMask[i] = make([]uint8, math.MaxUint8+1)
		} else {
			copyMask[i] = make([]uint8, len(row))
		}
		copy(copyMask[i], row)
	}

	return &DotsMask{
		mask: copyMask,
	}
}

func LocationToDotsMask(location Location) *DotsMask {
	if location.Empty() {
		return &DotsMask{
			mask: [][]uint8{},
		}
	}

	if location.DotCount() == 1 {
		return &DotsMask{
			mask: [][]uint8{
				{1},
			},
		}
	}

	firstDot := location.Dot(0)
	leftX := firstDot.X
	rightX := firstDot.X
	topY := firstDot.Y
	bottomY := firstDot.Y

	for i := uint16(0); i < location.DotCount(); i++ {
		dot := location.Dot(i)
		if leftX > dot.X {
			leftX = dot.X
		}
		if rightX < dot.X {
			rightX = dot.X
		}
		if topY > dot.Y {
			topY = dot.Y
		}
		if bottomY < dot.Y {
			bottomY = dot.Y
		}
	}

	if leftX == rightX && topY == bottomY {
		return &DotsMask{
			mask: [][]uint8{
				{1},
			},
		}
	}

	dm := NewZeroDotsMask(rightX-leftX+1, bottomY-topY+1)

	for i := uint16(0); i < location.DotCount(); i++ {
		dot := location.Dot(i)
		dm.mask[dot.Y-topY][dot.X-leftX] = 1
	}

	return dm
}

func NewZeroDotsMask(width, height uint8) *DotsMask {
	mask := make([][]uint8, height)
	for i := range mask {
		mask[i] = make([]uint8, width)
	}
	return &DotsMask{
		mask: mask,
	}
}

func (dm *DotsMask) Copy() *DotsMask {
	copyMask := make([][]uint8, len(dm.mask))

	for i, row := range dm.mask {
		if len(row) > math.MaxUint8 {
			copyMask[i] = make([]uint8, math.MaxUint8+1)
		} else {
			copyMask[i] = make([]uint8, len(row))
		}
		copy(copyMask[i], row)
	}

	return &DotsMask{
		mask: copyMask,
	}
}

func (dm *DotsMask) Width() uint8 {
	width := 0
	for _, row := range dm.mask {
		if width < len(row) {
			width = len(row)
		}
	}
	return uint8(width)
}

func (dm *DotsMask) Height() uint8 {
	return uint8(len(dm.mask))
}

func (dm *DotsMask) TurnOver() *DotsMask {
	newMask := NewZeroDotsMask(dm.Width(), dm.Height())
	for i := range dm.mask {
		copy(newMask.mask[len(dm.mask)-1-i], dm.mask[i])
	}
	return newMask
}

func (dm *DotsMask) TurnRight() *DotsMask {
	newMask := NewZeroDotsMask(dm.Height(), dm.Width())
	for i := 0; i < len(dm.mask); i++ {
		for j := 0; j < len(dm.mask[i]); j++ {
			newMask.mask[j][len(dm.mask)-1-i] = dm.mask[i][j]
		}
	}
	return newMask
}

func (dm *DotsMask) TurnLeft() *DotsMask {
	newMask := NewZeroDotsMask(dm.Height(), dm.Width())
	for i := 0; i < len(dm.mask); i++ {
		for j := 0; j < len(dm.mask[i]); j++ {
			newMask.mask[len(newMask.mask)-1-j][i] = dm.mask[i][j]
		}
	}
	return newMask
}

func (dm *DotsMask) TurnRandom() *DotsMask {
	const (
		caseReturnCopy = iota
		caseReturnTurnRight
		caseReturnTurnLeft
		caseReturnTurnOver
		turnReturnCasesCount
	)

	switch rand.Intn(turnReturnCasesCount) {
	case caseReturnCopy:
		return dm.Copy()
	case caseReturnTurnRight:
		return dm.TurnRight()
	case caseReturnTurnLeft:
		return dm.TurnLeft()
	case caseReturnTurnOver:
		return dm.TurnOver()
	}

	return nil
}

func (dm *DotsMask) Location(x, y uint8) Location {
	location := make(Location, 0)
	for i := 0; i < len(dm.mask); i++ {
		for j := 0; j < len(dm.mask[i]); j++ {
			if dm.mask[i][j] > 0 {
				location = append(location, Dot{
					X: x + uint8(j),
					Y: y + uint8(i),
				})
			}
		}
	}
	return location
}

func (dm *DotsMask) Empty() bool {
	for i := range dm.mask {
		for j := range dm.mask[i] {
			if dm.mask[i][j] > 0 {
				return false
			}
		}
	}
	return true
}

func (dm *DotsMask) DotCount() (count uint16) {
	for i := range dm.mask {
		for j := range dm.mask[i] {
			if dm.mask[i][j] > 0 {
				count++
			}
		}
	}
	return
}
