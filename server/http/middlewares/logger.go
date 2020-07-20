package middlewares

import (
	"net/http"

	"github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const headerSnakeClient = "X-Snake-Client"

func NewLogger(logger *logrus.Logger, name string) negroni.Handler {
	m := negronilogrus.NewMiddlewareFromLogger(logger, name)
	m.Before = before
	return m
}

func before(entry *logrus.Entry, req *http.Request, remoteAddr string) *logrus.Entry {
	var client = req.Header.Get(headerSnakeClient)

	if len(client) == 0 {
		client = req.UserAgent()
	}

	return entry.WithFields(logrus.Fields{
		"client":  client,
		"request": req.RequestURI,
		"method":  req.Method,
		"remote":  remoteAddr,
	})
}
