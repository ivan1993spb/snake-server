package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteCreateGame = "/game"

const MethodCreateGame = http.MethodPost

type createGameHandler struct {
	logger *logrus.Logger
}

type ErrCreateGameHandler string

func NewCreateGameHandler(logger *logrus.Logger) http.Handler {
	return &createGameHandler{
		logger: logger,
	}
}

func (h *createGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler start")

	h.logger.Info("game handler end")
}
