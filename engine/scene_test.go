package engine

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scene_Locate(t *testing.T) {
	scene := &Scene{
		area: &Area{
			width:  100,
			height: 100,
		},
		locations:      []Location{},
		locationsMutex: &sync.RWMutex{},
	}

	location := Location{&Dot{1, 1}, &Dot{1, 2}}
	err := scene.Locate(location)
	require.Nil(t, err)
	require.Equal(t, []Location{location}, scene.locations)
}
