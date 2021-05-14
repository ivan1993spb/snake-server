package handlers

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/connections"
)

func Test_getObjectsHandler_unsafeDiscardOutdatedTimers_discards(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	var (
		timer1 = time.Now().Add(-waitRequestTimeout / 2)
		timer2 = time.Now().Add(-waitRequestTimeout / 3)
		timer3 = time.Now().Add(-waitRequestTimeout / 4)
		timer4 = time.Now().Add(-waitRequestTimeout * 10)
		timer5 = time.Now().Add(-waitRequestTimeout * 2)
		timer6 = time.Now().Add(-(waitRequestTimeout + time.Minute))
		timer7 = time.Now().Add(-(waitRequestTimeout + time.Hour))
		timer8 = time.Now().Add(-waitRequestTimeout / 10)
		timer9 = time.Now().Add(-waitRequestTimeout)
	)

	handler := &getObjectsHandler{
		logger:       logger,
		groupManager: groupManager,

		timers: map[int]time.Time{
			1: timer1,
			2: timer2,
			3: timer3,
			4: timer4,
			5: timer5,
			6: timer6,
			7: timer7,
			8: timer8,
			9: timer9,
		},
		timersMux: &sync.Mutex{},
	}

	expected := map[int]time.Time{
		1: timer1,
		2: timer2,
		3: timer3,
		8: timer8,
	}

	handler.unsafeDiscardOutdatedTimers()

	require.Equal(t, expected, handler.timers)
}

func Test_getObjectsHandler_unsafeSetupTimer_setupTimers(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	handler := &getObjectsHandler{
		logger:       logger,
		groupManager: groupManager,

		timers:    map[int]time.Time{},
		timersMux: &sync.Mutex{},
	}

	for i := 1; i < 100; i++ {
		{
			_, exists := handler.timers[i]
			require.False(t, exists)
		}

		handler.unsafeSetupTimer(i)

		{
			_, exists := handler.timers[i]
			require.True(t, exists)
		}

		require.Len(t, handler.timers, i)
	}
}

func Test_getObjectsHandler_writeResponseJSON_writesCorrectResponseJSON(t *testing.T) {
	const groupsLimit = 5
	const connsLimit = 10

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	require.Nil(t, err)
	require.NotNil(t, groupManager)

	handler := &getObjectsHandler{
		logger:       logger,
		groupManager: groupManager,

		timers:    map[int]time.Time{},
		timersMux: &sync.Mutex{},
	}

	response := struct {
		Code   int     `json:"code"`
		Text   string  `json:"text"`
		Msg    string  `json:"msg"`
		Income float32 `json:"income"`
	}{
		Code:   312,
		Text:   "This is a description",
		Msg:    "Some message",
		Income: 31.2,
	}
	expectedJSON := `{"code":312,"text":"This is a description","msg":"Some message","income":31.2}` + "\n"

	recorder := httptest.NewRecorder()
	handler.writeResponseJSON(recorder, http.StatusCreated, response)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, expectedJSON, recorder.Body.String())
	require.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
}
