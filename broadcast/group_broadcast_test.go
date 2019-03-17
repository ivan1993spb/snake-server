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

func Test_GroupBroadcast_stop(t *testing.T) {
	gb := &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan Message, broadcastMainChanBufferSize),
		chs: []chan Message{
			make(chan Message, broadcastChanBufferSize),
			make(chan Message, broadcastChanBufferSize),
			make(chan Message, broadcastChanBufferSize),
			make(chan Message, broadcastChanBufferSize),
			make(chan Message, broadcastChanBufferSize),
		},
		chsMux: &sync.RWMutex{},

		flagStarted: true,
		startedMux:  &sync.Mutex{},
	}

	gb.stop()
	require.Len(t, gb.chs, 0)
}

func Test_GroupBroadcast_createChan(t *testing.T) {
	gb := &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan Message, broadcastMainChanBufferSize),
		chs:    make([]chan Message, 0),
		chsMux: &sync.RWMutex{},

		flagStarted: true,
		startedMux:  &sync.Mutex{},
	}

	gb.createChan()
	gb.createChan()
	gb.createChan()
	gb.createChan()

	require.Len(t, gb.chs, 4)

	for _, ch := range gb.chs {
		close(ch)
	}
}
