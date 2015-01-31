// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type Stream struct {
	cxt   context.Context
	conns []*WebsocketWrapper
	src   chan *OutputMessage
	wg    *sync.WaitGroup
}

func NewStream(cxt context.Context, src <-chan *OutputMessage) (
	*Stream, error) {
	if err := cxt.Err(); err != nil {
		return nil, fmt.Errorf("cannot create stream: %s", err)
	}

	wg := new(sync.WaitGroup)

	scxt, cancel := context.WithCancel(cxt)

	s := &Stream{scxt, make([]*WebsocketWrapper, 0),
		make(chan *OutputMessage), wg}
	s.AddSource(src)

	go func() {
		wg.Wait()
		cancel()
	}()

	go func() {
		if glog.V(INFOLOG_LEVEL_POOLS) {
			defer glog.Infoln("common game stream finished")
		}

		for msg := range s.src {
			select {
			case <-scxt.Done():
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

func (s *Stream) AddConn(ww *WebsocketWrapper) error {
	for i := range s.conns {
		if s.conns[i] == ww {
			return errors.New("cannot create connection to common" +
				" pool stream: passed connection already exists")
		}
	}

	s.conns = append(s.conns, ww)

	return nil
}

func (s *Stream) DelConn(ww *WebsocketWrapper) error {
	for i := range s.conns {
		if s.conns[i] == ww {
			s.conns = append(s.conns[:i], s.conns[i+1:]...)
			return nil
		}
	}

	return errors.New("cannot delete connection from common pool" +
		" stream: passed connection was not found")
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
