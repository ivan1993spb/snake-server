// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

type InputMessageHandler func(m *InputMessage)

type WebsocketWrapper struct {
	*websocket.Conn
	handlers map[string]InputMessageHandler
	Closed   <-chan struct{}
}

type errReceiveMessage struct {
	err error
}

func (e *errReceiveMessage) Error() string {
	return "receiving message error: " + e.err.Error()
}

func WrapWebsocket(ws *websocket.Conn) *WebsocketWrapper {
	var (
		handlers = make(map[string]InputMessageHandler)
		closec   = make(chan struct{})
	)

	go func() {
		for {
			var msg *InputMessage

			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				if err != io.EOF {
					glog.Errorln(&errReceiveMessage{err})
				}
				break
			}

			if len(msg.Header) == 0 {
				glog.Errorln(&errReceiveMessage{
					errors.New("empty message header")},
				)
				break
			}

			if _, exists := handlers[msg.Header]; !exists {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln("unexpected header")
				}
				break
			}

			go handlers[msg.Header](msg)
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("stoping message receiver")
		}

		close(closec)
	}()

	return &WebsocketWrapper{ws, handlers, closec}
}

func (ww *WebsocketWrapper) Send(header string, data interface{},
) error {
	return ww.SendMessage(&OutputMessage{header, data})
}

func (ww *WebsocketWrapper) SendMessage(msg *OutputMessage) error {
	err := websocket.JSON.Send(ww.Conn, msg)
	if err != nil {
		return fmt.Errorf("sending message error: %s", err)
	}

	return nil
}

func (ww *WebsocketWrapper) BindHandler(header string,
	handler InputMessageHandler) {
	ww.handlers[header] = handler
}
