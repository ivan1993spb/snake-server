package broadcast

import (
	"sync"
	"time"
)

const (
	broadcastMainChanBufferSize = 128
	broadcastChanBufferSize     = 128

	broadcastSendTimeout = time.Millisecond
)

type Message string

type GroupBroadcast struct {
	chStop chan struct{}
	chMain chan Message
	chs    []chan Message
	chsMux *sync.RWMutex

	flagStarted bool
	startedMux  *sync.Mutex
}

func NewGroupBroadcast() *GroupBroadcast {
	return &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan Message, broadcastMainChanBufferSize),
		chs:    make([]chan Message, 0),
		chsMux: &sync.RWMutex{},

		flagStarted: false,
		startedMux:  &sync.Mutex{},
	}
}

func (gb *GroupBroadcast) BroadcastMessageTimeout(message Message, timeout time.Duration) bool {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case gb.chMain <- message:
		return true
	case <-gb.chStop:
	case <-timer.C:
	}

	return false
}

func (gb *GroupBroadcast) BroadcastMessage(message Message) {
	select {
	case gb.chMain <- message:
	case <-gb.chStop:
	}
}

func (gb *GroupBroadcast) Start(stop <-chan struct{}) {
	gb.startedMux.Lock()
	defer gb.startedMux.Unlock()

	if gb.flagStarted {
		return
	}
	gb.flagStarted = true

	go gb.listenStopChan(stop)

	go gb.listenBroadcastMessage()
}

func (gb *GroupBroadcast) listenBroadcastMessage() {
	for {
		select {
		case message, ok := <-gb.chMain:
			if !ok {
				return
			}
			gb.broadcast(message)
		case <-gb.chStop:
			return
		}
	}
}

func (gb *GroupBroadcast) listenStopChan(stop <-chan struct{}) {
	select {
	case <-stop:
	}
	gb.stop()
}

func (gb *GroupBroadcast) broadcast(message Message) {
	gb.chsMux.RLock()
	defer gb.chsMux.RUnlock()

	for _, ch := range gb.chs {
		select {
		case ch <- message:
		case <-gb.chStop:
		}
	}
}

func (gb *GroupBroadcast) createChan() chan Message {
	ch := make(chan Message, broadcastChanBufferSize)

	gb.chsMux.Lock()
	gb.chs = append(gb.chs, ch)
	gb.chsMux.Unlock()

	return ch
}

func (gb *GroupBroadcast) deleteChan(ch chan Message) {
	go func() {
		for range ch {
		}
	}()

	gb.chsMux.Lock()
	for i := range gb.chs {
		if gb.chs[i] == ch {
			gb.chs = append(gb.chs[:i], gb.chs[i+1:]...)
			close(ch)
			break
		}
	}
	gb.chsMux.Unlock()
}

func (gb *GroupBroadcast) ListenMessages(stop <-chan struct{}, buffer uint) <-chan Message {
	ch := gb.createChan()
	chOut := make(chan Message, buffer)

	go func() {
		defer close(chOut)
		defer gb.deleteChan(ch)

		for {
			select {
			case <-stop:
				return
			case <-gb.chStop:
				return
			case message, ok := <-ch:
				if !ok {
					return
				}
				gb.sendMessageTimeout(chOut, message, stop, broadcastSendTimeout)
			}
		}
	}()

	return chOut
}

func (gb *GroupBroadcast) sendMessageTimeout(ch chan Message, message Message, stop <-chan struct{}, timeout time.Duration) {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- message:
	case <-gb.chStop:
	case <-stop:
	case <-timer.C:
	}
}

func (gb *GroupBroadcast) stop() {
	close(gb.chStop)
	close(gb.chMain)

	gb.chsMux.Lock()
	defer gb.chsMux.Unlock()

	for _, ch := range gb.chs {
		close(ch)
	}

	gb.chs = gb.chs[:0]
}
