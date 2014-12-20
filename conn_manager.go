package main

import (
	"errors"

	"bitbucket.org/pushkin_ivan/clever-snake/game"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/websocket"
)

type GameData struct {
	Game *game.Game
	chWs <-chan *websocket.Conn
}

type ConnManager struct{}

func NewConnManager() pwshandler.ConnManager {
	return &ConnManager{}
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) Handle(ws *websocket.Conn,
	data pwshandler.Environment) error {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("Websocket handler started for new connection")
	}
	defer func() {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("Websocket handler was finished")
		}
	}()

	if _, ok := data.(*GameData); ok {
		// if glog.V(INFOLOG_LEVEL_CONNS) {
		// 	glog.Infoln("Subscribe connection to game stream")
		// }

		// if glog.V(INFOLOG_LEVEL_CONNS) {
		// 	glog.Infoln("Starting listening for player commands")
		// }

		return nil
	}

	return errors.New("Game data was not received")
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(ws *websocket.Conn, err error) {
	if err == nil {
		err = errors.New("Passed nil errer to reporting")
	}

	// Log error message
	if ws != nil {
		glog.Errorln("IP:", ws.Request().RemoteAddr, ", Error:", err)
	} else {
		glog.Errorln("Error:", err)
	}
}
