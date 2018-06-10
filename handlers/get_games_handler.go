package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetGames = "/games"

const MethodGetGames = http.MethodGet

type responseGetGamesEntity struct {
	ID     int `json:"id"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type responseGetGamesHandler struct {
	Games []responseGetGamesEntity `json:"games"`
	Limit int                      `json:"limit"`
	Count int                      `json:"count"`
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
	groupCount := h.groupManager.GroupCount()

	response := responseGetGamesHandler{
		Games: make([]responseGetGamesEntity, 0, groupCount),
		Limit: h.groupManager.GroupLimit(),
		Count: groupCount,
	}

	for id, group := range h.groupManager.Groups() {
		response.Games = append(response.Games, responseGetGamesEntity{
			ID:     id,
			Limit:  group.GetLimit(),
			Count:  group.GetCount(),
			Width:  int(group.GetWorldWidth()),
			Height: int(group.GetWorldHeight()),
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGetGamesHandler(err.Error()))
	}
}
