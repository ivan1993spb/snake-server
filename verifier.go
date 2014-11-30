package main

import (
	"net/http"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
)

// Implementing pwshandler.RequestVerifier interface
type RequestVerifier struct{}

func NewVerifier(HashSalt string) pwshandler.RequestVerifier {
	return &RequestVerifier{}
}

func (*RequestVerifier) Verify(*http.Request) error {
	if glog.V(1) {
		glog.Infoln("Request was verified")
	}
	return nil
}
