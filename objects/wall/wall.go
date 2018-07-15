package wall

import (
	"fmt"
	"sync"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/satori/go.uuid"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const wallTypeLabel = "wall"

const wallMinBreakForce = 1000

type Wall struct {
	uuid     string
	world    *world.World
	location engine.Location
	mux      *sync.RWMutex
}

type ErrCreateWall string

func (e ErrCreateWall) Error() string {
	return "cannot create wall: " + string(e)
}

func NewWall(world *world.World, dm *engine.DotsMask) (*Wall, error) {
	wall := &Wall{
		uuid:  uuid.Must(uuid.NewV4()).String(),
		world: world,
		mux:   &sync.RWMutex{},
	}

	location, err := world.CreateObjectRandomByDotsMask(wall, dm)
	if err != nil {
		return nil, ErrCreateWall(err.Error())
	}

	wall.mux.Lock()
	wall.location = location
	wall.mux.Unlock()

	return wall, nil
}

func NewWallLocation(world *world.World, location engine.Location) (*Wall, error) {
	wall := &Wall{
		uuid:  uuid.Must(uuid.NewV4()).String(),
		world: world,
		mux:   &sync.RWMutex{},
	}

	wall.mux.Lock()
	defer wall.mux.Unlock()

	location, err := world.CreateObjectAvailableDots(wall, location)
	if err != nil {
		return nil, ErrCreateWall(err.Error())
	}

	wall.location = location

	return wall, nil
}

type errWallBreak string

func (e errWallBreak) Error() string {
	return "wall break error: " + string(e)
}

func (w *Wall) Break(dot engine.Dot, force float32) (success bool, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.location.Contains(dot) {
		if force < wallMinBreakForce {
			return false, nil
		}

		location := w.location.Delete(dot)

		if location.DotCount() > 0 {
			if err := w.world.UpdateObject(w, w.location, location); err != nil {
				return false, errWallBreak(err.Error())
			} else {
				w.location = location
			}
		} else {
			if err := w.world.DeleteObject(w, w.location); err != nil {
				return false, errWallBreak(err.Error())
			}

			w.location = w.location[:0]
		}

		return true, nil
	}

	return false, errWallBreak("wall does not contain dot")
}

func (w *Wall) String() string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return fmt.Sprintf("wall %d", len(w.location))
}

func (w *Wall) MarshalJSON() ([]byte, error) {
	w.mux.RLock()
	defer w.mux.RUnlock()
	return ffjson.Marshal(&wall{
		UUID: w.uuid,
		Dots: w.location,
		Type: wallTypeLabel,
	})
}

type wall struct {
	UUID string          `json:"uuid"`
	Dots engine.Location `json:"dots,omitempty"`
	Type string          `json:"type"`
}
