package engine

import (
	"fmt"
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

func Test_Rect_Width_ReturnsRectWidth(t *testing.T) {
	tests := []struct {
		rect Rect
	}{
		{Rect{0, 0, 10, 3}},
		{Rect{0, 0, 22, 4}},
		{Rect{0, 0, 123, 5}},
		{Rect{0, 0, 0, 233}},
	}

	for i, test := range tests {
		require.Equal(t, test.rect.w, test.rect.Width(), fmt.Sprintf("number: %d", i))
	}
}

func Test_Rect_Height_ReturnsRectHeight(t *testing.T) {
	tests := []struct {
		rect Rect
	}{
		{Rect{0, 0, 10, 3}},
		{Rect{0, 0, 22, 4}},
		{Rect{0, 0, 123, 5}},
		{Rect{0, 0, 0, 233}},
	}

	for i, test := range tests {
		require.Equal(t, test.rect.h, test.rect.Height(), fmt.Sprintf("number: %d", i))
	}
}

func Test_Rect_X_ReturnsRectX(t *testing.T) {
	tests := []struct {
		rect Rect
	}{
		{Rect{3, 23, 10, 3}},
		{Rect{32, 1, 22, 4}},
		{Rect{1, 32, 123, 5}},
		{Rect{42, 231, 0, 233}},
	}

	for i, test := range tests {
		require.Equal(t, test.rect.x, test.rect.X(), fmt.Sprintf("number: %d", i))
	}
}

func Test_Rect_Y_ReturnsRectY(t *testing.T) {
	tests := []struct {
		rect Rect
	}{
		{Rect{213, 231, 10, 3}},
		{Rect{23, 32, 22, 4}},
		{Rect{123, 132, 123, 5}},
		{Rect{22, 3, 0, 233}},
	}

	for i, test := range tests {
		require.Equal(t, test.rect.y, test.rect.Y(), fmt.Sprintf("number: %d", i))
	}
}

func Test_Rect_Dots_ReturnsDotList(t *testing.T) {
	tests := []struct {
		rect Rect
		dots []Dot
	}{
		{Rect{213, 231, 10, 3}, []Dot{
			// 0
			{213, 231}, {214, 231}, {215, 231}, {216, 231}, {217, 231},
			{218, 231}, {219, 231}, {220, 231}, {221, 231}, {222, 231},
			// 1
			{213, 232}, {214, 232}, {215, 232}, {216, 232}, {217, 232},
			{218, 232}, {219, 232}, {220, 232}, {221, 232}, {222, 232},
			// 2
			{213, 233}, {214, 233}, {215, 233}, {216, 233}, {217, 233},
			{218, 233}, {219, 233}, {220, 233}, {221, 233}, {222, 233},
		}},
		{Rect{23, 32, 3, 3}, []Dot{
			{23, 32}, {24, 32}, {25, 32},
			{23, 33}, {24, 33}, {25, 33},
			{23, 34}, {24, 34}, {25, 34},
		}},
		{Rect{123, 132, 1, 1}, []Dot{
			{123, 132},
		}},
		{Rect{22, 3, 0, 233}, []Dot{}},
		{Rect{22, 3, 123, 0}, []Dot{}},
		{Rect{22, 3, 0, 0}, []Dot{}},
	}

	for i, test := range tests {
		require.Equal(t, test.dots, test.rect.Dots(), fmt.Sprintf("number: %d", i))
	}
}
