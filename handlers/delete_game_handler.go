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

type responseDeleteGameHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
	ID   int    `json:"id"`
}

type deleteGameHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrDeleteGameHandler string

func (e ErrDeleteGameHandler) Error() string {
	return "delete game handler error: " + string(e)
}

func NewDeleteGameHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &deleteGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *deleteGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseDeleteGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
			ID:   id,
		})
		return
	}

	h.logger.Infoln("group id to delete:", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseDeleteGameHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
				ID:   id,
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseDeleteGameHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
				ID:   id,
			})
		}
		return
	}

	if !group.IsEmpty() {
		h.logger.Warn(ErrDeleteGameHandler("try to delete not empty group"))
		h.logger.Warnf("there is %d opened connections in group %d", group.GetCount(), id)
		h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseDeleteGameHandlerError{
			Code: http.StatusServiceUnavailable,
			Text: "cannot delete not empty game",
			ID:   id,
		})
		return
	}

	if err := h.groupManager.Delete(group); err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))

		switch err {
		case connections.ErrDeleteNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseDeleteGameHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
				ID:   id,
			})
		case connections.ErrDeleteNotEmptyGroup:
			h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseDeleteGameHandlerError{
				Code: http.StatusServiceUnavailable,
				Text: "cannot delete not empty game",
				ID:   id,
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseDeleteGameHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
				ID:   id,
			})
		}
		return
	}

	h.logger.Info("stop group")
	group.Stop()

	h.logger.Infoln("group deleted:", id)

	h.writeResponseJSON(w, http.StatusOK, responseDeleteGameHandler{
		ID: id,
	})
}

func (h *deleteGameHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrDeleteGameHandler(err.Error()))
	}
}
