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

const appleNutritionalValue uint16 = 1

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

type errAppleBite string

func (e errAppleBite) Error() string {
	return "apple bite error: " + string(e)
}

func (a *Apple) Bite(dot engine.Dot) (nv uint16, success bool, err error) {
	a.mux.RLock()
	defer a.mux.RUnlock()

	if a.dot.Equals(dot) {
		if err := a.world.DeleteObject(a, engine.Location{a.dot}); err != nil {
			return 0, false, errAppleBite(err.Error())
		}
		return appleNutritionalValue, true, nil
	}

	return 0, false, errAppleBite("apple does not contain dot")
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
