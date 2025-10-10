package apple

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const appleTypeLabel = "apple"

type Apple struct {
	id    world.Identifier
	world world.Interface
	dot   engine.Dot
	mux   *sync.RWMutex
}

const appleNutritionalValue uint16 = 1

type errCreateApple string

func (e errCreateApple) Error() string {
	return "cannot create apple: " + string(e)
}

// NewApple creates and locates new apple
func NewApple(w world.Interface) (*Apple, error) {
	apple := &Apple{
		id:  w.IdentifierRegistry().Obtain(),
		mux: &sync.RWMutex{},
	}

	apple.mux.Lock()
	defer apple.mux.Unlock()

	location, err := w.CreateObjectRandomDot(apple)
	if err != nil {
		w.IdentifierRegistry().Release(apple.id)

		return nil, errCreateApple(err.Error())
	}

	if location.Empty() {
		w.IdentifierRegistry().Release(apple.id)

		if err := w.DeleteObject(apple, location); err != nil {
			return nil, errCreateApple("no location located and cannot delete apple")
		}
		return nil, errCreateApple("no location located")
	}

	apple.dot = location.Dot(0)
	apple.world = w

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
		a.world.IdentifierRegistry().Release(a.id)
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
	return json.Marshal(&apple{
		ID:   a.id,
		Dot:  a.dot,
		Type: appleTypeLabel,
	})
}

type apple struct {
	ID   world.Identifier `json:"id"`
	Dot  engine.Dot       `json:"dot"`
	Type string           `json:"type"`
}
