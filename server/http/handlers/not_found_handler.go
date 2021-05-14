package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

var notFoundJSONResponse = []byte(`{"code":404,"text":"not found"}`)

type notFoundHandler struct {
	logger logrus.FieldLogger
}

type ErrNotFoundHandler string

func (e ErrNotFoundHandler) Error() string {
	return "not found handler error: " + string(e)
}

func NewNotFoundHandler(logger logrus.FieldLogger) http.Handler {
	return &notFoundHandler{
		logger: logger,
	}
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	if _, err := w.Write(notFoundJSONResponse); err != nil {
		h.logger.Error(ErrNotFoundHandler(err.Error()))
	}
}
