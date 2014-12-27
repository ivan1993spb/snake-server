package main

import (
	"encoding/json"
	"errors"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type PoolFeatures struct {
	startStreamConn StartStreamConnFunc
	stopStreamConn  StopStreamConnFunc
	// startPlayer starts player
	startPlayer game.StartPlayerFunc
	cxt         context.Context
}

type errConnProcessing struct {
	err error
}

func (e *errConnProcessing) Error() string {
	return "error of connection processing in connection manager: " +
		e.err.Error()
}

type ConnManager struct{}

func NewConnManager() pwshandler.ConnManager {
	return &ConnManager{}
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
			var msg *InputMessage
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				if err != io.EOF {
					glog.Errorln("connection error:", err)
				}
				break
			}

			if len(msg.Header) == 0 {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln("input message with empty header")
				}
				continue
			}

			if msg.Header != HEADER_GAME {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln(
						"input message with unexpected header:",
						msg.Header,
					)
				}
				continue
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

	// Starting private game stream

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

				buffer, err := json.Marshal(&OutputMessage{
					HEADER_GAME, data,
				})
				if err != nil {
					glog.Errorln(
						"cannot marshal private game data:",
						err,
					)
					continue
				}

				_, err = ws.Write(buffer)
				if err != nil {
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

	<-cxt.Done()

	return nil
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(ws *websocket.Conn, err error) {
	if err == nil {
		err = errors.New("passed nil errer for reporting")
	}

	glog.Errorln(err)

	err = websocket.JSON.Send(ws, &OutputMessage{HEADER_ERROR, err})

	if err != nil {
		glog.Error(err)
	}
}
