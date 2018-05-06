package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGameWebSocketByID = "/games/{id}/ws"

const MethodGame = http.MethodGet

type gameWebSocketHandler struct {
	logger       *logrus.Logger
	groupManager *connections.ConnectionGroupManager
}

type ErrGameWebSocketHandler string

func (e ErrGameWebSocketHandler) Error() string {
	return "game websocket handler error: " + string(e)
}

func NewGameWebSocketHandler(logger *logrus.Logger, groupManager *connections.ConnectionGroupManager) http.Handler {
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infoln("group id", id)

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))

		switch err {
		case connections.ErrNotFoundGroup:
			http.NotFound(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if group.IsFull() {
		h.logger.Warn(ErrGameWebSocketHandler("group is full"))
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	connection, err := connections.NewConnection(w, r)
	if err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if err := group.Handle(connection); err != nil {
		h.logger.Error(ErrGameWebSocketHandler(err.Error()))

		switch err.Err {
		case connections.ErrGroupIsFull:
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
}
