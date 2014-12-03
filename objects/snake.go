package objects

import (
	"errors"
	"math"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"bitbucket.org/pushkin_ivan/clever-snake/logic"
	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

const (
	_SNAKE_START_LENGTH    uint16  = 3
	_SNAKE_START_SPEED             = time.Second * 4
	_SNAKE_SPEED_FACTOR    float64 = 1.02
	_SNAKE_STRENGTH_FACTOR float32 = 1
)

// Snake object
type Snake struct {
	pg *playground.Playground

	dots   playground.DotList
	length uint16

	// Next motion direction
	nextDirection playground.Direction

	// Time of last movement
	lastMove time.Time

	parentCxt context.Context

	stop context.CancelFunc
}

// CreateSnake creates new snake
func CreateSnake(pg *playground.Playground, cxt context.Context,
) (*Snake, error) {

	if pg != nil {
		return nil, errors.New("Passed nil playground")
	}

	var (
		dir  = playground.RandomDirection()
		dots playground.DotList
		err  error
	)

	switch dir {
	case playground.DIR_NORTH, playground.DIR_SOUTH:
		dots, err = pg.GetEmptyField(1, uint8(_SNAKE_START_LENGTH))
	case playground.DIR_EAST, playground.DIR_WEST:
		dots, err = pg.GetEmptyField(uint8(_SNAKE_START_LENGTH), 1)
	}

	if err != nil {
		return nil, err
	}

	if dir == playground.DIR_SOUTH || dir == playground.DIR_EAST {
		dots = dots.Reverse()
	}

	// Parent context stores in snake to pass it to corpse when snake
	// will be died. Snakes context are passed in run func
	scxt, cncl := context.WithCancel(cxt)

	snake := &Snake{pg, dots, _SNAKE_START_LENGTH, dir, time.Now(),
		cxt, cncl}

	if err = pg.Locate(snake); err != nil {
		return nil, err
	}

	snake.run(scxt)

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

// Implementing playground.Object interface
func (s *Snake) Pack() string {
	return strconv.Itoa(int(s.length)) + "%" + s.dots.Pack()
}

// Implementing playground.Shifting interface
func (s *Snake) Updated() time.Time {
	return s.lastMove
}

// Implementing playground.Shifting interface
func (s *Snake) PackChanges() string {
	return strconv.Itoa(int(s.length)) + "%" + s.dots[0].Pack() +
		"-" + s.dots[len(s.dots)-1].Pack()
}

// Implementing logic.Living interface
func (s *Snake) Die() {
	s.pg.Delete(s)
	if s.stop != nil {
		s.stop()
	}
	CreateCorpse(s.pg, s.parentCxt, s.dots)
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

func (s *Snake) run(cxt context.Context) {
	go func() {
		for {
			select {
			case <-cxt.Done():
				return
			case <-time.After(s.calculateDelay()):
			}

			if s.pg.Located(s) {
				return
			}

			// Calculate next position
			dot, err := s.GetNextHeadDot()
			if err != nil {
				return
			}

			if object := s.pg.GetObjectByDot(dot); object != nil {
				if err = logic.Clash(s, object, dot); err != nil {
					return
				}
			}

			if s.pg.Located(s) {
				return
			}

			tmpDots := make(playground.DotList, len(s.dots)+1)
			copy(tmpDots[1:], s.dots)
			tmpDots[0] = dot
			s.dots = tmpDots

			if s.length < s.DotCount() {
				s.dots = s.dots[:len(s.dots)-1]
			}

			s.lastMove = time.Now()
		}
	}()
}

func (s *Snake) calculateDelay() time.Duration {
	k := math.Pow(_SNAKE_SPEED_FACTOR, float64(s.length))
	return time.Duration(k * float64(_SNAKE_START_SPEED))
}

// GetNextHeadDot calculates new position of snake's head by its
// direction and current head position
func (s *Snake) GetNextHeadDot() (*playground.Dot, error) {
	return s.pg.Navigate(s.dots[0], s.nextDirection, 1)
}

// Implementing logic.Controlled interface
func (s *Snake) Command(cmd string) error {
	switch cmd {
	case "n":
		s.SetMovementDirection(playground.DIR_NORTH)
	case "e":
		s.SetMovementDirection(playground.DIR_EAST)
	case "s":
		s.SetMovementDirection(playground.DIR_SOUTH)
	case "w":
		s.SetMovementDirection(playground.DIR_WEST)
	default:
		return errors.New("Cannot execute command")
	}
	return nil
}

func (s *Snake) SetMovementDirection(nextDir playground.Direction) {
	if playground.ValidDirection(nextDir) {
		currDir := playground.CalculateDirection(s.dots[1], s.dots[0])
		// Next direction cannot be opposite to current direction
		if playground.ReverseDirection(nextDir) != currDir {
			s.nextDirection = nextDir
		}
	}
}
