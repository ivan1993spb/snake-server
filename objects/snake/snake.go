package snake

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/objects/apple"
	"github.com/ivan1993spb/snake-server/objects/corpse"
	"github.com/ivan1993spb/snake-server/objects/mouse"
	"github.com/ivan1993spb/snake-server/objects/wall"
	"github.com/ivan1993spb/snake-server/objects/watermelon"
)

const (
	snakeStartLength    = 3
	snakeStartSpeed     = time.Second * 4
	snakeSpeedFactor    = 1.02
	snakeStrengthFactor = 1
)

type Command string

const (
	CommandToNorth Command = "n"
	CommandToEast  Command = "e"
	CommandToSouth Command = "s"
	CommandToWest  Command = "w"
)

var snakeCommands = map[Command]engine.Direction{
	CommandToNorth: engine.DirectionNorth,
	CommandToEast:  engine.DirectionEast,
	CommandToSouth: engine.DirectionSouth,
	CommandToWest:  engine.DirectionWest,
}

// Snake object
type Snake struct {
	world game.World

	location engine.Location
	length   uint16

	// Motion direction
	direction engine.Direction
}

// CreateSnake creates new snake
func CreateSnake(world game.World) (*Snake, error) {
	var (
		dir  = engine.RandomDirection()
		err  error
		dots engine.Location
	)

	snake := &Snake{}

	switch dir {
	case engine.DirectionNorth, engine.DirectionSouth:
		dots, err = world.CreateObjectRandomRect(snake, 1, uint8(snakeStartLength))
	case engine.DirectionEast, engine.DirectionWest:
		dots, err = world.CreateObjectRandomRect(snake, uint8(snakeStartLength), 1)
	}
	if err != nil {
		// TODO: Create error
		return nil, err
	}

	if dir == engine.DirectionSouth || dir == engine.DirectionEast {
		reversedDots := dots.Reverse()
		dots = reversedDots
	}

	snake.world = world
	snake.location = dots
	snake.length = snakeStartLength
	snake.direction = dir

	return snake, nil
}

func (s *Snake) Die() {
	s.world.DeleteObject(s, s.location)
}

func (s *Snake) Feed(f int8) {
	if f > 0 {
		s.length += uint16(f)
	}
}

func (s *Snake) Strength() float32 {
	return snakeStrengthFactor * float32(s.length)
}

func (s *Snake) Run(ch <-chan game.Event) error {
	go func() {
		var ticker = time.NewTicker(s.calculateDelay())

		for {
			select {
			case <-ticker.C:
			}

			// Calculate next position
			dot, err := s.getNextHeadDot()
			if err != nil {
				// TODO How to emit error ?
				//s.p.OccurredError(s, err)
				return
			}

			// TODO: Delete this logic
			if object := s.world.GetObjectByDot(dot); object != nil {
				switch object := object.(type) {
				case *apple.Apple:
					object.NutritionalValue(dot)
				case *corpse.Corpse:
					object.NutritionalValue(dot)
				case *mouse.Mouse:
				case *Snake:
				case *wall.Wall:
				case *watermelon.Watermelon:
					object.NutritionalValue(dot)
				}

				ticker = time.NewTicker(s.calculateDelay())
			}

			tmpLocation := make(engine.Location, len(s.location)+1)
			copy(tmpLocation[1:], s.location)
			tmpLocation[0] = dot
			s.world.UpdateObject(s, s.location, tmpLocation)
			s.location = tmpLocation

			if s.length < s.location.DotCount() {
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
		return s.world.Navigate(s.location[0], s.direction, 1)
	}

	return nil, fmt.Errorf("cannot get next head location: errEmptyDotList")
}

func (s *Snake) Command(cmd Command) error {
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
