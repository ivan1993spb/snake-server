package handlers

import (
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

const URLRouteGetInfo = "/"

const MethodGetInfo = http.MethodGet

type responseGetInfoHandler struct {
	Author  string `json:"author"`
	License string `json:"license"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

type getInfoHandler struct {
	logger logrus.FieldLogger
	info   []byte
}

type ErrGetInfoHandler string

func (e ErrGetInfoHandler) Error() string {
	return "get info handler error: " + string(e)
}

func NewGetInfoHandler(logger logrus.FieldLogger, author, license, version, build string) http.Handler {
	info, err := ffjson.Marshal(&responseGetInfoHandler{
		Author:  author,
		License: license,
		Version: version,
		Build:   build,
	})
	if err != nil {
		logger.WithError(err).Error("error on create info handler")
		panic(err)
	}
	return &getInfoHandler{
		logger: logger,
		info:   info,
	}
}

func (h *getInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(h.info); err != nil {
		h.logger.Error(ErrGetInfoHandler(err.Error()))
	}
}
