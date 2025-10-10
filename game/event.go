package game

import "github.com/ivan1993spb/snake-server/world"

type EventType uint8

const (
	EventTypeError EventType = iota
	EventTypeObjectCreate
	EventTypeObjectDelete
	EventTypeObjectDeleteV2
	EventTypeObjectUpdate
	EventTypeObjectUpdateV2
	EventTypeObjectChecked
)

var eventsLabels = map[EventType]string{
	EventTypeError:          "error",
	EventTypeObjectCreate:   "create",
	EventTypeObjectDelete:   "delete",
	EventTypeObjectDeleteV2: "delete_v2",
	EventTypeObjectUpdate:   "update",
	EventTypeObjectUpdateV2: "update_v2",
	EventTypeObjectChecked:  "checked",
}

func (event EventType) String() string {
	if label, ok := eventsLabels[event]; ok {
		return label
	}
	return "unknown"
}

var eventTypesJSONs = map[EventType][]byte{
	EventTypeError:          []byte(`"error"`),
	EventTypeObjectCreate:   []byte(`"create"`),
	EventTypeObjectDelete:   []byte(`"delete"`),
	EventTypeObjectDeleteV2: []byte(`"delete_v2"`),
	EventTypeObjectUpdate:   []byte(`"update"`),
	EventTypeObjectUpdateV2: []byte(`"update_v2"`),
	EventTypeObjectChecked:  []byte(`"checked"`),
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
	world.EventTypeError:          EventTypeError,
	world.EventTypeObjectCreate:   EventTypeObjectCreate,
	world.EventTypeObjectDelete:   EventTypeObjectDelete,
	world.EventTypeObjectDeleteV2: EventTypeObjectDeleteV2,
	world.EventTypeObjectUpdate:   EventTypeObjectUpdate,
	world.EventTypeObjectUpdateV2: EventTypeObjectUpdateV2,
	world.EventTypeObjectChecked:  EventTypeObjectChecked,
}

func worldEventTypeToGameEventType(worldEventType world.EventType) EventType {
	if gameEventType, ok := eventTypesCasting[worldEventType]; ok {
		return gameEventType
	}
	return 0
}
