// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// Stream is common pool data stream
type Stream struct {
	cxt   context.Context
	conns []*WebsocketWrapper
	src   chan *OutputMessage
	wg    *sync.WaitGroup
}

func NewStream(cxt context.Context) (*Stream, error) {
	if err := cxt.Err(); err != nil {
		return nil, fmt.Errorf("cannot create stream: %s", err)
	}

	s := &Stream{
		cxt,
		make([]*WebsocketWrapper, 0),
		make(chan *OutputMessage),
		new(sync.WaitGroup),
	}

	go func() {
		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("starting pool stream")
			defer glog.Infoln("finishing pool stream")
		}

		for msg := range s.src {
			select {
			case <-cxt.Done():
				return
			default:
			}

			if len(s.conns) == 0 {
				continue
			}

			for i := 0; i < len(s.conns); {
				if err := s.conns[i].SendMessage(msg); err != nil {
					// Remove connection on error
					glog.Errorln("cannot send common pool data:", err)

					if glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Infoln("removing connection from stream")
					}

					s.conns = append(s.conns[:i], s.conns[i+1:]...)
				} else {
					i++
				}
			}
		}
	}()

	return s, nil
}

func (s *Stream) AddConn(ww *WebsocketWrapper) {
	for i := range s.conns {
		if s.conns[i] == ww {
			return
		}
	}

	s.conns = append(s.conns, ww)
}

func (s *Stream) DelConn(ww *WebsocketWrapper) {
	for i := range s.conns {
		if s.conns[i] == ww {
			s.conns = append(s.conns[:i], s.conns[i+1:]...)
			return
		}
	}
}

func (s *Stream) AddSource(src <-chan *OutputMessage) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for data := range src {
			select {
			case <-s.cxt.Done():
				return
			default:
				s.src <- data
			}
		}
	}()
}

func (s *Stream) AddSourceHeader(header string,
	src <-chan interface{}) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for data := range src {
			select {
			case <-s.cxt.Done():
				return
			default:
				s.src <- &OutputMessage{header, data}
			}
		}
	}()
}

func (s *Stream) Stop() {
	s.wg.Wait()
	close(s.src)
}
