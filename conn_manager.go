package main

import (
	"errors"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

const INPUT_MAX_LENGTH = 512

type PoolFeatures struct {
	startStream StartStreamFunc
	startPlayer game.StartPlayerFunc
}

type errConnHandling struct {
	err error
}

func (e *errConnHandling) Error() string {
	return "Error of connection handling: " + e.err.Error()
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

	if game, ok := data.(*PoolFeatures); ok {
		if err := game.startStream(ws); err != nil {
			return &errConnHandling{err}
		}

		if err := game.startPlayer(StartListen(ws)); err != nil {
			return &errConnHandling{err}
		}

		return nil
	}

	return &errConnHandling{errors.New("Pool data was not received")}
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(_ *websocket.Conn, err error) {
	if err == nil {
		err = errors.New("Passed nil errer for reporting")
	}
	glog.Errorln(err)
}

func StartListen(ws *websocket.Conn) <-chan []byte {
	input := make(chan []byte)

	go func() {
		buffer := make([]byte, INPUT_MAX_LENGTH)
		for {
			n, err := ws.Read(buffer)
			if err != nil {
				if err != io.EOF {
					glog.Errorln(&errConnHandling{err})
				}
				close(input)
				return
			}
			input <- buffer[:n]
		}
	}()

	return input
}
