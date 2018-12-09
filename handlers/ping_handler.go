package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRoutePing = "/ping"

const MethodPing = http.MethodGet

var pingResponseBody = []byte(`{"pong":1}`)

type pingHandler struct {
	logger logrus.FieldLogger
}

type ErrPingHandler string

func (e ErrPingHandler) Error() string {
	return "ping handler error: " + string(e)
}

func NewPingHandler(logger logrus.FieldLogger) http.Handler {
	return &pingHandler{
		logger: logger,
	}
}

func (h *pingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(pingResponseBody); err != nil {
		h.logger.Error(ErrPingHandler(err.Error()))
	}
}
