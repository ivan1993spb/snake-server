package engine

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// storeObject stores object only if pointer p is empty
func storeObject(p *unsafe.Pointer, object *Object) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(uintptr(0)), unsafe.Pointer(object))
}

// emptyObject deletes certain object from pointer p
func emptyObject(p *unsafe.Pointer, object *Object) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(object), unsafe.Pointer(uintptr(0)))
}

// empty clears specified pointer p
func empty(p *unsafe.Pointer) {
	atomic.SwapPointer(p, unsafe.Pointer(uintptr(0)))
}

// fieldIsEmpty returns true if pointer p is empty
func fieldIsEmpty(p unsafe.Pointer) bool {
	return uintptr(p) == uintptr(0)
}

// Map structure represents core map
type Map struct {
	fields [][]*unsafe.Pointer
	area   Area
}

// NewMap creates and returns new empty Map with area a
func NewMap(a Area) *Map {
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

// Print prints a map
func (m *Map) Print() {
	fmt.Println("Map size:", m.area)

	for y := uint8(0); y < m.area.height; y++ {
		fmt.Printf("%4d |", y)
		for x := uint8(0); x < m.area.width; x++ {

			if p := atomic.LoadPointer(m.fields[y][x]); fieldIsEmpty(p) {
				fmt.Print(" .")
			} else {
				fmt.Print(" x")
			}
		}
		fmt.Println()
	}
}

// Has returns true if there is an object under the dot dot otherwise returns false
func (m *Map) Has(dot Dot) bool {
	p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
	return !fieldIsEmpty(p)
}

// Set sets given object under specified dot
func (m *Map) Set(dot Dot, object *Object) {
	if m.area.ContainsDot(dot) {
		atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(object))
	}
}

// Get returns an object by given dot
func (m *Map) Get(dot Dot) (*Object, bool) {
	if !m.area.ContainsDot(dot) {
		return nil, false
	}

	p := atomic.LoadPointer(m.fields[dot.Y][dot.X])

	if fieldIsEmpty(p) {
		return nil, false
	}

	object := (*Object)(p)

	return object, true
}

// SetIfAbsent sets given object under a dot only if the dot is empty.
// Returns true if the object has been set under the dot
func (m *Map) SetIfAbsent(dot Dot, object *Object) bool {
	if !m.area.ContainsDot(dot) {
		return false
	}
	return storeObject(m.fields[dot.Y][dot.X], object)
}

// Remove removes object under the specified dot
func (m *Map) Remove(dot Dot) {
	if m.area.ContainsDot(dot) {
		empty(m.fields[dot.Y][dot.X])
	}
}

// RemoveObject removes certain passed object
func (m *Map) RemoveObject(dot Dot, object *Object) {
	if m.area.ContainsDot(dot) {
		emptyObject(m.fields[dot.Y][dot.X], object)
	}
}

// HasAny returns true if at least one of the dots in the slice has been
// linked with an object.
func (m *Map) HasAny(dots []Dot) bool {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
			if !fieldIsEmpty(p) {
				return true
			}
		}
	}

	return false
}

// HasAll returns true all the dots in the slice has been linked with an
// object.
func (m *Map) HasAll(dots []Dot) bool {
	for _, dot := range dots {
		if !m.area.ContainsDot(dot) {
			return false
		}

		p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
		if p == unsafe.Pointer(uintptr(0)) {
			return false
		}
	}

	return true
}

// MGet returns a map in which keys are dots and values are linked objects
func (m *Map) MGet(dots []Dot) map[Dot]*Object {
	items := make(map[Dot]*Object)

	for _, dot := range dots {
		if !m.area.ContainsDot(dot) {
			continue
		}

		p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
		if !fieldIsEmpty(p) {
			object := (*Object)(p)
			items[dot] = object
		}
	}

	return items
}

// MRemove removes all objects under the specified dots in passed slice
func (m *Map) MRemove(dots []Dot) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			empty(m.fields[dot.Y][dot.X])
		}
	}
}

// MRemoveObject removes certain passed object under the specified dots in slice
func (m *Map) MRemoveObject(dots []Dot, object *Object) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			emptyObject(m.fields[dot.Y][dot.X], object)
		}
	}
}

// MSet sets given object under specified dots
func (m *Map) MSet(dots []Dot, object *Object) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(object))
		}
	}
}

// MSetIfAllAbsent sets given object under specified dots and returns true only
// if all the dots are empty otherwise the function does nothing and returns false
func (m *Map) MSetIfAllAbsent(dots []Dot, object *Object) bool {
	var i int

	for ; i < len(dots); i++ {
		dot := dots[i]
		if !m.area.ContainsDot(dot) {
			continue
		}

		if !storeObject(m.fields[dot.Y][dot.X], object) {
			break
		}
	}

	if i == len(dots) {
		return true
	}

	i -= 1

	// Rollback
	for ; i >= 0; i-- {
		dot := dots[i]
		if !m.area.ContainsDot(dot) {
			continue
		}
		emptyObject(m.fields[dot.Y][dot.X], object)
	}

	return false
}

// MSetIfAbsent sets given object under specified dots. If a dot in the slice is engaged
// the function skips the dot. Returns a list of dot which was eventually linked with the
// passed object
func (m *Map) MSetIfAbsent(dots []Dot, object *Object) []Dot {
	result := make([]Dot, 0, len(dots))

	for _, dot := range dots {
		if !m.area.ContainsDot(dot) {
			continue
		}

		if storeObject(m.fields[dot.Y][dot.X], object) {
			result = append(result, dot)
		}
	}

	return result
}
