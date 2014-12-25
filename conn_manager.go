package main

import (
	"encoding/json"
	"errors"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

type PoolFeatures struct {
	startStreamConn StartStreamConnFunc
	stopStreamConn  StopStreamConnFunc
	// startPlayer starts player
	startPlayer game.StartPlayerFunc
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

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting player")
	}

	// output is channel for transferring private game information
	// that is useful only for current player
	output, err := poolFeatures.startPlayer(acceptPlayerCommands(ws))
	if err != nil {
		return &errConnProcessing{err}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("starting private game stream")
	}

	startPrivateStream(ws, output)

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

func acceptPlayerCommands(ws *websocket.Conn) <-chan *game.Command {
	input := make(chan *game.Command)

	go func() {
		defer close(input)
		for {
			var msg *InputMessage
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				if err != io.EOF {
					glog.Errorln("cannot read data:", err)
				}
				break
			}

			if len(msg.Header) == 0 {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln("empty header")
				}
				continue
			}

			if msg.Header != HEADER_GAME {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln("unexpected header:", msg.Header)
				}
				continue
			}

			var cmd *game.Command
			if err := json.Unmarshal(msg.Data, &cmd); err != nil {
				glog.Errorln("cannot parse command:", err)
				continue
			}

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("accepted command:", cmd.Command)
			}

			input <- cmd
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("player listener finished")
		}
	}()

	return input
}

func startPrivateStream(ws *websocket.Conn,
	output <-chan interface{}) {
	for data := range output {
		buffer, err := json.Marshal(&OutputMessage{
			HEADER_GAME, data,
		})
		if err != nil {
			glog.Errorln("cannot marshal private game data:", err)
			continue
		}

		_, err = ws.Write(buffer)
		if err != nil {
			glog.Errorln("cannot send private game data:", err)
			break
		}
	}

	// Wait for closing output channel
	for range output {
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("private game stream finished")
	}
}
