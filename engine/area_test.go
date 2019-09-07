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

		{Dot{99, 0}, DirectionEast, 1, Dot{0, 0}, nil},
		{Dot{0, 99}, DirectionSouth, 1, Dot{0, 0}, nil},

		{Dot{0, 99}, DirectionSouth, 150, Dot{0, 49}, nil},
		{Dot{0, 99}, DirectionNorth, 150, Dot{0, 49}, nil},
		{Dot{0, 0}, DirectionEast, 150, Dot{50, 0}, nil},
		{Dot{0, 0}, DirectionWest, 150, Dot{50, 0}, nil},

		{Dot{30, 30}, DirectionWest, 5, Dot{25, 30}, nil},

		// Error case: invalid direction
		{Dot{0, 0}, 21, 150, Dot{}, &ErrNavigation{
			Err: &ErrInvalidDirection{
				Direction: 21,
			},
		}},
		// Error case: area doesn't contains dot
		{Dot{250, 25}, DirectionWest, 150, Dot{}, &ErrNavigation{
			Err: &ErrAreaNotContainsDot{
				Dot: Dot{250, 25},
			},
		}},
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

func Test_Area_Navigate_SquareMaxArea255x255(t *testing.T) {
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

		{Dot{0, 0}, DirectionWest, 1, Dot{254, 0}, nil},
		{Dot{0, 0}, DirectionEast, 1, Dot{1, 0}, nil},
		{Dot{0, 0}, DirectionNorth, 1, Dot{0, 254}, nil},
		{Dot{0, 0}, DirectionSouth, 1, Dot{0, 1}, nil},

		{Dot{254, 0}, DirectionEast, 1, Dot{0, 0}, nil},
		{Dot{0, 254}, DirectionSouth, 1, Dot{0, 0}, nil},
	}

	area := Area{
		width:  255,
		height: 255,
	}

	for i, test := range tests {
		actualDot, actualErr := area.Navigate(test.inputDot, test.inputDir, test.inputDis)
		require.Equal(t, test.expectedDot, actualDot, fmt.Sprintf("number %d", i))
		require.Equal(t, test.expectedErr, actualErr, fmt.Sprintf("number %d", i))
	}
}

func Test_Area_MarshalJSON(t *testing.T) {
	tests := []struct {
		area Area
		json []byte
	}{
		{
			Area{10, 10},
			[]byte("[10,10]"),
		},
		{
			Area{255, 255},
			[]byte("[255,255]"),
		},
		{
			Area{0, 0},
			[]byte("[0,0]"),
		},
		{
			Area{0, 1},
			[]byte("[0,1]"),
		},
		{
			Area{2, 1},
			[]byte("[2,1]"),
		},
		{
			Area{255, 1},
			[]byte("[255,1]"),
		},
		{
			Area{255, 100},
			[]byte("[255,100]"),
		},
		{
			Area{0, 255},
			[]byte("[0,255]"),
		},
	}

	for i, test := range tests {
		actualJSON, err := test.area.MarshalJSON()
		require.Nil(t, err, "test %d", i)
		require.Equal(t, test.json, actualJSON, "test %d", i)
	}
}
