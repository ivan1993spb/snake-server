package main

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// Playground wrapper
type Playground interface {
	Pack() string
	Updated() bool
}

type stream struct {
	playground  Playground
	subscribers []*websocket.Conn
}

func newStream(pg Playground, firstWs *websocket.Conn) *stream {
	return &stream{pg, []*websocket.Conn{firstWs}}
}

func (s *stream) addSubscriber(ws *websocket.Conn) {
	if !s.connExists(ws) {
		s.subscribers = append(s.subscribers, ws)
	}
}

func (s *stream) delSubscriber(i int) {
	if i > -1 && len(s.subscribers) > i {
		s.subscribers = append(
			s.subscribers[:i],
			s.subscribers[i+1:]...,
		)
	}
}

func (s *stream) connExists(ws *websocket.Conn) bool {
	return s.connIndex(ws) > -1
}

func (s *stream) connIndex(ws *websocket.Conn) int {
	for i := range s.subscribers {
		if s.subscribers[i] == ws {
			return i
		}
	}
	return -1
}

func (s *stream) push() {
	if s.playground.Updated() {
		data := "playground:" + s.playground.Pack()
		for i := 0; i < len(s.subscribers); {
			err := websocket.Message.Send(s.subscribers[i], data)
			if err != nil {
				s.delSubscriber(i)
			} else {
				i++
			}
		}
	}
}

type Streamer struct {
	delay    time.Duration
	streams  []*stream
	pingPong chan chan struct{}
	parCxt   context.Context    // Parent context
	cancel   context.CancelFunc // Cancel func of child context
}

func NewStreamer(cxt context.Context, delay time.Duration,
) (*Streamer, error) {
	if err := cxt.Err(); err != nil {
		return nil, err
	}
	if delay <= 0 {
		return nil, errors.New("Invalid delay")
	}

	return &Streamer{
		delay:    delay,
		streams:  make([]*stream, 0),
		pingPong: make(chan chan struct{}),
		parCxt:   cxt,
	}, nil
}

func (s *Streamer) delStream(i int) {
	if i > -1 && len(s.streams) > i {
		s.streams = append(s.streams[:i], s.streams[i+1:]...)
	}
}

func (s *Streamer) Subscribe(pg Playground, ws *websocket.Conn) {
	if ws == nil {
		return
	}

	defer func() {
		if !s.running() {
			if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
				glog.Infoln("Starting streamer")
			}
			s.start()
		}
	}()

	for _, stream := range s.streams {
		if stream.playground == pg {
			if !stream.connExists(ws) {
				if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
					glog.Infoln("Creating new subscriber to stream")
				}
				stream.addSubscriber(ws)
			}
			return
		}
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Creating new subscriber to new stream")
	}
	s.streams = append(
		s.streams,
		newStream(pg, ws),
	)
}

func (s *Streamer) Unsubscribe(pg Playground, ws *websocket.Conn) {
	for i := range s.streams {
		if s.streams[i].playground == pg {
			if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
				glog.Infoln("Necessary stream was found")
			}

			if j := s.streams[i].connIndex(ws); j > -1 {
				if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
					glog.Infoln("Subscriber was found")
					glog.Infoln("Removing subscriber from stream")
				}

				s.streams[i].delSubscriber(j)

				if len(s.streams[i].subscribers) == 0 {
					if glog.V(INFOLOG_LEVEL_ABOUT_STREAMS) {
						glog.Infoln(
							"Stream has no subscribers.",
							"Removing stream",
						)
					}
					s.delStream(i)
				}

				if len(s.streams) == 0 && s.running() {
					if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
						glog.Infoln(
							"Streamer is empty.",
							"Stoping streamer",
						)
					}
					s.stop()
				}
				return
			}
		}
	}
}

func (s *Streamer) start() {
	if !s.running() {
		var cxt context.Context
		cxt, s.cancel = context.WithCancel(s.parCxt)
		s.run(cxt)
	}
}

func (s *Streamer) stop() {
	if s.running() && s.cancel != nil {
		s.cancel()
	}
}

func (s *Streamer) running() bool {
	var ch = make(chan struct{})
	go func() {
		s.pingPong <- ch
	}()
	select {
	case <-ch:
		return true
	case <-time.After(s.delay):
		<-s.pingPong
	}
	close(ch)
	return false
}

func (s *Streamer) run(cxt context.Context) {
	if len(s.streams) == 0 {
		return
	}

	if s.running() {
		return
	}

	go func() {
		var t = time.Tick(s.delay)

		for {
			select {
			case <-cxt.Done():
				if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
					glog.Infoln(
						"Stopping streamer:",
						"context was canceled",
					)
				}
				return
			case ch := <-s.pingPong:
				ch <- struct{}{}
				continue
			case <-t:
			}

			if len(s.streams) == 0 {
				if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
					glog.Infoln(
						"Stopping streamer:",
						"there is no one stream",
					)
				}
				return
			}

			for i := 0; i < len(s.streams); {
				if len(s.streams[i].subscribers) == 0 {
					s.delStream(i)
					continue
				}

				if s.streams[i].playground.Updated() {
					s.streams[i].push()
				}

				i++
			}
		}
	}()
}
