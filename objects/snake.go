package objects

import (
	"errors"
	"math"
	"strconv"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/logic"
	"bitbucket.org/pushkin_ivan/clever-snake/playground"
	"golang.org/x/net/context"
)

const (
	_SNAKE_START_LENGTH uint16 = 3
	_SNAKE_START_SPEED         = time.Second * 4

	// Magnification of speed factor
	_SNAKE_SPEED_FACTOR float64 = 1.02

	_SNAKE_LOC_RETRIES_NUMBER = 15

	_SNAKE_STRENGTH_FACTOR float32 = 0.5
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

	dead chan struct{}
}

// NewSnake creates new snake
func NewSnake(pg *playground.Playground) (*Snake, error) {
	if pg != nil {
		if w, h := pg.GetSize(); uint16(w) < _SNAKE_START_LENGTH ||
			uint16(h) < _SNAKE_START_LENGTH {
			return nil, errors.New("Playground is too small")
		}

		for i := 0; i < _SNAKE_LOC_RETRIES_NUMBER; i++ {
			direction, dots, err := findPlace(pg)
			if err == nil {
				snake := &Snake{pg, dots, _SNAKE_START_LENGTH,
					direction, time.Now(), make(chan struct{})}

				if pg.Locate(snake) == nil {
					return snake, nil
				}
			}
		}
	}

	return nil, errors.New("Cannot create snake")
}

// findPlace calculates movement direction and dots for snake
func findPlace(pg *playground.Playground,
) (playground.Direction, playground.DotList, error) {
	if pg == nil {
		return 0, nil, errors.New("Passed nil playground")
	}

	direction := playground.RandomDirection()
	dots := make(playground.DotList, _SNAKE_START_LENGTH)

	dots[_SNAKE_START_LENGTH-1] = pg.RandomDot()

	var i int16 = -1 * (int16(_SNAKE_START_LENGTH) - 2)
	for ; i < int16(_SNAKE_START_LENGTH+1); i++ {
		dot, err := pg.Navigate(dots[_SNAKE_START_LENGTH-1],
			direction, i)
		if err == nil {
			return 0, nil, err
		}
		if pg.Occupied(dot) {
			return 0, nil, errors.New("Occupied dot")
		}
		if i <= 0 {
			dots[-i] = dot
		}
	}

	return direction, dots, nil
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
	close(s.dead)
	NewCorpse(s.pg, s.dots)
}

// Implementing logic.Living interface
func (s *Snake) Feed(f int8) {
	if f > 0 {
		s.length += uint16(f)
	} else {
		f *= -1
		if s.length > uint16(f) {
			s.length -= uint16(f)
		} else {
			s.length = 0
		}
	}
}

// Implementing logic.Resistant interface
func (s *Snake) Strength() float32 {
	return _SNAKE_STRENGTH_FACTOR * float32(s.length)
}

// Implementing logic.Runnable interface
func (s *Snake) Run(cxt context.Context) {
	for {

		select {
		case <-cxt.Done():
			return
		case <-s.dead:
			return
		case <-time.After(calculateDelay(s.length)):
		}

		if s.pg.Located(s) {
			return
		}

		// Calculate next position
		dot, err := s.GetNextHeadDot()
		if err != nil {
			return
		}

		// If this dot is occupied run clash handler
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
}

// calculateDelay calculates delay by snake length
func calculateDelay(length uint16) time.Duration {
	k := math.Pow(_SNAKE_SPEED_FACTOR, float64(length))
	// Delay in nano secunds
	delay := k * float64(_SNAKE_START_SPEED)
	return time.Duration(delay)
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
		s.SetDirection(playground.DIR_NORTH)
	case "e":
		s.SetDirection(playground.DIR_EAST)
	case "s":
		s.SetDirection(playground.DIR_SOUTH)
	case "w":
		s.SetDirection(playground.DIR_WEST)
	}
	return errors.New("Cannot execute command")
}

func (s *Snake) SetDirection(dir playground.Direction) {
	if playground.ValidDirection(dir) {
		// Difference between opposite directions equals two if
		// constants of directions were defined in correct sequence!
		direction := playground.CalculateDirection(s.dots[1],
			s.dots[0])

		if direction != dir && (direction-dir)%2 != 0 {
			s.nextDirection = dir
		}
	}
}
