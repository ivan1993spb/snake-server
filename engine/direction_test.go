package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Direction_CalculateDirection(t *testing.T) {
	tests := []struct {
		from        Dot
		to          Dot
		expectedDir Direction
	}{
		{Dot{0, 0}, Dot{0, 1}, DirectionSouth},
		{Dot{0, 0}, Dot{0, 2}, DirectionSouth},
		{Dot{0, 10}, Dot{0, 20}, DirectionSouth},
		{Dot{0, 10}, Dot{0, 8}, DirectionNorth},
		{Dot{0, 50}, Dot{0, 8}, DirectionNorth},
		{Dot{0, 50}, Dot{10, 50}, DirectionEast},
		{Dot{20, 50}, Dot{22, 50}, DirectionEast},
		{Dot{20, 50}, Dot{10, 50}, DirectionWest},
		{Dot{20, 50}, Dot{19, 50}, DirectionWest},
	}

	for i, test := range tests {
		actualDir := CalculateDirection(test.from, test.to)
		require.Equal(t, test.expectedDir, actualDir, fmt.Sprintf("number %d", i))
	}
}

func Test_Direction_CalculateDirection_ReturnsValidDirectionForEqualDots(t *testing.T) {
	require.True(t, ValidDirection(CalculateDirection(Dot{}, Dot{})))
	require.True(t, ValidDirection(CalculateDirection(Dot{10, 10}, Dot{10, 10})))
}

func Test_ValidDirection_ValidatesDirectionsCorrectly(t *testing.T) {
	tests := []struct {
		direction   Direction
		expectedRes bool
	}{
		{DirectionNorth, true},
		{DirectionEast, true},
		{DirectionSouth, true},
		{DirectionWest, true},
		{22, false},
		{44, false},
	}

	for i, test := range tests {
		require.Equal(t, test.expectedRes, ValidDirection(test.direction), fmt.Sprintf("number %d", i))
	}
}

func Test_Direction_Reverse_ReversesDirection(t *testing.T) {
	tests := []struct {
		directionInput    Direction
		directionExpected Direction
		expectError       bool
	}{
		{DirectionNorth, DirectionSouth, false},
		{DirectionEast, DirectionWest, false},
		{DirectionSouth, DirectionNorth, false},
		{DirectionWest, DirectionEast, false},
		{22, 0, true},
		{44, 0, true},
	}

	for i, test := range tests {
		actualDirection, err := test.directionInput.Reverse()
		if test.expectError {
			require.NotNil(t, err, fmt.Sprintf("number %d not error", i))
		} else {
			require.Nil(t, err, fmt.Sprintf("number %d error", i))
		}
		require.Equal(t, test.directionExpected, actualDirection, fmt.Sprintf("number %d", i))
	}
}

func Test_Direction_MarshalJSON(t *testing.T) {
	tests := []struct {
		direction    Direction
		expectedJSON []byte
		expectedErr  error
	}{
		{DirectionNorth, directionsJSON[DirectionNorth], nil},
		{DirectionEast, directionsJSON[DirectionEast], nil},
		{DirectionSouth, directionsJSON[DirectionSouth], nil},
		{DirectionWest, directionsJSON[DirectionWest], nil},
		{22, unknownDirectionJSON, &ErrDirectionMarshal{
			Err: &ErrInvalidDirection{
				Direction: 22,
			},
		}},
		{12, unknownDirectionJSON, &ErrDirectionMarshal{
			Err: &ErrInvalidDirection{
				Direction: 12,
			},
		}},
	}

	for i, test := range tests {
		json, err := test.direction.MarshalJSON()
		require.Equal(t, test.expectedJSON, json, fmt.Sprintf("number %d", i))
		require.Equal(t, test.expectedErr, err, fmt.Sprintf("number %d", i))
	}
}
