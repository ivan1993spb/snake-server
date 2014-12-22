package main

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// StartStreamConnFunc starts stream for passed websocket connection
type StartStreamConnFunc func(*websocket.Conn) error

// StopStreamConnFunc stops stream for passed websocket connection
type StopStreamConnFunc func(*websocket.Conn) error

// StartStream starts common pool game stream
func StartStream(cxt context.Context, chByte <-chan []byte,
) (StartStreamConnFunc, StopStreamConnFunc, error) {
	if err := cxt.Err(); err != nil {
		return nil, nil, fmt.Errorf("Cannot start stream: %s", err)
	}

	conns := make([]*websocket.Conn, 0)

	go func() {
		for {
			select {
			case <-cxt.Done():
				if glog.V(INFOLOG_LEVEL_POOLS) {
					glog.Infoln("Common game stream stops")
				}
				return
			case data := <-chByte:
				// Send data for each websocket connection in conns
				for i := 0; i < len(conns); {
					if _, err := conns[i].Write(data); err != nil {
						// Remove connection on error

						if glog.V(INFOLOG_LEVEL_CONNS) {
							glog.Warningln(
								"Cannot send common game data:",
								err,
							)
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
		}, nil
}
