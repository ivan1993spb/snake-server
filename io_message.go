/*
Input and output data is JSON objects:

	{"header": HEADER, "data": DATA}
*/
package main

import "encoding/json"

// Output headers
const (
	HEADER_ERROR = "error" // Header for error reporting
	HEADER_INFO  = "info"  // Header for info messages
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
