package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_NewExperimentalPlayground_CreatesPlayground(t *testing.T) {
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
		pg, err := NewExperimentalPlayground(test.width, test.height)
		require.Nil(t, err, "test %d", i)
		require.NotNil(t, pg, "test %d", i)
		require.Equal(t, test.width, pg.Area().Width(), "test %d", i)
		require.Equal(t, test.height, pg.Area().Height(), "test %d", i)
		require.Equal(t, uint16(test.width)*uint16(test.height), pg.Area().Size(), "test %d", i)
	}
}

func Test_ExperimentalPlayground_CreateObject(t *testing.T) {
	area := engine.MustArea(100, 100)

	object := &struct {
		a int
	}{
		10,
	}

	location := engine.Location{engine.Dot{0, 0}}

	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	err := pg.CreateObject(object, location)
	require.Nil(t, err)
	require.Len(t, pg.objectsContainers, 1)

	require.NotEmpty(t, pg.objectsContainers[object])

	actual, ok := pg.gameMap.Get(engine.Dot{0, 0})
	require.True(t, ok)
	require.Equal(t, pg.objectsContainers[object], actual)
}

func Test_ExperimentalPlayground_CreateObjectRandomRect(t *testing.T) {
	area := engine.MustArea(100, 100)

	object := &struct {
		str string
	}{
		"ok",
	}

	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 100)

	for _, dot := range location {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, pg.objectsContainers[object], actual)
	}
}

func Test_ExperimentalPlayground_UpdateObject(t *testing.T) {
	t.SkipNow()
}

func Benchmark_ExperimentalPlayground_UpdateObject(b *testing.B) {
	const (
		areaWidth  = 150
		areaHeight = 100
	)

	type TestStructure struct {
		a, b, c, d int64
		e, f, g, h float64
	}

	const (
		first = iota
		second
		count
	)

	pg, err := NewExperimentalPlayground(areaWidth, areaHeight)

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

func Test_ExperimentalPlayground_DeleteObject(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct{}{}
	container := engine.NewContainer(object)

	pg.gameMap.MSet(engine.Location{{1, 1}}, container)
	pg.objectsContainers[object] = container

	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))

	err := pg.DeleteObject(object, engine.Location{{1, 1}})
	require.Nil(t, err)

	require.Empty(t, pg.objectsContainers)
	require.Empty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))
}

func Test_ExperimentalPlayground_CreateObjectAvailableDots_EmptySquareScene(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
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

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object], actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectAvailableDots_NotEmptyScene(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

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
	container2 := engine.NewContainer(object2)
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.gameMap.MSet(location2, container2)
	pg.objectsContainers[object2] = container2
	for _, dot := range location2 {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, container2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.Nil(t, err)
	require.True(t, location1.Equals(location1Actual))

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location1.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object1], actualContainer, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, container2, actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectAvailableDots_LocationNotAvailable(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

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
	container2 := engine.NewContainer(object2)
	location2 := location1.Copy()

	pg.gameMap.MSet(location2, container2)
	pg.objectsContainers[object2] = container2
	for _, dot := range location2 {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, container2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.NotNil(t, err)
	require.Nil(t, location1Actual)

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, container2, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectAvailableDots_LocationsIntersects(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	// Object to locate
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
	container2 := engine.NewContainer(object2)
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.gameMap.MSet(location2, container2)
	pg.objectsContainers[object2] = container2
	for _, dot := range location2 {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, container2, actual)
	}

	location1Actual, err := pg.CreateObjectAvailableDots(object1, location1)
	require.Nil(t, err)
	require.True(t, location1Expected.Equals(location1Actual))

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location1Expected.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object1], actualContainer, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, container2, actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_UpdateObjectAvailableDots_EmptyMap(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct{}{}
	container := engine.NewContainer(object)

	pg.gameMap.MSet(engine.Location{{1, 1}}, container)
	pg.objectsContainers[object] = container

	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))

	location, err := pg.UpdateObjectAvailableDots(object, engine.Location{{1, 1}}, engine.Location{{1, 2}})
	require.Nil(t, err)
	require.Equal(t, engine.Location{{1, 2}}, location)
	require.Empty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))
	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 2}}))
}

func Test_ExperimentalPlayground_UpdateObjectAvailableDots_NotEmptyMap(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object1 := &struct {
		data string
	}{"first"}
	container1 := engine.NewContainer(object1)

	pg.gameMap.MSet(engine.Location{{1, 1}}, container1)
	pg.objectsContainers[object1] = container1

	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))

	object2 := &struct {
		data string
	}{"second"}
	container2 := engine.NewContainer(object2)

	pg.gameMap.MSet(engine.Location{{1, 3}}, container2)
	pg.objectsContainers[object2] = container2

	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 3}}))

	location, err := pg.UpdateObjectAvailableDots(object1, engine.Location{{1, 1}}, engine.Location{{1, 2}})
	require.Nil(t, err)
	require.Equal(t, engine.Location{{1, 2}}, location)
	require.Empty(t, pg.gameMap.MGet(engine.Location{{1, 1}}))
	require.NotEmpty(t, pg.gameMap.MGet(engine.Location{{1, 2}}))
}

func Test_ExperimentalPlayground_UpdateObjectAvailableDots_NotEmptyMap_BigObjects(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object1 := &struct {
		data string
	}{"first"}
	container1 := engine.NewContainer(object1)
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

	pg.gameMap.MSet(location1Old, container1)
	pg.objectsContainers[object1] = container1
	for _, dot := range location1Old {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, container1, actual)
	}

	object2 := &struct {
		data string
	}{"second"}
	container2 := engine.NewContainer(object2)
	location2 := engine.Location{
		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 5},
		{2, 6},
	}

	pg.gameMap.MSet(location2, container2)
	pg.objectsContainers[object2] = container2
	for _, dot := range location2 {
		require.True(t, pg.gameMap.Has(dot))
		actual, ok := pg.gameMap.Get(dot)
		require.True(t, ok)
		require.Equal(t, container2, actual)
	}

	location1Actual, err := pg.UpdateObjectAvailableDots(object1, location1Old, location1New)
	require.Nil(t, err)
	require.True(t, location1Expected.Equals(location1Actual))

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location1Expected.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, container1, actualContainer, "dot %s", dot)
		} else if location2.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, container2, actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectRandomDot_SquareEmptyPlayground(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomDot(object)
	require.Nil(t, err)
	require.Len(t, location, 1)

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object], actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectRandomDot_EmptyPlayground(t *testing.T) {
	area := engine.MustArea(200, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomDot(object)
	require.Nil(t, err)
	require.Len(t, location, 1)

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object], actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectRandomRect_SquareEmptyPlayground_SquareRect(t *testing.T) {
	area := engine.MustArea(100, 100)
	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomRect(object, 10, 10)
	require.Nil(t, err)
	require.Len(t, location, 10*10)

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object], actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}

func Test_ExperimentalPlayground_CreateObjectRandomRect_SquareEmptyPlayground(t *testing.T) {
	area := engine.MustArea(100, 100)

	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[interface{}]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	object := &struct {
		data string
	}{"first"}

	location, err := pg.CreateObjectRandomRect(object, 10, 8)
	require.Nil(t, err)
	require.Len(t, location, 10*8)

	for _, dot := range pg.Area().Dots() {
		actualContainer, ok := pg.gameMap.Get(dot)

		if location.Contains(dot) {
			require.True(t, ok, "dot %s", dot)
			require.Equal(t, pg.objectsContainers[object], actualContainer, "dot %s", dot)
		} else {
			require.False(t, ok, "dot %s", dot)
			require.Nil(t, actualContainer, "dot %s", dot)
		}
	}
}
