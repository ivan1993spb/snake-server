package broadcast

import (
	"sync"
	"time"
)

const (
	broadcastMainChanBufferSize = 64
	broadcastChanBufferSize     = 64
	broadcastSendTimeout        = time.Millisecond * 100
)

type BroadcastMessage string

type GroupBroadcast struct {
	chStop chan struct{}
	chMain chan BroadcastMessage
	chs    []chan BroadcastMessage
	chsMux *sync.RWMutex

	flagStarted bool
}

func NewGroupBroadcast() *GroupBroadcast {
	return &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan BroadcastMessage, broadcastMainChanBufferSize),
		chs:    make([]chan BroadcastMessage, 0),
		chsMux: &sync.RWMutex{},
	}
}

func (gb *GroupBroadcast) BroadcastMessageTimeout(message BroadcastMessage, timeout time.Duration) bool {
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

func (gb *GroupBroadcast) BroadcastMessage(message BroadcastMessage) {
	select {
	case gb.chMain <- message:
	case <-gb.chStop:
	}
}

func (gb *GroupBroadcast) Start(stop <-chan struct{}) {
	if gb.flagStarted {
		return
	}
	gb.flagStarted = true

	go func() {
		select {
		case <-stop:
		}
		gb.stop()
	}()

	go func() {
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
	}()
}

func (gb *GroupBroadcast) broadcast(message BroadcastMessage) {
	gb.chsMux.RLock()
	defer gb.chsMux.RUnlock()

	for _, ch := range gb.chs {
		select {
		case ch <- message:
		case <-gb.chStop:
		}
	}
}

func (gb *GroupBroadcast) createChan() chan BroadcastMessage {
	ch := make(chan BroadcastMessage, broadcastChanBufferSize)

	gb.chsMux.Lock()
	gb.chs = append(gb.chs, ch)
	gb.chsMux.Unlock()

	return ch
}

func (gb *GroupBroadcast) deleteChan(ch chan BroadcastMessage) {
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

func (gb *GroupBroadcast) ListenMessages(stop <-chan struct{}, buffer uint) <-chan BroadcastMessage {
	ch := gb.createChan()
	chOut := make(chan BroadcastMessage, buffer)

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
				gb.send(chOut, message, stop, broadcastSendTimeout)
			}
		}
	}()

	return chOut
}

func (gb *GroupBroadcast) send(ch chan BroadcastMessage, message BroadcastMessage, stop <-chan struct{}, timeout time.Duration) {
	const tickSize = 5

	var timer = time.NewTimer(timeout)
	defer timer.Stop()

	var ticker = time.NewTicker(timeout / tickSize)
	defer ticker.Stop()

	if cap(ch) == 0 {
		select {
		case ch <- message:
		case <-gb.chStop:
		case <-stop:
		case <-timer.C:
		}
	} else {
		for {
			select {
			case ch <- message:
				return
			case <-gb.chStop:
				return
			case <-stop:
				return
			case <-timer.C:
				return
			case <-ticker.C:
				if len(ch) == cap(ch) {
					select {
					case <-ch:
					case ch <- message:
						return
					case <-stop:
						return
					case <-gb.chStop:
						return
					case <-timer.C:
						return
					}
				}
			}
		}
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
