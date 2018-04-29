package engine

// DotList is set of dots
type DotList []*Dot

// contains returns true if object contains passed dot
func (dl DotList) Contains(dot *Dot) bool {
	if len(dl) > 0 {
		for _, d := range dl {
			if d.Equals(dot) {
				return true
			}
		}
	}

	return false
}

// Delete deletes dot from object
func (dl DotList) Delete(dot *Dot) DotList {
	if len(dl) > 0 {
		for i := range dl {
			if dl[i].Equals(dot) {
				return append(dl[:i], dl[i+1:]...)
			}
		}
	}
	return DotList{}
}

func (dl DotList) Add(dot *Dot) DotList {
	return append(dl, dot)
}

// Reverse reverses dot sequence in object
func (dl DotList) Reverse() DotList {
	if len(dl) > 0 {
		ro := make(DotList, 0, len(dl))
		for i := len(dl) - 1; i >= 0; i-- {
			ro = append(ro, dl[i])
		}

		return ro
	}

	return DotList{}
}

func (dl DotList) Dot(i uint16) *Dot {
	return dl[i]
}

func (dl DotList) DotCount() uint16 {
	return uint16(len(dl))
}

func (dl DotList) Copy() DotList {
	newList := make(DotList, 0, len(dl))
	copy(newList, dl)
	return newList
}

func (dl DotList) Equals(o2 DotList) bool {
	if len(dl) == 0 && len(o2) == 0 {
		return true
	}
	if len(dl) != len(o2) {
		return false
	}

	for i := 0; i < len(dl); i++ {
		if !dl[i].Equals(o2[i]) {
			return false
		}
	}

	return true
}
