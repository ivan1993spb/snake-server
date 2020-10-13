package playground

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
)

func Test_NewExperimentalPlayground(t *testing.T) {
	tests := []struct {
		width, height uint8

		err bool
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
		{
			width:  0,
			height: 0,
			err:    true,
		},
		{
			width:  200,
			height: 0,
			err:    true,
		},
	}

	for i, test := range tests {
		pg, err := NewExperimentalPlayground(test.width, test.height)

		if test.err {
			require.NotNil(t, err, "test %d", i)
			require.Nil(t, pg, "test %d", i)
		} else {
			require.Nil(t, err, "test %d", i)
			require.NotNil(t, pg, "test %d", i)
			require.Equal(t, test.width, pg.Area().Width(), "test %d", i)
			require.Equal(t, test.height, pg.Area().Height(), "test %d", i)
		}
	}
}

func Test_ExperimentalPlayground_CreateObject(t *testing.T) {
	const (
		AreaWidth  = 100
		AreaHeight = 100
	)

	area := engine.MustArea(AreaWidth, AreaHeight)

	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[engine.Object]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	// Literally anything
	object := &struct {
		a int
	}{
		10,
	}

	location1 := engine.Location{engine.Dot{0, 0}}
	location2 := engine.Location{engine.Dot{1, 1}}

	// Register the object correctly
	{
		err := pg.CreateObject(object, location1)
		require.Nil(t, err)
		require.Len(t, pg.objectsContainers, 1)

		require.NotEmpty(t, pg.objectsContainers[object])

		actual, ok := pg.gameMap.Get(engine.Dot{0, 0})
		require.True(t, ok)
		require.Equal(t, pg.objectsContainers[object], actual)
	}

	// Try to register the same object second time
	{
		err := pg.CreateObject(object, location2)
		require.NotNil(t, err)
	}
}

func Test_ExperimentalPlayground_GetObjectByDot_WorksCorrectly(t *testing.T) {
	const (
		messageFormat               = "test %d"
		messageObjectLocationFormat = "test %d, object %v, location %s"
		messageObjectDotFormat      = "test %d, object %v, dot %s"
	)

	var (
		object1 = "object1"
		object2 = "object2"
		object3 = "object3"
		object4 = "object4"
		object5 = "object5"
		object6 = "object6"
	)

	tests := []struct {
		areaWidth, areaHeight uint8
		// Map of objects and their locations
		objects map[engine.Object]engine.Location
		// Map of dots and expected objects located at the certain dot
		checks map[engine.Dot]engine.Object
	}{
		// Test 1
		{
			areaWidth:  100,
			areaHeight: 100,
			objects: map[engine.Object]engine.Location{
				object1: {
					{10, 20},
					{11, 24},
					{90, 30},
					{21, 34},
				},
				object2: {
					{33, 20},
					{21, 24},
					{31, 30},
				},
			},
			checks: map[engine.Dot]engine.Object{
				{90, 30}: object1,
				{21, 24}: object2,
				{21, 25}: nil,
				{99, 99}: nil,
			},
		},
		// Test 2
		{
			areaWidth:  150,
			areaHeight: 150,
			objects: map[engine.Object]engine.Location{
				object3: {
					{149, 20},
					{15, 44},
					{10, 43},
				},
				object4: {
					{133, 120},
					{121, 124},
					{131, 130},
				},
			},
			checks: map[engine.Dot]engine.Object{
				{15, 44}:   object3,
				{10, 43}:   object3,
				{11, 43}:   nil,
				{133, 120}: object4,
				{131, 130}: object4,
				{131, 131}: nil,
			},
		},
		// Test 3
		{
			areaWidth:  151,
			areaHeight: 203,
			objects: map[engine.Object]engine.Location{
				object5: {
					{98, 199},
					{76, 52},
					{65, 76},
				},
				object6: {
					{76, 77},
					{32, 87},
					{43, 74},
				},
			},
			checks: map[engine.Dot]engine.Object{
				{0, 0}:     nil,
				{200, 200}: nil,
				{98, 199}:  object5,
				{65, 76}:   object5,
				{43, 74}:   object6,
				{70, 70}:   nil,
			},
		},
	}

	for i, test := range tests {
		number := i + 1

		// Init playground
		area, err := engine.NewArea(test.areaWidth, test.areaHeight)
		require.Nil(t, err, messageFormat, number)

		pg := &ExperimentalPlayground{
			gameMap: engine.NewMap(area),

			objectsContainers:    make(map[engine.Object]*engine.Container),
			objectsContainersMux: &sync.RWMutex{},
		}

		// Add objects manually
		for object, location := range test.objects {
			container := engine.NewContainer(object)

			pg.objectsContainers[object] = container
			ok := pg.gameMap.MSetIfAllVacant(location, container)
			require.True(t, ok, messageObjectLocationFormat, number, object, location)
		}

		// Check objects presence at the certain dots
		for dot, objectExpect := range test.checks {
			actual := pg.GetObjectByDot(dot)
			require.Equal(t, objectExpect, actual, messageObjectDotFormat, number, objectExpect, dot)
		}
	}
}

func Test_ExperimentalPlayground_GetObjectsByDots(t *testing.T) {
	const (
		objectAddMessage  = "problem adding the object manually: %s at %s"
		testNumberMessage = "test number %d"
	)

	const (
		areaWidth  = 210
		areaHeight = 158
	)

	var (
		object1   = "object1"
		location1 = engine.Location{{1, 3}, {1, 2}, {6, 5}, {7, 3}, {9, 2}}

		object2   = "object2"
		location2 = engine.Location{{22, 32}, {32, 43}, {43, 27}, {12, 65}, {12, 66}}

		object3   = "object3"
		location3 = engine.Location{{46, 3}, {52, 1}, {51, 1}, {50, 2}, {52, 2}, {52, 4}}

		object4   = "object4"
		location4 = engine.Location{{2, 32}, {1, 32}, {5, 54}, {1, 66}, {7, 90}}

		object5   = "object5"
		location5 = engine.Location{{102, 54}, {123, 34}, {101, 35}, {115, 34}}

		object6   = "object6"
		location6 = engine.Location{{0, 0}, {201, 0}, {21, 200}}
	)

	var objects = map[engine.Object]engine.Location{
		object1: location1,
		object2: location2,
		object3: location3,
		object4: location4,
		object5: location5,
		object6: location6,
	}

	area, err := engine.NewArea(areaWidth, areaHeight)
	require.Nil(t, err)

	pg := &ExperimentalPlayground{
		gameMap: engine.NewMap(area),

		objectsContainers:    make(map[engine.Object]*engine.Container),
		objectsContainersMux: &sync.RWMutex{},
	}

	// Add objects manually
	for object, location := range objects {
		container := engine.NewContainer(object)

		pg.objectsContainers[object] = container
		ok := pg.gameMap.MSetIfAllVacant(location, container)
		require.True(t, ok, objectAddMessage, object, location)
	}

	tests := []struct {
		dots            []engine.Dot
		expectedObjects []engine.Object
	}{
		// Test 1
		{
			dots:            nil,
			expectedObjects: nil,
		},
		// Test 2
		{
			dots:            nil,
			expectedObjects: []engine.Object{},
		},
		// Test 3
		{
			dots:            []engine.Dot{{1, 1}},
			expectedObjects: nil,
		},
		// Test 4
		{
			dots:            []engine.Dot{{1, 1}},
			expectedObjects: []engine.Object{},
		},
		// Test 5
		{
			dots:            []engine.Dot{{0, 0}},
			expectedObjects: []engine.Object{object6},
		},
		// Test 6
		{
			dots:            []engine.Dot{{0, 0}, {201, 0}},
			expectedObjects: []engine.Object{object6},
		},
		// Test 7
		{
			dots:            []engine.Dot{{1, 3}, {1, 2}, {0, 0}, {201, 0}, {102, 54}},
			expectedObjects: []engine.Object{object1, object5, object6},
		},
		// Test 8
		{
			dots:            []engine.Dot{{1, 3}, {1, 2}, {6, 5}, {7, 3}, {9, 2}},
			expectedObjects: []engine.Object{object1},
		},
		// Test 9
		{
			dots:            []engine.Dot{{22, 3}, {22, 2}, {22, 5}, {2, 3}, {22, 5}},
			expectedObjects: []engine.Object{},
		},
		// Test 10
		{
			dots: []engine.Dot{{1, 3}, {1, 2}, {22, 32}, {32, 43},
				{46, 3}, {52, 1}, {2, 32}, {1, 32}, {102, 54},
				{123, 34}, {0, 0}, {201, 0}},
			expectedObjects: []engine.Object{object1, object2, object3, object4, object5, object6},
		},
	}

	for i, test := range tests {
		number := i + 1

		actualObjects := pg.GetObjectsByDots(test.dots)
		require.Subset(t, test.expectedObjects, actualObjects, testNumberMessage, number)
		require.Subset(t, actualObjects, test.expectedObjects, testNumberMessage, number)
	}
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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
	require.Nil(t, err)
	require.Empty(t, location1Actual)

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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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

		objectsContainers:    make(map[engine.Object]*engine.Container),
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
