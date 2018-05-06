package game

type EventType uint8

const (
	EventTypeError EventType = iota
	EventTypeObjectCreate
	EventTypeObjectDelete
	EventTypeObjectUpdate
	EventTypeObjectChecked
)

type Event struct {
	Type    EventType
	Payload interface{}
}
