package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteCreateGame = "/game"

const MethodCreateGame = http.MethodPost

const (
	postFieldConnectionLimit = "limit"
	postFieldMapWidth        = "width"
	postFieldMapHeight       = "height"
)

type responseCreateGameHandler struct {
	ID     int   `json:"id"`
	Limit  int   `json:"limit"`
	Width  uint8 `json:"width"`
	Height uint8 `json:"height"`
}

type createGameHandler struct {
	logger       *logrus.Logger
	groupManager *connections.ConnectionGroupManager
}

type ErrCreateGameHandler string

func (e ErrCreateGameHandler) Error() string {
	return "create game handler error: " + string(e)
}

func NewCreateGameHandler(logger *logrus.Logger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &createGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *createGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("create game handler start")
	defer h.logger.Info("create game handler end")

	connectionLimit, err := strconv.Atoi(r.PostFormValue(postFieldConnectionLimit))
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mapWidth, err := strconv.ParseUint(r.PostFormValue(postFieldMapWidth), 10, 8)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mapHeight, err := strconv.ParseUint(r.PostFormValue(postFieldMapHeight), 10, 8)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infof("create group: limit=%d, width=%d, height=%d", connectionLimit, mapWidth, mapHeight)

	group, err := connections.NewConnectionGroup(h.logger, connectionLimit, uint8(mapWidth), uint8(mapHeight))
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.logger.Info("initialized group")

	id, err := h.groupManager.Add(group)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	h.logger.Infoln("created group with id:", id)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(responseCreateGameHandler{
		ID:     id,
		Limit:  connectionLimit,
		Width:  uint8(mapWidth),
		Height: uint8(mapHeight),
	})
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
