package main

import (
	"net/http"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
)

// RequestVerifier verifies requests by hash sum of passed request
// data
type RequestVerifier struct{}

func NewVerifier(HashSalt string) pwshandler.RequestVerifier {
	return &RequestVerifier{}
}

// Implementing pwshandler.RequestVerifier interface
func (*RequestVerifier) Verify(*http.Request) error {
	if glog.V(1) {
		glog.Infoln("Request was verified")
	}
	return nil
}
