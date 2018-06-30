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

var upgrader = websocket.Upgrader{
	ReadBufferSize:    wsReadBufferSize,
	WriteBufferSize:   wsWriteBufferSize,
	EnableCompression: false,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type responseGameWebSocketHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type gameWebSocketHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrGameWebSocketHandler string

func (e ErrGameWebSocketHandler) Error() string {
	return "game web-socket handler error: " + string(e)
}

func NewGameWebSocketHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &gameWebSocketHandler{
		logger:       logger,
		groupManager: groupManager,
	}
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

	h.logger.Infoln("try to connect to game group id", id)

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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusInternalServerError, &responseGameWebSocketHandlerError{
			Code: http.StatusInternalServerError,
			Text: "web-socket upgrade connection error",
		})
		return
	}

	conn.SetReadLimit(wsReadMessageLimit)

	h.logger.Info("start connection worker")

	if err := group.Handle(connections.NewConnectionWorker(conn, h.logger)); err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		return
	}
}

func (h *gameWebSocketHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
	}
}
