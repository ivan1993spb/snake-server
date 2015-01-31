// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Input and output data in websocket is JSON objects with format:

	{"header": HEADER, "data": DATA}

*/

package main

import "encoding/json"

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
	Header string          `json:"header"`
	Data   json.RawMessage `json:"data"`
}
