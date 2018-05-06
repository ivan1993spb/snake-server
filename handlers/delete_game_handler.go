package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteDeleteGameByID = "/games/{id}"

const MethodDeleteGame = http.MethodDelete

type responseDeleteGameHandler struct {
	ID int `json:"id"`
}

type deleteGameHandler struct {
	logger       *logrus.Logger
	groupManager *connections.ConnectionGroupManager
}

type ErrDeleteGameHandler string

func (e ErrDeleteGameHandler) Error() string {
	return "delete game handler error: " + string(e)
}

func NewDeleteGameHandler(logger *logrus.Logger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &deleteGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *deleteGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("delete game handler start")
	defer h.logger.Info("delete game handler end")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infoln("group id to delete:", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))

		switch err {
		case connections.ErrNotFoundGroup:
			http.NotFound(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if !group.IsEmpty() {
		h.logger.Warnln("try to delete not empty group:", id)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	if err := h.groupManager.Delete(group); err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))

		switch err {
		case connections.ErrDeleteNotFoundGroup:
			http.NotFound(w, r)
		case connections.ErrDeleteNotEmptyGroup:
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Infoln("group deleted:", id)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(responseDeleteGameHandler{
		ID: id,
	})
	if err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
