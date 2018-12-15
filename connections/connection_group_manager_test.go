package connections

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func Test_ConnectionGroupManager_Add_GeneratesValidIDs(t *testing.T) {
	const groupLimit = 10
	const connsLimit = 100

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	m := &ConnectionGroupManager{
		groups:      map[int]*ConnectionGroup{},
		groupsMutex: &sync.RWMutex{},
		groupLimit:  groupLimit,
		connsLimit:  connsLimit,
		logger:      logger,
	}

	for i := 0; i < groupLimit; i++ {
		id, err := m.Add(&ConnectionGroup{
			limit:      connsLimit / groupLimit,
			counterMux: &sync.RWMutex{},
			game:       nil,
			broadcast:  nil,
			logger:     logger,
			chs:        nil,
			chsMux:     nil,
			stop:       nil,
			stopper:    nil,
		})
		require.Equal(t, i+firstGroupId, id, "unexpected group id")
		require.Nil(t, err, "unexpected error")
	}
}

func Test_ConnectionGroupManager_Add_GetErrGroupLimitReached(t *testing.T) {
	const groupLimit = 2
	const connsLimit = 100

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	m := &ConnectionGroupManager{
		groups: map[int]*ConnectionGroup{
			1: {
				limit:      connsLimit / groupLimit,
				counterMux: &sync.RWMutex{},
				game:       nil,
				broadcast:  nil,
				logger:     logger,
				chs:        nil,
				chsMux:     nil,
				stop:       nil,
				stopper:    nil,
			},
			2: {
				limit:      connsLimit / groupLimit,
				counterMux: &sync.RWMutex{},
				game:       nil,
				broadcast:  nil,
				logger:     logger,
				chs:        nil,
				chsMux:     nil,
				stop:       nil,
				stopper:    nil,
			},
		},
		groupsMutex: &sync.RWMutex{},
		groupLimit:  groupLimit,
		connsLimit:  connsLimit,
		logger:      logger,
	}

	id, err := m.Add(&ConnectionGroup{
		limit:      1,
		counterMux: &sync.RWMutex{},
		game:       nil,
		broadcast:  nil,
		logger:     logger,
		chs:        nil,
		chsMux:     nil,
		stop:       nil,
		stopper:    nil,
	})

	require.Zero(t, id)
	require.Equal(t, ErrGroupLimitReached, err)
}

func Test_ConnectionGroupManager_Add_GetErrConnsLimitReached(t *testing.T) {
	const groupLimit = 10
	const connsLimit = 100

	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	m := &ConnectionGroupManager{
		groups: map[int]*ConnectionGroup{
			1: {
				limit:      50,
				counterMux: &sync.RWMutex{},
				game:       nil,
				broadcast:  nil,
				logger:     logger,
				chs:        nil,
				chsMux:     nil,
				stop:       nil,
				stopper:    nil,
			},
			2: {
				limit:      50,
				counterMux: &sync.RWMutex{},
				game:       nil,
				broadcast:  nil,
				logger:     logger,
				chs:        nil,
				chsMux:     nil,
				stop:       nil,
				stopper:    nil,
			},
		},
		groupsMutex: &sync.RWMutex{},
		groupLimit:  groupLimit,
		connsLimit:  connsLimit,
		connsCount:  connsLimit,
		logger:      logger,
	}

	id, err := m.Add(&ConnectionGroup{
		limit:      10,
		counterMux: &sync.RWMutex{},
		game:       nil,
		broadcast:  nil,
		logger:     logger,
		chs:        nil,
		chsMux:     nil,
		stop:       nil,
		stopper:    nil,
	})

	require.Zero(t, id)
	require.Equal(t, ErrConnsLimitReached, err)
}
