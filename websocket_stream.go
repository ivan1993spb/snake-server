package main

import (
	"errors"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

// StartStreamConnFunc starts stream for passed websocket connection
type StartStreamConnFunc func(*websocket.Conn) error

// StopStreamConnFunc stops stream for passed websocket connection
type StopStreamConnFunc func(*websocket.Conn) error

//  StartGameStream starts common pool game stream
func StartGameStream(chByte <-chan []byte,
) (StartStreamConnFunc, StopStreamConnFunc) {

	conns := make([]*websocket.Conn, 0)

	go func() {
		for data := range chByte {
			if len(data) == 0 || len(conns) == 0 {
				continue
			}

			// Send data for each websocket connection in conns
			for i := 0; i < len(conns); {
				if _, err := conns[i].Write(data); err != nil {
					// Remove connection on error
					glog.Warningln(
						"Cannot send common game data:",
						err,
					)

					if glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Infoln(
							"Removing wsconn from game stream",
						)
					}

					conns = append(conns[:i], conns[i+1:]...)
				} else {
					i++
				}
			}
		}

		if len(conns) < 0 {
			conns = conns[:0]
		}

		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("Common game stream finished")
		}
	}()

	return func(ws *websocket.Conn) error {
			// Check if passed websocket connection already exists
			for i := range conns {
				if conns[i] == ws {
					return errors.New("Cannot create connection to " +
						"common pool game stream: Passed connection" +
						" already exists")
				}
			}

			conns = append(conns, ws)

			return nil
		}, func(ws *websocket.Conn) (err error) {
			for i := range conns {
				if conns[i] == ws {
					conns = append(conns[:i], conns[i+1:]...)
					return nil
				}
			}

			return errors.New("Cannot remove connection from common" +
				" pool game stream: Passed connection was not found")
		}
}
