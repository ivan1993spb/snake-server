package game

import (
	// "bitbucket.org/pushkin_ivan/clever-snake/game/playground"
	"encoding/json"
	"golang.org/x/net/context"
)

//
type StartPlayerFunc func(<-chan interface{}) (<-chan interface{}, error)

type errStartingGame struct {
	err error
}

func (e *errStartingGame) Error() string {
	return "Starting game error: " + e.err.Error()
}

type Object json.Marshaler

type ObjectSet map[uint16]Object

func StartGame(cxt context.Context, pgW, pgH uint8,
) (<-chan interface{}, StartPlayerFunc, error) {
	if err := cxt.Err(); err != nil {
		return nil, nil, &errStartingGame{err}
	}

	// pg, err := playground.NewPlayground(pgW, pgH)
	// if err != nil {
	// 	return nil, nil, &errStartingGame{err}
	// }

	// objects := make(map[uint16]Object)

	output := make(chan interface{})
	go func() {
		<-cxt.Done()
		close(output)
	}()

	return output,
		func(input <-chan interface{}) (<-chan interface{}, error) {
			output := make(chan interface{})
			go func() {
				for range input {
				}
				close(output)
			}()
			return output, nil
		}, nil
}
