package game

type EventType uint8

const (
	EventTypeError EventType = iota
	EventTypeObjectCreate
	EventTypeObjectDelete
	EventTypeObjectUpdate
)

type Event struct {
	Type    EventType
	Payload interface{}
}
