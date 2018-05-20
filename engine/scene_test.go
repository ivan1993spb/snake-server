package engine

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scene_Locate(t *testing.T) {
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

func Benchmark_Scene_Locate(b *testing.B) {
	// TODO: Implement benchmark.
}

func Benchmark_Scene_Relocate(b *testing.B) {
	// TODO: Implement benchmark.
}

func Test_Scene_LocateRandomRect(t *testing.T) {
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

func Test_Scene_LocateAvailableDots_EmptyScene(t *testing.T) {
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
