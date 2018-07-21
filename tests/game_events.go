package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/player"
)

const (
	defaultWidth  = 150
	defaultHeight = 150

	defaultPlayerCount = 100

	defaultTestDuration = time.Minute
)

const chanGameEventsBuffer = 8192

var (
	width  uint
	height uint

	playerCount int

	testDuration time.Duration
)

func init() {
	flag.UintVar(&width, "width", defaultWidth, "map width")
	flag.UintVar(&height, "height", defaultHeight, "map height")
	flag.IntVar(&playerCount, "players", defaultPlayerCount, "players count")
	flag.DurationVar(&testDuration, "duration", defaultTestDuration, "test duration")
	flag.Parse()
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	logger.WithFields(logrus.Fields{
		"map_width":     width,
		"map_height":    height,
		"player_count":  playerCount,
		"test_duration": testDuration,
	}).Info("start test")

	stop := make(chan struct{})

	g, err := game.NewGame(logger, uint8(width), uint8(height))
	if err != nil {
		logger.Fatal("cannot create game")
	}

	g.Start(stop)

	world := g.World()
	chins := make([]chan string, playerCount)

	for i := 0; i < playerCount; i++ {
		chins[i] = make(chan string)
		p := player.NewPlayer(logger, world)
		p.Start(stop, chins[i])
	}

	chEvents := g.ListenEvents(stop, chanGameEventsBuffer)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(testDuration)
	defer timer.Stop()

	var (
		counter int
		average float32
	)

	for {
		select {
		case <-chEvents:
			counter++
		case <-ticker.C:
			average = (average + float32(counter)) / 2
			logger.Infof("count %d events/second, average %f events/second", counter, average)
			counter = 0
		case <-timer.C:
			logger.Info("stop")
			return
		}
	}
}
