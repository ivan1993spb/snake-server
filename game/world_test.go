package game

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

func Test_World_Events(t *testing.T) {
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")

	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := playground.NewPlayground(scene)
	require.NotNil(t, pg, "cannot initialize playground")

	world := &World{
		pg:          pg,
		chMain:      make(chan Event, worldEventsChanMainBufferSize),
		chsProxy:    make([]chan Event, 0),
		chsProxyMux: &sync.RWMutex{},
		stopGlobal:  make(chan struct{}, 0),
	}

	stopWorld := make(chan struct{})
	world.start(stopWorld)
	defer close(stopWorld)

	stop := make(chan struct{})

	chEventsFirst := world.Events(stop, 0)
	chEventsSecond := world.Events(stop, 4)

	object := &struct{}{}

	require.Nil(t, world.CreateObject(object, engine.Location{engine.NewDot(0, 0)}))

	require.Equal(t, EventTypeObjectCreate, (<-chEventsFirst).Type)

	require.Nil(t, world.UpdateObject(object, engine.Location{engine.NewDot(0, 0)}, engine.Location{engine.NewDot(1, 1)}))
	require.Equal(t, EventTypeObjectUpdate, (<-chEventsFirst).Type)

	require.NotNil(t, world.GetObjectByDot(engine.NewDot(1, 1)))
	require.Equal(t, EventTypeObjectChecked, (<-chEventsFirst).Type)

	require.Nil(t, world.DeleteObject(object, engine.Location{engine.NewDot(1, 1)}))
	require.Equal(t, EventTypeObjectDelete, (<-chEventsFirst).Type)

	require.Equal(t, EventTypeObjectCreate, (<-chEventsSecond).Type)
	require.Equal(t, EventTypeObjectUpdate, (<-chEventsSecond).Type)
	require.Equal(t, EventTypeObjectChecked, (<-chEventsSecond).Type)
	require.Equal(t, EventTypeObjectDelete, (<-chEventsSecond).Type)

	close(stop)

}

func Test_World_UpdateObject(t *testing.T) {
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")

	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := playground.NewPlayground(scene)
	require.NotNil(t, pg, "cannot initialize playground")

	object := &struct{}{}
	err = pg.CreateObject(object, engine.Location{engine.NewDot(0, 0)})
	require.Nil(t, err, "cannot create object to playground")

	world := &World{
		pg:          pg,
		chMain:      make(chan Event, worldEventsChanMainBufferSize),
		chsProxy:    make([]chan Event, 0),
		chsProxyMux: &sync.RWMutex{},
		stopGlobal:  make(chan struct{}, 0),
	}
	stop := make(chan struct{})
	world.start(stop)
	defer close(stop)

	err = world.UpdateObject(object, engine.Location{engine.NewDot(0, 0)}, engine.Location{engine.NewDot(1, 1)})
	require.Nil(t, err)
	require.True(t, pg.EntityExists(object, engine.Location{engine.NewDot(1, 1)}))
	require.False(t, pg.EntityExists(object, engine.Location{engine.NewDot(0, 0)}))
	require.False(t, pg.LocationExists(engine.Location{engine.NewDot(0, 0)}))
	require.True(t, pg.LocationExists(engine.Location{engine.NewDot(1, 1)}))

	err = world.UpdateObject(object, engine.Location{engine.NewDot(1, 1)}, engine.Location{engine.NewDot(3, 3)})
	require.Nil(t, err)
	require.True(t, pg.EntityExists(object, engine.Location{engine.NewDot(3, 3)}))
	require.False(t, pg.EntityExists(object, engine.Location{engine.NewDot(1, 1)}))
	require.False(t, pg.LocationExists(engine.Location{engine.NewDot(1, 1)}))
	require.True(t, pg.LocationExists(engine.Location{engine.NewDot(3, 3)}))

	err = world.UpdateObject(object, engine.Location{engine.NewDot(3, 3)}, engine.Location{engine.NewDot(0, 5)})
	require.Nil(t, err)
	require.True(t, pg.EntityExists(object, engine.Location{engine.NewDot(0, 5)}))
	require.False(t, pg.EntityExists(object, engine.Location{engine.NewDot(3, 3)}))
	require.False(t, pg.LocationExists(engine.Location{engine.NewDot(3, 3)}))
	require.True(t, pg.LocationExists(engine.Location{engine.NewDot(0, 5)}))
}

func Benchmark_World_UpdateObject(b *testing.B) {
	// TODO: Implement benchmark.
}
