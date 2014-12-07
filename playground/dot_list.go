package playground

// DotList represents list of dots and provides list packing
type DotList []*Dot

// Pack packs dot list to string in accordance with standard ST_1
func (dl DotList) Pack() (output string) {
	if len(dl) > 0 {
		for _, dot := range dl {
			output += ";"
			if dot == nil {
				output += NewDefaultDot().Pack()
			} else {
				output += dot.Pack()
			}
		}

		output = output[1:]
	}

	return
}

// Contains returns true if dot list contains passed dot
func (dl DotList) Contains(dot *Dot) bool {
	for i := range dl {
		if dl[i].Equals(dot) {
			return true
		}
	}

	return false
}

// Delete deletes dot from dot list
func (dl DotList) Delete(dot *Dot) {
	for i := range dl {
		if dl[i].Equals(dot) {
			dl = append(dl[:i], dl[i+1:]...)
			return
		}
	}
}

func (dl DotList) Reverse() (rdl DotList) {
	if len(dl) > 0 {
		rdl = make(DotList, 0, len(dl))
		for i := len(dl) - 1; i >= 0; i-- {
			rdl = append(rdl, dl[i])
		}
	}
	return
}
