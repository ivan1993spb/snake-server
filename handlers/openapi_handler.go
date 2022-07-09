package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteOpenAPI = "/openapi.yaml"

type openAPIHandler struct {
	logger logrus.FieldLogger
	spec   []byte
}

type ErrOpenAPIHandler string

func (e ErrOpenAPIHandler) Error() string {
	return "openapi handler error: " + string(e)
}

func NewOpenAPIHandler(logger logrus.FieldLogger, spec []byte) http.Handler {
	return &openAPIHandler{
		logger: logger,
		spec:   spec,
	}
}

func (h *openAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(h.spec); err != nil {
		h.logger.Error(ErrOpenAPIHandler(err.Error()))
	}
}
