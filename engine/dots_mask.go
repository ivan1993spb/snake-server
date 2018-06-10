package engine

import "math"

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

var DotsMaskHome = NewDotsMask([][]uint8{
	{1, 1, 0, 1, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 0, 0, 0, 1},
	{1, 1, 1, 1, 1},
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

func NewDotsMask(mask [][]uint8) *DotsMask {
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
