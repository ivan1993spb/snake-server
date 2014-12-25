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

// StartGameStream starts common pool game stream
func StartGameStream(stream <-chan interface{},
) (StartStreamConnFunc, StopStreamConnFunc) {

	conns := make([]*websocket.Conn, 0)

	go func() {
		for data := range stream {
			if len(conns) == 0 || data == nil {
				continue
			}

			// Send data for each websocket connection in conns
			for i := 0; i < len(conns); {
				err := websocket.JSON.Send(conns[i], &OutputMessage{
					HEADER_GAME, data,
				})
				if err != nil {
					// Remove connection on error
					glog.Errorln(
						"cannot send common game data:", err,
					)

					if glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Infoln(
							"removing connection from game stream",
						)
					}

					conns = append(conns[:i], conns[i+1:]...)
				} else {
					i++
				}
			}
		}

		if glog.V(INFOLOG_LEVEL_POOLS) {
			glog.Infoln("common game stream finished")
		}
	}()

	return func(ws *websocket.Conn) error {
			// Check if passed websocket connection already exists
			for i := range conns {
				if conns[i] == ws {
					return errors.New("cannot create connection to " +
						"common pool game stream: passed connection" +
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

			return errors.New("cannot remove connection from common" +
				" pool game stream: passed connection was not found")
		}
}
