package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteWelcome = "/"

const MethodWelcome = http.MethodGet

var welcomeMessage = []byte(`Welcome to Snake-Server!`)

type welcomeHandler struct {
	logger logrus.FieldLogger
}

type ErrWelcomeHandler string

func (e ErrWelcomeHandler) Error() string {
	return "welcome handler error: " + string(e)
}

func NewWelcomeHandler(logger logrus.FieldLogger) http.Handler {
	return &welcomeHandler{
		logger: logger,
	}
}

func (h *welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(welcomeMessage); err != nil {
		h.logger.Error(ErrWelcomeHandler(err.Error()))
	}
}
