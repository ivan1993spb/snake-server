package world

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
