package world

import (
	"fmt"
	"sync"
	"time"

	playground "github.com/ivan1993spb/snake-server/cplayground"
	"github.com/ivan1993spb/snake-server/engine"
)

const (
	worldEventsChanMainBufferSize  = 512
	worldEventsChanProxyBufferSize = 128
)

const worldEventsSendTimeout = time.Millisecond * 50

type World struct {
	pg          *playground.Playground
	chMain      chan Event
	chsProxy    []chan Event
	chsProxyMux *sync.RWMutex
	stopGlobal  chan struct{}
	flagStarted bool
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
	}, nil
}

func (w *World) event(event Event) {
	select {
	case w.chMain <- event:
	case <-w.stopGlobal:
	}
}

func (w *World) Start(stop <-chan struct{}) {
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
				w.sendEvent(chOut, event, stop, worldEventsSendTimeout)
			}
		}
	}()

	return chOut
}

func (w *World) sendEvent(ch chan Event, event Event, stop <-chan struct{}, timeout time.Duration) {
	const tickSize = 5

	var timer = time.NewTimer(timeout)
	defer timer.Stop()

	var ticker = time.NewTicker(timeout / tickSize)
	defer ticker.Stop()

	if cap(ch) == 0 {
		select {
		case ch <- event:
		case <-w.stopGlobal:
		case <-stop:
		case <-timer.C:
		}
	} else {
		for {
			select {
			case ch <- event:
				return
			case <-w.stopGlobal:
				return
			case <-stop:
				return
			case <-timer.C:
				return
			case <-ticker.C:
				if len(ch) == cap(ch) {
					select {
					case <-ch:
					case ch <- event:
						return
					case <-w.stopGlobal:
						return
					case <-stop:
						return
					case <-timer.C:
						return
					}
				}
			}
		}
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

func (w *World) ObjectExists(object interface{}) bool {
	return w.pg.ObjectExists(object)
}

func (w *World) LocationExists(location engine.Location) bool {
	return w.pg.LocationExists(location)
}

func (w *World) EntityExists(object interface{}, location engine.Location) bool {
	return w.pg.EntityExists(object, location)
}

func (w *World) GetObjectByLocation(location engine.Location) interface{} {
	if object := w.pg.GetObjectByLocation(location); object != nil {
		w.event(Event{
			Type:    EventTypeObjectChecked,
			Payload: object,
		})
		return object
	}
	return nil

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

func (w *World) GetEntityByDot(dot engine.Dot) (interface{}, engine.Location) {
	if object, location := w.pg.GetEntityByDot(dot); object != nil && !location.Empty() {
		w.event(Event{
			Type:    EventTypeObjectChecked,
			Payload: object,
		})
		return object, location
	}
	return nil, nil
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
	return location, err
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
	return err
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
	return location, err
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
	return location, err
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
	return location, err
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
	return location, err
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
	return location, err
}

func (w *World) LocationOccupied(location engine.Location) bool {
	return w.pg.LocationOccupied(location)
}

func (w *World) Navigate(dot engine.Dot, dir engine.Direction, dis uint8) (engine.Dot, error) {
	return w.pg.Navigate(dot, dir, dis)
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

func (w *World) GetObjects() []interface{} {
	return w.pg.GetObjects()
}
