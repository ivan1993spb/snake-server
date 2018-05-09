package connections

import "encoding/json"

type InputMessageType uint8

const (
	InputMessageTypeSnakeCommand InputMessageType = iota
)

type InputMessage struct {
	Type    InputMessageType
	Payload json.RawMessage
}
