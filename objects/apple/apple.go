package apple

import (
	"fmt"
	"sync"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const appleTypeLabel = "apple"

type Apple struct {
	uuid  string
	world *world.World
	dot   engine.Dot
	mux   *sync.RWMutex
}

type ErrCreateApple string

func (e ErrCreateApple) Error() string {
	return "cannot create apple: " + string(e)
}

// NewApple creates and locates new apple
func NewApple(world *world.World) (*Apple, error) {
	apple := &Apple{
		uuid: uuid.Must(uuid.NewV4()).String(),
		mux:  &sync.RWMutex{},
	}

	location, err := world.CreateObjectRandomDot(apple)
	if err != nil {
		return nil, ErrCreateApple(err.Error())
	}
	if len(location) == 0 {
		return nil, ErrCreateApple("created empty location")
	}

	apple.mux.Lock()
	apple.dot = location.Dot(0)
	apple.world = world
	apple.mux.Unlock()

	return apple, nil
}

func (a *Apple) String() string {
	a.mux.RLock()
	defer a.mux.RUnlock()
	return fmt.Sprintf("apple %s", a.dot)
}

func (a *Apple) NutritionalValue(dot engine.Dot) uint16 {
	a.mux.RLock()
	defer a.mux.RUnlock()

	if a.dot.Equals(dot) {
		// TODO: Handle error?
		a.world.DeleteObject(a, engine.Location{a.dot})
		return 1
	}

	return 0
}

func (a *Apple) MarshalJSON() ([]byte, error) {
	a.mux.RLock()
	defer a.mux.RUnlock()
	return ffjson.Marshal(&apple{
		UUID: a.uuid,
		Dot:  a.dot,
		Type: appleTypeLabel,
	})
}

type apple struct {
	UUID string     `json:"uuid"`
	Dot  engine.Dot `json:"dot"`
	Type string     `json:"type"`
}
