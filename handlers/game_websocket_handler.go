package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGameWebSocketByID = "/games/{id}"

const MethodGame = http.MethodGet

const wsReadMessageLimit = 128

const wsReadBufferSize = 2048

const wsWriteBufferSize = 20480

const messageUpgradeConnectionError = "web-socket upgrade connection error"

type responseGameWebSocketHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type gameWebSocketHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
	upgrader     *websocket.Upgrader
}

type ErrGameWebSocketHandler string

func (e ErrGameWebSocketHandler) Error() string {
	return "game web-socket handler error: " + string(e)
}

func NewGameWebSocketHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:    wsReadBufferSize,
		WriteBufferSize:   wsWriteBufferSize,
		EnableCompression: false,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	handler := &gameWebSocketHandler{
		logger:       logger,
		groupManager: groupManager,
		upgrader:     upgrader,
	}

	upgrader.Error = handler.errorUpgradeConnection

	return handler
}

func (h *gameWebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler start")
	defer h.logger.Info("game handler end")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGameWebSocketHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
		})
		return
	}

	h.logger.WithField("game", id).Info("try to connect to game group")

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseGameWebSocketHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseGameWebSocketHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
			})
		}
		return
	}

	if group.IsFull() {
		h.logger.Warn(ErrGameWebSocketHandler("group is full"))
		h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseGameWebSocketHandlerError{
			Code: http.StatusServiceUnavailable,
			Text: "group is full",
		})
		return
	}

	h.logger.Info("upgrade connection")

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		// Response is written by failed upgrader
		return
	}

	conn.SetReadLimit(wsReadMessageLimit)

	h.logger.Info("start connection worker")

	if err := group.Handle(connections.NewConnectionWorker(conn, h.logger)); err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		return
	}
}

func (h *gameWebSocketHandler) errorUpgradeConnection(w http.ResponseWriter, _ *http.Request, status int, _ error) {
	// Composing error message for upgrade failure case
	w.Header().Set("Sec-Websocket-Version", "13")

	h.writeResponseJSON(w, status, &responseGameWebSocketHandlerError{
		Code: status,
		Text: messageUpgradeConnectionError,
	})
}

func (h *gameWebSocketHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
	}
}
