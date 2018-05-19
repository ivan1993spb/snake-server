package snake

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/objects"
	"github.com/ivan1993spb/snake-server/objects/corpse"
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
	id    uint64
	world *world.World

	dots   []engine.Dot
	length uint16

	mux *sync.RWMutex

	// Motion direction
	direction engine.Direction
}

// NewSnake creates new snake
func NewSnake(world *world.World) (*Snake, error) {
	var (
		dir      = engine.RandomDirection()
		err      error
		location engine.Location
	)

	snake := &Snake{}

	switch dir {
	case engine.DirectionNorth, engine.DirectionSouth:
		location, err = world.CreateObjectRandomRect(snake, 1, uint8(snakeStartLength))
	case engine.DirectionEast, engine.DirectionWest:
		location, err = world.CreateObjectRandomRect(snake, uint8(snakeStartLength), 1)
	}
	if err != nil {
		// TODO: Create error.
		return nil, err
	}

	if dir == engine.DirectionSouth || dir == engine.DirectionEast {
		// TODO: Reverse?
		reversedDots := location.Reverse()
		location = reversedDots
	}

	snake.id = *(*uint64)(unsafe.Pointer(&snake))
	snake.world = world
	snake.dots = []engine.Dot(location)
	snake.length = snakeStartLength
	snake.direction = dir
	snake.mux = &sync.RWMutex{}

	return snake, nil
}

func (s *Snake) String() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return fmt.Sprintf("snake %s", s.dots)
}

func (s *Snake) Die() {
	s.mux.Lock()
	s.world.DeleteObject(s, engine.Location(s.dots))
	s.mux.Unlock()
	corpse.NewCorpse(s.world, s.dots)
}

func (s *Snake) Feed(f int8) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if f > 0 {
		s.length += uint16(f)
	}
}

func (s *Snake) Strength() float32 {
	s.mux.Lock()
	defer s.mux.Unlock()
	return snakeStrengthFactor * float32(s.length)
}

func (s *Snake) Run(stop <-chan struct{}) {
	go func() {
		var ticker = time.NewTicker(s.calculateDelay())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.move()
			case <-stop:
				return
			}
		}
	}()
}

func (s *Snake) move() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Calculate next position
	dot, err := s.getNextHeadDot()
	if err != nil {
		// TODO How to emit error ?
		//s.p.OccurredError(s, err)
		return err
	}

	if object := s.world.GetObjectByDot(dot); object != nil {
		// TODO: Use interfaces to interact objects.
		if food, ok := object.(objects.Food); ok {
			s.length += food.NutritionalValue(dot)
		} else {
			s.Die()
			return nil
		}

		// TODO: Reload ticker.
		//ticker = time.NewTicker(s.calculateDelay())
	}

	tmpDots := make([]engine.Dot, len(s.dots)+1)
	copy(tmpDots[1:], s.dots)
	tmpDots[0] = dot

	if s.length < uint16(len(tmpDots)) {
		tmpDots = tmpDots[:len(tmpDots)-1]
	}

	// TODO: Handle error.
	if err := s.world.UpdateObject(s, engine.Location(s.dots), tmpDots); err != nil {
		return err
	}

	s.dots = tmpDots

	return nil
}

func (s *Snake) calculateDelay() time.Duration {
	return time.Duration(math.Pow(snakeSpeedFactor, float64(s.length)) * float64(snakeStartSpeed))
}

// getNextHeadDot calculates new position of snake's head by its
// direction and current head position
func (s *Snake) getNextHeadDot() (engine.Dot, error) {
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
			s.direction = nextDir
			return nil
		}
	}

	return errors.New("invalid direction")
}

func (s *Snake) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(&snake{
		ID:   s.id,
		Dots: s.dots,
	})
}

type snake struct {
	ID   uint64       `json:"id"`
	Dots []engine.Dot `json:"dots"`
}
