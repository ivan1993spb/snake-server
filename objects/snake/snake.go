package snake

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/olebedev/emitter"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

const (
	snakeStartLength    = 3
	snakeStartSpeed     = time.Second * 4
	snakeSpeedFactor    = 1.02
	snakeStrengthFactor = 1
)

type SnakeCommand string

const (
	SnakeCommandToNorth SnakeCommand = "n"
	SnakeCommandToEast  SnakeCommand = "e"
	SnakeCommandToSouth SnakeCommand = "s"
	SnakeCommandToWest  SnakeCommand = "w"
)

var snakeCommands = map[SnakeCommand]engine.Direction{
	SnakeCommandToNorth: engine.DirectionNorth,
	SnakeCommandToEast:  engine.DirectionEast,
	SnakeCommandToSouth: engine.DirectionSouth,
	SnakeCommandToWest:  engine.DirectionWest,
}

// Snake object
type Snake struct {
	pg *playground.Playground

	location engine.Location
	length   uint16

	// Motion direction
	direction engine.Direction
}

// CreateSnake creates new snake
func CreateSnake(pg *playground.Playground) (*Snake, error) {
	var (
		dir  = engine.RandomDirection()
		err  error
		dots engine.Location
	)

	snake := &Snake{}

	switch dir {
	case engine.DirectionNorth, engine.DirectionSouth:
		dots, err = pg.CreateObjectRandomRect(snake, 1, uint8(snakeStartLength))
	case engine.DirectionEast, engine.DirectionWest:
		dots, err = pg.CreateObjectRandomRect(snake, uint8(snakeStartLength), 1)
	}
	if err != nil {
		// TODO: Create error
		return nil, err
	}

	if dir == engine.DirectionSouth || dir == engine.DirectionEast {
		reversedDots := dots.Reverse()
		dots = reversedDots
	}

	snake.pg = pg
	snake.location = dots
	snake.length = snakeStartLength
	snake.direction = dir

	return snake, nil
}

// Implementing playground.Location interface
func (s *Snake) DotCount() uint16 {
	return uint16(len(s.location))
}

// Implementing playground.Location interface
func (s *Snake) Dot(i uint16) *engine.Dot {
	if uint16(len(s.location)) > i {
		return s.location[i]
	}

	return nil
}

// Implementing logic.Living interface
func (s *Snake) Die() {
	s.pg.DeleteObject(s, s.location)
}

// Implementing logic.Living interface
func (s *Snake) Feed(f int8) {
	if f > 0 {
		s.length += uint16(f)
	}
}

// Implementing logic.Resistant interface
func (s *Snake) Strength() float32 {
	return snakeStrengthFactor * float32(s.length)
}

func (s *Snake) Run(emitter *emitter.Emitter) error {
	go func() {
		var ticker = time.NewTicker(s.calculateDelay())

		for {
			select {
			case <-ticker.C:
			}

			if !s.pg.Located(s) {
				return
			}

			// Calculate next position
			dot, err := s.getNextHeadDot()
			if err != nil {
				s.p.OccurredError(s, err)
				return
			}

			// TODO: Delete this logic
			if object := s.pg.GetObjectByDot(dot); object != nil {
				if err = logic.Clash(s, object, dot); err != nil {
					s.p.OccurredError(s, err)
					return
				}

				if !s.pg.Located(s) {
					return
				}

				ticker = time.NewTicker(s.calculateDelay())
			}

			tmpLocation := make(engine.Location, len(s.location)+1)
			copy(tmpLocation[1:], s.location)
			tmpLocation[0] = dot
			s.pg.UpdateObject(s, s.location, tmpLocation)
			s.location = tmpLocation

			if s.length < s.DotCount() {
				s.location = s.location[:len(s.location)-1]
			}
		}
	}()

	return nil
}

func (s *Snake) calculateDelay() time.Duration {
	return time.Duration(math.Pow(snakeSpeedFactor, float64(s.length)) * float64(snakeStartSpeed))
}

// getNextHeadDot calculates new position of snake's head by its
// direction and current head position
func (s *Snake) getNextHeadDot() (*engine.Dot, error) {
	if len(s.location) > 0 {
		return s.pg.Navigate(s.location[0], s.direction, 1)
	}

	return nil, fmt.Errorf("cannot get next head location: %s", errEmptyDotList)
}

// Implementing logic.Controlled interface
func (s *Snake) Command(cmd SnakeCommand) error {
	if direction, ok := snakeCommands[cmd]; ok {
		s.setMovementDirection(direction)
		return nil
	}
	return errors.New("cannot execute command")
}

func (s *Snake) setMovementDirection(nextDir engine.Direction) error {
	if engine.ValidDirection(nextDir) {
		currDir := engine.CalculateDirection(s.location[1], s.location[0])
		rNextDir, err := nextDir.Reverse()
		if err != nil {
			return fmt.Errorf("cannot set movement direction: %s", err)
		}

		// Next direction cannot be opposite to current direction
		if rNextDir == currDir {
			return errors.New("next direction cannot be opposite to current direction")
		} else {
			s.direction = nextDir
		}
	}

	return nil
}
