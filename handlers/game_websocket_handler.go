package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGameWebSocket = "/game/{id}/ws"

const MethodGame = http.MethodGet

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

type gameWebSocketHandler struct {
	logger       *logrus.Logger
	groupManager *connections.ConnectionGroupManager
}

type ErrGameHandler string

func NewGameWebSocketHandler(logger *logrus.Logger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &gameWebSocketHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *gameWebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler start")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars[routeVarGroupID])
	if err != nil {
		// TODO: Create custom error
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infoln("group id", id)

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

	h.logger.Info(group)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err)
	}
	defer conn.Close()

	group.Add(conn)

	// TODO: Implement handler.

	group.Delete(conn)

	h.logger.Info("game handler end")
}
