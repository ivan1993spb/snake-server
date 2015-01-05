// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type ConnManager struct{}

func NewConnManager() pwshandler.ConnManager {
	return new(ConnManager)
}

type errConnProcessing struct {
	err error
}

func (e *errConnProcessing) Error() string {
	return "error of connection processing in connection manager: " +
		e.err.Error()
}

// Implementing pwshandler.ConnManager interface
func (*ConnManager) Handle(ws *websocket.Conn,
	data pwshandler.Environment) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("connection handler was started")
		defer glog.Infoln("connection handler was finished")
	}

	poolFeatures, ok := data.(*PoolFeatures)
	if !ok || poolFeatures == nil {
		return &errConnProcessing{
			errors.New("pool data was not received"),
		}
	}

	// Setup game stream
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("creating connection to common game stream")
	}
	if err := poolFeatures.startStreamConn(ws); err != nil {
		return &errConnProcessing{err}
	}
	defer func() {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("removing connection from common game stream")
		}
		if err := poolFeatures.stopStreamConn(ws); err != nil {
			glog.Errorln(&errConnProcessing{err})
		}
	}()

	// Pool context is parent of connection context
	cxt, cancel := context.WithCancel(poolFeatures.cxt)

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                   BEGIN COMMAND ACCEPTER                    *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Channel for player commands
	input := make(chan *game.Command)

	// Starting command accepter

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting command accepter")
	}

	go func() {
		for {
			msg, err := ReceiveMessage(ws, HEADER_GAME)
			if err != nil {
				if err != io.EOF {
					glog.Errorln("cannot receive player command:",
						err)
				}
				break
			}

			var cmd *game.Command
			if err := json.Unmarshal(msg.Data, &cmd); err != nil {
				glog.Errorln("cannot parse player command:", err)
				continue
			}

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("accepted command:", cmd.Command)
			}

			input <- cmd
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("command accepter finished")
		}

		close(input)
		// Canceling connection context
		cancel()
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                    END COMMAND ACCEPTER                     *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Starting player

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting player")
	}

	// output is channel for transferring private game information
	// that is useful only for current player
	output, err := poolFeatures.startPlayer(cxt, input)
	if err != nil {
		return &errConnProcessing{err}
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                   BEGIN PRIVATE STREAM                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Starting private stream

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting private game stream")
	}

	go func() {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			defer glog.Infoln("private game stream finished")
		}
		for {
			select {
			case <-cxt.Done():
				return
			case data := <-output:
				if data == nil {
					continue
				}

				if err :=
					SendMessage(ws, HEADER_GAME, data); err != nil {
					glog.Errorln(&errConnProcessing{fmt.Errorf(
						"cannot send private game data: %s", err)})
					return
				}
			}
		}
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                     END PRIVATE STREAM                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	<-cxt.Done()

	return nil
}

type errErrorHandling struct {
	err error
}

func (e *errErrorHandling) Error() string {
	return "error of error handling: " + e.err.Error()
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(ws *websocket.Conn, err error) {
	if err == nil {
		err = &errErrorHandling{
			errors.New("passed nil errer for reporting"),
		}
	}

	glog.Errorln(err)

	if err = SendMessage(ws, HEADER_ERROR, err.Error()); err != nil {
		glog.Errorln(&errErrorHandling{err})
	}
}
