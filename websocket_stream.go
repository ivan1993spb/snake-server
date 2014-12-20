package main

import (
	"io"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func StartStream(cxt context.Context, chByte <-chan []byte,
) <-chan *websocket.Conn {
	if glog.V(INFOLOG_LEVEL_POOLS) {
		glog.Infoln("Starting stream")
	}

	chWs := make(chan *websocket.Conn)

	go func() {
		var webSocks = make([]*websocket.Conn, 0)

	loop:
		select {
		case <-cxt.Done():
			if glog.V(INFOLOG_LEVEL_POOLS) {
				glog.Infoln("Stopping stream")
			}
		case ws := <-chWs:
			for i := range webSocks {
				if webSocks[i] == ws {
					goto loop
				}
			}
			webSocks = append(webSocks, ws)
			goto loop
		case data := <-chByte:
			for i := 0; i < len(webSocks); {
				if _, err := webSocks[i].Write(data); err != nil {
					if err != io.EOF && glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Warningln("Connection error:", err)
					}

					webSocks = append(webSocks[:i], webSocks[i+1:]...)
				} else {
					i++
				}
			}
			goto loop
		}

		close(chWs)
	}()

	return chWs
}
