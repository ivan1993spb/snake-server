package engine

import (
	"strconv"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type SampleObject struct {
	a, b, c, d, e int
	f, g          float64
}

func getSampleContainerFirst() *Container {
	return &Container{
		object: &SampleObject{
			a: 4,
			b: 43,
			c: 1,
			d: 4,
			e: 54,
			f: 5.212,
			g: 12341.21,
		},
	}
}

func getSampleContainerSecond() *Container {
	return &Container{
		object: &SampleObject{
			a: 2,
			b: 1,
			c: 4,
			d: 4,
			e: 0x32,
			f: 545.55,
			g: 1123.4,
		},
	}
}

func getSampleContainerThird() *Container {
	return &Container{
		object: &SampleObject{
			a: 3,
			b: 123,
			c: 4,
			d: 44,
			e: 0xffffff,
			f: 6.55,
			g: 36.4,
		},
	}
}

func getSampleMapArea(a Area) *Map {
	field := make(map[Dot]*unsafe.Pointer)

	for _, dot := range a.Dots() {
		var emptyFieldPointer = unsafe.Pointer(uintptr(0))
		field[dot] = &emptyFieldPointer
	}

	return &Map{
		field: field,
		area:  a,
	}
}

func Test_storeContainer_storesContainer(t *testing.T) {
	container := getSampleContainerFirst()

	var emptyFieldPointer = unsafe.Pointer(uintptr(0))
	storePointer := &emptyFieldPointer

	require.True(t, storeContainer(storePointer, container))
	require.Equal(t, uintptr(unsafe.Pointer(container)), uintptr(*storePointer))
}

func Test_emptyContainer_emptiesPointer(t *testing.T) {
	container := getSampleContainerFirst()

	storePointer := unsafe.Pointer(container)

	require.True(t, emptyContainer(&storePointer, container))
	require.Equal(t, uintptr(0), uintptr(storePointer))
}

func Test_storeContainer_returnsFalseIfPointerIsEngaged(t *testing.T) {
	// First
	containerFirst := getSampleContainerFirst()

	// Second
	containerSecond := getSampleContainerSecond()

	// Store cell
	storePointer := unsafe.Pointer(containerFirst)

	require.False(t, storeContainer(&storePointer, containerSecond))
	require.NotEqual(t, uintptr(0), uintptr(storePointer))
	require.Equal(t, uintptr(unsafe.Pointer(containerFirst)), uintptr(storePointer))
}

func Test_emptyContainer_returnsFalseIfMismatching(t *testing.T) {
	// First
	containerFirst := getSampleContainerFirst()

	// Second
	containerSecond := getSampleContainerSecond()

	// Store cell
	storePointer := unsafe.Pointer(containerFirst)

	require.False(t, emptyContainer(&storePointer, containerSecond))
	require.NotEqual(t, uintptr(0), uintptr(storePointer))
	require.Equal(t, uintptr(unsafe.Pointer(containerFirst)), uintptr(storePointer))
}

func Test_fieldIsEmpty(t *testing.T) {
	// First
	containerFirst := getSampleContainerFirst()

	// Second
	containerSecond := getSampleContainerSecond()

	tests := []struct {
		pointer  unsafe.Pointer
		expected bool
	}{
		{
			pointer:  unsafe.Pointer(uintptr(0)),
			expected: true,
		},
		{
			pointer:  unsafe.Pointer(containerFirst),
			expected: false,
		},
		{
			pointer:  unsafe.Pointer(uintptr(0)),
			expected: true,
		},
		{
			pointer:  unsafe.Pointer(containerSecond),
			expected: false,
		},
		{
			pointer:  unsafe.Pointer(uintptr(0)),
			expected: true,
		},
	}

	for i, test := range tests {
		require.Equal(t, test.expected, fieldIsEmpty(test.pointer), "number "+strconv.Itoa(i))
	}
}

func Test_NewMap_CreatesEmptyMap(t *testing.T) {
	area := MustArea(53, 20)
	m := NewMap(area)

	require.Equal(t, area, m.area)

	for y := uint8(0); y < area.height; y++ {
		for x := uint8(0); x < area.width; x++ {
			p := atomic.LoadPointer(m.field[Dot{x, y}])
			require.True(t, uintptr(0) == uintptr(p))
		}
	}
}

func Test_Map_Area_ReturnsArea(t *testing.T) {
	area := MustArea(53, 20)
	m := Map{
		area: area,
	}
	require.Equal(t, area, m.Area())
}

func Test_Map_Print_Prints(t *testing.T) {
	area := MustArea(23, 15)
	m := getSampleMapArea(area)

	// First
	pointerToContainerFirst := unsafe.Pointer(getSampleContainerFirst())
	dotFirst := Dot{3, 4}

	// Second
	pointerToContainerSecond := unsafe.Pointer(getSampleContainerSecond())
	dotSecond := Dot{22, 4}

	// Second
	pointerToContainerThird := unsafe.Pointer(getSampleContainerThird())
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	m.Print()
}

func Test_Map_Has_ReturnsValidIndicator(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	// First
	pointerToContainerFirst := unsafe.Pointer(getSampleContainerFirst())
	dotFirst := Dot{3, 4}

	// Second
	pointerToContainerSecond := unsafe.Pointer(getSampleContainerSecond())
	dotSecond := Dot{22, 4}

	// Second
	pointerToContainerThird := unsafe.Pointer(getSampleContainerThird())
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	for _, dot := range area.Dots() {
		if dot.Equals(dotFirst) || dot.Equals(dotSecond) || dot.Equals(dotThird) {
			require.True(t, m.Has(dot))
		} else {
			require.False(t, m.Has(dot))
		}
	}
}

func Test_Map_Set(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	containerFirst := getSampleContainerFirst()
	dotFirst := Dot{1, 3}
	m.Set(dotFirst, containerFirst)

	for _, dot := range area.Dots() {
		pointer := *m.field[dot]

		if dot.Equals(dotFirst) {
			require.Equal(t, uintptr(unsafe.Pointer(containerFirst)), uintptr(pointer))
		} else {
			require.Equal(t, uintptr(0), uintptr(pointer))
		}
	}
}

func Test_Map_Get_ReturnsValidContainer(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{3, 4}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	for _, dot := range area.Dots() {
		result, ok := m.Get(dot)

		switch {
		case dot.Equals(dotFirst):
			require.True(t, ok)
			require.Equal(t, containerFirst, result)
		case dot.Equals(dotSecond):
			require.True(t, ok)
			require.Equal(t, containerSecond, result)
		case dot.Equals(dotThird):
			require.True(t, ok)
			require.Equal(t, containerThird, result)
		default:
			require.False(t, ok)
			require.Nil(t, result)
		}
	}
}

func Test_Map_Get_ReturnsNilFalse(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	dot1 := Dot{122, 33}
	dot2 := Dot{0, 32}
	dot3 := Dot{12, 211}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		result, ok := m.Get(dot)
		require.False(t, ok, "number "+strconv.Itoa(i))
		require.Nil(t, result, "number "+strconv.Itoa(i))
	}
}

func Test_Map_SetIfVacant_OnEmptyMap(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	containerFirst := getSampleContainerFirst()

	dot1 := Dot{1, 21}
	dot2 := Dot{0, 12}
	dot3 := Dot{12, 11}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfVacant(dot, containerFirst)
		require.True(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_SetIfVacant_InvalidDots(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	containerFirst := getSampleContainerFirst()

	dot1 := Dot{1, 211}
	dot2 := Dot{0, 121}
	dot3 := Dot{12, 111}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfVacant(dot, containerFirst)
		require.False(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_SetIfVacant_OccupiedDots(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	containerFirst := getSampleContainerFirst()
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{11, 5}

	m.field[dot1] = &pointerToContainerSecond
	m.field[dot2] = &pointerToContainerSecond
	m.field[dot3] = &pointerToContainerSecond

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfVacant(dot, containerFirst)
		require.False(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_Remove_RemovesContainers(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	m.Remove(dotFirst)
	m.Remove(dotSecond)
	m.Remove(dotThird)

	for _, dot := range area.Dots() {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}
}

func Test_Map_Remove_DoesNothingInCaseOfInvalidDots(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	m.Remove(Dot{201, 0})
	m.Remove(Dot{202, 105})
	m.Remove(Dot{203, 0})

	for _, dot := range area.Dots() {
		pointer := *m.field[dot]

		switch {
		case dot.Equals(dotFirst):
			require.Equal(t, uintptr(pointerToContainerFirst), uintptr(pointer))
		case dot.Equals(dotSecond):
			require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer))
		case dot.Equals(dotThird):
			require.Equal(t, uintptr(pointerToContainerThird), uintptr(pointer))
		default:
			require.Equal(t, uintptr(0), uintptr(pointer))
		}
	}
}

func Test_Map_RemoveContainer_IgnoresInvalidDot(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.field[dot1] = &pointerToContainerFirst
	m.field[dot2] = &pointerToContainerFirst
	m.field[dot3] = &pointerToContainerFirst

	m.RemoveContainer(Dot{124, 212}, containerFirst)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(pointerToContainerFirst), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_RemoveContainer_RemovesContainer(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.field[dot1] = &pointerToContainerFirst
	m.field[dot2] = &pointerToContainerFirst
	m.field[dot3] = &pointerToContainerFirst

	m.RemoveContainer(dot1, containerFirst)
	m.RemoveContainer(dot2, containerFirst)
	m.RemoveContainer(dot3, containerFirst)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(0), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_RemoveContainer_DoesNotRemoveMismatchedContainer(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)

	containerSecond := getSampleContainerSecond()

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.field[dot1] = &pointerToContainerFirst
	m.field[dot2] = &pointerToContainerFirst
	m.field[dot3] = &pointerToContainerFirst

	m.RemoveContainer(dot1, containerSecond)
	m.RemoveContainer(dot2, containerSecond)
	m.RemoveContainer(dot3, containerSecond)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(pointerToContainerFirst), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_HasAny_ReturnsCorrectResult(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.field[dot1] = &pointerToContainerFirst
	m.field[dot2] = &pointerToContainerFirst
	m.field[dot3] = &pointerToContainerFirst

	location1 := []Dot{
		// Dot to be skipped
		{
			X: 150,
			Y: 220,
		},
		//
		{
			X: 1,
			Y: 2,
		},
		{
			X: 1,
			Y: 3,
		},
		{
			X: 1,
			Y: 4,
		},
		dot2,
	}

	location2 := []Dot{
		// Dot to be skipped
		{
			X: 123,
			Y: 1,
		},
		//
		{
			X: 1,
			Y: 3,
		},
		{
			X: 1,
			Y: 4,
		},
		dot3,
	}

	location3 := []Dot{
		// Dot to be skipped
		{
			X: 129,
			Y: 1,
		},
		//
		{
			X: 1,
			Y: 3,
		},
		{
			X: 1,
			Y: 4,
		},
	}

	tests := []struct {
		location []Dot
		expected bool
	}{
		{
			location: location1,
			expected: true,
		},
		{
			location: location2,
			expected: true,
		},
		{
			location: location3,
			expected: false,
		},
	}

	for i, test := range tests {
		result := m.HasAny(test.location)
		require.Equal(t, test.expected, result, "number "+strconv.Itoa(i))
	}
}

func Test_Map_HasAll_returnsCorrectResult(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	tests := []struct {
		dots     []Dot
		expected bool
	}{
		{
			dots:     []Dot{},
			expected: true,
		},
		{
			dots:     []Dot{dotFirst, dotFirst},
			expected: true,
		},
		{
			dots:     []Dot{dotFirst, dotSecond, dotThird},
			expected: true,
		},
		{
			dots:     []Dot{dotFirst, dotSecond, dotThird, {1, 1}},
			expected: false,
		},
		{
			dots:     []Dot{dotFirst, {201, 202}},
			expected: false,
		},
	}

	for i, test := range tests {
		result := m.HasAll(test.dots)
		require.Equal(t, test.expected, result, "test number "+strconv.Itoa(i))
	}
}

func Test_Map_MGet_returnsValidDotContainerMap(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	dots := []Dot{
		{200, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	result := m.MGet(dots)

	require.Len(t, result, 3)

	// Check the containers

	{
		resultContainerFirst, ok := result[dotFirst]
		require.True(t, ok)
		require.Equal(t, containerFirst, resultContainerFirst)
	}

	{
		resultContainerSecond, ok := result[dotSecond]
		require.True(t, ok)
		require.Equal(t, containerSecond, resultContainerSecond)
	}

	{
		resultContainerThird, ok := result[dotThird]
		require.True(t, ok)
		require.Equal(t, containerThird, resultContainerThird)
	}

	{
		result, ok := result[Dot{20, 4}]
		require.False(t, ok)
		require.Nil(t, result)
	}
}

func Test_Map_MRemove_removes(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	m.MRemove(dots)

	for _, dot := range area.Dots() {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}
}

func Test_Map_MRemoveContainer_removesCertainContainers(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	m.MRemoveContainer(dots, containerSecond)

	{
		pointer := *m.field[dotFirst]
		require.Equal(t, uintptr(pointerToContainerFirst), uintptr(pointer))
	}

	{
		pointer := *m.field[dotSecond]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}

	{
		pointer := *m.field[dotThird]
		require.Equal(t, uintptr(pointerToContainerThird), uintptr(pointer))
	}
}

func Test_Map_MSet(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	m.MSet(dots, containerSecond)

	for i, dot := range dots {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]
			require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}
}

func Test_Map_MSetIfAllVacant_SetsContainerOnMapCorrectly(t *testing.T) {
	area := MustArea(56, 45)
	m := getSampleMapArea(area)
	containerFirst := getSampleContainerFirst()
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{11, 5}

	m.field[dot1] = &pointerToContainerSecond
	m.field[dot2] = &pointerToContainerSecond
	m.field[dot3] = &pointerToContainerSecond

	location := []Dot{
		// Dot to be skipped
		{
			X: 100,
			Y: 200,
		},
		//
		{
			X: 1,
			Y: 2,
		},
		{
			X: 1,
			Y: 3,
		},
		{
			X: 1,
			Y: 4,
		},
		{
			X: 1,
			Y: 5,
		},
	}

	ok := m.MSetIfAllVacant(location, containerFirst)

	require.True(t, ok)

	for i, dot := range location {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]
			require.NotEqual(t, uintptr(0), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.field[dot]
		require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
	}
}

func Test_Map_MSetIfAllVacant_RollbacksChangesAndReturnsFalse(t *testing.T) {
	area := MustArea(56, 45)
	m := getSampleMapArea(area)
	containerFirst := getSampleContainerFirst()
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.field[dot1] = &pointerToContainerSecond
	m.field[dot2] = &pointerToContainerSecond
	m.field[dot3] = &pointerToContainerSecond

	location := []Dot{
		// Dot to be skipped
		{
			X: 100,
			Y: 200,
		},
		//
		{
			X: 1,
			Y: 2,
		},
		{
			X: 1,
			Y: 3,
		},
		{
			X: 1,
			Y: 4,
		},
		dot3,
	}

	ok := m.MSetIfAllVacant(location, containerFirst)

	require.False(t, ok)

	for i, dot := range location {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]

			if dot.Equals(dot3) {
				require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
			} else {
				require.Equal(t, uintptr(0), uintptr(pointer), "number "+strconv.Itoa(i))
			}
		}
	}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]
			require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}
}

func Test_Map_MSetIfVacant(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	containerFirst := getSampleContainerFirst()
	pointerToContainerFirst := unsafe.Pointer(containerFirst)
	dotFirst := Dot{120, 66}

	// Second
	containerSecond := getSampleContainerSecond()
	pointerToContainerSecond := unsafe.Pointer(containerSecond)
	dotSecond := Dot{22, 4}

	// Second
	containerThird := getSampleContainerThird()
	pointerToContainerThird := unsafe.Pointer(containerThird)
	dotThird := Dot{12, 0}

	m.field[dotFirst] = &pointerToContainerFirst
	m.field[dotSecond] = &pointerToContainerSecond
	m.field[dotThird] = &pointerToContainerThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	resultDots := m.MSetIfVacant(dots, containerSecond)

	require.Len(t, resultDots, 3)

	t.Log("Check result dot set")
	for i, dot := range resultDots {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]
			require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	t.Log("Check if containers has been linked correclty on the map")
	for i, dot := range []Dot{{11, 66}, {4, 2}, {20, 4}, dotSecond} {
		if area.ContainsDot(dot) {
			pointer := *m.field[dot]
			require.Equal(t, uintptr(pointerToContainerSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	t.Log("Check other containers")
	{
		pointer := *m.field[dotFirst]
		require.Equal(t, uintptr(pointerToContainerFirst), uintptr(pointer))
	}

	{
		pointer := *m.field[dotThird]
		require.Equal(t, uintptr(pointerToContainerThird), uintptr(pointer))
	}
}
