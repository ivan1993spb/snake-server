package objects

import "bitbucket.org/pushkin_ivan/clever-snake/game/playground"

type Apple struct {
	p   GameProcessor
	pg  *playground.Playground
	dot *playground.Dot
}

// CreateApple creates and locates new apple
func CreateApple(p GameProcessor, pg *playground.Playground,
) (*Apple, error) {
	if p == nil {
		return nil, &errCreateObject{errNilGameProcessor}
	}
	if pg == nil {
		return nil, &errCreateObject{errNilPlayground}
	}

	dot, err := pg.GetRandomEmptyDot()
	if err != nil {
		return nil, &errCreateObject{err}
	}

	apple := &Apple{p, pg, dot}

	if err := pg.Locate(apple, true); err != nil {
		return nil, &errCreateObject{err}
	}

	p.OccurredCreating(apple)

	return apple, nil
}

// Implementing playground.Object interface
func (*Apple) DotCount() uint16 {
	return 1
}

// Implementing playground.Object interface
func (a *Apple) Dot(i uint16) *playground.Dot {
	if i == 0 {
		return a.dot
	}
	return nil
}

// Implementing logic.Food interface
func (a *Apple) NutritionalValue(dot *playground.Dot) int8 {
	if a.dot.Equals(dot) {
		a.pg.Delete(a)

		a.p.OccurredDeleting(a)

		return 1
	}

	return 0
}
