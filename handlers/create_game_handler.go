package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteCreateGame = "/games"

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

type responseCreateGameHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type createGameHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrCreateGameHandler string

func (e ErrCreateGameHandler) Error() string {
	return "create game handler error: " + string(e)
}

func NewCreateGameHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &createGameHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *createGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connectionLimit, err := strconv.Atoi(r.PostFormValue(postFieldConnectionLimit))
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid limit",
		})
		return
	}
	if connectionLimit <= 0 {
		h.logger.Warnln(ErrCreateGameHandler("invalid connection limit"), connectionLimit)
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid limit",
		})
		return
	}

	mapWidth, err := strconv.ParseUint(r.PostFormValue(postFieldMapWidth), 10, 8)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid width",
		})
		return
	}
	if mapWidth == 0 {
		h.logger.Warnln(ErrCreateGameHandler("invalid map width"), mapWidth)
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid width",
		})
		return
	}

	mapHeight, err := strconv.ParseUint(r.PostFormValue(postFieldMapHeight), 10, 8)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid height",
		})
		return
	}
	if mapHeight == 0 {
		h.logger.Warnln(ErrCreateGameHandler("invalid map height"), mapHeight)
		h.writeResponseJSON(w, http.StatusBadRequest, &responseCreateGameHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid height",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"width":            mapWidth,
		"height":           mapHeight,
		"connection_limit": connectionLimit,
	}).Debug("create game group")

	group, err := connections.NewConnectionGroup(h.logger, connectionLimit, uint8(mapWidth), uint8(mapHeight))
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
		h.writeResponseJSON(w, http.StatusInternalServerError, &responseCreateGameHandlerError{
			Code: http.StatusInternalServerError,
			Text: "cannot create game",
		})
		return
	}

	id, err := h.groupManager.Add(group)
	if err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))

		switch err {
		case connections.ErrGroupLimitReached:
			h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseCreateGameHandlerError{
				Code: http.StatusServiceUnavailable,
				Text: "groups limit reached",
			})
		case connections.ErrConnsLimitReached:
			h.writeResponseJSON(w, http.StatusServiceUnavailable, &responseCreateGameHandlerError{
				Code: http.StatusServiceUnavailable,
				Text: "connections limit reached",
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseCreateGameHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
			})
		}
		return
	}

	h.logger.Info("start group")
	group.Start()

	h.logger.WithField("group_id", id).Infoln("created group")

	h.writeResponseJSON(w, http.StatusCreated, &responseCreateGameHandler{
		ID:     id,
		Limit:  group.GetLimit(),
		Width:  uint8(mapWidth),
		Height: uint8(mapHeight),
	})
}

func (h *createGameHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrCreateGameHandler(err.Error()))
	}
}
