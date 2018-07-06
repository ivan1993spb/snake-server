package snake

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/objects"
	"github.com/ivan1993spb/snake-server/world"
)

const (
	snakeTypeLabel = "snake"

	snakeStartSpeed  = time.Millisecond * 500
	snakeSpeedFactor = 1

	snakeStartLength = 3
	snakeStartMargin = 1

	snakeMaxInteractionRetries = 5

	snakeForceBaby  = 1
	snakeForceAdult = 2
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

	stopper *sync.Once
	stop    chan struct{}
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
		stopper:   &sync.Once{},
		stop:      make(chan struct{}),
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

func (s *Snake) die() error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if err := s.world.DeleteObject(s, s.location); err != nil {
		return fmt.Errorf("die snake error: %s", err)
	}

	// Do not empty location to pass it for corpse creation.

	return nil
}

func (s *Snake) feed(f uint16) {
	if f > 0 {
		s.mux.Lock()
		defer s.mux.Unlock()
		s.length += f
	}
}

type errSnakeHit string

func (e errSnakeHit) Error() string {
	return "snake hit error: " + string(e)
}

func (s *Snake) Hit(dot engine.Dot, force float32) (success bool, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.location.Contains(dot) {
		if force > s.unsafeGetForce() {
			s.stopper.Do(func() {
				close(s.stop)
			})

			newLocation := s.location.Delete(dot)
			if err := s.world.UpdateObject(s, s.location, newLocation); err != nil {
				return false, errSnakeHit(err.Error())
			}

			s.location = newLocation

			return true, nil
		}

		return false, nil
	}

	return false, errSnakeHit("snake does not contain dot")
}

func (s *Snake) unsafeGetForce() float32 {
	if s.length > snakeStartLength {
		return snakeForceAdult
	}
	return snakeForceBaby
}

func (s *Snake) getForce() float32 {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.unsafeGetForce()
}

func (s *Snake) Run(stop <-chan struct{}, logger logrus.FieldLogger) <-chan struct{} {
	snakeStop := make(chan struct{})
	logger = logger.WithField("uuid", s.uuid)

	go func() {
		var ticker = time.NewTicker(s.calculateDelay())
		defer ticker.Stop()
		defer close(snakeStop)
		defer func() {
			if err := s.die(); err != nil {
				logger.WithError(err).Error("die snake error")
			}
		}()
		defer s.stopper.Do(func() {
			close(s.stop)
		})

		for {
			select {
			case <-ticker.C:
				if err := s.move(); err != nil {
					if err != errUnsuccessfulInteraction {
						logger.WithError(err).Error("snake move error")
					}
					return
				}
			case <-stop:
				// Global stop
				return
			case <-s.stop:
				// Local snake stop
				return
			}
		}
	}()

	return snakeStop
}

type errSnakeMove string

func (e errSnakeMove) Error() string {
	return "move snake error: " + string(e)
}

var errUnsuccessfulInteraction = errSnakeMove("unsuccessful interaction")

func (s *Snake) move() error {
	// Calculate next position
	dot, err := s.getNextHeadDot()
	if err != nil {
		return errSnakeMove(err.Error())
	}

	retries := 0

	for {
		if object := s.world.GetObjectByDot(dot); object != nil {
			if success, err := s.interactObject(object, dot); err != nil {
				return errSnakeMove(err.Error())
			} else if !success {
				return errUnsuccessfulInteraction
			}
		} else {
			break
		}

		if retries >= snakeMaxInteractionRetries {
			return errSnakeMove("interaction retries limit reached")
		}

		retries++
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

type errInteractObject string

func (e errInteractObject) Error() string {
	return "object interaction error: " + string(e)
}

func (s *Snake) interactObject(object interface{}, dot engine.Dot) (success bool, err error) {
	if food, ok := object.(objects.Food); ok {
		nv, success, err := food.Bite(dot)
		if err != nil {
			return false, errInteractObject(err.Error())
		}
		if success {
			s.feed(nv)
		}
		return success, nil
	}

	if alive, ok := object.(objects.Alive); ok {
		success, err := alive.Hit(dot, s.getForce())
		if err != nil {
			return false, errInteractObject(err.Error())
		}
		return success, nil
	}

	if object, ok := object.(objects.Object); ok {
		success, err := object.Break(dot, s.getForce())
		if err != nil {
			return false, errInteractObject(err.Error())
		}
		return success, nil
	}

	return false, nil
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
		if err := s.setMovementDirection(direction); err != nil {
			return fmt.Errorf("cannot execute command: %s", err)
		}
		return nil
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
	Dots []engine.Dot `json:"dots,omitempty"`
	Type string       `json:"type"`
}
