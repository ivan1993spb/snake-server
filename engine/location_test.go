package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Location_ContainsDot(t *testing.T) {
	tests := []struct {
		location Location
		dot      Dot
		contains bool
	}{
		{Location{Dot{0, 0}}, Dot{0, 0}, true},
		{Location{Dot{0, 0}, Dot{0, 1}, Dot{0, 2}}, Dot{0, 1}, true},
		{Location{Dot{1, 2}, Dot{1, 3}, Dot{1, 4}, Dot{2, 3}}, Dot{5, 1}, false},
		{Location{Dot{5, 2}, Dot{5, 1}, Dot{5, 0}}, Dot{5, 3}, false},
		{Location{Dot{20, 20}, Dot{21, 21}, Dot{20, 21}}, Dot{21, 20}, false},
	}

	for i, test := range tests {
		test.location.Contains(test.dot)
		require.Equal(t, test.contains, test.location.Contains(test.dot), fmt.Sprintf("number: %d", i))
	}
}

func Test_Location_Copy(t *testing.T) {
	locations := []Location{
		{Dot{0, 0}, Dot{0, 1}, Dot{0, 2}},
		{Dot{1, 0}, Dot{2, 1}, Dot{3, 2}},
		{Dot{1, 3}, Dot{5, 2}, Dot{3, 2}, Dot{1, 0}, Dot{2, 1}},
	}

	for i, location := range locations {
		require.Equal(t, location, location.Copy(), fmt.Sprintf("number: %d", i))
	}
}

func Test_Location_Equals(t *testing.T) {
	tests := []struct {
		first  Location
		second Location
		equals bool
	}{
		{
			first:  Location{Dot{0, 0}, Dot{0, 1}, Dot{0, 2}},
			second: Location{Dot{0, 0}, Dot{0, 1}, Dot{0, 2}},
			equals: true,
		},
		{
			first:  Location{Dot{0, 0}, Dot{0, 1}, Dot{0, 2}},
			second: Location{Dot{1, 0}, Dot{0, 1}, Dot{0, 2}},
			equals: false,
		},
		{
			first:  Location{Dot{0, 0}, Dot{1, 0}, Dot{0, 2}},
			second: Location{Dot{1, 0}, Dot{0, 0}, Dot{0, 2}},
			equals: true,
		},
	}

	for i, test := range tests {
		if test.equals {
			require.True(t, test.first.Equals(test.second), fmt.Sprintf("number: %d", i))
			require.True(t, test.second.Equals(test.first), fmt.Sprintf("number: %d", i))
		} else {
			require.False(t, test.first.Equals(test.second), fmt.Sprintf("number: %d", i))
			require.False(t, test.second.Equals(test.first), fmt.Sprintf("number: %d", i))
		}
	}
}
