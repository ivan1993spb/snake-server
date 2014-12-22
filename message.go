package main

import "encoding/json"

// Header represents message type
type Header string

const (
	// Used only by client
	HEADER_AUTH = "auth" // Header for auth data

	// Used only by server
	HEADER_PROTOCOL = "protocol" // Header for protocol declaration
	HEADER_ERROR    = "error"    // Header for error reporting
	HEADER_INFO     = "info"     // Header for info messages

	// Used by server and client
	HEADER_GAME = "game" // Header for game data

)

// Message is input or output message
type Message struct {
	Header Header
	Data   json.RawMessage
}
