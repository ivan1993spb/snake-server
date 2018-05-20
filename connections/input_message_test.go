package connections

import (
	"encoding/json"
	"testing"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/testify/require"
)

func Test_InputMessageType_UnmarshalJSON_CorrectData(t *testing.T) {
	data := []byte(`{"type": "snake", "payload": "north"}`)
	expected := InputMessage{
		Type:    InputMessageTypeSnakeCommand,
		Payload: json.RawMessage(`"north"`),
	}
	var inputMessage InputMessage
	err := ffjson.Unmarshal(data, &inputMessage)
	require.Nil(t, err)
	require.Equal(t, expected, inputMessage)
}

func Test_InputMessageType_UnmarshalJSON_InvalidMessageType(t *testing.T) {
	data := []byte(`{"type": "invalid", "payload": "north"}`)
	var inputMessage InputMessage
	err := ffjson.Unmarshal(data, &inputMessage)
	require.Equal(t, ErrUnknownInputMessageType, err)
}
