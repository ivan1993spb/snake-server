package engine

// Location is set of dots
type Location []*Dot

// contains returns true if object contains passed dot
func (l Location) Contains(dot *Dot) bool {
	if len(l) > 0 {
		for _, d := range l {
			if d.Equals(dot) {
				return true
			}
		}
	}

	return false
}

// Delete deletes dot from object
func (l Location) Delete(dot *Dot) Location {
	if len(l) > 0 {
		for i := range l {
			if l[i].Equals(dot) {
				return append(l[:i], l[i+1:]...)
			}
		}
	}
	return Location{}
}

func (l Location) Add(dot *Dot) Location {
	return append(l, dot)
}

// Reverse reverses dot sequence in object
func (l Location) Reverse() Location {
	if len(l) > 0 {
		ro := make(Location, 0, len(l))
		for i := len(l) - 1; i >= 0; i-- {
			ro = append(ro, l[i])
		}

		return ro
	}

	return Location{}
}

func (l Location) Dot(i uint16) *Dot {
	return l[i]
}

func (l Location) DotCount() uint16 {
	return uint16(len(l))
}

func (l Location) Copy() Location {
	newLocation := make(Location, 0, len(l))
	copy(newLocation, l)
	return newLocation
}

func (l1 Location) Equals(l2 Location) bool {
	if len(l1) == 0 && len(l2) == 0 {
		return true
	}

	if len(l1) != len(l2) {
		return false
	}

	if len(l1.Difference(l2)) > 0 {
		return false
	}

	return true
}

func (l1 Location) Difference(l2 Location) Location {
	var diff Location

	for i := 0; i < 2; i++ {
		for _, dot1 := range l1 {
			found := false
			for _, dot2 := range l2 {
				if dot1.Equals(dot2) {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, dot1)
			}
		}
		if i == 0 {
			l1, l2 = l2, l1
		}
	}

	return diff
}

func (l1 Location) Intersection(l2 Location) (intersection Location) {
	low, high := l1, l2
	if len(l1) > len(l2) {
		low = l2
		high = l1
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			f1 := i + 1
			f2 := j + 1
			if l.Equals(h) {
				intersection = append(intersection, h)
				if f1 < len(low) && f2 < len(high) {
					if !low[f1].Equals(high[f2]) {
						done = true
					}
				}
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		if done {
			break
		}
	}
	return
}
