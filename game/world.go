package game

import (
	"time"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/playground"
)

type World struct {
	pg       *playground.Playground
	chEvents chan *Event
	stop     chan struct{}
	timeout  time.Duration
}

func NewWorld(pg *playground.Playground, buffer uint, timeout time.Duration) *World {
	return &World{
		pg:       pg,
		chEvents: make(chan *Event, buffer),
		stop:     make(chan struct{}, 0),
		timeout:  timeout,
	}
}

// TODO: Create func ListenEvents with out channel <-chan *Event to provide multi listeners

func (w *World) event(event *Event) {
	go func() {
		var timer = time.NewTimer(w.timeout)
		defer timer.Stop()
		select {
		case w.chEvents <- event:
		case <-w.stop:
		case <-timer.C:
		}
	}()
}

func (w *World) Stop() {
	close(w.stop)
	close(w.chEvents)
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
	return w.pg.GetObjectByLocation(location)
}

func (w *World) GetObjectByDot(dot *engine.Dot) interface{} {
	return w.pg.GetObjectByDot(dot)
}

func (w *World) GetEntityByDot(dot *engine.Dot) (interface{}, engine.Location) {
	return w.pg.GetEntityByDot(dot)
}

func (w *World) GetObjectsByDots(dots []*engine.Dot) []interface{} {
	return w.pg.GetObjectsByDots(dots)
}

func (w *World) CreateObject(object interface{}, location engine.Location) error {
	if err := w.pg.CreateObject(object, location); err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(&Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return nil
}

func (w *World) CreateObjectAvailableDots(object interface{}, location engine.Location) (engine.Location, *playground.ErrCreateObjectAvailableDots) {
	location, err := w.pg.CreateObjectAvailableDots(object, location)
	if err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(&Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, err
}

func (w *World) DeleteObject(object interface{}, location engine.Location) *playground.ErrDeleteObject {
	err := w.pg.DeleteObject(object, location)
	if err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(&Event{
		Type:    EventTypeObjectDelete,
		Payload: object,
	})
	return err
}

func (w *World) UpdateObject(object interface{}, old, new engine.Location) *playground.ErrUpdateObject {
	if err := w.pg.UpdateObject(object, old, new); err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return err
	}
	w.event(&Event{
		Type:    EventTypeObjectUpdate,
		Payload: object,
	})
	return nil
}

func (w *World) UpdateObjectAvailableDots(object interface{}, old, new engine.Location) (engine.Location, *playground.ErrUpdateObjectAvailableDots) {
	location, err := w.pg.UpdateObjectAvailableDots(object, old, new)
	if err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(&Event{
		Type:    EventTypeObjectUpdate,
		Payload: object,
	})
	return location, err
}

func (w *World) CreateObjectRandomDot(object interface{}) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomDot(object)
	if err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(&Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, err
}

func (w *World) CreateObjectRandomRect(object interface{}, rw, rh uint8) (engine.Location, error) {
	location, err := w.pg.CreateObjectRandomRect(object, rw, rh)
	if err != nil {
		w.event(&Event{
			Type:    EventTypeError,
			Payload: err,
		})
		return nil, err
	}
	w.event(&Event{
		Type:    EventTypeObjectCreate,
		Payload: object,
	})
	return location, err
}

func (w *World) Navigate(dot *engine.Dot, dir engine.Direction, dis uint8) (*engine.Dot, error) {
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
