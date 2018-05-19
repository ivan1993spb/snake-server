package connections

import (
	"sync"
	"time"
)

const (
	broadcastMainChanBufferSize = 64
	broadcastChanBufferSize     = 32
	broadcastSendTimeout        = time.Millisecond * 100
)

type GroupBroadcast struct {
	chStop chan struct{}
	chMain chan OutputMessage
	chs    []chan OutputMessage
	chsMux *sync.RWMutex

	flagStarted bool
}

func NewGroupBroadcast() *GroupBroadcast {
	return &GroupBroadcast{
		chStop: make(chan struct{}),
		chMain: make(chan OutputMessage, broadcastMainChanBufferSize),
		chs:    make([]chan OutputMessage, 0),
		chsMux: &sync.RWMutex{},
	}
}

func (gb *GroupBroadcast) BroadcastMessage(message OutputMessage) {
	select {
	case gb.chMain <- message:
	case <-gb.chStop:
	}
}

func (gb *GroupBroadcast) Start() {
	if gb.flagStarted {
		return
	}
	gb.flagStarted = true

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

func (gb *GroupBroadcast) broadcast(message OutputMessage) {
	gb.chsMux.RLock()
	defer gb.chsMux.RUnlock()

	for _, ch := range gb.chs {
		select {
		case ch <- message:
		case <-gb.chStop:
		}
	}
}

func (gb *GroupBroadcast) createChan() chan OutputMessage {
	ch := make(chan OutputMessage, broadcastChanBufferSize)

	gb.chsMux.Lock()
	gb.chs = append(gb.chs, ch)
	gb.chsMux.Unlock()

	return ch
}

func (gb *GroupBroadcast) deleteChan(ch chan OutputMessage) {
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

func (gb *GroupBroadcast) OutputMessages(stop <-chan struct{}, buffer uint) <-chan OutputMessage {
	ch := gb.createChan()
	chOut := make(chan OutputMessage, buffer)

	go func() {
		defer close(chOut)
		defer gb.deleteChan(ch)

		for {
			select {
			case <-stop:
				return
			case <-gb.chStop:
				return
			case message := <-ch:
				gb.send(chOut, message, stop, broadcastSendTimeout)
			}
		}
	}()

	return chOut
}

func (gb *GroupBroadcast) send(ch chan OutputMessage, message OutputMessage, stop <-chan struct{}, timeout time.Duration) {
	var timer = time.NewTimer(timeout)
	defer timer.Stop()
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
			default:
				if len(ch) == cap(ch) {
					<-ch
				}
			}
		}
	}
}

func (gb *GroupBroadcast) Stop() {
	close(gb.chStop)
	close(gb.chMain)

	gb.chsMux.Lock()
	defer gb.chsMux.Unlock()

	for _, ch := range gb.chs {
		close(ch)
	}

	gb.chs = gb.chs[:0]
}
