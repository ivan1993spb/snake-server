package main

import (
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

type errConnVerifying struct {
	err error
}

func (e *errConnVerifying) Error() string {
	return "Verifying connection error: " + e.err.Error()
}

// RequestVerifier verifies requests by hash sum of passed request
// data
type RequestVerifier struct{}

func NewRequestVerifier(HashSalt string) pwshandler.RequestVerifier {
	return &RequestVerifier{}
}

// Implementing pwshandler.RequestVerifier interface
func (*RequestVerifier) Verify(ws *websocket.Conn) (err error) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Verifying accepted connection")
	}

	err = websocket.JSON.Send(ws, &Message{
		HEADER_INFO, "Verifying connection",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	// Check received hash

	// ...

	err = websocket.JSON.Send(ws, &Message{
		HEADER_INFO, "Connection was verified",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	return nil
}
