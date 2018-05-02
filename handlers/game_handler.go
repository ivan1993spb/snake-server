package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const URLRouteGameByID = "/game/{id}"

const MethodGame = http.MethodDelete

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type gameHandler struct {
	logger *logrus.Logger
}

type ErrGameHandler string

func NewGameHandler(logger *logrus.Logger) http.Handler {
	return &gameHandler{
		logger: logger,
	}
}

func (h *gameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler start")

	vars := mux.Vars(r)
	h.logger.Infoln("vars", vars)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err)
	}

	// TODO: Implement handler.

	conn.Close()

	h.logger.Info("game handler end")
}
