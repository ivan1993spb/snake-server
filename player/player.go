package player

import (
	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/objects/snake"
)

type Player struct {
	game *game.Game
}

func NewPlayer(game *game.Game) *Player {
	return &Player{
		game: game,
	}
}

func (p *Player) Start() {
	s, _ := snake.NewSnake(p.game.World())
	// TODO: Pass stop channel?
	s.Run(p.game.Events(make(chan struct{})))
}
