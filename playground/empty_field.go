package playground

import "errors"

const _RETRIES_NUMBER = 10

// GetEmptyField finds empty field with passed width and height
func (pg *Playground) GetEmptyField(w, h uint8) (DotList, error) {
	var pgW, pgH = pg.GetSize()

	if w*h == 0 || w > pgW || h > pgH {
		return nil, errors.New("Invalid field size")
	}

	var (
		x0, y0 uint8
		dots   DotList
		count  int
	)

mainLoop:

	x0 = uint8(random.Intn(int(pgW - w)))
	y0 = uint8(random.Intn(int(pgH - h)))
	dots = dots[:0]

	for x := x0; x < x0+w; x++ {
		for y := y0; y < y0+h; y++ {
			if dot := NewDot(x, y); !pg.Occupied(dot) {
				dots = append(dots, dot)
			} else if count < _RETRIES_NUMBER {
				count++
				goto mainLoop
			} else {
				return nil, errors.New("Cannot find empty field")
			}
		}
	}

	return dots, nil

}

// GetEmptyDot finds empty random dot
func (pg *Playground) GetEmptyDot() (*Dot, error) {
	for count := 0; count < _RETRIES_NUMBER; count++ {
		if dot := pg.RandomDot(); !pg.Occupied(dot) {
			return dot, nil
		}
	}

	return nil, errors.New("Cannot find empty dot")
}
