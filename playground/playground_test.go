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

	objectExists := &struct{}{}

	err = scene.Locate(engine.Location{engine.NewDot(0, 0)})
	require.Nil(t, err)

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{
				object:   objectExists,
				location: engine.Location{engine.NewDot(0, 0)},
			},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	require.True(t, pg.ObjectExists(objectExists))

	objectNotExists := &struct{}{}

	require.False(t, pg.ObjectExists(objectNotExists))
}

func Test_Playground_CreateObject(t *testing.T) {
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")
	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}
	location := engine.Location{engine.NewDot(0, 0)}

	pg := &Playground{
		scene:         scene,
		entities:      []entity{},
		entitiesMutex: &sync.RWMutex{},
	}

	err = pg.CreateObject(object, location)
	require.Nil(t, err)
}