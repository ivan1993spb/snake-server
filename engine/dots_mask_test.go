package engine

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewDotsMask(t *testing.T) {
	dm1 := NewDotsMask([][]uint8{})
	require.Equal(t, [][]uint8{}, dm1.mask)

	dm2 := NewDotsMask([][]uint8{
		{1},
		{1},
		{1},
		{1},
	})
	require.Equal(t, [][]uint8{
		{1},
		{1},
		{1},
		{1},
	}, dm2.mask)

	dm3 := NewDotsMask([][]uint8{
		{1},
		{1},
		{1, 1, 1, 1, 1},
		{1},
	})
	require.Equal(t, [][]uint8{
		{1},
		{1},
		{1, 1, 1, 1, 1},
		{1},
	}, dm3.mask)

	dm4 := NewDotsMask([][]uint8{
		{1, 0, 1},
		{1},
		{1, 1, 1, 1, 1},
		{1},
	})
	require.Equal(t, [][]uint8{
		{1, 0, 1},
		{1},
		{1, 1, 1, 1, 1},
		{1},
	}, dm4.mask)
}

func Test_NewDotsMask_TheBiggestMask(t *testing.T) {
	rawMap := make([][]uint8, math.MaxUint8+1)
	for i := range rawMap {
		rawMap[i] = make([]uint8, math.MaxUint8+1)
		for j := range rawMap[i] {
			rawMap[i][j] = 1
		}
	}

	dm := NewDotsMask(rawMap)
	require.Equal(t, rawMap, dm.mask)
}

func Test_NewZeroDotsMask(t *testing.T) {
	dm1 := NewZeroDotsMask(0, 0)
	require.Equal(t, [][]uint8{}, dm1.mask)

	dm2 := NewZeroDotsMask(5, 3)
	require.Equal(t, [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}, dm2.mask)

	dm3 := NewZeroDotsMask(3, 5)
	require.Equal(t, [][]uint8{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}, dm3.mask)
}

func Test_DotsMask_Copy(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	dm1Copy := dm1.Copy()
	require.Equal(t, [][]uint8{
		{1, 0, 1},
		{1},
		{1, 1, 1, 1, 1},
		{1},
	}, dm1Copy.mask)

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 1, 1},
		},
	}
	dm2Copy := dm2.Copy()
	require.Equal(t, [][]uint8{
		{1, 0, 0, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 0, 0, 1, 1},
	}, dm2Copy.mask)
}

func Test_DotsMask_Width(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	require.Equal(t, uint8(5), dm1.Width())

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 1, 1},
		},
	}
	require.Equal(t, uint8(5), dm2.Width())

	dm3 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
	}
	require.Equal(t, uint8(3), dm3.Width())

	dm4 := &DotsMask{
		mask: [][]uint8{
			{1},
			{1},
			{1},
			{1},
		},
	}
	require.Equal(t, uint8(1), dm4.Width())
}

func Test_DotsMask_Height(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	require.Equal(t, uint8(4), dm1.Height())

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 1, 1},
		},
	}
	require.Equal(t, uint8(3), dm2.Height())

	dm3 := &DotsMask{
		mask: [][]uint8{
			{1},
			{1},
			{1},
			{1},
		},
	}
	require.Equal(t, uint8(4), dm3.Height())
}

func Test_DotsMask_TurnOver(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{1, 0, 0, 0, 0},
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0},
			{1, 0, 1, 0, 0},
		},
	}, dm1.TurnOver())

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1, 1, 1},
			{0, 1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{0, 1, 0},
			{1, 1, 1},
			{1, 0, 1},
		},
	}, dm2.TurnOver())
}

func Test_DotsMask_TurnRight(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{1, 1, 1, 1},
			{0, 1, 0, 0},
			{0, 1, 0, 1},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		},
	}, dm1.TurnRight())

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1, 1, 1},
			{0, 1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{0, 1, 1},
			{1, 1, 0},
			{0, 1, 1},
		},
	}, dm2.TurnRight())
}

func Test_DotsMask_TurnLeft(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1},
			{1, 1, 1, 1, 1},
			{1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{1, 0, 1, 0},
			{0, 0, 1, 0},
			{1, 1, 1, 1},
		},
	}, dm1.TurnLeft())

	dm2 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 1},
			{1, 1, 1},
			{0, 1},
		},
	}
	require.Equal(t, &DotsMask{
		mask: [][]uint8{
			{1, 1, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
	}, dm2.TurnLeft())
}

func Test_DotsMask_Location(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 1, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
	}
	require.Equal(t, Location{
		Dot{
			X: 4,
			Y: 3,
		},
		Dot{
			X: 5,
			Y: 3,
		},
		Dot{
			X: 5,
			Y: 4,
		},
		Dot{
			X: 6,
			Y: 4,
		},
		Dot{
			X: 4,
			Y: 5,
		},
		Dot{
			X: 5,
			Y: 5,
		},
	}, dm1.Location(4, 3))
}

func Test_DotsMask_Empty(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 1, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
	}
	require.False(t, dm1.Empty())

	dm2 := &DotsMask{
		mask: [][]uint8{},
	}
	require.True(t, dm2.Empty())

	dm3 := &DotsMask{
		mask: [][]uint8{
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
		},
	}
	require.True(t, dm3.Empty())

	dm4 := &DotsMask{
		mask: [][]uint8{
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0},
		},
	}
	require.False(t, dm4.Empty())
}

func Test_LocationToDotsMask(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 1, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
	}
	location1 := Location{
		Dot{
			X: 4,
			Y: 3,
		},
		Dot{
			X: 5,
			Y: 3,
		},
		Dot{
			X: 5,
			Y: 4,
		},
		Dot{
			X: 6,
			Y: 4,
		},
		Dot{
			X: 4,
			Y: 5,
		},
		Dot{
			X: 5,
			Y: 5,
		},
	}
	require.Equal(t, dm1, LocationToDotsMask(location1))
}

func Test_LocationToDotsMask_Second(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 0},
			{0, 1},
		},
	}
	location1 := Location{
		Dot{
			X: 5,
			Y: 4,
		},
		Dot{
			X: 4,
			Y: 3,
		},
	}
	require.Equal(t, dm1, LocationToDotsMask(location1))
}

func Test_LocationToDotsMask_ReturnsOnePointMaskForOneDotLocation(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1},
		},
	}
	location1 := Location{
		Dot{
			X: 4,
			Y: 3,
		},
	}
	require.Equal(t, dm1, LocationToDotsMask(location1))
}

func Test_LocationToDotsMask_ReturnsEmptyMaskForEmptyLocation(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{},
	}
	location1 := Location{}
	require.Equal(t, dm1, LocationToDotsMask(location1))
}

func Test_LocationToDotsMask_ReturnsOnePointMaskForLocationWithTheSameDots(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1},
		},
	}
	location1 := Location{
		Dot{
			X: 4,
			Y: 3,
		},
		Dot{
			X: 4,
			Y: 3,
		},
		Dot{
			X: 4,
			Y: 3,
		},
		Dot{
			X: 4,
			Y: 3,
		},
	}
	require.Equal(t, dm1, LocationToDotsMask(location1))
}

func Test_DotsMask_DotCount(t *testing.T) {
	dm1 := &DotsMask{
		mask: [][]uint8{
			{1, 1, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
	}
	require.Equal(t, uint16(6), dm1.DotCount())

	dm2 := &DotsMask{
		mask: [][]uint8{},
	}
	require.Equal(t, uint16(0), dm2.DotCount())

	dm3 := &DotsMask{
		mask: [][]uint8{
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
		},
	}
	require.Equal(t, uint16(0), dm3.DotCount())

	dm4 := &DotsMask{
		mask: [][]uint8{
			{1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 1},
			{0, 0, 0, 0, 0, 0},
		},
	}
	require.Equal(t, uint16(3), dm4.DotCount())
}
