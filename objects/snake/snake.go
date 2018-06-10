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
	uuid uuid.UUID

	world *world.World

	dots   []engine.Dot
	length uint16

	direction engine.Direction

	mux *sync.RWMutex
}

// NewSnake creates new snake
func NewSnake(world *world.World) (*Snake, error) {
	var (
		err      error
		location engine.Location
	)

	snake := newDefaultSnake(world)
	dir := engine.RandomDirection()

	switch dir {
	case engine.DirectionNorth, engine.DirectionSouth:
		location, err = world.CreateObjectRandomRect(snake, 1, uint8(snakeStartLength))
	case engine.DirectionEast, engine.DirectionWest:
		location, err = world.CreateObjectRandomRect(snake, uint8(snakeStartLength), 1)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot create snake: %s", err)
	}

	if dir == engine.DirectionSouth || dir == engine.DirectionEast {
		// TODO: Reverse?
		reversedDots := location.Reverse()
		location = reversedDots
	}

	snake.relocate([]engine.Dot(location))
	snake.setDirection(dir)

	return snake, nil
}

func newDefaultSnake(world *world.World) *Snake {
	return &Snake{
		uuid:      uuid.Must(uuid.NewV4()),
		world:     world,
		dots:      make([]engine.Dot, snakeStartLength),
		length:    snakeStartLength,
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}
}

func (s *Snake) relocate(dots []engine.Dot) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.dots = dots
}

func (s *Snake) GetUUID() uuid.UUID {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.uuid
}

func (s *Snake) setDirection(dir engine.Direction) {
	s.mux.Lock()
	defer s.mux.Unlock()
	// TODO: Use atomic.
	s.direction = dir
}

func (s *Snake) String() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return fmt.Sprintf("snake %s", s.dots)
}

func (s *Snake) Die() {
	s.mux.RLock()
	s.world.DeleteObject(s, engine.Location(s.dots))
	s.mux.RUnlock()
	//corpse.NewCorpse(s.world, s.dots)
}

func (s *Snake) feed(f uint16) {
	if f > 0 {
		s.mux.Lock()
		defer s.mux.Unlock()
		// TODO: Use atomic.
		s.length += f
	}
}

func (s *Snake) strength() float32 {
	s.mux.Lock()
	defer s.mux.Unlock()
	// TODO: Use atomic.

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
		// TODO How to emit error ?
		//s.p.OccurredError(s, err)
		return err
	}

	if object := s.world.GetObjectByDot(dot); object != nil {
		if food, ok := object.(objects.Food); ok {
			s.feed(food.NutritionalValue(dot))
		} else {
			s.Die()
			return errors.New("die collusion")
		}

		// TODO: Reload ticker.
		//ticker = time.NewTicker(s.calculateDelay())
	}

	s.mux.RLock()
	tmpDots := make([]engine.Dot, len(s.dots)+1)
	copy(tmpDots[1:], s.dots)
	s.mux.RUnlock()
	tmpDots[0] = dot

	if s.length < uint16(len(tmpDots)) {
		tmpDots = tmpDots[:len(tmpDots)-1]
	}

	// TODO: Handle error.
	if err := s.world.UpdateObject(s, engine.Location(s.dots), tmpDots); err != nil {
		return fmt.Errorf("update snake error: %s", err)
	}

	s.relocate(tmpDots)

	return nil
}

func (s *Snake) calculateDelay() time.Duration {
	s.mux.RLock()
	defer s.mux.RUnlock()
	// TODO: Use atomic.
	return time.Duration(math.Pow(snakeSpeedFactor, float64(s.length)) * float64(snakeStartSpeed))
}

// getNextHeadDot calculates new position of snake's head by its
// direction and current head position
func (s *Snake) getNextHeadDot() (engine.Dot, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.dots) > 0 {
		return s.world.Navigate(s.dots[0], s.direction, 1)
	}

	return engine.Dot{}, fmt.Errorf("cannot get next head dots: errEmptyDotList")
}

func (s *Snake) Command(cmd Command) error {
	if direction, ok := snakeCommands[cmd]; ok {
		// TODO: Handle err.
		s.setMovementDirection(direction)
		return nil
	}
	return errors.New("cannot execute command")
}

func (s *Snake) setMovementDirection(nextDir engine.Direction) error {
	if engine.ValidDirection(nextDir) {
		currDir := engine.CalculateDirection(s.dots[1], s.dots[0])
		rNextDir, err := nextDir.Reverse()
		if err != nil {
			return fmt.Errorf("cannot set movement direction: %s", err)
		}

		// Next direction cannot be opposite to current direction
		if rNextDir == currDir {
			return errors.New("next direction cannot be opposite to current direction")
		} else {
			s.setDirection(nextDir)
			return nil
		}
	}

	return errors.New("invalid direction")
}

func (s *Snake) GetLocation() engine.Location {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return engine.Location(s.dots).Copy()
}

func (s *Snake) MarshalJSON() ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return ffjson.Marshal(&snake{
		UUID: s.uuid.String(),
		Dots: s.dots,
		Type: "snake",
	})
}

type snake struct {
	UUID string       `json:"uuid"`
	Dots []engine.Dot `json:"dots"`
	Type string       `json:"type"`
}
