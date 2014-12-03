package main

import (
	"errors"
	"net/http"
	// "time"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"

	"bitbucket.org/pushkin_ivan/clever-snake/playground"
)

type GameData struct {
	Context    context.Context
	Playground *playground.Playground
}

type ConnManager struct {
	streamer *Streamer
}

func NewConnManager(s *Streamer) (pwshandler.ConnManager, error) {
	if s == nil {
		return nil, errors.New("Passed nil streamer")
	}

	return &ConnManager{s}, nil
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) Handle(conn *websocket.Conn,
	env pwshandler.Environment) error {

	if /*gameData*/ _, ok := env.(*GameData); ok {

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                  BEGIN INIT PLAYER                  *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                   END INIT PLAYER                   *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		return nil
	}

	return errors.New("Game data was not received")
}

// Implementing pwshandler.ConnManager interface
func (m *ConnManager) HandleError(_ http.ResponseWriter,
	r *http.Request, err error) {
	if err == nil {
		err = errors.New("Passed nil errer to reporting")
	}

	// Log error message
	if r != nil {
		glog.Infoln("IP:", r.RemoteAddr, ", Error:", err)
	} else {
		glog.Infoln("Error:", err)
	}
}
