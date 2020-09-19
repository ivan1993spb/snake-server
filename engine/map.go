package engine

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// storeContainer stores a container only if the pointer p is empty
func storeContainer(p *unsafe.Pointer, container *Container) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(uintptr(0)), unsafe.Pointer(container))
}

// emptyContainer deletes a certain container from the pointer p
func emptyContainer(p *unsafe.Pointer, container *Container) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(container), unsafe.Pointer(uintptr(0)))
}

// empty cleans a specified pointer p
func empty(p *unsafe.Pointer) {
	atomic.SwapPointer(p, unsafe.Pointer(uintptr(0)))
}

// fieldIsEmpty returns true if the pointer p is empty
func fieldIsEmpty(p unsafe.Pointer) bool {
	return uintptr(p) == uintptr(0)
}

// Map structure represents core map
type Map struct {
	fields [][]*unsafe.Pointer
	area   Area
}

// NewMap creates and returns a new empty Map with area a
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

// Area returns area of a Map
func (m *Map) Area() Area {
	return m.area
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

// Has returns true if there is a container under the dot dot, otherwise returns false
func (m *Map) Has(dot Dot) bool {
	p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
	return !fieldIsEmpty(p)
}

// Set sets given container under specified dot
func (m *Map) Set(dot Dot, container *Container) {
	if m.area.ContainsDot(dot) {
		atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(container))
	}
}

// Get returns a container by the given dot
func (m *Map) Get(dot Dot) (*Container, bool) {
	if !m.area.ContainsDot(dot) {
		return nil, false
	}

	p := atomic.LoadPointer(m.fields[dot.Y][dot.X])

	if fieldIsEmpty(p) {
		return nil, false
	}

	container := (*Container)(p)

	return container, true
}

// SetIfVacant  sets the given container under the dot only if the dot is empty.
// Returns true if the container has been set under the dot.
func (m *Map) SetIfVacant(dot Dot, container *Container) bool {
	if !m.area.ContainsDot(dot) {
		return false
	}
	return storeContainer(m.fields[dot.Y][dot.X], container)
}

// Remove removes a container under the specified dot
func (m *Map) Remove(dot Dot) {
	if m.area.ContainsDot(dot) {
		empty(m.fields[dot.Y][dot.X])
	}
}

// RemoveContainer removes a certain passed container
func (m *Map) RemoveContainer(dot Dot, container *Container) {
	if m.area.ContainsDot(dot) {
		emptyContainer(m.fields[dot.Y][dot.X], container)
	}
}

// HasAny returns true if at least one of the dots in the slice has been
// linked with a container.
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

// HasAll returns true if all dots in passed slice has been linked with a
// container.
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

// MGet returns a map in which keys are dots and values are linked with them containers
func (m *Map) MGet(dots []Dot) map[Dot]*Container {
	items := make(map[Dot]*Container)

	for _, dot := range dots {
		if !m.area.ContainsDot(dot) {
			continue
		}

		p := atomic.LoadPointer(m.fields[dot.Y][dot.X])
		if !fieldIsEmpty(p) {
			container := (*Container)(p)
			items[dot] = container
		}
	}

	return items
}

// MRemove removes all containers under the specified dots in the passed slice
func (m *Map) MRemove(dots []Dot) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			empty(m.fields[dot.Y][dot.X])
		}
	}
}

// MRemoveContainer removes the certain passed container under the specified dots in the slice
func (m *Map) MRemoveContainer(dots []Dot, container *Container) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			emptyContainer(m.fields[dot.Y][dot.X], container)
		}
	}
}

// MSet sets the given container under specified dots
func (m *Map) MSet(dots []Dot, container *Container) {
	for _, dot := range dots {
		if m.area.ContainsDot(dot) {
			atomic.SwapPointer(m.fields[dot.Y][dot.X], unsafe.Pointer(container))
		}
	}
}

// MSetIfAllVacant  sets the given container under specified dots and returns true only
// if all the dots has been set, otherwise the function does nothing and returns false.
func (m *Map) MSetIfAllVacant(dots []Dot, container *Container) bool {
	var i int

	for ; i < len(dots); i++ {
		dot := dots[i]
		if !m.area.ContainsDot(dot) {
			continue
		}

		if !storeContainer(m.fields[dot.Y][dot.X], container) {
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
		emptyContainer(m.fields[dot.Y][dot.X], container)
	}

	return false
}

// MSetIfVacant  sets the given container under specified dots. If a dot in the slice has already
// been linked, the function skips the dot. MSetIfVacant  returns a list of dots which was eventually
// linked with the passed container
func (m *Map) MSetIfVacant(dots []Dot, container *Container) []Dot {
	result := make([]Dot, 0, len(dots))

	for _, dot := range dots {
		if !m.area.ContainsDot(dot) {
			continue
		}

		if storeContainer(m.fields[dot.Y][dot.X], container) {
			result = append(result, dot)
		}
	}

	return result
}
