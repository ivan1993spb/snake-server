package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const URLRouteDeleteGame = "/game"

const MethodDeleteGame = http.MethodDelete

type deleteGameHandler struct {
	logger *logrus.Logger
}

type ErrDeleteGameHandler string

func NewDeleteGameHandler(logger *logrus.Logger) http.Handler {
	return &createGameHandler{
		logger: logger,
	}
}

func (h *deleteGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("game handler start")

	h.logger.Info("game handler end")
}
