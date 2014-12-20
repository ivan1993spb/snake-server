package game

// import (
// 	"golang.org/x/net/context"
// 	"golang.org/x/net/websocket"
// )

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

type Game struct{}

func NewGame(...interface{}) *Game {
	return &Game{}
}

func (g *Game) Start() {

}

func (g *Game) GetStream() <-chan []byte {
	return make(chan []byte)
}
