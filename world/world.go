package world

import (
	"fmt"
	"sync"
	"time"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

const (
	worldEventsChanMainBufferSize  = 4096
	worldEventsChanProxyBufferSize = 4096
)

const worldEventsSendTimeout = time.Millisecond

type World struct {
	pg          *playground.Playground
	chMain      chan Event
	chsProxy    []chan Event
	chsProxyMux *sync.RWMutex
	stopGlobal  chan struct{}
	flagStarted bool
	startedMux  *sync.Mutex

	identifierRegistry *IdentifierRegistry
}

func NewWorld(width, height uint8) (*World, error) {
	pg, err := playground.NewPlayground(width, height)
	if err != nil {
		return nil, fmt.Errorf("cannot create world: %s", err)
	}

	return &World{
		pg:          pg,
		chMain:      make(chan Event, worldEventsChanMainBufferSize),
		chsProxy:    make([]chan Event, 0),
		chsProxyMux: &sync.RWMutex{},
		stopGlobal:  make(chan struct{}),

		flagStarted: false,
		startedMux:  &sync.Mutex{},

		identifierRegistry: NewIdentifierRegistry(),
	}, nil
}

func (w *World) event(event Event) {
	select {
	case w.chMain <- event:
	case <-w.stopGlobal:
	}
}

func (w *World) Start(stop <-chan struct{}) {
	w.startedMux.Lock()
	defer w.startedMux.Unlock()

	if w.flagStarted {
		return
	}
	w.flagStarted = true

	go func() {
		select {
		case <-stop:
		}
		w.stop()
	}()

	go func() {
		for {
			select {
			case event, ok := <-w.chMain:
				if !ok {
					return
				}
				w.broadcast(event)
			case <-w.stopGlobal:
				return
			}
		}
	}()
}

func (w *World) broadcast(event Event) {
	w.chsProxyMux.RLock()
	defer w.chsProxyMux.RUnlock()

	for _, chProxy := range w.chsProxy {
		select {
		case chProxy <- event:
		case <-w.stopGlobal:
		}
	}
}

func (w *World) createChanProxy() chan Event {
	chProxy := make(chan Event, worldEventsChanProxyBufferSize)

	w.chsProxyMux.Lock()
	w.chsProxy = append(w.chsProxy, chProxy)
	w.chsProxyMux.Unlock()

	return chProxy
}

func (w *World) deleteChanProxy(chProxy chan Event) {
	go func() {
		for range chProxy {
		}
	}()

	w.chsProxyMux.Lock()
	for i := range w.chsProxy {
		if w.chsProxy[i] == chProxy {
			w.chsProxy = append(w.chsProxy[:i], w.chsProxy[i+1:]...)
			close(chProxy)
			break
		}
	}
	w.chsProxyMux.Unlock()
}

func (w *World) Events(stop <-chan struct{}, buffer uint) <-chan Event {
	chProxy := w.createChanProxy()
	chOut := make(chan Event, buffer)

	go func() {
		defer close(chOut)
		defer w.deleteChanProxy(chProxy)

		for {
			select {
			case <-stop:
				return
			case <-w.stopGlobal:
				return
			case event := <-chProxy:
				w.sendEvent(chOut, event, stop)
			}
		}
	}()

	return chOut
}

func (w *World) sendEvent(ch chan Event, event Event, stop <-chan struct{}) {
	if event.Type == EventTypeObjectUpdate || event.Type == EventTypeObjectChecked {
		w.sendEventTimeout(ch, event, stop, worldEventsSendTimeout)
	} else {
		w.sendEventStrict(ch, event, stop)
	}
}

func (w *World) sendEventStrict(ch chan Event, event Event, stop <-chan struct{}) {
	select {
	case ch <- event:
	case <-w.stopGlobal:
	case <-stop:
	}
}

func (w *World) sendEventTimeout(ch chan Event, event Event, stop <-chan struct{}, timeout time.Duration) {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- event:
	case <-w.stopGlobal:
	case <-stop:
	case <-timer.C:
	}
}

func (w *World) stop() {
	close(w.stopGlobal)
	close(w.chMain)

	w.chsProxyMux.Lock()
	defer w.chsProxyMux.Unlock()

	for _, ch := range w.chsProxy {
		close(ch)
	}

	w.chsProxy = w.chsProxy[:0]
}

func (w *World) GetObjectByDot(dot engine.Dot) interface{} {
	if object := w.pg.GetObjectByDot(dot); object != nil {
		w.event(Event{
			Type:    EventTypeObjectChecked,
			Payload: object,
		})
		return object
	}
	return nil
}

func (w *World) GetObjectsByDots(dots []engine.Dot) []interface{} {
	if objects := w.pg.GetObjectsByDots(dots); len(objects) > 0 {
		for _, object := range objects {
			w.event(Event{
				Type:    EventTypeObjectChecked,
				Payload: object,
			})
		}
		return objects
	}
	return nil
}

func (w *World) CreateObject(object interface{}, location engine.Location) error {
	if err := w.pg.CreateObject(object, location); err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return nil
}

func (w *World) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, error) {
	location, err := w.pg.CreateObjectAvailableDots(object, location)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, nil
}

func (w *World) DeleteObject(object interface{}, location engine.Location) error {
	err := w.pg.DeleteObject(object, location)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(Event{
		Type:    EventTypeObjectDelete,
		Payload: object,
	})
	return nil
}

func (w *World) UpdateObject(object interface{}, old, new engine.Location) error {
	if err := w.pg.UpdateObject(object, old, new); err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(Event{
		Type:    EventTypeObjectUpdate,
		Payload: object,
	})
	return nil
}

func (w *World) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, error) {
	location, err := w.pg.UpdateObjectAvailableDots(object, old, new)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectUpdate,
		Payload: object,
	})
	return location, nil
}

func (w *World) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomDot(object)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, nil
}

func (w *World) CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomRect(object, rw, rh)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, nil
}

func (w *World) CreateObjectRandomRectMargin(object interface{}, rw, rh, margin uint8) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomRectMargin(object, rw, rh, margin)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, nil
}

func (w *World) CreateObjectRandomByDotsMask(object interface{}, dm *engine.DotsMask) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomByDotsMask(object, dm)
	if err != nil {
		w.event(Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, nil
}

func (w *World) LocationOccupied(location engine.Location) bool {
	return w.pg.LocationOccupied(location)
}

func (w *World) Size() uint16 {
	return w.pg.Size()
}

func (w *World) Width() uint8 {
	return w.pg.Width()
}

func (w *World) Height() uint8 {
	return w.pg.Height()
}

func (w *World) Area() engine.Area {
	return w.pg.Area()
}

func (w *World) GetObjects() []interface{} {
	return w.pg.GetObjects()
}

func (w *World) ObtainIdentifier() Identifier {
	return w.identifierRegistry.Obtain()
}

func (w *World) ReleaseIdentifier(id Identifier) {
	w.identifierRegistry.Release(id)
}
