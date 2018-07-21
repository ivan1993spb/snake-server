package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetObjects = "/games/{id}/objects"

const MethodGetObjects = http.MethodGet

type responseGetObjectsHandler struct {
	Objects []interface{} `json:"objects"`
}

type responseGetObjectsHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type getObjectsHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrGetObjectsHandler string

func (e ErrGetObjectsHandler) Error() string {
	return "get objects handler error: " + string(e)
}

func NewGetObjectsHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getObjectsHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *getObjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("parse game id error")
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGetObjectsHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
		})
		return
	}

	h.logger.WithField("game", id).Infoln("game id received")

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("cannot get game group")

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseGetObjectsHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseGetObjectsHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
			})
		}
		return
	}

	h.writeResponseJSON(w, http.StatusOK, &responseGetObjectsHandler{
		Objects: group.GetObjects(),
	})
}

func (h *getObjectsHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("encode response error")
	}
}
