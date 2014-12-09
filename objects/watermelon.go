package objects

import (
	"math/rand"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
	"golang.org/x/net/context"
)

const (
	_WATERMELON_W = 2 // Width
	_WATERMELON_H = 2 // Height

	_WATERMELON_AREA = _WATERMELON_W * _WATERMELON_H

	_WATERMELON_MIN_NUTR_VALUE = 1
	_WATERMELON_MAX_NUTR_VALUE = 10

	_WATERMELON_NUTR_VAR = _WATERMELON_MAX_NUTR_VALUE -
		_WATERMELON_MIN_NUTR_VALUE

	_WATERMELON_MAX_EXPERIENCE = time.Second * 30
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type Watermelon struct {
	pg      *playground.Playground
	dots    [_WATERMELON_AREA]*playground.Dot
	updated time.Time
	cancel  context.CancelFunc
}

func CreateWatermelon(pg *playground.Playground, cxt context.Context,
) (*Watermelon, error) {

	if pg == nil {
		return nil, &errCreateObject{playground.ErrNilPlayground}
	}
	if err := cxt.Err(); err != nil {
		return nil, &errCreateObject{err}
	}

	dots, err := pg.GetEmptyField(_WATERMELON_W, _WATERMELON_H)
	if err != nil {
		return nil, &errCreateObject{err}
	}

	wcxt, cncl := context.WithCancel(cxt)
	watermelon := &Watermelon{
		pg:      pg,
		updated: time.Now(),
		cancel:  cncl,
	}
	copy(watermelon.dots[:], dots[:_WATERMELON_AREA])

	if err := pg.Locate(watermelon); err != nil {
		return nil, &errCreateObject{err}
	}

	if err := watermelon.run(wcxt); err != nil {
		pg.Delete(watermelon)
		return nil, &errCreateObject{err}
	}

	return watermelon, nil
}

// Implementing playground.Object interface
func (w *Watermelon) DotCount() (c uint16) {
	for _, dot := range w.dots {
		if dot != nil {
			c++
		}
	}
	return
}

// Implementing playground.Object interface
func (w *Watermelon) Dot(i uint16) (dot *playground.Dot) {
	if i < w.DotCount() {
		var j uint16
		for _, dot = range w.dots {
			if dot != nil {
				if i == j {
					break
				}
				j++
			}
		}
	}
	return
}

// Implementing playground.Object interface
func (w *Watermelon) Pack() (res string) {
	for _, dot := range w.dots {
		if dot != nil {
			res += ";" + dot.Pack()
		} else {
			res += ";_"
		}
	}
	if len(res) > 0 {
		res = res[1:]
	}
	return
}

// Implementing playground.Shifting interface
func (w *Watermelon) PackChanges() string {
	var output string
	for _, dot := range w.dots {
		if dot != nil {
			output += "1"
		} else {
			output += "0"
		}
	}
	return output
}

// Implementing playground.Shifting interface
func (w *Watermelon) Updated() time.Time {
	return w.updated
}

// Implementing logic.Food interface
func (w *Watermelon) NutritionalValue(dot *playground.Dot) int8 {
	for i := range w.dots {
		if w.dots[i].Equals(dot) {
			w.updated = time.Now()
			w.dots[i] = nil

			return _WATERMELON_MIN_NUTR_VALUE +
				int8(random.Intn(_WATERMELON_NUTR_VAR))
		}
	}

	return 0
}

func (w *Watermelon) run(cxt context.Context) error {
	if err := cxt.Err(); err != nil {
		return err
	}

	go func() {
		select {
		case <-cxt.Done():
		case <-time.After(_WATERMELON_MAX_EXPERIENCE):
		}
		if w.pg.Located(w) {
			w.pg.Delete(w)
		}
	}()

	return nil
}
