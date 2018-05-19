package game

import "github.com/ivan1993spb/snake-server/world"

type EventType uint8

const (
	EventTypeError EventType = iota
	EventTypeObjectCreate
	EventTypeObjectDelete
	EventTypeObjectUpdate
	EventTypeObjectChecked
)

var eventsLabels = map[EventType]string{
	EventTypeError:         "error",
	EventTypeObjectCreate:  "create",
	EventTypeObjectDelete:  "delete",
	EventTypeObjectUpdate:  "update",
	EventTypeObjectChecked: "checked",
}

func (event EventType) String() string {
	if label, ok := eventsLabels[event]; ok {
		return label
	}
	return "unknown"
}

type Event struct {
	Type    EventType
	Payload interface{}
}

var eventTypesCasting = map[world.EventType]EventType{
	world.EventTypeError:         EventTypeError,
	world.EventTypeObjectCreate:  EventTypeObjectCreate,
	world.EventTypeObjectDelete:  EventTypeObjectDelete,
	world.EventTypeObjectUpdate:  EventTypeObjectUpdate,
	world.EventTypeObjectChecked: EventTypeObjectChecked,
}

func worldEventTypeToGameEventType(worldEventType world.EventType) EventType {
	if gameEventType, ok := eventTypesCasting[worldEventType]; ok {
		return gameEventType
	}
	return 0
}
