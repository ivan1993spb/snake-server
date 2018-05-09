package connections

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/objects/snake"
)

type ConnectionWorker struct {
	conn      *websocket.Conn
	logger    *logrus.Logger
	chStop    chan struct{}
	chStopErr chan error
	chOutput  chan OutputMessage
	chInput   chan InputMessageType

	flagStarted bool
	flagStopped bool
}

func NewConnectionWorker(conn *websocket.Conn, logger *logrus.Logger) *ConnectionWorker {
	return &ConnectionWorker{
		conn:      conn,
		logger:    logger,
		chStopErr: make(chan error, 0),
	}
}

func (cw *ConnectionWorker) Start(game *game.Game) error {
	//return cw.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"), time.Now().Add(time.Second))

	if cw.flagStarted {
		// Return error
		return nil
	}

	cw.startRead()
	cw.startWrite()

	// Start connection read
	// Start connection write
	// Listen game events
	// Parse input messages (?) and send game commands

	s, _ := snake.CreateSnake(game.World())
	go game.RunObserver(s, 16)
	//game.RunObserver(cw, 16)

	cw.flagStarted = true

	return <-cw.chStopErr
}

func (cw *ConnectionWorker) startRead() {
	go func() {
		for {
			//messageType, r, err := cw.conn.NextReader()
			_, _, err := cw.conn.NextReader()
			if err != nil {
				cw.chStopErr <- err
				return
			}
		}
	}()
}

func (cw *ConnectionWorker) startWrite() {
	go func() {
		for {
			select {
			case message := <-cw.chOutput:
				cw.logger.Info(message)
				//err := cw.writeMessage(message)
				// TODO: Handler error.
			case <-cw.chStop:
				return
			}

		}
	}()
}

func (cw *ConnectionWorker) writeMessage(message *OutputMessage) error {
	w, err := cw.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return nil
	}

	// TODO: Write message.

	// TODO: Return error.

	w.Close()
	return nil
}

func (cw *ConnectionWorker) Run(ch <-chan game.Event) {
	for {
		select {
		case event := <-ch:
			cw.logger.Info(event)
		case <-cw.chStop:
			return
		}
	}
}
