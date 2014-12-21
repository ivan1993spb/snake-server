package main

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// StartStreamFunc starts stream for passed websocket connection
type StartStreamFunc func(*websocket.Conn) error

// StartStream starts game stream. Func receives websockets from
// returned StartStreamFunc and saves it. When bytes received from
// passed channel StartStream sends bytes to all saved websockets
func StartStream(cxt context.Context, chByte <-chan []byte,
) (StartStreamFunc, error) {
	if err := cxt.Err(); err != nil {
		return nil, fmt.Errorf("Cannot start stream: %s", err)
	}

	// Channel for creation new websockets
	chWs := make(chan *websocket.Conn)

	go func() {
		defer close(chWs)
		var webSocks = make([]*websocket.Conn, 0)

	loop:
		select {
		case <-cxt.Done():
			if glog.V(INFOLOG_LEVEL_POOLS) {
				glog.Infoln("Common game stream stops")
			}
			return
		case ws := <-chWs:
			for i := range webSocks {
				if webSocks[i] == ws {
					goto loop
				}
			}
			webSocks = append(webSocks, ws)
		case data := <-chByte:
			for i := 0; i < len(webSocks); {
				if _, err := webSocks[i].Write(data); err != nil {
					if glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Errorln(
							"Cannot send common game data:",
							err,
						)
					}
					webSocks = append(webSocks[:i], webSocks[i+1:]...)
				} else {
					i++
				}
			}
		}

		goto loop
	}()

	return func(ws *websocket.Conn) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf(
					"Cannot create connection to game stream: %s", r,
				)
			}
		}()

		chWs <- ws
		return
	}, nil
}
