package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

// Form keys that are used for request verifying
const (
	FORM_KEY_SUM  = "sum"
	FORM_KEY_PART = "part"
)

type errConnVerifying struct {
	err error
}

func (e *errConnVerifying) Error() string {
	return "cannot verify connection: " + e.err.Error()
}

// RequestVerifier verifies requests
type RequestVerifier struct {
	hashSalt string
}

func NewRequestVerifier(hashSalt string) pwshandler.RequestVerifier {
	return &RequestVerifier{hashSalt}
}

// Implementing pwshandler.RequestVerifier interface
func (rv *RequestVerifier) Verify(ws *websocket.Conn) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("verifying accepted connection")
		defer glog.Infoln("connection was verified")
	}

	err := websocket.JSON.Send(ws, &OutputMessage{
		HEADER_INFO, "verifying connection",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	sum, err := hex.DecodeString(ws.Request().FormValue(FORM_KEY_SUM))
	if len(sum) != sha256.Size {
		return &errConnVerifying{errors.New("invalid sum size")}
	}

	controlSum := sha256.Sum256([]byte(
		rv.hashSalt + ws.Request().FormValue(FORM_KEY_PART),
	))

	for i := 0; i < sha256.Size; i++ {
		if controlSum[i] != sum[i] {
			return &errConnVerifying{
				errors.New("request is not trusted"),
			}
		}
	}

	err = websocket.JSON.Send(ws, &OutputMessage{
		HEADER_INFO, "connection was verified",
	})
	if err != nil {
		return &errConnVerifying{err}
	}

	return nil
}
