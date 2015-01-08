// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playground

// DotList represents list of dots
type DotList []*Dot

// Contains returns true if dot list contains passed dot
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

// Delete deletes dot from dot list
func (dl DotList) Delete(dot *Dot) {
	if len(dl) > 0 {
		for i := range dl {
			if dl[i].Equals(dot) {
				dl = append(dl[:i], dl[i+1:]...)
				return
			}
		}
	}
}

// Reverse reverses dot list
func (dl DotList) Reverse() DotList {
	if len(dl) > 0 {
		rdl := make(DotList, 0, len(dl))
		for i := len(dl) - 1; i >= 0; i-- {
			rdl = append(rdl, dl[i])
		}

		return rdl
	}

	return DotList{}
}

// Implementing Entity interface
func (dl DotList) Dot(i uint16) *Dot {
	return dl[i]
}

// Implementing Entity interface
func (dl DotList) DotCount() uint16 {
	return uint16(len(dl))
}
