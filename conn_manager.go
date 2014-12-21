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
	startStream StartStreamFunc
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

	if game, ok := data.(*PoolFeatures); ok {
		if err := game.startStream(ws); err != nil {
			return &errConnProcessing{err}
		}

		cxt, cancel := context.WithCancel(game.poolContext)
		input := make(chan []byte)

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("Starting new player")
		}
		output, err := game.startPlayer(cxt, input)

		if err == nil {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Starting private game stream")
			}
			go func() {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					defer glog.Infoln("Private game stream stops")
				}
				for {
					select {
					case <-cxt.Done():
						return
					case data := <-output:
						if _, err := ws.Write(data); err != nil {
							if glog.V(INFOLOG_LEVEL_CONNS) {
								glog.Errorln(
									"Cannot send private game data:",
									err,
								)
							}
							return
						}
					}
				}
			}()

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Starting player listener")
			}
			go func() {
				buffer := make([]byte, INPUT_MAX_LENGTH)
				for {
					n, err := ws.Read(buffer)
					if err != nil {
						if err != io.EOF {
							glog.Errorln("Cannot read data:", err)
						}
						if glog.V(INFOLOG_LEVEL_CONNS) {
							glog.Infoln("Player listener stops")
						}

						cancel()
						return
					}
					input <- buffer[:n]
				}
			}()

			<-cxt.Done()

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("Finishing player")
			}

			close(input)

			return nil

		}

		return &errConnProcessing{err}
	}

	return &errConnProcessing{
		errors.New("Pool data was not received"),
	}
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(_ *websocket.Conn, err error) {
	if err == nil {
		err = errors.New("Passed nil errer for reporting")
	}
	glog.Errorln(err)
}
