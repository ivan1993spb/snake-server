package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_Playground_ObjectExists(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	objectExists := &struct{}{}

	err = scene.Locate(engine.Location{engine.Dot{0, 0}})
	require.Nil(t, err)

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{
				object:   objectExists,
				location: engine.Location{engine.Dot{0, 0}},
			},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	require.True(t, pg.ObjectExists(objectExists))

	objectNotExists := &struct{}{}

	require.False(t, pg.ObjectExists(objectNotExists))
}

func Test_Playground_CreateObject(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}
	location := engine.Location{engine.Dot{0, 0}}

	pg := &Playground{
		scene:         scene,
		entities:      []entity{},
		entitiesMutex: &sync.RWMutex{},
	}

	err = pg.CreateObject(object, location)
	require.Nil(t, err)
}

func Test_Playground_CreateObjectRandomRect(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}

	pg := &Playground{
		scene:         scene,
		entities:      []entity{},
		entitiesMutex: &sync.RWMutex{},
	}

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 100)
}

func Test_Playground_UpdateObject(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{object, engine.Location{engine.Dot{0, 0}}},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	err = scene.Locate(engine.Location{engine.Dot{0, 0}})
	require.Nil(t, err)

	err = pg.UpdateObject(object, engine.Location{engine.Dot{0, 0}}, engine.Location{engine.Dot{1, 1}})
	require.Nil(t, err)
	require.True(t, scene.Located(engine.Location{engine.Dot{1, 1}}))
	require.False(t, scene.Located(engine.Location{engine.Dot{0, 0}}))

	err = pg.UpdateObject(object, engine.Location{engine.Dot{1, 1}}, engine.Location{engine.Dot{2, 2}})
	require.Nil(t, err)
	require.True(t, scene.Located(engine.Location{engine.Dot{2, 2}}))
	require.False(t, scene.Located(engine.Location{engine.Dot{1, 1}}))

	err = pg.UpdateObject(object, engine.Location{engine.Dot{2, 2}}, engine.Location{engine.Dot{0, 0}})
	require.Nil(t, err)
	require.True(t, scene.Located(engine.Location{engine.Dot{0, 0}}))
	require.False(t, scene.Located(engine.Location{engine.Dot{2, 2}}))
}

func Benchmark_Playground_UpdateObject(b *testing.B) {
	// TODO: Implement benchmark.
}

func Test_Playground_DeleteObject(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	object := &struct{}{}

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{object, engine.Location{engine.Dot{0, 0}}},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	err = scene.Locate(engine.Location{engine.Dot{0, 0}})
	require.Nil(t, err)

	err = pg.DeleteObject(object, engine.Location{engine.Dot{0, 0}})
	require.Nil(t, err)
	require.False(t, scene.Located(engine.Location{engine.Dot{0, 0}}))
	require.Len(t, pg.entities, 0)
}

func Test_Playground_CreateObjectAvailableDots_EmptyScene(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := &Playground{
		scene:         scene,
		entities:      []entity{},
		entitiesMutex: &sync.RWMutex{},
	}

	object := &struct{}{}
	location := engine.Location{
		engine.Dot{0, 0},
		engine.Dot{0, 1},
		engine.Dot{0, 2},
	}

	actualLocation, err := pg.CreateObjectAvailableDots(object, location)
	require.Nil(t, err)
	require.Equal(t, location, actualLocation)
}

func Test_Playground_CreateObjectAvailableDots_LocationNotAvailable(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	locatedObject := &struct{ int }{0}

	pg := &Playground{
		scene: scene,
		entities: []entity{
			{locatedObject, engine.Location{
				engine.Dot{0, 0},
				engine.Dot{0, 1},
				engine.Dot{0, 2},
			}},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	require.Nil(t, scene.Locate(engine.Location{
		engine.Dot{0, 0},
		engine.Dot{0, 1},
		engine.Dot{0, 2},
	}))

	object := &struct{}{}

	location := engine.Location{
		engine.Dot{0, 0},
		engine.Dot{0, 1},
		engine.Dot{0, 2},
	}

	actualLocation, err := pg.CreateObjectAvailableDots(object, location)
	require.NotNil(t, err)
	require.Nil(t, actualLocation)
}

func Test_Playground_CreateObjectAvailableDots_LocationsIntersects(t *testing.T) {
	scene, err := engine.NewScene(100, 100)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	locatedObject := &struct{ int }{0}
	pg := &Playground{
		scene: scene,
		entities: []entity{
			{locatedObject, engine.Location{
				engine.Dot{0, 0},
				engine.Dot{0, 1},
				engine.Dot{0, 2},
			}},
		},
		entitiesMutex: &sync.RWMutex{},
	}

	require.Nil(t, scene.Locate(engine.Location{
		engine.Dot{0, 0},
		engine.Dot{0, 1},
		engine.Dot{0, 2},
	}))

	object := &struct{}{}

	location := engine.Location{
		engine.Dot{0, 0},
		engine.Dot{1, 1},
		engine.Dot{1, 2},
	}

	actualLocation, err := pg.CreateObjectAvailableDots(object, location)
	require.Nil(t, err)
	require.Equal(t, engine.Location{
		engine.Dot{1, 1},
		engine.Dot{1, 2},
	}, actualLocation)
}
