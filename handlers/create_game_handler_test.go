package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/middlewares"
)

func Test_CreateGameHandler_ServeHTTP_CreatesGroup(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	handler := NewCreateGameHandler(logger, groupManager)

	r := mux.NewRouter()
	r.Path(URLRouteCreateGame).Methods(MethodCreateGame).Handler(handler)

	n := negroni.New(middlewares.NewRecovery(logger), middlewares.NewLogger(logger, "api"))
	n.UseHandler(r)

	data := &url.Values{}
	data.Add(postFieldConnectionLimit, "10")
	data.Add(postFieldMapWidth, "100")
	data.Add(postFieldMapHeight, "100")

	request := httptest.NewRequest(MethodCreateGame, URLRouteCreateGame, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()

	n.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusCreated, recorder.Code)

	require.Len(t, groupManager.Groups(), 1)
	group, err := groupManager.Get(1)
	require.Nil(t, err)
	require.NotNil(t, group)
	require.Equal(t, 10, group.GetLimit())
	require.Nil(t, groupManager.Delete(group))

	hook.Reset()
}
