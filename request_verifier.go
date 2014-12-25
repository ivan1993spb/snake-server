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
	return "cannot verify connection: " + e.err.Error()
}

// RequestVerifier verifies requests
type RequestVerifier struct{}

func NewRequestVerifier(HashSalt string) pwshandler.RequestVerifier {
	return &RequestVerifier{}
}

// Implementing pwshandler.RequestVerifier interface
func (*RequestVerifier) Verify(ws *websocket.Conn) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("verifying accepted connection")
		defer glog.Infoln("connection was verified")
	}

	var err error

	err = websocket.JSON.Send(ws, &OutputMessage{
		HEADER_INFO, "verifying connection",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	// Check received hash

	// ...

	err = websocket.JSON.Send(ws, &OutputMessage{
		HEADER_INFO, "connection was verified",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	return nil
}
