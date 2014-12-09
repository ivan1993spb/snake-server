package objects

import (
	"errors"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
	"golang.org/x/net/context"
)

// Time for which corpse will be lie on playground
const _CORPSE_MAX_EXPERIENCE = time.Second * 15

// Snakes can eat corpses
type Corpse struct {
	pg *playground.Playground

	dots playground.DotList

	updated time.Time

	// last nipped piece
	nippedPiece *playground.Dot

	stop context.CancelFunc
}

// Corpses are created when a snake dies
func CreateCorpse(pg *playground.Playground, cxt context.Context,
	dots playground.DotList) (*Corpse, error) {

	if pg == nil {
		return nil, errors.New("Passed nil playground")
	}
	if len(dots) == 0 {
		return nil, errors.New("Passed empty dot list")
	}
	if err := cxt.Err(); err != nil {
		return nil, err
	}

	ccxt, cancel := context.WithCancel(cxt)
	corpse := &Corpse{pg, dots, time.Now(), nil, cancel}

	if err := pg.Locate(corpse); err != nil {
		return nil, err
	}

	if err := corpse.run(ccxt); err != nil {
		pg.Delete(corpse)
		return nil, err
	}

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

// Implementing playground.Object interface
func (c *Corpse) Pack() string {
	return c.dots.Pack()
}

// Implementing playground.Shifting interface. Updated returns last
// time when a piece of corpse was eaten
func (c *Corpse) Updated() time.Time {
	return c.updated
}

// Implementing playground.Shifting interface
func (c *Corpse) PackChanges() string {
	if c.nippedPiece != nil {
		return "nip" + c.nippedPiece.Pack()
	}
	return c.Pack()
}

// Implementing logic.Food interface
func (c *Corpse) NutritionalValue(dot *playground.Dot) int8 {
	if c.dots.Contains(dot) {
		c.dots.Delete(dot)

		if len(c.dots) > 0 {
			c.nippedPiece = dot
			c.updated = time.Now()
		} else {
			// Remove corpse if it was eaten
			c.pg.Delete(c)
			if c.stop != nil {
				c.stop()
			}
		}

		return 2
	}

	return 0
}

func (c *Corpse) run(cxt context.Context) error {
	if err := cxt.Err(); err != nil {
		return err
	}

	go func() {
		select {
		case <-cxt.Done():
			// If pool are closed or corpse was eaten
		case <-time.After(_CORPSE_MAX_EXPERIENCE):
			// If corpse lies too long
		}
		if c.pg.Located(c) {
			c.pg.Delete(c)
		}
	}()

	return nil
}
