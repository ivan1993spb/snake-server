package player

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/objects/snake"
	"github.com/ivan1993spb/snake-server/world"
)

const countdown = 5

const chanMessageBuffer = 128

const chanErrorBuffer = 32

type Player struct {
	world  world.Interface
	logger logrus.FieldLogger
}

func NewPlayer(logger logrus.FieldLogger, world world.Interface) *Player {
	return &Player{
		logger: logger,
		world:  world,
	}
}

func (p *Player) Start(stop <-chan struct{}, chin <-chan string) <-chan Message {
	chout := make(chan Message, chanMessageBuffer)
	localStopper := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		<-stop
		close(localStopper)

		wg.Wait()
		close(chout)
	}()

	go func() {
		defer wg.Done()

		chout <- NewMessageNotice("welcome to snake-server!")
		chout <- NewMessageSize(p.world.Width(), p.world.Height())
		chout <- NewMessageObjects(p.world.GetObjects())

		for {
			chout <- NewMessageCountdown(countdown)

			timer := time.NewTimer(time.Second * countdown)
			select {
			case <-timer.C:
				timer.Stop()
			case <-localStopper:
				timer.Stop()
				return
			}

			chout <- NewMessageNotice("start")

			p.emptyInputChan(localStopper, chin)

			s, err := snake.NewSnake(p.world)
			if err != nil {
				chout <- NewMessageError("cannot create snake")
				p.logger.Errorln("cannot create snake to player:", err)
				continue
			}
			snakeStop := s.Run(localStopper, p.logger)

			chout <- NewMessageSnake(s.GetID())

			wg.Add(1)
			go func() {
				defer wg.Done()
				errch := p.processSnakeCommands(snakeStop, chin, s)

				for {
					select {
					case <-snakeStop:
						return
					case err, ok := <-errch:
						if !ok {
							return
						}
						chout <- NewMessageError(err.Error())
					}
				}
			}()

			select {
			case <-snakeStop:
			case <-localStopper:
				return
			}
		}
	}()

	return chout
}

func (p *Player) processSnakeCommands(stop <-chan struct{}, chin <-chan string, s *snake.Snake) <-chan error {
	errch := make(chan error, chanErrorBuffer)

	go func() {
		defer close(errch)
		for {
			select {
			case <-stop:
				return
			case command, ok := <-chin:
				if !ok {
					return
				}

				p.logger.WithField("command", command).Debug("received snake command")

				if err := s.Command(snake.Command(command)); err != nil {
					errch <- err
				}
			}
		}
	}()

	return errch
}

func (p *Player) emptyInputChan(stop <-chan struct{}, chin <-chan string) {
	for {
		if len(chin) > 0 {
			select {
			case <-chin:
			case <-stop:
				return
			default:
				return
			}
		} else {
			return
		}
	}
}
