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
