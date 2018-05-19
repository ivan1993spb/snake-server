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

var eventTypesJSONs = map[EventType][]byte{
	EventTypeError:         []byte(`"error"`),
	EventTypeObjectCreate:  []byte(`"create"`),
	EventTypeObjectDelete:  []byte(`"delete"`),
	EventTypeObjectUpdate:  []byte(`"update"`),
	EventTypeObjectChecked: []byte(`"checked"`),
}

func (event EventType) MarshalJSON() ([]byte, error) {
	if json, ok := eventTypesJSONs[event]; ok {
		return json, nil
	}
	return []byte(`"unknown"`), nil
}

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
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
