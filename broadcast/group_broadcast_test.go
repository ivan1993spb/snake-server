package broadcast

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GroupBroadcast_Start(t *testing.T) {
	gb := &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan Message, broadcastMainChanBufferSize),
		chs:    make([]chan Message, 0),
		chsMux: &sync.RWMutex{},

		flagStarted: false,
		startedMux:  &sync.Mutex{},
	}

	stop := make(chan struct{})
	defer close(stop)

	gb.Start(stop)
	require.True(t, gb.flagStarted)
}
