package engine

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scene_Locate_SquareScene(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{Dot{1, 1}, Dot{1, 2}}
	err := scene.Locate(location)
	require.Nil(t, err)
	require.Equal(t, []Location{location}, scene.locations)
}

func Test_Scene_Locate_RectScene(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 200,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{Dot{1, 1}, Dot{1, 2}}
	err := scene.Locate(location)
	require.Nil(t, err)
	require.Equal(t, []Location{location}, scene.locations)
}

func Benchmark_Scene_Locate(b *testing.B) {
	scene := func() *Scene {
		return &Scene{
			area: Area{
				width:  100,
				height: 100,
			},
			locations: []Location{
				{
					{0, 1},
					{0, 2},
					{0, 3},
					{0, 4},
					{0, 5},
					{0, 6},
				},
				{
					{4, 1},
					{4, 2},
					{4, 3},
					{4, 4},
					{4, 5},
					{4, 6},
				},
				{
					{5, 2},
					{6, 2},
					{7, 2},
					{8, 2},
					{9, 2},
					{10, 2},
				},
			},
			locationsMutex: &sync.RWMutex{},
		}
	}

	for n := 0; n < b.N; n++ {
		scene().Locate(Location{
			{10, 11},
			{10, 12},
			{10, 13},
			{10, 14},
			{10, 15},
			{10, 16},
			{10, 17},
			{10, 18},
			{10, 19},
			{10, 20},
			{10, 21},
		})
	}
}

func Benchmark_Scene_Relocate(b *testing.B) {
	// TODO: Implement benchmark.
}

func Test_Scene_LocateRandomRect_SquareScene(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location, err := scene.LocateRandomRect(1, 5)
	require.Nil(t, err)
	require.Equal(t, []Location{location}, scene.locations)
}

func Test_Scene_LocateRandomRect_RectScene(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  150,
			height: 99,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location, err := scene.LocateRandomRect(1, 5)
	require.Nil(t, err)
	require.Equal(t, []Location{location}, scene.locations)
}

func Test_Scene_LocateAvailableDots_EmptySquareScene(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{Dot{1, 1}, Dot{1, 2}}
	locationActual := scene.LocateAvailableDots(location)
	require.Equal(t, []Location{location}, scene.locations)
	require.Equal(t, []Location{locationActual}, scene.locations)
	require.Equal(t, location, locationActual)
}

func Test_Scene_LocateAvailableDots_LocationNotAvailable(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations: []Location{
			{Dot{1, 1}, Dot{1, 2}},
		},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{Dot{1, 1}, Dot{1, 2}}
	locationActual := scene.LocateAvailableDots(location)
	require.Equal(t, Location{}, locationActual)
}

func Test_Scene_LocateAvailableDots_LocationsIntersects(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations: []Location{
			{Dot{1, 1}, Dot{1, 2}, Dot{1, 3}, Dot{1, 4}},
		},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{Dot{1, 1}, Dot{1, 0}}
	locationActual := scene.LocateAvailableDots(location)
	require.Equal(t, Location{Dot{1, 0}}, locationActual)
	require.Equal(t, []Location{
		{Dot{1, 1}, Dot{1, 2}, Dot{1, 3}, Dot{1, 4}},
		{Dot{1, 0}},
	}, scene.locations)
}

func Test_Scene_LocateRandomRectMargin_LocatesValidRectWithMargin(t *testing.T) {
	scene := &Scene{
		area: Area{
			width:  100,
			height: 100,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location, err := scene.LocateRandomRectMargin(2, 3, 2)
	require.Nil(t, err)
	require.Len(t, location, 6)
	require.Len(t, scene.locations, 1)
	require.Len(t, scene.locations[0], 6)
}
