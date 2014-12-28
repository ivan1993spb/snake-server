package game

import (
	"encoding/json"
	"time"

	// "bitbucket.org/pushkin_ivan/clever-snake/game/playground"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

//
type StartPlayerFunc func(cxt context.Context, input <-chan *Command,
) (<-chan interface{}, error)

type errStartingGame struct {
	err error
}

func (e *errStartingGame) Error() string {
	return "starting game error: " + e.err.Error()
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

		defer close(output)
		defer glog.Infoln("finishing game")

		// all running objects work like this code:
		for {
			select {
			case <-cxt.Done():
				return
			case <-time.Tick(time.Second * 3):
				output <- "test"
			}
		}
	}()

	return output,
		func(pcxt context.Context, input <-chan *Command) (
			<-chan interface{}, error) {
			if pcxt.Err() != nil {
				return nil, nil
			}
			output := make(chan interface{})
			go func() {

				defer close(output)
				defer glog.Infoln("finishing player")

				select {
				case <-pcxt.Done():
				case <-time.After(time.Second):
				}

				for {
					select {
					case <-pcxt.Done():
						return
					case cmd := <-input:
						if cmd == nil {
							return
						}
						output <- "received cmd =)"
					}
				}
			}()
			return output, nil
		}, nil
}

type Command struct {
	Command string          `json:"command"`
	Params  json.RawMessage `json:"params"`
}

type Notice struct {
}
