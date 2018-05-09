package connections

import "encoding/json"

type InputMessageType uint8

const (
	InputMessageTypeGameCommand InputMessageType = iota
)

type InputMessage struct {
	Type    InputMessageType
	Payload json.RawMessage
}
