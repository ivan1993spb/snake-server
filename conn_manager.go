package main

import (
	"errors"
	"io"

	"bitbucket.org/pushkin_ivan/clever-snake/objects"
	"bitbucket.org/pushkin_ivan/clever-snake/playground"
	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
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
func (m *ConnManager) Handle(ws *websocket.Conn,
	data pwshandler.Environment) error {
	if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
		glog.Infoln("Websocket handler started for new connection")
	}
	defer func() {
		if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
			glog.Infoln("Websocket handler was finished")
		}
	}()

	if gameData, ok := data.(*GameData); ok {
		if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
			glog.Infoln("Subscribe connection to game stream")
		}
		// Starting game stream
		m.streamer.Subscribe(gameData.Playground, ws)

		// Defer unsubscribing
		defer func() {
			if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
				glog.Infoln("Unsubscribe connection from stream")
			}
			m.streamer.Unsubscribe(gameData.Playground, ws)
		}()

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                  BEGIN INIT PLAYER                  *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
			glog.Infoln("Start player init")
		}

		snake, err := objects.CreateSnake(
			gameData.Playground,
			gameData.Context,
		)
		if err != nil {
			return err
		}

		if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
			glog.Infoln("Snake was created")
		}

		/* * * * * * * * * * * * * * * * * * * * * * * * * * * *
		 *                   END INIT PLAYER                   *
		 * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
			glog.Infoln("Starting listening for player commands")
		}

		var (
			input = make(chan string) // Input commands
			errch = make(chan error)  // Errors
		)

		go func() {
			for {
				var data string
				var err = websocket.Message.Receive(ws, &data)
				if err != nil {
					errch <- err
					return
				}
				input <- string(data)
			}
		}()

		for {
			select {
			case <-gameData.Context.Done():
				if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
					glog.Infoln(
						"Parent context was canceled:",
						"finishing handler",
					)
				}
				return nil
			case err := <-errch:
				if err != io.EOF {
					if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
						glog.Infoln("Error with connection:", err)
					}
					return err
				} else if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
					glog.Infoln("Connection was closed")
				}
				return nil
			case cmd := <-input:
				if glog.V(INFOLOG_LEVEL_ABOUT_CONNS) {
					glog.Infoln("Accepted player command")
				}
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
