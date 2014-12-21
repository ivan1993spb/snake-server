package game

import (
	"bitbucket.org/pushkin_ivan/clever-snake/game/playground"
	"golang.org/x/net/context"
)

type StartPlayerFunc func(<-chan []byte) error

type errStartingGame struct {
	err error
}

func (e *errStartingGame) Error() string {
	return "Starting game error: " + e.err.Error()
}

func StartGame(cxt context.Context, pgW, pgH uint8,
) (<-chan []byte, StartPlayerFunc, error) {
	if err := cxt.Err(); err != nil {
		return nil, nil, &errStartingGame{err}
	}

	pg, err := playground.NewPlayground(pgW, pgH)
	if err != nil {
		return nil, nil, &errStartingGame{err}
	}

	objects := make(map[uint16]Object)

	go func() {
	}()

	return make(chan []byte), func(ch <-chan []byte) error {
		<-ch
		return nil
	}, nil
}
