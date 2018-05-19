package player

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/objects/snake"
)

type Player struct {
	game   *game.Game
	logger logrus.FieldLogger
}

func NewPlayer(logger logrus.FieldLogger, game *game.Game) *Player {
	return &Player{
		game:   game,
		logger: logger,
	}
}

func (p *Player) Start(stop <-chan struct{}) <-chan interface{} {
	s, _ := snake.NewSnake(p.game.World())
	s.Run(stop)

	chout := make(chan interface{})

	go func() {
		<-stop
		s.Die()
		close(chout)
	}()

	return chout

}
