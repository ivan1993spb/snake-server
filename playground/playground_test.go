package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/concurrent-map"
	"github.com/ivan1993spb/snake-server/engine"
)

func Test_Playground_CreateObject(t *testing.T) {
	object := &struct{}{}
	location := engine.Location{engine.Dot{0, 0}}

	pg, err := NewPlayground(100, 100)
	require.Nil(t, err)

	err = pg.CreateObject(object, location)
	require.Nil(t, err)
}

func Test_Playground_CreateObjectRandomRect(t *testing.T) {
	object := &struct{}{}

	pg, err := NewPlayground(100, 100)
	require.Nil(t, err)

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 100)
}

func Test_Playground_UpdateObject(t *testing.T) {
	t.SkipNow()
}

func Benchmark_Playground_UpdateObject(b *testing.B) {
	// TODO: Implement benchmark.
}

func Test_Playground_DeleteObject(t *testing.T) {
	pg := &Playground{
		cMap:       cmap.NewDefault(),
		objects:    make([]interface{}, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct{}{}

	pg.cMap.MSet(prepareMap(object, engine.Location{{1, 1}}))
	pg.objects = append(pg.objects, object)

	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))

	err := pg.DeleteObject(object, engine.Location{{1, 1}})
	require.Nil(t, err)

	require.Empty(t, pg.objects)
	require.Empty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))
}

func Test_Playground_CreateObjectAvailableDots_EmptyScene(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_CreateObjectAvailableDots_LocationNotAvailable(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_CreateObjectAvailableDots_LocationsIntersects(t *testing.T) {
	t.SkipNow()
	// TODO: Implement test.
}

func Test_Playground_UpdateObjectAvailableDots_EmptyMap(t *testing.T) {
	pg := &Playground{
		cMap:       cmap.NewDefault(),
		objects:    make([]interface{}, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct{}{}

	pg.cMap.MSet(prepareMap(object, engine.Location{{1, 1}}))
	pg.objects = append(pg.objects, object)

	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))

	location, err := pg.UpdateObjectAvailableDots(object, engine.Location{{1, 1}}, engine.Location{{1, 2}})
	require.Nil(t, err)
	require.Equal(t, engine.Location{{1, 2}}, location)
	require.Empty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))
	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 2}}.Hash()))
}

func Test_Playground_UpdateObjectAvailableDots_NotEmptyMap(t *testing.T) {
	pg := &Playground{
		cMap:       cmap.NewDefault(),
		objects:    make([]interface{}, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object1 := &struct {
		data string
	}{"first"}

	pg.cMap.MSet(prepareMap(object1, engine.Location{{1, 1}}))
	pg.objects = append(pg.objects, object1)

	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))

	object2 := &struct {
		data string
	}{"second"}

	pg.cMap.MSet(prepareMap(object2, engine.Location{{1, 3}}))
	pg.objects = append(pg.objects, object2)

	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 3}}.Hash()))

	location, err := pg.UpdateObjectAvailableDots(object1, engine.Location{{1, 1}}, engine.Location{{1, 2}})
	require.Nil(t, err)
	require.Equal(t, engine.Location{{1, 2}}, location)
	require.Empty(t, pg.cMap.MGet(engine.Location{{1, 1}}.Hash()))
	require.NotEmpty(t, pg.cMap.MGet(engine.Location{{1, 2}}.Hash()))
}
