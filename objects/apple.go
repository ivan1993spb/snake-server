package objects

import "bitbucket.org/pushkin_ivan/clever-snake/playground"

type Apple struct {
	pg  *playground.Playground
	dot *playground.Dot
}

// CreateApple creates and locates new apple
func CreateApple(pg *playground.Playground) (*Apple, error) {
	if pg == nil {
		return nil, playground.ErrNilPlayground
	}

	dot, err := pg.GetEmptyDot()
	if err != nil {
		return nil, err
	}

	apple := &Apple{pg, dot}

	if err := pg.Locate(apple); err != nil {
		return nil, err

	}

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

// Implementing playground.Object interface
func (a *Apple) Pack() string {
	return a.dot.Pack()
}

// Implementing logic.Food interface
func (a *Apple) NutritionalValue(dot *playground.Dot) int8 {
	if a.dot.Equals(dot) {
		a.pg.Delete(a)
		CreateApple(a.pg)
		return 1
	}
	return 0
}
