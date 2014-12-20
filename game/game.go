package game

import (
	"golang.org/x/net/context"
)

type StartPlayerFunc func(<-chan []byte) error

// type Game struct {
// 	chError    chan error
// 	chCreating chan interface{}
// 	chUpdating chan interface{}
// 	chDeleting chan interface{}
// }

// func (game *Game) Start() {

// }

// func (game *Game) Play(ws *websocket.Conn) {

// }

func StartGame(cxt context.Context, pgW, pgH uint8,
) (<-chan []byte, StartPlayerFunc, error) {
	return make(chan []byte), func(ch <-chan []byte) error {
		<-ch
		return nil
	}, nil
}
