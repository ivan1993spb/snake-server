package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetCapacity = "/capacity"

const MethodGetCapacity = http.MethodGet

type responseGetCapacityHandler struct {
	Capacity float32 `json:"capacity"`
}

type getCapacityHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager
}

type ErrGetCapacityHandler string

func (e ErrGetCapacityHandler) Error() string {
	return "get capacity handler error: " + string(e)
}

func NewGetCapacityHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getCapacityHandler{
		logger:       logger,
		groupManager: groupManager,
	}
}

func (h *getCapacityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(responseGetCapacityHandler{
		Capacity: h.groupManager.Capacity(),
	})
	if err != nil {
		h.logger.Error(ErrGetGameHandler(err.Error()))
	}
}
