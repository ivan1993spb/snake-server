/*
Input and output data is JSON objects:

	{"header": HEADER, "data": DATA}
*/
package main

import (
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/net/websocket"
)

// Output headers
const (
	HEADER_ERROR   = "error"   // Header for error reporting
	HEADER_INFO    = "info"    // Header for info messages
	HEADER_POOL_ID = "pool_id" // Header for sending pool ids
	HEADER_CONN_ID = "conn_id" // Header for sending connection ids
)

type OutputMessage struct {
	Header string      `json:"header"`
	Data   interface{} `json:"data"`
}

// Input headers
const (
	HEADER_AUTH = "auth" // Header for auth data
)

type InputMessage struct {
	Header string `json:"header"`
	// Do not parse data while header is unknown
	Data json.RawMessage `json:"data"`
}

// Input/output headers
const (
	HEADER_GAME = "game" // Header for game data
)

type errReceiveMessage struct {
	err error
}

func (e *errReceiveMessage) Error() string {
	return "receiving message error: " + e.err.Error()
}

var (
	ErrConnStop  = &errReceiveMessage{io.EOF}
	ErrConnFatal = &errReceiveMessage{
		errors.New("fatal connection error")}
)

func ReceiveMessage(ws *websocket.Conn, headers ...string,
) (*InputMessage, error) {
	if len(headers) == 0 {
		return nil, &errReceiveMessage{errors.New("no headers")}
	}

	var data []byte
	if err := websocket.Message.Receive(ws, &data); err != nil {
		if err == io.EOF {
			return nil, ErrConnStop
		}
		return nil, ErrConnFatal
	}

	var msg *InputMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, &errReceiveMessage{err}
	}

	if len(msg.Header) == 0 {
		return nil, &errReceiveMessage{
			errors.New("empty message header"),
		}
	}

	for i := 0; i < len(headers); i++ {
		if headers[i] == msg.Header {
			return msg, nil
		}
	}

	return nil, &errReceiveMessage{errors.New("unexpected header")}
}

type errSendMessage struct {
	err error
}

func (e *errSendMessage) Error() string {
	return "sending message error: " + e.err.Error()
}

func SendMessage(ws *websocket.Conn, header string, data interface{},
) error {
	data, err := json.Marshal(&OutputMessage{header, data})
	if err != nil {
		return &errSendMessage{err}
	}

	err = websocket.Message.Send(ws, data)
	if err != nil {
		return &errSendMessage{err}
	}

	return nil
}
