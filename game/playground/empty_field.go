// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package playground

import "errors"

const _RETRIES_NUMBER = 35

var errRetriesLimit = errors.New("retries limit was reached")

type errEmptyField struct {
	err error
}

func (e *errEmptyField) Error() string {
	return "cannot find empty field: " + e.err.Error()
}

// GetEmptyField finds empty field with passed width and height
func (pg *Playground) GetRandomEmptyRect(rw, rh uint8,
) (*Rect, error) {

	var count = 0

loop:

	if rect, err := pg.RandomRect(rw, rh); err == nil {
		for i := uint16(0); i < rect.DotCount(); i++ {
			if pg.Occupied(rect.Dot(i)) {
				goto rewind
			}
		}
		return rect, nil
	} else {
		return nil, &errEmptyField{err}
	}

rewind:

	if count < _RETRIES_NUMBER {
		count++
		goto loop
	}

	return nil, &errEmptyField{errRetriesLimit}
}

// GetRandomEmptyDot returns random empty dot
func (pg *Playground) GetRandomEmptyDot() (*Dot, error) {
	for count := 0; count < _RETRIES_NUMBER; count++ {
		if dot := pg.RandomDot(); !pg.Occupied(dot) {
			return dot, nil
		}
	}

	return nil, &errEmptyField{errRetriesLimit}
}
