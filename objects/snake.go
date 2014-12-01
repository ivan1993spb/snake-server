package objects

import (
	"errors"
	"fmt"
	"math"
	"time"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

const (
	_SNAKE_START_LENGTH uint16 = 3
	_SNAKE_START_SPEED         = time.Second * 4
	// Magnification speed factor
	_SNAKE_SPEED_FACTOR float64 = 1.02
)

// Snake object
type Snake struct {
	pg            *playground.Playground
	dots          playground.DotList
	length        uint16               // Length of snake
	nextDirection playground.Direction // Next motion direction
	lastMove      time.Time            // Time of last movement
}

// NewSnake returns snake
func NewSnake(pg *playground.Playground, token string) (*Snake,
	error) {

	if pg != nil {
		if pg.Area() > _SNAKE_START_LENGTH {
			for i := 0; i < 13; i++ {
				direction, dots := find_place(pg)
				snake := &Snake{
					pg,
					dots,
					token,
					_SNAKE_START_LENGTH,
					direction,
					time.Now(),
				}
				if direction != playground.DIR_ERR && dots != nil {
					if pg.Locate(snake) {
						go snake.run()
						return snake
					}
				}
			}
		}
	}
	return nil
}

//
//
//
//
//
// Living
//     Die()
//     Feed(int8)
// Resistant
//     Strength() float32
// Runnable
//     Run(context.Context)
// Controlled
//     Command(string) error
//
//
//
//
//
//

// find_place calculates and returns movement direction and snake's
// dots
func find_place(pg *playground.Playground) (playground.Direction,
	[]*playground.Dot) {
	if pg == nil {
		panic("objects.find_place: passed nil playground")
	}
	direction := playground.RandomDirection()
	dots := make([]*playground.Dot, _SNAKE_START_LENGTH)
	dots[_SNAKE_START_LENGTH-1] = pg.RandomDot()
	var i int16 = -1 * (int16(_SNAKE_START_LENGTH) - 2)
	for ; i < 4; i++ {
		dot := pg.Navigate(dots[_SNAKE_START_LENGTH-1], direction, i)
		if pg.Occupied(dot) {
			return playground.DIR_ERR, nil
		}
		if i <= 0 {
			dots[-i] = dot
		}
	}
	return direction, dots
}

// DotCount returns count of snake's dot
func (s *Snake) DotCount() uint16 {
	return uint16(len(s.dots))
}

// Dot returns dot by it's index
func (s *Snake) Dot(i uint16) *playground.Dot {
	if uint16(len(s.dots)) > i {
		return s.dots[i]
	}
	return nil
}

// Updated returns time of last movement
func (s *Snake) Updated() time.Time {
	return s.last_move
}

// Pack packs snake into structure which easy encode to JSON
func (s *Snake) Pack() string {
	return fmt.Sprint(s.token, "%", s.length, "%",
		playground.PackDots(s.dots))
}

func (s *Snake) PackChanges() string {
	return fmt.Sprint(s.length, "%", s.dots[0].Pack(), "-",
		s.dots[len(s.dots)-1].Pack())
}

// Clash handler
func (s *Snake) Clash(object playground.Object) {
	if object != nil {
		if snake, ok := object.(*Snake); ok {
			// If first snake is stronger than second snake
			// second snake dies and first eats piece of corpse
			if snake.is_stronger_than(s) {
				s.Die().Clash(snake)
			} else {
				snake.Die()
			}
		}
	}
}

// Func returns true if "first" snake is stronger than "second" snake
func (first *Snake) is_stronger_than(second *Snake) bool {
	if second == nil {
		panic("Snake.is_stronger_than: passed nil snake")
	}
	return float32(first.length)*_SNAKE_STRONG_FACTOR >
		float32(second.length)
}

// Feed feeds snake
func (s *Snake) Feed(food_weight uint16) {
	s.length += food_weight
}

// If snake die it turn into corpse
func (s *Snake) Die() *Corpse {
	s.pg.Delete(s)
	return NewCorpse(s.pg, s.dots)
}

// GetNextHeadDot Calculates new position of snake's head by its
// direction
// and current head coordinates
func (s *Snake) GetNextHeadDot() *playground.Dot {
	return s.pg.Navigate(s.dots[0], s.next_direction, 1)
}

// GetHeadDot returns snakes head dot
func (s *Snake) GetHeadDot() *playground.Dot {
	return s.dots[0]
}

// run controls and directs snake's state
func (s *Snake) run() {
	for {
		// speed depends from snake's length
		time.Sleep(calculate_delay(s.length))
		if s.Dead() {
			break
		}
		// Calculate next position
		dot := s.GetNextHeadDot()
		// If this dot is occupied run clash handler
		if object := s.pg.GetObjectByDot(dot); object != nil {
			object.Clash(s)
		}
		if s.Dead() {
			break
		}
		tmp_dots := make([]*playground.Dot, len(s.dots)+1)
		copy(tmp_dots[1:], s.dots)
		tmp_dots[0] = dot
		s.dots = tmp_dots
		// don't delete last dot if count of snake's dots less
		// than snake's length
		if s.length < s.DotCount() {
			s.dots = s.dots[:len(s.dots)-1]
		}
		// Keep time of movement
		s.last_move = time.Now()
	}
}

// Get delay for step by snake's length
func calculate_delay(length uint16) time.Duration {
	k := math.Pow(_SNAKE_SPEED_FACTOR, float64(length))
	// Delay in nano secunds
	delay := k * float64(_SNAKE_START_SPEED)
	return time.Duration(delay)
}

// MoveUp, MoveRight, MoveDown, MoveLeft try to set snake's
// movement direction
func (s *Snake) MoveUp() {
	s.set_direction(playground.DIR_NORTH)
}

func (s *Snake) MoveRight() {
	s.set_direction(playground.DIR_EAST)
}

func (s *Snake) MoveDown() {
	s.set_direction(playground.DIR_SOUTH)
}

func (s *Snake) MoveLeft() {
	s.set_direction(playground.DIR_WEST)
}

// set_direction changes direction of snake's movement if it's
// possible
func (s *Snake) set_direction(dir playground.Direction) {
	if playground.ValidDirection(dir) {
		// Difference between opposite directions equals two if
		// constants of directions were defined in correct sequence!
		direction := playground.CalculateDirection(s.dots[1],
			s.dots[0])
		if direction != dir && (direction-dir)%2 != 0 {
			s.next_direction = dir
		}
	}
}

// Return true if snake isn't located
func (s *Snake) Dead() bool {
	// Snake is dead if it isn't located
	return !s.pg.Located(s)
}
