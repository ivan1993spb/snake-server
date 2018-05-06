package connections

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
)

type ConnectionWorker struct {
	conn   *websocket.Conn
	logger *logrus.Logger
}

func NewConnectionWorker(conn *websocket.Conn, logger *logrus.Logger) *ConnectionWorker {
	return &ConnectionWorker{
		conn:   conn,
		logger: logger,
	}
}

func (c *ConnectionWorker) Run(game *game.Game) error {
	return c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"), time.Now().Add(time.Second))
}
