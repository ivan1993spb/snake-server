package engine

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

func storeObject(p *unsafe.Pointer, object *Object) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(uintptr(0)), unsafe.Pointer(object))
}

func emptyObject(p *unsafe.Pointer, object *Object) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(object), unsafe.Pointer(uintptr(0)))
}

func fieldIsEmpty(p unsafe.Pointer) bool {
	return uintptr(p) == uintptr(0)
}

type Map struct {
	fields [][]*unsafe.Pointer
	area   Area
}

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

func (m *Map) Has(dot Dot) bool {
	p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
	return !fieldIsEmpty(p)
}

func (m *Map) Set(dot Dot, object *Object) {
	if m.area.ContainsDot(dot) {
		atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(object))
	}
}

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

func (m *Map) SetIfAbsent(dot Dot, object *Object) bool {
	if !m.area.ContainsDot(dot) {
		return false
	}
	return storeObject(m.fields[dot.Y][dot.X], object)
}

func (m *Map) Remove(dot Dot) {
	if m.area.ContainsDot(dot) {
		atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(uintptr(0)))
	}
}

func (m *Map) RemoveObject(dot Dot, object *Object) {
	if m.area.ContainsDot(dot) {
		emptyObject(m.fields[dot.Y][dot.X], object)
	}
}

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

func (m *Map) MRemove(dots []Dot) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(uintptr(0)))
		}
	}
}

func (m *Map) MRemoveObject(dots []Dot, object *Object) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			emptyObject(m.fields[dot.Y][dot.X], object)
		}
	}
}

func (m *Map) MSet(dots []Dot, object *Object) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(object))
		}
	}
}

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
