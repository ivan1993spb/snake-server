package player

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

const countdown = 5
const chanMessageBuffer = 16

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

func (p *Player) Start(stop <-chan struct{}) <-chan Message {
	chout := make(chan Message, chanMessageBuffer)
	localStopper := make(chan struct{})

	go func() {
		<-stop
		close(localStopper)
		close(chout)
	}()

	go func() {
		chout <- Message{
			Type:    MessageTypeNotice,
			Payload: MessageNotice("welcome to snake server!"),
		}

		chout <- Message{
			Type: MessageTypeSize,
			Payload: MessageSize{
				Width:  p.world.Width(),
				Height: p.world.Height(),
			},
		}

		for {
			chout <- Message{
				Type:    MessageTypeCountdown,
				Payload: MessageCountdown(countdown),
			}
			timer := time.NewTimer(time.Second * countdown)
			select {
			case <-timer.C:
				timer.Stop()
			case <-localStopper:
				timer.Stop()
				return
			}

			s, err := snake.NewSnake(p.world)
			if err != nil {
				chout <- Message{
					Type:    MessageTypeError,
					Payload: MessageError("cannot create snake"),
				}
				p.logger.Errorln("cannot create snake to player:", err)
				continue
			}
			snakeStop := s.Run(localStopper)

			chout <- Message{
				Type:    MessageTypeSnake,
				Payload: MessageSnake(s.GetID()),
			}

			select {
			case <-snakeStop:
			case <-localStopper:
				s.Die()
				return
			}
		}
	}()

	return chout
}

func (p *Player) SnakeCommand() {
}
