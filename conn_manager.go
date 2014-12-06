package main

import (
	"errors"
	"net/http"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"

	"bitbucket.org/pushkin_ivan/clever-snake/objects"
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
	if glog.V(3) {
		glog.Infoln("Websocket handler started for new connection")
	}
	defer func() {
		if glog.V(3) {
			glog.Infoln("Websocket handler stops")
		}
	}()

	if gameData, ok := env.(*GameData); ok {
		if glog.V(4) {
			glog.Infoln("Handler receive game data")
		}

		if glog.V(3) {
			glog.Infoln("Subscribe connection to stream")
		}
		// Starting game stream
		m.streamer.Subscribe(gameData.Playground, conn)
		// Defer unsubscribing
		defer func() {
			if glog.V(3) {
				glog.Infoln("Handler are finishing")
			}
			m.streamer.Unsubscribe(gameData.Playground, conn)
		}()

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                  BEGIN INIT PLAYER                  *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		if glog.V(4) {
			glog.Infoln("Start player init")
		}

		snake, err := objects.CreateSnake(
			gameData.Playground,
			gameData.Context,
		)
		if err != nil {
			return err
		}

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                   END INIT PLAYER                   *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		if glog.V(3) {
			glog.Infoln("Starting listening for player commands")
		}

		var (
			input = make(chan string) // Input commands
			errch = make(chan error)  // Errors
		)

		go func() {
			for {
				if ty, cmd, err := conn.ReadMessage(); err != nil {
					errch <- err
					return
				} else if ty == websocket.TextMessage {
					input <- string(cmd)
				}
			}
		}()

		for {
			select {
			case <-gameData.Context.Done():
				return nil
			case err := <-errch:
				return err
			case cmd := <-input:
				if err := snake.Command(cmd); err != nil {
					return err
				}
			}
		}

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
