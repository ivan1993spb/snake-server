package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_Playground_ObjectExists(t *testing.T) {
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")
	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{
				object:   object,
				location: engine.Location{engine.NewDot(0, 0)},
			},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	require.True(t, pg.ObjectExists(object))
}
