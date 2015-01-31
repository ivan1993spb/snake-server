// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
)

type GameConnManager struct{}

func (*GameConnManager) Handle(ww *WebsocketWrapper,
	poolFeatures *PoolFeatures) error {

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("connection handler was started")
		defer glog.Infoln("connection handler was finished")
	}

	// Setup game stream

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("creating connection to common game stream")
	}
	poolFeatures.stream.AddConn(ww)

	defer func() {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("removing connection from common game stream")
		}
		poolFeatures.stream.DelConn(ww)
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                   BEGIN COMMAND ACCEPTER                    *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Channel for player commands
	input := make(chan *game.Command)
	defer close(input)

	// Starting command accepter

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting command accepter")
	}

	ww.BindHandler(HEADER_GAME, func(msg *InputMessage) {
		var cmd *game.Command
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			glog.Errorln("cannot parse player command:", err)
			return
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("accepted command")
		}

		input <- cmd
	})

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                    END COMMAND ACCEPTER                     *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Starting player

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting player")
	}

	// output is channel for transferring private game information
	// that is useful only for current player
	output, err := poolFeatures.startPlayer(poolFeatures.cxt, input)
	if err != nil {
		return fmt.Errorf("connection processing error: %s", err)
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
			case <-poolFeatures.cxt.Done():
				return
			case data := <-output:
				if data == nil {
					continue
				}

				if err := ww.Send(HEADER_GAME, data); err != nil {
					glog.Errorln(
						"cannot send private game data:",
						err,
					)
					return
				}
			}
		}
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                     END PRIVATE STREAM                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	select {
	case <-ww.Closed: // Waiting for connection closing
	case <-poolFeatures.cxt.Done(): // Waiting for pool finishing
	}

	return nil
}
