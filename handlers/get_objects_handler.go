package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
)

const URLRouteGetObjects = "/games/{id}/objects"

const MethodGetObjects = http.MethodGet

const waitRequestTimeout = time.Second * 10

type responseGetObjectsHandler struct {
	Objects []interface{} `json:"objects"`
}

type responseGetObjectsHandlerError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type getObjectsHandler struct {
	logger       logrus.FieldLogger
	groupManager *connections.ConnectionGroupManager

	timers    map[int]time.Time
	timersMux *sync.Mutex
}

type ErrGetObjectsHandler string

func (e ErrGetObjectsHandler) Error() string {
	return "get objects handler error: " + string(e)
}

func NewGetObjectsHandler(logger logrus.FieldLogger, groupManager *connections.ConnectionGroupManager) http.Handler {
	return &getObjectsHandler{
		logger:       logger,
		groupManager: groupManager,

		timers:    map[int]time.Time{},
		timersMux: &sync.Mutex{},
	}
}

func (h *getObjectsHandler) unsafeDiscardOutdatedTimers() {
	for groupId, timer := range h.timers {
		if time.Since(timer) > waitRequestTimeout {
			delete(h.timers, groupId)
		}
	}
}

func (h *getObjectsHandler) unsafeSetupTimer(groupId int) {
	h.timers[groupId] = time.Now()
}

func (h *getObjectsHandler) getRetryAfterHeaderSeconds(groupId int) int {
	h.timersMux.Lock()
	defer h.timersMux.Unlock()

	h.unsafeDiscardOutdatedTimers()

	if timer, ok := h.timers[groupId]; ok {
		since := time.Since(timer)
		if since < waitRequestTimeout {
			return int(math.Ceil((waitRequestTimeout - since).Seconds()))
		}
	}

	h.unsafeSetupTimer(groupId)

	return 0
}

func (h *getObjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("parse game id error")
		h.writeResponseJSON(w, http.StatusBadRequest, &responseGetObjectsHandlerError{
			Code: http.StatusBadRequest,
			Text: "invalid game id",
		})
		return
	}

	h.logger.WithField("game", id).Infoln("game id received")

	group, err := h.groupManager.Get(id)
	if err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("cannot get game group")

		switch err {
		case connections.ErrNotFoundGroup:
			h.writeResponseJSON(w, http.StatusNotFound, &responseGetObjectsHandlerError{
				Code: http.StatusNotFound,
				Text: "game not found",
			})
		default:
			h.writeResponseJSON(w, http.StatusInternalServerError, &responseGetObjectsHandlerError{
				Code: http.StatusInternalServerError,
				Text: "unknown error",
			})
		}
		return
	}

	if retryAfterSeconds := h.getRetryAfterHeaderSeconds(id); retryAfterSeconds > 0 {
		seconds := strconv.Itoa(retryAfterSeconds)
		w.Header().Set("Retry-After", seconds)
		h.writeResponseJSON(w, http.StatusTooManyRequests, &responseGetObjectsHandlerError{
			Code: http.StatusTooManyRequests,
			Text: "retry after " + seconds + " second(s)",
		})
		return
	}

	h.writeResponseJSON(w, http.StatusOK, &responseGetObjectsHandler{
		Objects: group.GetObjects(),
	})
}

func (h *getObjectsHandler) writeResponseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(ErrGetObjectsHandler(err.Error())).Error("encode response error")
	}
}
