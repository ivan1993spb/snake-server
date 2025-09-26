package world

type EventType uint8

const (
	EventTypeError EventType = iota
	EventTypeObjectCreate
	EventTypeObjectDelete
	EventTypeObjectUpdate
	EventTypeObjectUpdateV2
	EventTypeObjectChecked
)

var eventsLabels = map[EventType]string{
	EventTypeError:          "error",
	EventTypeObjectCreate:   "create",
	EventTypeObjectDelete:   "delete",
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
	Type    EventType
	Payload interface{}
}
