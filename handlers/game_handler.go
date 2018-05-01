package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const URLRouteGame = "/game/{id}"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type gameHandler struct {
	logger *logrus.Logger
}

func NewGameHandler(logger *logrus.Logger) http.Handler {
	return &gameHandler{
		logger: logger,
	}
}

func (h *gameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err)
	}

	// TODO: Implement handler.
}
