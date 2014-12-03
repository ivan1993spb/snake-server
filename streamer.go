package main

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type ShiftingObject interface {
	Pack() string
	Updated() bool
}

type Streamer struct {
	delay          time.Duration
	subscriptions  map[ShiftingObject][]*websocket.Conn
	parentCxt, cxt context.Context
	cancel         context.CancelFunc
}

func NewStreamer(cxt context.Context, delay time.Duration,
) (*Streamer, error) {
	if err := cxt.Err(); err != nil {
		return nil, err
	}
	if delay == 0 {
		return nil, errors.New("Invalid delay")
	}

	return &Streamer{delay: delay,
		subscriptions: make(map[ShiftingObject][]*websocket.Conn),
		parentCxt:     cxt}, nil
}

func (s *Streamer) Subscribe(o ShiftingObject, c *websocket.Conn,
) error {
	if c == nil {
		return errors.New("Passed nil connection")
	}

	if !s.running() {
		s.start()
	}

	if _, ok := s.subscriptions[o]; ok {
		for _, conn := range s.subscriptions[o] {
			if conn == c {
				return errors.New("Connection already subscribed")
			}
		}
		s.subscriptions[o] = append(s.subscriptions[o], c)
	} else {
		s.subscriptions[o] = []*websocket.Conn{c}
	}

	return nil
}

func (s *Streamer) Unsubscribe(o ShiftingObject, c *websocket.Conn) {
	if _, ok := s.subscriptions[o]; ok {
		for i := range s.subscriptions[o] {
			if s.subscriptions[o][i] == c {
				s.subscriptions[o] = append(s.subscriptions[o][:i],
					s.subscriptions[o][i+1:]...)
				break
			}
		}

		if len(s.subscriptions[o]) == 0 {
			delete(s.subscriptions, o)
		}

		if len(s.subscriptions) == 0 && s.running() {
			s.stop()
		}
	}
}

func (s *Streamer) start() {
	if !s.running() {
		s.cxt, s.cancel = context.WithCancel(s.parentCxt)
		s.run()
	}
}

func (s *Streamer) stop() {
	if s.running() && s.cancel != nil {
		s.cancel()
		s.cxt, s.cancel = nil, nil
	}
}

func (s *Streamer) running() bool {
	return s.cxt == nil || s.cxt.Err() != nil
}

func (s *Streamer) run() {
	if !s.running() && len(s.subscriptions) > 0 {
		go func() {
			for {
				select {
				case <-s.cxt.Done():
					return
				case <-time.After(s.delay):
				}
				if len(s.subscriptions) > 0 {
					for o, conns := range s.subscriptions {
						if o.Updated() {
							for _, conn := range conns {
								err := conn.WriteMessage(
									websocket.TextMessage,
									[]byte(o.Pack()),
								)
								if err != nil {
									s.Unsubscribe(o, conn)
								}
							}
						}
					}
				}
			}
		}()
	}
}
