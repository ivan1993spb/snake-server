package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetGames = "/games"

const MethodGetGames = http.MethodGet

// TODO: Create width and height?
type responseGetGamesEntity struct {
	ID    int `json:"id"`
	Limit int `json:"limit"`
	Count int `json:"count"`
}

type responseGetGamesHandler struct {
	Games []responseGetGamesEntity `json:"games"`
}

type getGamesHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrGetGamesHandler string

func (e ErrGetGamesHandler) Error() string {
	return "get game handler error: " + string(e)
}

func NewGetGamesHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getGamesHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *getGamesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get games handler start")
	defer h.logger.Info("get games handler end")

	response := responseGetGamesHandler{
		Games: []responseGetGamesEntity{},
	}
	for id, group := range h.groupManager.Groups() {
		response.Games = append(response.Games, responseGetGamesEntity{
			ID:    id,
			Limit: group.GetLimit(),
			Count: group.GetCount(),
		})
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGetGamesHandler(err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
