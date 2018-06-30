package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteBroadcast = "/games/{id}/broadcast"

const MethodBroadcast = http.MethodPost

const broadcastTimeout = time.Millisecond

const broadcastMaxBodySize = 128

const postFieldBroadcast = "message"

type responseBroadcastHandler struct {
	Success bool `json:"success"`
}

type responseBroadcastHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type broadcastHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrBroadcastHandler string

func (e ErrBroadcastHandler) Error() string {
	return "broadcast handler error: " + string(e)
}

func NewBroadcastHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &broadcastHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *broadcastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, broadcastMaxBodySize)
	if err := r.ParseForm(); err != nil {
		h.logger.Error(ErrBroadcastHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseBroadcastHandlerError{
			Code: http.StatusBadRequest,
			Text: "bad request",
		})
		return
	}

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(ErrBroadcastHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseBroadcastHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
		})
		return
	}

	h.logger.Infoln("group id to broadcast:", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Errorln("cannot get group:", err.Error())

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseBroadcastHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseBroadcastHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
			})
		}
		return
	}

	if group.IsEmpty() {
		h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseBroadcastHandlerError{
			Code: http.StatusServiceUnavailable,
			Text: "group is empty",
		})
		return
	}

	message := r.PostFormValue(postFieldBroadcast)
	if len(message) == 0 {
		h.writeResponseJSON(w, http.StatusBadRequest, &responseBroadcastHandlerError{
			Code: http.StatusBadRequest,
			Text: "message is empty",
		})
		return
	}

	h.writeResponseJSON(w, http.StatusOK, &responseBroadcastHandler{
		Success: group.BroadcastMessageTimeout(message, broadcastTimeout),
	})
}

func (h *broadcastHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrBroadcastHandler(err.Error()))
	}
}
