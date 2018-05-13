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
	world.start()

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

	world.stop()
}
