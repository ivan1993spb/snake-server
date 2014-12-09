package main

import (
	"errors"
	"io"
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

func newStream(pg Playground) *stream {
	return &stream{pg, make([]*websocket.Conn, 0, 1)}
}

func (s *stream) addSubscriber(ws *websocket.Conn) error {
	if !s.hasSubscriber(ws) {
		s.subscribers = append(s.subscribers, ws)
		return nil
	}

	return errors.New(
		"Passed connection already subscribed to this stream")
}

func (s *stream) delSubscriber(ws *websocket.Conn) error {
	if s.hasSubscriber(ws) {
		for i := range s.subscribers {
			if s.subscribers[i] == ws {
				s.subscribers = append(
					s.subscribers[:i],
					s.subscribers[i+1:]...,
				)

				return nil
			}
		}
	}

	return errors.New("Cannot remove not found subscriber")
}

func (s *stream) hasSubscriber(ws *websocket.Conn) bool {
	for i := range s.subscribers {
		if s.subscribers[i] == ws {
			return true
		}
	}
	return false
}

func (s *stream) isEmpty() bool {
	return len(s.subscribers) == 0
}

func (s *stream) pushData() {
	if s.playground.Updated() {
		// Data for streaming
		var data = s.playground.Pack()

		for _, ws := range s.subscribers {
			if err := websocket.Message.Send(ws, data); err != nil {
				if err != io.EOF {
					if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
						glog.Error("Connection error:", err)
					}
				}

				s.delSubscriber(ws)
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

func (s *Streamer) getStreamWithPlayground(pg Playground) *stream {
	var stm = s.getStreamByPlayground(pg)
	if stm != nil {
		return stm
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Creating new stream")
	}

	stm = newStream(pg)
	s.streams = append(s.streams, stm)

	return stm
}

func (s *Streamer) getStreamByPlayground(pg Playground) *stream {
	for i := range s.streams {
		if s.streams[i].playground == pg {
			return s.streams[i]
		}
	}

	return nil
}

func (s *Streamer) delStream(sm *stream) error {
	if !sm.isEmpty() {
		return errors.New("Passed to removing stream is not empty")
	}

	for i := range s.streams {
		if s.streams[i] == sm {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return nil
		}
	}

	return errors.New("Passed to removing stream was not found")
}

func (s *Streamer) isEmpty() bool {
	return len(s.streams) == 0
}

func (s *Streamer) Subscribe(pg Playground, ws *websocket.Conn,
) error {
	if ws == nil {
		return errors.New("Passed nil connection")
	}

	defer func() {
		if !s.isEmpty() && !s.running() {
			if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
				glog.Infoln("Starting streamer")
			}
			if err := s.start(); err != nil {
				glog.Errorln("Cannot start stream:", err)
			}
		}
	}()

	var stm = s.getStreamWithPlayground(pg)
	if stm.hasSubscriber(ws) {
		return errors.New("Connection already subscribed to stream")
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Creating subscriber to stream")
	}

	stm.addSubscriber(ws)

	return nil
}

func (s *Streamer) Unsubscribe(pg Playground, ws *websocket.Conn,
) error {
	if pg == nil {
		return errors.New("Cannot subscribe to nil playground")
	}
	if ws == nil {
		return errors.New("Passed nil connection")
	}

	var stm = s.getStreamByPlayground(pg)
	if stm == nil {
		return errors.New("Stream with passed playground not found")
	}

	if !stm.hasSubscriber(ws) {
		return errors.New("Subscriber was not found")
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Subscriber was found")
	}
	if err := stm.delSubscriber(ws); err != nil {
		return err
	}
	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Subscriber was removed")
	}

	if stm.isEmpty() {
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln(
				"Stream has no subscribers. Removing stream",
			)
		}
		if err := s.delStream(stm); err != nil {
			return err
		}
	}

	if s.isEmpty() && s.running() {
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln("Streamer is empty. Stoping streamer")
		}
		if err := s.stop(); err != nil {
			glog.Errorln("Cannot stop streamer:", err)
		}
	}

	return nil
}

func (s *Streamer) start() error {
	if s.running() {
		return errors.New("Streamer already started")
	}
	var cxt context.Context
	cxt, s.cancel = context.WithCancel(s.parCxt)
	return s.run(cxt)
}

func (s *Streamer) stop() error {
	if s.cancel == nil {
		return errors.New("CancelFunc is nil")
	}
	s.cancel()
	return nil
}

func (s *Streamer) running() bool {
	var ch = make(chan struct{})
	defer close(ch)

	go func() { s.pingPong <- ch }()

	select {
	case <-ch:
		return true
	case <-time.After(s.delay):
		<-s.pingPong
	}

	return false
}

func (s *Streamer) run(cxt context.Context) error {

	if s.running() {
		return errors.New("Streamer already started")
	}

	go func() {
		var ticker = time.Tick(s.delay)
		for {
			select {
			case <-cxt.Done():
				if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
					glog.Infoln("Streamer was stopped")
				}
				return
			case ch := <-s.pingPong:
				ch <- struct{}{}
				time.Sleep(s.delay / 2)
			case <-ticker:
			}

			if !s.isEmpty() {
				for _, stm := range s.streams {
					if stm.isEmpty() {
						s.delStream(stm)
						continue
					}

					stm.pushData()
				}
			}
		}
	}()

	return nil
}
