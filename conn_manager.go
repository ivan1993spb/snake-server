package main

import (
	"errors"
	"net/http"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"bitbucket.org/pushkin_ivan/simple-2d-playground"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type GameData struct {
	Context    context.Context
	Playground *playground.Playground
}

type ConnManager struct{}

func NewConnManager() pwshandler.ConnManager {
	return &ConnManager{}
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) Handle(conn *websocket.Conn,
	env pwshandler.Environment) error {

	if gameData, ok := env.(*GameData); ok {

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *          GAME LOGIC IS HERE. INIT PLAYER            *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		return nil
	}

	return errors.New("Game data was not received")
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(_ http.ResponseWriter,
	_ *http.Request, err error) {
	// Write error message to log
	glog.Exitln(err)
}
