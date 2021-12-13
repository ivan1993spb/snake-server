package mouse

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

const mouseTypeLabel = "mouse"

type Mouse struct {
	id world.Identifier

	dot       engine.Dot
	direction engine.Direction

	world world.Interface
	mux   *sync.RWMutex

	once sync.Once
	stop chan struct{}
}

type errCreateMouse string

func (e errCreateMouse) Error() string {
	return "cannot create mouse: " + string(e)
}

func NewMouse(world world.Interface) (*Mouse, error) {
	mouse := &Mouse{
		id: world.IdentifierRegistry().Obtain(),

		world: world,
		mux:   &sync.RWMutex{},

		stop: make(chan struct{}),
	}

	mouse.mux.Lock()
	defer mouse.mux.Unlock()

	location, err := world.CreateObjectRandomDot(mouse)
	if err != nil {
		world.IdentifierRegistry().Release(mouse.id)

		return nil, errCreateMouse(err.Error())
	}

	if location.Empty() {
		world.IdentifierRegistry().Release(mouse.id)

		if err := world.DeleteObject(mouse, location); err != nil {
			return nil, errCreateMouse("no location located and cannot delete mouse")
		}
		return nil, errCreateMouse("no location located")
	}

	mouse.dot = location.Dot(0)
	mouse.direction = engine.RandomDirection()

	return mouse, nil
}

const mouseNutritionalValue uint16 = 15

type errMouseBite string

func (e errMouseBite) Error() string {
	return "mouse bite error: " + string(e)
}

func (m *Mouse) Bite(dot engine.Dot) (nv uint16, success bool, err error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if m.dot.Equals(dot) {
		m.die()

		if err := m.world.DeleteObject(m, engine.Location{m.dot}); err != nil {
			return 0, false, errMouseBite(err.Error())
		}
		return mouseNutritionalValue, true, nil
	}

	return 0, false, errMouseBite("mouse does not contain dot")
}

func (m *Mouse) String() string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return fmt.Sprintf("mouse %s", m.dot)
}

const mouseMarshalBufferSize = 72

func (m *Mouse) MarshalJSON() ([]byte, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	buff := bytes.NewBuffer(make([]byte, 0, mouseMarshalBufferSize))
	buff.WriteString(`{"type":"`)
	buff.WriteString(mouseTypeLabel)
	buff.WriteString(`","id":`)
	buff.WriteString(m.id.String())
	buff.WriteString(`,"dot":`)
	if dotJSON, err := m.dot.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buff.Write(dotJSON)
	}
	buff.WriteString(`,"direction":`)
	if directionJSON, err := m.direction.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buff.Write(directionJSON)
	}
	buff.WriteByte('}')

	return buff.Bytes(), nil
}

func (m *Mouse) die() {
	m.once.Do(func() {
		close(m.stop)
	})
}

const mouseStepDistance = 1

func (m *Mouse) move() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	dir := engine.RandomDirection()
	dot, err := m.world.Area().Navigate(m.dot, dir, mouseStepDistance)
	if err != nil {
		return err
	}

	if object := m.world.GetObjectByDot(dot); object != nil {
		return fmt.Errorf("dot is occupied")
	}

	if err := m.world.UpdateObject(m, engine.Location{m.dot}, engine.Location{dot}); err != nil {
		return err
	}

	m.direction = dir
	m.dot = dot
	return nil
}

const (
	mouseTickDurationMin = time.Second
	mouseTickDurationMax = time.Second * 3
)

func genMouseTickDuration() time.Duration {
	return mouseTickDurationMin + time.Duration(rand.Int63n(int64(mouseTickDurationMax-mouseTickDurationMin)))
}

func (m *Mouse) Run(stop <-chan struct{}) {
	var ticker = time.NewTicker(genMouseTickDuration())

	go func() {
		select {
		case <-stop:
			m.die()
		case <-m.stop:
		}

		ticker.Stop()
		m.world.IdentifierRegistry().Release(m.id)
	}()

	go func() {
		for {
			select {
			case <-m.stop:
				return
			case <-ticker.C:
				m.move()
			}
		}
	}()
}
