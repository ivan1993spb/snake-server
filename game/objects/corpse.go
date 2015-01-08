package objects

import (
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
	"golang.org/x/net/context"
)

// Time for which corpse will be lie on playground
const _CORPSE_MAX_EXPERIENCE = time.Second * 15

// Snakes can eat corpses
type Corpse struct {
	p  GameProcessor
	pg *playground.Playground

	dots        playground.DotList
	nippedPiece *playground.Dot // last nipped piece

	stop context.CancelFunc
}

// Corpses are created when a snake dies
func CreateCorpse(p GameProcessor, pg *playground.Playground,
	cxt context.Context, dots playground.DotList) (*Corpse, error) {
	if p == nil {
		return nil, &errCreateObject{errNilGameProcessor}
	}
	if pg == nil {
		return nil, &errCreateObject{errNilPlayground}
	}
	if len(dots) == 0 {
		return nil, &errCreateObject{errEmptyDotList}
	}
	if err := cxt.Err(); err != nil {
		return nil, &errCreateObject{err}
	}

	ccxt, cancel := context.WithTimeout(cxt, _CORPSE_MAX_EXPERIENCE)
	corpse := &Corpse{p, pg, dots, nil, cancel}

	if err := pg.Locate(corpse, true); err != nil {
		return nil, &errCreateObject{err}
	}

	if err := corpse.run(ccxt); err != nil {
		pg.Delete(corpse)
		return nil, &errCreateObject{err}
	}

	p.OccurredCreating(corpse)

	return corpse, nil
}

// Implementing playground.Object interface
func (c *Corpse) DotCount() uint16 {
	return uint16(len(c.dots))
}

// Implementing playground.Object interface
func (c *Corpse) Dot(i uint16) *playground.Dot {
	if c.DotCount() > i {
		return c.dots[i]
	}
	return nil
}

// Implementing logic.Food interface
func (c *Corpse) NutritionalValue(dot *playground.Dot) int8 {
	if c.dots.Contains(dot) {
		c.dots.Delete(dot)

		if len(c.dots) > 0 {
			c.nippedPiece = dot
			c.p.OccurredUpdating(c)
		} else {
			c.stop()
		}

		return 2
	}

	return 0
}

func (c *Corpse) run(cxt context.Context) error {
	if err := cxt.Err(); err != nil {
		return &errStartingObject{err}
	}

	go func() {
		<-cxt.Done()

		if c.pg.Located(c) {
			c.pg.Delete(c)
		}

		c.p.OccurredDeleting(c)
	}()

	return nil
}
