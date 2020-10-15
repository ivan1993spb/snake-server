package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/concurrent-map"
	"github.com/ivan1993spb/snake-server/engine"
)

func Test_NewPlaygroundCMap_CreatesPlaygroundCMap(t *testing.T) {
	tests := []struct {
		width, height uint8
	}{
		{
			width:  200,
			height: 100,
		},
		{
			width:  80,
			height: 100,
		},
		{
			width:  30,
			height: 22,
		},
		{
			width:  255,
			height: 200,
		},
	}

	for i, test := range tests {
		pg, err := NewPlaygroundCMap(test.width, test.height)
		require.Nil(t, err, "test %d", i)
		require.NotNil(t, pg, "test %d", i)
		require.Equal(t, test.width, pg.area.Width(), "test %d", i)
		require.Equal(t, test.height, pg.area.Height(), "test %d", i)
		require.Equal(t, uint16(test.width)*uint16(test.height), pg.area.Size(), "test %d", i)
	}
}

func Test_PlaygroundCMap_CreateObject(t *testing.T) {
	object := &struct {
		a int
	}{
		10,
	}
	location := engine.Location{engine.Dot{0, 0}}

	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	err := pg.CreateObject(object, location)
	require.Nil(t, err)
	require.Len(t, pg.objects, 1)
	require.Equal(t, object, pg.objects[0])
	actual, ok := pg.cMap.Get(engine.Dot{0, 0}.Hash())
	require.True(t, ok)
	require.Equal(t, object, actual)
}

func Test_PlaygroundCMap_CreateObjectRandomRect(t *testing.T) {
	object := &struct {
		str string
	}{
		"ok",
	}

	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 100)

	for _, dot := range location {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object, actual)
	}
}

func Test_PlaygroundCMap_UpdateObject(t *testing.T) {
	t.SkipNow()
}

type TestStructure struct {
	a, b, c, d int64
	e, f, g, h float64
}

func Benchmark_PlaygroundCMap_UpdateObject(b *testing.B) {
	const (
		areaWidth  = 150
		areaHeight = 100
	)

	const (
		first = iota
		second
		count
	)

	pg, err := NewPlaygroundCMap(areaWidth, areaHeight)

	if err != nil {
		b.Fatal(err)
	}

	object := &TestStructure{}

	locations := [count]engine.Location{
		first:  engine.NewRect(0, 0, areaWidth, areaHeight/2).Location(),
		second: engine.NewRect(0, 0, areaWidth/2, areaHeight).Location(),
	}

	location := locations[first]
	pg.CreateObject(object, location)

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := pg.UpdateObject(object, location, locations[i%count]); err != nil {
			b.Fatal(err)
		}
		location = locations[i%count]
	}
}

func Test_PlaygroundCMap_DeleteObject(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
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

func Test_PlaygroundCMap_CreateObjectAvailableDots_EmptySquareScene(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct {
		data string
	}{"first"}

	location := engine.Location{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
		{2, 6},
		{3, 6},
		{4, 6},
		{5, 6},
	}

	actualLocation, err := pg.CreateObjectAvailableDots(object, location)
	require.Nil(t, err)
	require.True(t, location.Equals(actualLocation))

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectAvailableDots_NotEmptyScene(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	// Object to create
	object1 := &struct {
		data string
	}{"first"}
	location1 := engine.Location{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
	}

	// Located object
	object2 := &struct {
		data string
	}{"second"}
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.cMap.MSet(prepareMap(object2, location2))
	pg.objects = append(pg.objects, object2)
	for _, dot := range location2 {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.Nil(t, err)
	require.True(t, location1.Equals(location1Actual))

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location1.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object1, actualObject, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object2, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectAvailableDots_LocationNotAvailable(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	// Object to create
	object1 := &struct {
		data string
	}{"first"}
	location1 := engine.Location{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
	}

	// Located object
	object2 := &struct {
		data string
	}{"second"}
	location2 := location1.Copy()

	pg.cMap.MSet(prepareMap(object2, location2))
	pg.objects = append(pg.objects, object2)
	for _, dot := range location2 {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.NotNil(t, err)
	require.Nil(t, location1Actual)

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object2, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectAvailableDots_LocationsIntersects(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	// Object to create
	object1 := &struct {
		data string
	}{"first"}
	location1 := engine.Location{
		{2, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{2, 5},
		{2, 6},
	}
	location1Expected := engine.Location{
		{1, 2},
		{1, 3},
		{1, 4},
	}

	// Located object
	object2 := &struct {
		data string
	}{"second"}
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.cMap.MSet(prepareMap(object2, location2))
	pg.objects = append(pg.objects, object2)
	for _, dot := range location2 {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.Nil(t, err)
	require.True(t, location1Expected.Equals(location1Actual))

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location1Expected.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object1, actualObject, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object2, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_UpdateObjectAvailableDots_EmptyMap(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
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

func Test_PlaygroundCMap_UpdateObjectAvailableDots_NotEmptyMap(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
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

func Test_PlaygroundCMap_UpdateObjectAvailableDots_NotEmptyMap_BigObjects(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object1 := &struct {
		data string
	}{"first"}
	location1Old := engine.Location{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{1, 6},
	}
	location1New := engine.Location{
		{2, 0},
		{2, 1},
		{2, 2},
		{1, 3},
		{1, 4},
		{1, 5},
		{3, 6},
		{3, 6},
	}
	location1Expected := engine.Location{
		{2, 0},
		{1, 3},
		{1, 4},
		{1, 5},
		{3, 6},
	}

	pg.cMap.MSet(prepareMap(object1, location1Old))
	pg.objects = append(pg.objects, object1)
	for _, dot := range location1Old {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object1, actual)
	}

	object2 := &struct {
		data string
	}{"second"}
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.cMap.MSet(prepareMap(object2, location2))
	pg.objects = append(pg.objects, object2)
	for _, dot := range location2 {
		require.True(t, pg.cMap.Has(dot.Hash()))
		actual, ok := pg.cMap.Get(dot.Hash())
		require.True(t, ok)
		require.Equal(t, object2, actual)
	}

	location1Actual, err := pg.UpdateObjectAvailableDots(object1, location1Old, location1New)
	require.Nil(t, err)
	require.True(t, location1Expected.Equals(location1Actual))

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location1Expected.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object1, actualObject, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object2, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectRandomDot_SquareEmptyPlaygroundCMap(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomDot(object)
	require.Nil(t, err)
	require.Len(t, location, 1)

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectRandomDot_EmptyPlaygroundCMap(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(200, 100),
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomDot(object)
	require.Nil(t, err)
	require.Len(t, location, 1)

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectRandomRect_SquareEmptyPlaygroundCMap_SquareRect(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 10*10)

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}

func Test_PlaygroundCMap_CreateObjectRandomRect_SquareEmptyPlaygroundCMap(t *testing.T) {
	pg := &PlaygroundCMap{
		cMap:       cmap.NewDefault(),
		objects:    make([]engine.Object, 0),
		objectsMux: &sync.RWMutex{},
		area:       engine.MustArea(100, 100),
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomRect(object, 10, 8)
	require.Nil(t, err)
	require.Len(t, location, 10*8)

	for _, dot := range pg.area.Dots() {
		actualObject, ok := pg.cMap.Get(dot.Hash())

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, object, actualObject, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualObject, "dot %s", dot)
		}
	}
}
