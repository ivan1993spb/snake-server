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
	ID    int `json:"id"`
	Limit int `json:"limit"`
	Count int `json:"count"`
}

type getGameHandler struct {
	logger       *logrus.Logger
	groupManager *connections.ConnectionGroupManager
}

type ErrGetGameHandler string

func (e ErrGetGameHandler) Error() string {
	return "get game handler error: " + string(e)
}

func NewGetGameHandler(logger *logrus.Logger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *getGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get game handler start")
	defer h.logger.Info("get game handler end")

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
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(responseGetGameHandler{
		ID:    id,
		Limit: group.GetLimit(),
		Count: group.GetCount(),
	})
	if err != nil {
		h.logger.Error(ErrGetGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
