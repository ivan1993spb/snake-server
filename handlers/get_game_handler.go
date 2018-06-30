package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetGameByID = "/games/{id}"

const MethodGetGame = http.MethodGet

type responseGetGameHandler struct {
	ID     int `json:"id"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type responseGetGameHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
	ID   int    `json:"id"`
}

type getGameHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrGetGameHandler string

func (e ErrGetGameHandler) Error() string {
	return "get game handler error: " + string(e)
}

func NewGetGameHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *getGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(ErrGetGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGetGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
			ID:   id,
		})
		return
	}

	h.logger.Infoln("group id to get:", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Errorln("cannot get group:", err.Error())

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseGetGameHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
				ID:   id,
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseGetGameHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
				ID:   id,
			})
		}
		return
	}

	h.writeResponseJSON(w, http.StatusOK, &responseGetGameHandler{
		ID:     id,
		Limit:  group.GetLimit(),
		Count:  group.GetCount(),
		Width:  int(group.GetWorldWidth()),
		Height: int(group.GetWorldHeight()),
	})
}

func (h *getGameHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGetGameHandler(err.Error()))
	}
}
