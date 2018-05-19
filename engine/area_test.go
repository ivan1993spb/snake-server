package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewArea_InvalidSize(t *testing.T) {
	var area Area
	var err error

	area, err = NewArea(0, 0)
	require.NotNil(t, err)
	require.Equal(t, Area{}, area)

	area, err = NewArea(1, 0)
	require.NotNil(t, err)
	require.Equal(t, Area{}, area)

	area, err = NewArea(0, 1)
	require.NotNil(t, err)
	require.Equal(t, Area{}, area)
}

func Test_NewArea_ValidSize(t *testing.T) {
	var err error

	_, err = NewArea(1, 1)
	require.Nil(t, err)

	_, err = NewArea(100, 100)
	require.Nil(t, err)
}

func Test_Area_Navigate_SquareArea100x100(t *testing.T) {
	tests := []struct {
		inputDot    Dot
		inputDir    Direction
		inputDis    uint8
		expectedDot Dot
		expectedErr error
	}{
		{Dot{0, 0}, DirectionWest, 0, Dot{0, 0}, nil},
		{Dot{0, 0}, DirectionEast, 0, Dot{0, 0}, nil},
		{Dot{0, 0}, DirectionNorth, 0, Dot{0, 0}, nil},
		{Dot{0, 0}, DirectionSouth, 0, Dot{0, 0}, nil},

		{Dot{0, 0}, DirectionWest, 1, Dot{99, 0}, nil},
		{Dot{0, 0}, DirectionEast, 1, Dot{1, 0}, nil},
		{Dot{0, 0}, DirectionNorth, 1, Dot{0, 99}, nil},
		{Dot{0, 0}, DirectionSouth, 1, Dot{0, 1}, nil},
	}

	area := Area{
		width:  100,
		height: 100,
	}

	for i, test := range tests {
		actualDot, actualErr := area.Navigate(test.inputDot, test.inputDir, test.inputDis)
		require.Equal(t, test.expectedDot, actualDot, fmt.Sprintf("number %d", i))
		require.Equal(t, test.expectedErr, actualErr, fmt.Sprintf("number %d", i))
	}
}
