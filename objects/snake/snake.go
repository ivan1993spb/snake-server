package snake

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/objects"
	"github.com/ivan1993spb/snake-server/world"
)

const (
	snakeStartLength    = 3
	snakeStartSpeed     = time.Second
	snakeSpeedFactor    = 1.02
	snakeStrengthFactor = 1
	snakeStartMargin    = 1
	snakeTypeLabel      = "snake"
)

type Command string

const (
	CommandToNorth Command = "north"
	CommandToEast  Command = "east"
	CommandToSouth Command = "south"
	CommandToWest  Command = "west"
)

var snakeCommands = map[Command]engine.Direction{
	CommandToNorth: engine.DirectionNorth,
	CommandToEast:  engine.DirectionEast,
	CommandToSouth: engine.DirectionSouth,
	CommandToWest:  engine.DirectionWest,
}

// Snake object
type Snake struct {
	uuid string

	world *world.World

	location engine.Location
	length   uint16

	direction engine.Direction

	mux *sync.RWMutex
}

// NewSnake creates new snake
func NewSnake(world *world.World) (*Snake, error) {
	snake := newDefaultSnake(world)
	location, err := snake.locate()
	if err != nil {
		return nil, fmt.Errorf("cannot create snake: %s", err)
	}

	if snake.direction == engine.DirectionSouth || snake.direction == engine.DirectionEast {
		location = location.Reverse()
	}

	snake.setLocation(location)

	return snake, nil
}

func newDefaultSnake(world *world.World) *Snake {
	return &Snake{
		uuid:      uuid.Must(uuid.NewV4()).String(),
		world:     world,
		location:  make(engine.Location, snakeStartLength),
		length:    snakeStartLength,
		direction: engine.RandomDirection(),
		mux:       &sync.RWMutex{},
	}
}

func (s *Snake) locate() (engine.Location, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	switch s.direction {
	case engine.DirectionNorth, engine.DirectionSouth:
		return s.world.CreateObjectRandomRectMargin(s, 1, uint8(snakeStartLength), snakeStartMargin)
	case engine.DirectionEast, engine.DirectionWest:
		return s.world.CreateObjectRandomRectMargin(s, uint8(snakeStartLength), 1, snakeStartMargin)
	}
	return nil, errors.New("invalid direction")
}

func (s *Snake) setLocation(location engine.Location) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.location = location
}

func (s *Snake) GetUUID() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.uuid
}

func (s *Snake) setDirection(dir engine.Direction) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.direction = dir
}

func (s *Snake) String() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return fmt.Sprintf("snake %s", s.location)
}

func (s *Snake) Die() {
	s.mux.RLock()
	s.world.DeleteObject(s, engine.Location(s.location))
	s.mux.RUnlock()
}

func (s *Snake) feed(f uint16) {
	if f > 0 {
		s.mux.Lock()
		defer s.mux.Unlock()
		s.length += f
	}
}

func (s *Snake) strength() float32 {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return snakeStrengthFactor * float32(s.length)
}

func (s *Snake) Run(stop <-chan struct{}) <-chan struct{} {
	snakeStop := make(chan struct{})

	go func() {
		var ticker = time.NewTicker(s.calculateDelay())
		defer ticker.Stop()
		defer close(snakeStop)

		for {
			select {
			case <-ticker.C:
				if err := s.move(); err != nil {
					// TODO: Handle error.
					return
				}
			case <-stop:
				return
			}
		}
	}()

	return snakeStop
}

func (s *Snake) move() error {
	// Calculate next position
	dot, err := s.getNextHeadDot()
	if err != nil {
		return err
	}

	if object := s.world.GetObjectByDot(dot); object != nil {
		if food, ok := object.(objects.Food); ok {
			s.feed(food.NutritionalValue(dot))
		} else {
			s.Die()

			return errors.New("snake dies")
		}

		// TODO: Reload ticker.
		//ticker = time.NewTicker(s.calculateDelay())
	}

	s.mux.RLock()
	tmpLocation := make(engine.Location, len(s.location)+1)
	copy(tmpLocation[1:], s.location)
	s.mux.RUnlock()
	tmpLocation[0] = dot

	if s.length < uint16(len(tmpLocation)) {
		tmpLocation = tmpLocation[:len(tmpLocation)-1]
	}

	if err := s.world.UpdateObject(s, engine.Location(s.location), tmpLocation); err != nil {
		return fmt.Errorf("update snake error: %s", err)
	}

	s.setLocation(tmpLocation)

	return nil
}

func (s *Snake) calculateDelay() time.Duration {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return time.Duration(math.Pow(snakeSpeedFactor, float64(s.length)) * float64(snakeStartSpeed))
}

// getNextHeadDot calculates new position of snake's head by its direction and current head position
func (s *Snake) getNextHeadDot() (engine.Dot, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.location) > 0 {
		return s.world.Navigate(s.location[0], s.direction, 1)
	}

	return engine.Dot{}, errors.New("cannot get next head dots: empty location")
}

func (s *Snake) Command(cmd Command) error {
	if direction, ok := snakeCommands[cmd]; ok {
		return fmt.Errorf("cannot execute command: %s", s.setMovementDirection(direction))
	}

	return errors.New("cannot execute command: unknown command")
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
		}

		s.setDirection(nextDir)

		return nil
	}

	return errors.New("invalid direction")
}

func (s *Snake) GetLocation() engine.Location {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return engine.Location(s.location).Copy()
}

func (s *Snake) MarshalJSON() ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return ffjson.Marshal(&snake{
		UUID: s.uuid,
		Dots: s.location,
		Type: snakeTypeLabel,
	})
}

type snake struct {
	UUID string       `json:"uuid"`
	Dots []engine.Dot `json:"dots"`
	Type string       `json:"type"`
}
