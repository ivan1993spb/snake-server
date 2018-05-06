package connections

import (
	"github.com/gorilla/websocket"

	"github.com/ivan1993spb/snake-server/game"
)

type Connection struct {
}

func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{}
}

func (c *Connection) Run(game *game.Game) error {

}
