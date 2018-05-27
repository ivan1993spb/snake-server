package apple

import (
	"fmt"
	"sync"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

type Apple struct {
	world    *world.World
	location engine.Location
	mux      *sync.RWMutex
}

type ErrCreateApple string

func (e ErrCreateApple) Error() string {
	return "cannot create apple: " + string(e)
}

// NewApple creates and locates new apple
func NewApple(world *world.World) (*Apple, error) {
	apple := &Apple{
		mux: &sync.RWMutex{},
	}

	location, err := world.CreateObjectRandomDot(apple)
	if err != nil {
		return nil, ErrCreateApple(err.Error())
	}

	apple.mux.Lock()
	apple.location = location
	apple.world = world
	apple.mux.Unlock()

	return apple, nil
}

func (a *Apple) String() string {
	a.mux.RLock()
	defer a.mux.RUnlock()
	return fmt.Sprintf("apple %s", a.location)
}

func (a *Apple) NutritionalValue(dot engine.Dot) uint16 {
	a.mux.RLock()
	defer a.mux.RUnlock()

	if a.location.Equals(engine.Location{dot}) {
		// TODO: Handle error.
		a.world.DeleteObject(a, a.location)
		return 1
	}

	return 0
}

func (a *Apple) MarshalJSON() ([]byte, error) {
	a.mux.RLock()
	defer a.mux.RUnlock()
	return ffjson.Marshal(a.location)
}
