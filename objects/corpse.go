package objects

import (
	"errors"
	"time"

	"golang.org/x/net/context"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

// Time for which corpse will lie on playground
const _CORPSE_MAX_EXPERIENCE = time.Second * 15

// Snakes can eat corpses
type Corpse struct {
	// Playground on which corpse lies
	pg *playground.Playground

	// Dots on which corpse's pieces lie updated is last time when a
	// snake ate any piece of corpse
	dots playground.DotList

	updated time.Time

	// last nipped piece
	nippedPiece *playground.Dot

	eaten chan struct{}
}

// Corpses are created when a snake dies
func NewCorpse(pg *playground.Playground, dots playground.DotList,
) (*Corpse, error) {
	if pg != nil && len(dots) > 0 {
		c := &Corpse{pg, dots, time.Now(), nil, make(chan struct{})}
		if err := pg.Locate(c); err == nil {
			return c, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Cannot create corpse")
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

// Implementing logic.Runnable interface
func (c *Corpse) Run(cxt context.Context) {
	go func() {
		select {
		case <-cxt.Done():
			// If pool are closed
		case <-time.After(_CORPSE_MAX_EXPERIENCE):
			// If corpse lies too long
		case <-c.eaten:
			// If corpse was eaten
		}

		if c.pg.Located(c) {
			c.pg.Delete(c)
		}
	}()
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
	return ""
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
			close(c.eaten)
		}

		return 2
	}

	return 0
}
