package main

import (
	"errors"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

const INPUT_MAX_LENGTH = 512

type PoolFeatures struct {
	startStreamConn StartStreamConnFunc
	stopStreamConn  StopStreamConnFunc
	// startPlayer starts player
	startPlayer game.StartPlayerFunc
	poolContext context.Context
}

type errConnProcessing struct {
	err error
}

func (e *errConnProcessing) Error() string {
	return "Error of connection processing in connection manager: " +
		e.err.Error()
}

type ConnManager struct{}

func NewConnManager() pwshandler.ConnManager {
	return &ConnManager{}
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) Handle(ws *websocket.Conn,
	data pwshandler.Environment) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Websocket handler was started")
		defer glog.Infoln("Websocket handler was finished")
	}

	game, ok := data.(*PoolFeatures)
	if !ok {
		return &errConnProcessing{
			errors.New("Pool data was not received"),
		}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Creating connection to common game stream")
	}
	if err := game.startStreamConn(ws); err != nil {
		return &errConnProcessing{err}
	}

	// input is channel for transferring information from client to
	// player goroutine, for example: player commands
	input := make(chan []byte)

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Starting player")
	}

	// output is channel for transferring private game information
	// for only one player. This information are useful only for
	// current player
	output, err := game.startPlayer(input)
	if err != nil {
		return &errConnProcessing{err}
	}

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Starting private game stream")
	}
	// Send game data which are useful only for current player
	go func() {
		for data := range output {
			if _, err := ws.Write(data); err != nil {
				glog.Warningln("Cannot send private game data:", err)
				break
			}
		}

		// Wait for closing output channel
		for range output {
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("Private game stream finished")
		}
	}()

	stop := make(chan struct{})

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Starting player listener")
	}
	// Listen for player commands
	go func() {
		buffer := make([]byte, INPUT_MAX_LENGTH)

		for {
			n, err := ws.Read(buffer)
			if err != nil {
				if err != io.EOF {
					glog.Errorln("Cannot read data:", err)
				}
				break
			}
			input <- buffer[:n]
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("Player listener finished")
		}

		close(stop)
	}()

	select {
	case <-stop:
	case <-game.poolContext.Done():
	}

	close(input)

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Removing connection from common game stream")
	}
	if err := game.stopStreamConn(ws); err != nil {
		return &errConnProcessing{err}
	}

	return nil

}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(_ *websocket.Conn, err error) {
	if err == nil {
		err = errors.New("Passed nil errer for reporting")
	}
	glog.Errorln(err)
}
