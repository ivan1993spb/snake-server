package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Rect_Location(t *testing.T) {
	rect := Rect{
		x: 0,
		y: 0,
		w: 1,
		h: 5,
	}

	require.Equal(t, Location{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
		{0, 4},
	}, rect.Location())
}

func Test_Rect_Dot_HorizontalRect(t *testing.T) {
	rect := Rect{
		x: 0,
		y: 0,
		w: 1,
		h: 5,
	}

	require.Equal(t, Dot{0, 0}, rect.Dot(0))
	require.Equal(t, Dot{0, 1}, rect.Dot(1))
	require.Equal(t, Dot{0, 2}, rect.Dot(2))
	require.Equal(t, Dot{0, 3}, rect.Dot(3))
	require.Equal(t, Dot{0, 4}, rect.Dot(4))
}

func Test_Rect_Dot_VerticalRect(t *testing.T) {
	rect := Rect{
		x: 0,
		y: 0,
		w: 5,
		h: 1,
	}

	require.Equal(t, Dot{0, 0}, rect.Dot(0))
	require.Equal(t, Dot{1, 0}, rect.Dot(1))
	require.Equal(t, Dot{2, 0}, rect.Dot(2))
	require.Equal(t, Dot{3, 0}, rect.Dot(3))
	require.Equal(t, Dot{4, 0}, rect.Dot(4))
}

func Test_Rect_Dot_VerticalRectWithXY(t *testing.T) {
	rect := Rect{
		x: 5,
		y: 5,
		w: 5,
		h: 1,
	}

	require.Equal(t, Dot{5, 5}, rect.Dot(0))
	require.Equal(t, Dot{6, 5}, rect.Dot(1))
	require.Equal(t, Dot{7, 5}, rect.Dot(2))
	require.Equal(t, Dot{8, 5}, rect.Dot(3))
	require.Equal(t, Dot{9, 5}, rect.Dot(4))
}

func Test_Rect_Dot_HorizontalRectWithXY(t *testing.T) {
	rect := Rect{
		x: 5,
		y: 5,
		w: 1,
		h: 5,
	}

	require.Equal(t, Dot{5, 5}, rect.Dot(0))
	require.Equal(t, Dot{5, 6}, rect.Dot(1))
	require.Equal(t, Dot{5, 7}, rect.Dot(2))
	require.Equal(t, Dot{5, 8}, rect.Dot(3))
	require.Equal(t, Dot{5, 9}, rect.Dot(4))
}

func Test_Rect_Dot_SquareRectWithXY(t *testing.T) {
	rect := Rect{
		x: 5,
		y: 5,
		w: 3,
		h: 3,
	}

	require.Equal(t, Dot{5, 5}, rect.Dot(0))
	require.Equal(t, Dot{6, 5}, rect.Dot(1))
	require.Equal(t, Dot{7, 5}, rect.Dot(2))
	require.Equal(t, Dot{5, 6}, rect.Dot(3))
	require.Equal(t, Dot{6, 6}, rect.Dot(4))
	require.Equal(t, Dot{7, 6}, rect.Dot(5))
	require.Equal(t, Dot{5, 7}, rect.Dot(6))
	require.Equal(t, Dot{6, 7}, rect.Dot(7))
	require.Equal(t, Dot{7, 7}, rect.Dot(8))
}

func Test_Rect_MarshalJSON(t *testing.T) {
	tests := []struct {
		rect Rect
		json []byte
	}{
		{
			Rect{},
			[]byte("[0,0,0,0]"),
		},
		{
			Rect{1, 2, 3, 4},
			[]byte("[1,2,3,4]"),
		},
		{
			Rect{4, 3, 2, 1},
			[]byte("[4,3,2,1]"),
		},
		{
			Rect{255, 200, 160, 100},
			[]byte("[255,200,160,100]"),
		},
		{
			Rect{255, 200, 160, 255},
			[]byte("[255,200,160,255]"),
		},
	}

	for i, test := range tests {
		actualJSON, err := test.rect.MarshalJSON()
		require.Nil(t, err, "test %d", i)
		require.Equal(t, test.json, actualJSON, "test %d", i)
	}
}
