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
}

type ErrGetInfoHandler string

func (e ErrGetInfoHandler) Error() string {
	return "get info handler error: " + string(e)
}

func NewGetInfoHandler(logger logrus.FieldLogger) http.Handler {
	return &getInfoHandler{
		logger: logger,
	}
}

func (h *getInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Create server name, links, description, version, build tag.

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := fmt.Fprintln(w, "Wellcome to Snake-Server!"); err != nil {
		h.logger.Error(ErrGetInfoHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
