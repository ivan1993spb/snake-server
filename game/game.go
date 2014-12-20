package game

import (
	"golang.org/x/net/context"
)

type PlayFunc func(<-chan []byte) error

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
) (<-chan []byte, PlayFunc, error) {
	return make(chan []byte), func(<-chan []byte) error {
		return nil
	}, nil
}
