package objects

import (
	"errors"
	"fmt"
	"math"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/game/logic"
	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
	"golang.org/x/net/context"
)

const (
	_SNAKE_START_LENGTH    = 3
	_SNAKE_START_SPEED     = time.Second * 4
	_SNAKE_SPEED_FACTOR    = 1.02
	_SNAKE_STRENGTH_FACTOR = 1
)

// Snake object
type Snake struct {
	p  GameProcessor
	pg *playground.Playground

	dots   playground.DotList
	length uint16

	// Next motion direction
	nextDirection playground.Direction

	// stop is CancelFunc of child context of parentCxt which belonges
	// to snake. Calling stop() stops snake
	stop context.CancelFunc
}

// CreateSnake creates new snake
func CreateSnake(p GameProcessor, pg *playground.Playground,
	cxt context.Context) (*Snake, error) {
	if p == nil {
		return nil, &errCreateObject{errNilGameProcessor}
	}
	if pg == nil {
		return nil, &errCreateObject{errNilPlayground}
	}
	if err := cxt.Err(); err != nil {
		return nil, &errCreateObject{err}
	}

	var (
		dir = playground.RandomDirection()
		err error
		e   playground.Entity
	)

	switch dir {
	case playground.DIR_NORTH, playground.DIR_SOUTH:
		e, err = pg.GetRandomEmptyRect(1, uint8(_SNAKE_START_LENGTH))
	case playground.DIR_EAST, playground.DIR_WEST:
		e, err = pg.GetRandomEmptyRect(uint8(_SNAKE_START_LENGTH), 1)
	}
	if err != nil {
		return nil, &errCreateObject{err}
	}

	dots := playground.EntityToDotList(e)

	if dir == playground.DIR_SOUTH || dir == playground.DIR_EAST {
		dots = dots.Reverse()
	}

	scxt, cncl := context.WithCancel(cxt)

	snake := &Snake{p, pg, dots, _SNAKE_START_LENGTH, dir, cncl}

	if err = pg.Locate(snake, true); err != nil {
		return nil, &errCreateObject{err}
	}

	if err := snake.run(scxt); err != nil {
		pg.Delete(snake)
		return nil, &errCreateObject{err}
	}

	p.OccurredCreating(snake)

	return snake, nil
}

// Implementing playground.Object interface
func (s *Snake) DotCount() uint16 {
	return uint16(len(s.dots))
}

// Implementing playground.Object interface
func (s *Snake) Dot(i uint16) *playground.Dot {
	if uint16(len(s.dots)) > i {
		return s.dots[i]
	}

	return nil
}

// Implementing logic.Living interface
func (s *Snake) Die() {
	s.pg.Delete(s)
	if s.stop != nil {
		s.stop()
	}
}

// Implementing logic.Living interface
func (s *Snake) Feed(f int8) {
	if f > 0 {
		s.length += uint16(f)
	}
}

// Implementing logic.Resistant interface
func (s *Snake) Strength() float32 {
	return _SNAKE_STRENGTH_FACTOR * float32(s.length)
}

func (s *Snake) run(cxt context.Context) error {
	if err := cxt.Err(); err != nil {
		return &errStartingObject{err}
	}

	go func() {
		var ticker = time.NewTicker(s.calculateDelay())

		defer func() {
			if s.pg.Located(s) {
				s.Die()
			}
			ticker.Stop()
			s.p.OccurredDeleting(s)
		}()

		for {
			select {
			case <-cxt.Done():
				return
			case <-ticker.C:
			}

			if !s.pg.Located(s) {
				return
			}

			// Calculate next position
			dot, err := s.getNextHeadDot()
			if err != nil {
				s.p.OccurredError(err)
				return
			}

			if object := s.pg.GetObjectByDot(dot); object != nil {
				if err = logic.Clash(s, object, dot); err != nil {
					s.p.OccurredError(err)
					return
				}

				if !s.pg.Located(s) {
					return
				}

				ticker = time.NewTicker(s.calculateDelay())
			}

			tmpDots := make(playground.DotList, len(s.dots)+1)
			copy(tmpDots[1:], s.dots)
			tmpDots[0] = dot
			s.dots = tmpDots

			if s.length < s.DotCount() {
				s.dots = s.dots[:len(s.dots)-1]
			}
		}
	}()

	return nil
}

func (s *Snake) calculateDelay() time.Duration {
	return time.Duration(math.Pow(_SNAKE_SPEED_FACTOR,
		float64(s.length)) * float64(_SNAKE_START_SPEED))
}

// getNextHeadDot calculates new position of snake's head by its
// direction and current head position
func (s *Snake) getNextHeadDot() (*playground.Dot, error) {
	if len(s.dots) > 0 {
		return s.pg.Navigate(s.dots[0], s.nextDirection, 1)
	}

	return nil, fmt.Errorf("cannot get next head dot: %s",
		errEmptyDotList)
}

// Implementing logic.Controlled interface
func (s *Snake) Command(cmd string) error {
	switch cmd {
	case "n":
		s.setMovementDirection(playground.DIR_NORTH)
	case "e":
		s.setMovementDirection(playground.DIR_EAST)
	case "s":
		s.setMovementDirection(playground.DIR_SOUTH)
	case "w":
		s.setMovementDirection(playground.DIR_WEST)
	default:
		return errors.New("cannot execute command")
	}
	return nil
}

func (s *Snake) setMovementDirection(nextDir playground.Direction,
) error {
	if playground.ValidDirection(nextDir) {
		currDir := playground.CalculateDirection(s.dots[1], s.dots[0])
		rNextDir, err := playground.ReverseDirection(nextDir)
		if err != nil {
			return fmt.Errorf("cannot set movement direction: %s",
				err)
		}

		// Next direction cannot be opposite to current direction
		if rNextDir == currDir {
			return errors.New("next direction cannot be opposite to" +
				" current direction")
		} else {
			s.nextDirection = nextDir
		}
	}

	return nil
}
