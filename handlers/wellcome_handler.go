package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteWellcome = "/"

const MethodWellcome = http.MethodGet

var wellcomeMessage = []byte(`Wellcome to Snake-Server!`)

type wellcomeHandler struct {
	logger logrus.FieldLogger
}

type ErrWellcomeHandler string

func (e ErrWellcomeHandler) Error() string {
	return "wellcome handler error: " + string(e)
}

func NewWellcomeHandler(logger logrus.FieldLogger) http.Handler {
	return &wellcomeHandler{
		logger: logger,
	}
}

func (h *wellcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := w.Write(wellcomeMessage); err != nil {
		h.logger.Error(ErrWellcomeHandler(err.Error()))
	}
}
