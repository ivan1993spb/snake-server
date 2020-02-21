package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/middlewares"
)

func Test_GetGamesHandler_ServeHTTP_ReturnsBadRequestErrorWithInvalidLimit(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	handler := &getGamesHandler{
		logger:       logger,
		groupManager: groupManager,
	}

	r := mux.NewRouter()
	r.Path(URLRouteGetGames).Methods(MethodGetGames).Handler(handler)

	n := negroni.New(middlewares.NewRecovery(logger), middlewares.NewLogger(logger, "api"))
	n.UseHandler(r)

	invalidLimits := []string{"-1", "test", "-9999"}

	for i, limit := range invalidLimits {
		request := httptest.NewRequest(MethodGetGames, URLRouteGetGames, nil)
		q := request.URL.Query()
		q.Add(getFieldGamesLimit, limit)
		request.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()

		n.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusBadRequest, recorder.Code, "case number"+strconv.Itoa(i))
	}

	hook.Reset()
}

func Test_GetGamesHandler_ServeHTTP_ReturnsBadRequestErrorWithInvalidSorting(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	handler := &getGamesHandler{
		logger:       logger,
		groupManager: groupManager,
	}

	r := mux.NewRouter()
	r.Path(URLRouteGetGames).Methods(MethodGetGames).Handler(handler)

	n := negroni.New(middlewares.NewRecovery(logger), middlewares.NewLogger(logger, "api"))
	n.UseHandler(r)

	request := httptest.NewRequest(MethodGetGames, URLRouteGetGames, nil)
	q := request.URL.Query()
	q.Add(getFieldGamesSorting, "invalid")
	request.URL.RawQuery = q.Encode()

	recorder := httptest.NewRecorder()

	n.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	hook.Reset()
}
