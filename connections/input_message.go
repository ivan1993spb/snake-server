package connections

import (
	"bytes"
	"errors"
)

type InputMessageType uint8

const (
	InputMessageTypeSnakeCommand InputMessageType = iota
	InputMessageTypeBroadcast
)

var inputMessageTypeJSONs = map[InputMessageType][]byte{
	InputMessageTypeSnakeCommand: []byte(`"snake"`),
	InputMessageTypeBroadcast:    []byte(`"broadcast"`),
}

var ErrUnknownInputMessageType = errors.New("unknown input message type")

func (t *InputMessageType) UnmarshalJSON(data []byte) error {
	for msgType, commandJSON := range inputMessageTypeJSONs {
		if bytes.Equal(commandJSON, data) {
			*t = msgType
			return nil
		}
	}

	return ErrUnknownInputMessageType
}

var inputMessageTypeLabels = map[InputMessageType]string{
	InputMessageTypeSnakeCommand: "snake",
	InputMessageTypeBroadcast:    "broadcast",
}

func (t InputMessageType) String() string {
	if label, ok := inputMessageTypeLabels[t]; ok {
		return label
	}
	return "unknown"
}

//go:generate ffjson $GOFILE

// ffjson: noencoder
type InputMessage struct {
	Type    InputMessageType `json:"type"`
	Payload string           `json:"payload"`
}
