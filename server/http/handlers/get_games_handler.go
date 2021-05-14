package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetGames = "/games"

const MethodGetGames = http.MethodGet

const (
	getFieldGamesLimit   = "limit"
	getFieldGamesSorting = "sorting"
)

const (
	gamesSortingLabelSmart  = "smart"
	gamesSortingLabelRandom = "random"
)

type gamesSorting uint8

const (
	gamesSortingUndefined gamesSorting = iota
	gamesSortingSmart
	gamesSortingRandom
)

const gamesSortingDefault = gamesSortingRandom

var errInvalidGamesSorting = errors.New("invalid sorting")

func parseGamesSorting(sorting string) (gamesSorting, error) {
	if len(sorting) == 0 {
		return gamesSortingDefault, nil
	}

	if sorting == gamesSortingLabelSmart {
		return gamesSortingSmart, nil
	}

	if sorting == gamesSortingLabelRandom {
		return gamesSortingRandom, nil
	}

	return gamesSortingUndefined, errInvalidGamesSorting
}

type responseGetGamesEntity struct {
	ID     int    `json:"id"`
	Limit  int    `json:"limit"`
	Count  int    `json:"count"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Rate   uint32 `json:"rate"`
}

type responseGetGamesHandler struct {
	Games []*responseGetGamesEntity `json:"games"`
	Limit int                       `json:"limit"`
	Count int                       `json:"count"`
}

type responseGetGamesHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
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
	sorting, err := parseGamesSorting(r.URL.Query().Get(getFieldGamesSorting))
	if err != nil {
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGetGamesHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid sorting",
		})
		return
	}

	limitLabel := r.URL.Query().Get(getFieldGamesLimit)
	flagUseLimit := len(limitLabel) > 0
	limit, err := strconv.Atoi(limitLabel)
	if flagUseLimit && (err != nil || limit < 0) {
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGetGamesHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid limit value",
		})
		return
	}

	groupCount := h.groupManager.GroupCount()

	if groupCount == 0 {
		h.writeResponseJSON(w, http.StatusOK, &responseGetGamesHandler{
			Games: []*responseGetGamesEntity{},
			Limit: h.groupManager.GroupLimit(),
			Count: groupCount,
		})
		return
	}

	if flagUseLimit && limit == 0 {
		h.writeResponseJSON(w, http.StatusOK, &responseGetGamesHandler{
			Games: []*responseGetGamesEntity{},
			Limit: h.groupManager.GroupLimit(),
			Count: groupCount,
		})
		return
	}

	entities := make([]*responseGetGamesEntity, 0, groupCount)

	for id, group := range h.groupManager.Groups() {
		entities = append(entities, &responseGetGamesEntity{
			ID:     id,
			Limit:  group.GetLimit(),
			Count:  group.GetCount(),
			Width:  int(group.GetWorldWidth()),
			Height: int(group.GetWorldHeight()),
			Rate:   group.GetRate(),
		})
	}

	entities = sortGameEntities(sorting, entities)

	if flagUseLimit && limit < len(entities) {
		entities = entities[:limit]
	}

	h.writeResponseJSON(w, http.StatusOK, &responseGetGamesHandler{
		Games: entities,
		Limit: h.groupManager.GroupLimit(),
		Count: groupCount,
	})
}

func (h *getGamesHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error(ErrGetGamesHandler(err.Error()))
	}
}

func sortGameEntities(sorting gamesSorting, entities []*responseGetGamesEntity) []*responseGetGamesEntity {
	if sorting != gamesSortingSmart {
		return entities
	}

	emptyGameEntities := filterGameEntities(entities, func(entity *responseGetGamesEntity) bool {
		return entity.Count == 0
	})

	fullGameEntities := filterGameEntities(entities, func(entity *responseGetGamesEntity) bool {
		return entity.Count == entity.Limit
	})

	relevantGameEntities := filterGameEntities(entities, func(entity *responseGetGamesEntity) bool {
		return entity.Count > 0 && entity.Count < entity.Limit
	})

	sort.Slice(emptyGameEntities, func(i, j int) bool {
		return emptyGameEntities[i].Rate < emptyGameEntities[j].Rate
	})

	sort.Slice(fullGameEntities, func(i, j int) bool {
		return fullGameEntities[i].Limit < emptyGameEntities[j].Limit
	})

	sort.Slice(relevantGameEntities, func(i, j int) bool {
		return relevantGameEntities[i].Count < relevantGameEntities[j].Count
	})

	copy(entities, relevantGameEntities)
	copy(entities[len(relevantGameEntities):], emptyGameEntities)
	copy(entities[len(relevantGameEntities)+len(emptyGameEntities):], fullGameEntities)

	return entities
}

type gamesEntityFilter func(entity *responseGetGamesEntity) bool

func filterGameEntities(entities []*responseGetGamesEntity, filter gamesEntityFilter) []*responseGetGamesEntity {
	resultGameEntities := make([]*responseGetGamesEntity, 0)
	for _, entity := range entities {
		if filter(entity) {
			resultGameEntities = append(resultGameEntities, entity)
		}
	}
	return resultGameEntities
}
