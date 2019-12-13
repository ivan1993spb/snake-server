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

func getSampleWrappedObjectFirst() *Object {
	return &Object{
		value: &SampleObject{
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

func getSampleWrappedObjectSecond() *Object {
	return &Object{
		value: &SampleObject{
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

func getSampleWrappedObjectThird() *Object {
	return &Object{
		value: &SampleObject{
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
	m := make([][]*unsafe.Pointer, a.height)

	for y := uint8(0); y < a.height; y++ {
		m[y] = make([]*unsafe.Pointer, a.width)

		for x := uint8(0); x < a.width; x++ {
			var emptyFieldPointer = unsafe.Pointer(uintptr(0))
			m[y][x] = &emptyFieldPointer
		}
	}

	return &Map{
		fields: m,
		area:   a,
	}
}

func Test_storeObject_storesObject(t *testing.T) {
	wrappedObject := getSampleWrappedObjectFirst()

	var emptyFieldPointer = unsafe.Pointer(uintptr(0))
	storePointer := &emptyFieldPointer

	require.True(t, storeObject(storePointer, wrappedObject))
	require.Equal(t, uintptr(unsafe.Pointer(wrappedObject)), uintptr(*storePointer))
}

func Test_emptyObject_emptiesPointer(t *testing.T) {
	wrappedObject := getSampleWrappedObjectFirst()

	storePointer := unsafe.Pointer(wrappedObject)

	require.True(t, emptyObject(&storePointer, wrappedObject))
	require.Equal(t, uintptr(0), uintptr(storePointer))
}

func Test_storeObject_returnsFalseIfPointerIsEngaged(t *testing.T) {
	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()

	// Store cell
	storePointer := unsafe.Pointer(wrappedObjectFirst)

	require.False(t, storeObject(&storePointer, wrappedObjectSecond))
	require.NotEqual(t, uintptr(0), uintptr(storePointer))
	require.Equal(t, uintptr(unsafe.Pointer(wrappedObjectFirst)), uintptr(storePointer))
}

func Test_emptyObject_returnsFalseIfMismatching(t *testing.T) {
	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()

	// Store cell
	storePointer := unsafe.Pointer(wrappedObjectFirst)

	require.False(t, emptyObject(&storePointer, wrappedObjectSecond))
	require.NotEqual(t, uintptr(0), uintptr(storePointer))
	require.Equal(t, uintptr(unsafe.Pointer(wrappedObjectFirst)), uintptr(storePointer))
}

func Test_fieldIsEmpty(t *testing.T) {
	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()

	tests := []struct {
		pointer  unsafe.Pointer
		expected bool
	}{
		{
			pointer:  unsafe.Pointer(uintptr(0)),
			expected: true,
		},
		{
			pointer:  unsafe.Pointer(wrappedObjectFirst),
			expected: false,
		},
		{
			pointer:  unsafe.Pointer(uintptr(0)),
			expected: true,
		},
		{
			pointer:  unsafe.Pointer(wrappedObjectSecond),
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
			p := atomic.LoadPointer(m.fields[y][x])
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
	pointerToWrappedObjectFirst := unsafe.Pointer(getSampleWrappedObjectFirst())
	dotFirst := Dot{3, 4}

	// Second
	pointerToWrappedObjectSecond := unsafe.Pointer(getSampleWrappedObjectSecond())
	dotSecond := Dot{22, 4}

	// Second
	pointerToWrappedObjectThird := unsafe.Pointer(getSampleWrappedObjectThird())
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	m.Print()
}

func Test_Map_Has_ReturnsValidIndicator(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	// First
	pointerToWrappedObjectFirst := unsafe.Pointer(getSampleWrappedObjectFirst())
	dotFirst := Dot{3, 4}

	// Second
	pointerToWrappedObjectSecond := unsafe.Pointer(getSampleWrappedObjectSecond())
	dotSecond := Dot{22, 4}

	// Second
	pointerToWrappedObjectThird := unsafe.Pointer(getSampleWrappedObjectThird())
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

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

	wrappedObjectFirst := getSampleWrappedObjectFirst()
	dotFirst := Dot{1, 3}
	m.Set(dotFirst, wrappedObjectFirst)

	for _, dot := range area.Dots() {
		pointer := *m.fields[dot.Y][dot.X]

		if dot.Equals(dotFirst) {
			require.Equal(t, uintptr(unsafe.Pointer(wrappedObjectFirst)), uintptr(pointer))
		} else {
			require.Equal(t, uintptr(0), uintptr(pointer))
		}
	}
}

func Test_Map_Get_ReturnsValidObject(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{3, 4}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	for _, dot := range area.Dots() {
		result, ok := m.Get(dot)

		switch {
		case dot.Equals(dotFirst):
			require.True(t, ok)
			require.Equal(t, wrappedObjectFirst, result)
		case dot.Equals(dotSecond):
			require.True(t, ok)
			require.Equal(t, wrappedObjectSecond, result)
		case dot.Equals(dotThird):
			require.True(t, ok)
			require.Equal(t, wrappedObjectThird, result)
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

func Test_Map_SetIfAbsent_OnEmptyMap(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	wrappedObjectFirst := getSampleWrappedObjectFirst()

	dot1 := Dot{1, 21}
	dot2 := Dot{0, 12}
	dot3 := Dot{12, 11}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfAbsent(dot, wrappedObjectFirst)
		require.True(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_SetIfAbsent_InvalidDots(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	wrappedObjectFirst := getSampleWrappedObjectFirst()

	dot1 := Dot{1, 211}
	dot2 := Dot{0, 121}
	dot3 := Dot{12, 111}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfAbsent(dot, wrappedObjectFirst)
		require.False(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_SetIfAbsent_OccupiedDots(t *testing.T) {
	area := MustArea(23, 31)
	m := getSampleMapArea(area)
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{11, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectSecond
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectSecond
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectSecond

	for i, dot := range []Dot{dot1, dot2, dot3} {
		ok := m.SetIfAbsent(dot, wrappedObjectFirst)
		require.False(t, ok, "number "+strconv.Itoa(i))
	}
}

func Test_Map_Remove_RemovesObjects(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	m.Remove(dotFirst)
	m.Remove(dotSecond)
	m.Remove(dotThird)

	for _, dot := range area.Dots() {
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}
}

func Test_Map_Remove_DoesNothingInCaseOfInvalidDots(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	m.Remove(Dot{201, 0})
	m.Remove(Dot{202, 105})
	m.Remove(Dot{203, 0})

	for _, dot := range area.Dots() {
		pointer := *m.fields[dot.Y][dot.X]

		switch {
		case dot.Equals(dotFirst):
			require.Equal(t, uintptr(pointerToWrappedObjectFirst), uintptr(pointer))
		case dot.Equals(dotSecond):
			require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer))
		case dot.Equals(dotThird):
			require.Equal(t, uintptr(pointerToWrappedObjectThird), uintptr(pointer))
		default:
			require.Equal(t, uintptr(0), uintptr(pointer))
		}
	}
}

func Test_Map_RemoveObject_IgnoresInvalidDot(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectFirst
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectFirst
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectFirst

	m.RemoveObject(Dot{124, 212}, wrappedObjectFirst)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(pointerToWrappedObjectFirst), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_RemoveObject_RemovesObject(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectFirst
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectFirst
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectFirst

	m.RemoveObject(dot1, wrappedObjectFirst)
	m.RemoveObject(dot2, wrappedObjectFirst)
	m.RemoveObject(dot3, wrappedObjectFirst)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(0), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_RemoveObject_DoesNotRemoveMismatchedObject(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)

	wrappedObjectSecond := getSampleWrappedObjectSecond()

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectFirst
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectFirst
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectFirst

	m.RemoveObject(dot1, wrappedObjectSecond)
	m.RemoveObject(dot2, wrappedObjectSecond)
	m.RemoveObject(dot3, wrappedObjectSecond)

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(pointerToWrappedObjectFirst), uintptr(pointer), "dot number "+strconv.Itoa(i))
	}
}

func Test_Map_HasAny_ReturnsCorrectResult(t *testing.T) {
	area := MustArea(123, 211)
	m := getSampleMapArea(area)

	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectFirst
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectFirst
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectFirst

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
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

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

func Test_Map_MGet_returnsValidDotObjectMap(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

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

	// Check the objects

	{
		resultObjectFirst, ok := result[dotFirst]
		require.True(t, ok)
		require.Equal(t, wrappedObjectFirst, resultObjectFirst)
	}

	{
		resultObjectSecond, ok := result[dotSecond]
		require.True(t, ok)
		require.Equal(t, wrappedObjectSecond, resultObjectSecond)
	}

	{
		resultObjectThird, ok := result[dotThird]
		require.True(t, ok)
		require.Equal(t, wrappedObjectThird, resultObjectThird)
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
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

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
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}
}

func Test_Map_MRemoveObject_removesCertainObjects(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	m.MRemoveObject(dots, wrappedObjectSecond)

	{
		pointer := *m.fields[dotFirst.Y][dotFirst.X]
		require.Equal(t, uintptr(pointerToWrappedObjectFirst), uintptr(pointer))
	}

	{
		pointer := *m.fields[dotSecond.Y][dotSecond.X]
		require.Equal(t, uintptr(0), uintptr(pointer))
	}

	{
		pointer := *m.fields[dotThird.Y][dotThird.X]
		require.Equal(t, uintptr(pointerToWrappedObjectThird), uintptr(pointer))
	}
}

func Test_Map_MSet(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	m.MSet(dots, wrappedObjectSecond)

	for i, dot := range dots {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]
			require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}
}

func Test_Map_MSetIfAllAbsent_SetsObjectOnMapCorrectly(t *testing.T) {
	area := MustArea(56, 45)
	m := getSampleMapArea(area)
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{11, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectSecond
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectSecond
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectSecond

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

	ok := m.MSetIfAllAbsent(location, wrappedObjectFirst)

	require.True(t, ok)

	for i, dot := range location {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]
			require.NotEqual(t, uintptr(0), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		pointer := *m.fields[dot.Y][dot.X]
		require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
	}
}

func Test_Map_MSetIfAllAbsent_RollbacksChangesAndReturnsFalse(t *testing.T) {
	area := MustArea(56, 45)
	m := getSampleMapArea(area)
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)

	dot1 := Dot{6, 12}
	dot2 := Dot{2, 30}
	dot3 := Dot{1, 5}

	m.fields[dot1.Y][dot1.X] = &pointerToWrappedObjectSecond
	m.fields[dot2.Y][dot2.X] = &pointerToWrappedObjectSecond
	m.fields[dot3.Y][dot3.X] = &pointerToWrappedObjectSecond

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

	ok := m.MSetIfAllAbsent(location, wrappedObjectFirst)

	require.False(t, ok)

	for i, dot := range location {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]

			if dot.Equals(dot3) {
				require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
			} else {
				require.Equal(t, uintptr(0), uintptr(pointer), "number "+strconv.Itoa(i))
			}
		}
	}

	for i, dot := range []Dot{dot1, dot2, dot3} {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]
			require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}
}

func Test_Map_MSetIfAbsent(t *testing.T) {
	area := MustArea(200, 100)
	m := getSampleMapArea(area)

	// First
	wrappedObjectFirst := getSampleWrappedObjectFirst()
	pointerToWrappedObjectFirst := unsafe.Pointer(wrappedObjectFirst)
	dotFirst := Dot{120, 66}

	// Second
	wrappedObjectSecond := getSampleWrappedObjectSecond()
	pointerToWrappedObjectSecond := unsafe.Pointer(wrappedObjectSecond)
	dotSecond := Dot{22, 4}

	// Second
	wrappedObjectThird := getSampleWrappedObjectThird()
	pointerToWrappedObjectThird := unsafe.Pointer(wrappedObjectThird)
	dotThird := Dot{12, 0}

	m.fields[dotFirst.Y][dotFirst.X] = &pointerToWrappedObjectFirst
	m.fields[dotSecond.Y][dotSecond.X] = &pointerToWrappedObjectSecond
	m.fields[dotThird.Y][dotThird.X] = &pointerToWrappedObjectThird

	dots := []Dot{
		{201, 66}, // Dot to be skipped

		dotFirst,
		{11, 66},
		dotSecond,
		{4, 2},
		dotThird,
		{20, 4},
	}

	resultDots := m.MSetIfAbsent(dots, wrappedObjectSecond)

	require.Len(t, resultDots, 3)

	t.Log("Check result dot set")
	for i, dot := range resultDots {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]
			require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	t.Log("Check if objects has been linked correclty on the map")
	for i, dot := range []Dot{{11, 66}, {4, 2}, {20, 4}, dotSecond} {
		if area.ContainsDot(dot) {
			pointer := *m.fields[dot.Y][dot.X]
			require.Equal(t, uintptr(pointerToWrappedObjectSecond), uintptr(pointer), "number "+strconv.Itoa(i))
		}
	}

	t.Log("Check another objects")
	{
		pointer := *m.fields[dotFirst.Y][dotFirst.X]
		require.Equal(t, uintptr(pointerToWrappedObjectFirst), uintptr(pointer))
	}

	{
		pointer := *m.fields[dotThird.Y][dotThird.X]
		require.Equal(t, uintptr(pointerToWrappedObjectThird), uintptr(pointer))
	}
}
