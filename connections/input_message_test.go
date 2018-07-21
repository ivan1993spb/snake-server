package connections

import (
	"testing"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/testify/require"
)

func Test_InputMessageType_UnmarshalJSON_CorrectData(t *testing.T) {
	data := []byte(`{"type": "snake", "payload": "north"}`)
	expected := InputMessage{
		Type:    InputMessageTypeSnakeCommand,
		Payload: "north",
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
	require.NotNil(t, err)
}

func Test_InputMessageType_UnmarshalJSON_BroadcastMessageTypes(t *testing.T) {
	data := []byte(`{"type": "broadcast", "payload": "hello"}`)
	expected := InputMessage{
		Type:    InputMessageTypeBroadcast,
		Payload: "hello",
	}
	var inputMessage InputMessage
	err := ffjson.Unmarshal(data, &inputMessage)
	require.Nil(t, err)
	require.Equal(t, expected, inputMessage)
}
