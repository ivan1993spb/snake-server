package handlers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteGetInfo = "/"

const MethodGetInfo = http.MethodGet

type getInfoHandler struct {
	logger logrus.FieldLogger
	info   string
}

type ErrGetInfoHandler string

func (e ErrGetInfoHandler) Error() string {
	return "get info handler error: " + string(e)
}

func NewGetInfoHandler(logger logrus.FieldLogger, version, build string) http.Handler {
	return &getInfoHandler{
		logger: logger,
		info:   fmt.Sprintf("Wellcome to Snake-Server!\nVersion: %s (build %s)\n", version, build),
	}
}

func (h *getInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := fmt.Fprint(w, h.info); err != nil {
		h.logger.Error(ErrGetInfoHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
