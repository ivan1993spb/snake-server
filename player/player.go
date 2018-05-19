package player

import (
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

type Player struct {
	world  *world.World
	logger logrus.FieldLogger
}

func NewPlayer(logger logrus.FieldLogger, world *world.World) *Player {
	return &Player{
		logger: logger,
		world:  world,
	}
}

func (p *Player) Start(stop <-chan struct{}) <-chan interface{} {
	s, _ := snake.NewSnake(p.world)
	s.Run(stop)

	chout := make(chan interface{})

	go func() {
		<-stop
		s.Die()
		close(chout)
	}()

	return chout

}
