package engine

import (
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
