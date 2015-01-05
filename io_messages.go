// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Input and output data is JSON objects:

	{"header": HEADER, "data": DATA}
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Headers
const (
	// Output headers
	HEADER_ERROR   = "error"   // Header for error reporting
	HEADER_INFO    = "info"    // Header for info messages
	HEADER_POOL_ID = "pool_id" // Header for sending pool ids
	HEADER_CONN_ID = "conn_id" // Header for sending connection ids

	// Input/output headers
	HEADER_GAME = "game" // Header for game data
)

type OutputMessage struct {
	Header string      `json:"header"`
	Data   interface{} `json:"data"`
}

type InputMessage struct {
	Header string `json:"header"`
	// Do not parse data while header is unknown
	Data json.RawMessage `json:"data"`
}

type errReceiveMessage struct {
	err error
}

func (e *errReceiveMessage) Error() string {
	return "receiving message error: " + e.err.Error()
}

func ReceiveMessage(ws *websocket.Conn, headers ...string,
) (*InputMessage, error) {
	if len(headers) == 0 {
		return nil, &errReceiveMessage{errors.New("no headers")}
	}

	var msg *InputMessage
	if err := websocket.JSON.Receive(ws, &msg); err != nil {
		if err != io.EOF {
			err = &errReceiveMessage{err}
		}
		return nil, err
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

func SendMessage(ws *websocket.Conn, header string, data interface{},
) error {
	err := websocket.JSON.Send(ws, &OutputMessage{header, data})
	if err != nil {
		return fmt.Errorf("sending message error: %s", err)
	}

	return nil
}
