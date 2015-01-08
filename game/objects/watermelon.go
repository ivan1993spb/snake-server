package objects

import (
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
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

type Watermelon struct {
	p  GameProcessor
	pg *playground.Playground

	dots [_WATERMELON_AREA]*playground.Dot
	stop context.CancelFunc
}

func CreateWatermelon(p GameProcessor, pg *playground.Playground,
	cxt context.Context) (*Watermelon, error) {
	if p == nil {
		return nil, &errCreateObject{errNilGameProcessor}
	}
	if pg == nil {
		return nil, &errCreateObject{errNilPlayground}
	}
	if err := cxt.Err(); err != nil {
		return nil, &errCreateObject{err}
	}

	e, err := pg.GetRandomEmptyRect(_WATERMELON_W, _WATERMELON_H)
	if err != nil {
		return nil, &errCreateObject{err}
	}
	dots := playground.EntityToDotList(e)

	wcxt, cncl := context.WithTimeout(cxt, _WATERMELON_MAX_EXPERIENCE)
	watermelon := &Watermelon{
		p:    p,
		pg:   pg,
		stop: cncl,
	}
	copy(watermelon.dots[:], dots[:_WATERMELON_AREA])

	if err := pg.Locate(watermelon, true); err != nil {
		return nil, &errCreateObject{err}
	}

	if err := watermelon.run(wcxt); err != nil {
		pg.Delete(watermelon)
		return nil, &errCreateObject{err}
	}

	p.OccurredCreating(watermelon)

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

// Implementing logic.Food interface
func (w *Watermelon) NutritionalValue(dot *playground.Dot) int8 {
	for i := range w.dots {
		if w.dots[i].Equals(dot) {
			w.dots[i] = nil

			w.p.OccurredUpdating(w)

			return _WATERMELON_MIN_NUTR_VALUE +
				int8(random.Intn(_WATERMELON_NUTR_VAR))
		}
	}

	if w.isEaten() {
		w.stop()
	}

	return 0
}

func (w *Watermelon) isEaten() bool {
	for i := 0; i < _WATERMELON_AREA; i++ {
		if w.dots[i] != nil {
			return false
		}
	}

	return true
}

func (w *Watermelon) run(cxt context.Context) error {
	if err := cxt.Err(); err != nil {
		return &errStartingObject{err}
	}

	go func() {
		<-cxt.Done()

		if w.pg.Located(w) {
			w.pg.Delete(w)
		}

		w.p.OccurredUpdating(w)
	}()

	return nil
}
