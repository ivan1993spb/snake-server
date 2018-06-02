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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infoln("group id to get:", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Errorln("cannot get group:", err.Error())

		switch err {
		case connections.ErrNotFoundGroup:
			http.NotFound(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(responseGetGameHandler{
		ID:     id,
		Limit:  group.GetLimit(),
		Count:  group.GetCount(),
		Width:  int(group.GetWorldWidth()),
		Height: int(group.GetWorldHeight()),
	})
	if err != nil {
		h.logger.Error(ErrGetGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
