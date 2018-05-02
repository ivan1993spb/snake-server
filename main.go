package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/handlers"
)

const description = "Snake server"

type errStartingServer struct {
	err error
}

func (e *errStartingServer) Error() string {
	return "starting server error: " + e.err.Error()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                  BEGIN PARSING PARAMETERS                   *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	var (
		// Networking
		// TODO: fix to one "addr" param
		addr string

		// Server limits
		poolLimit uint
		//emptyRoomExpire time.Duration // Create if users will be able to add rooms

		// Room properties
		connLimit, pgW, pgH uint
	)

	flag.StringVar(&addr, "addr", "", "addr")
	flag.UintVar(&poolLimit, "pool_limit", 10, "max pool count on server")
	flag.UintVar(&connLimit, "conn_limit", 4, "max connection count on pool")
	flag.UintVar(&pgW, "pg_w", 40, "playground width")
	flag.UintVar(&pgH, "pg_h", 28, "playground height")

	flag.Usage = func() {
		fmt.Println(description)
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	flag.Parse()

	logger := logrus.New()
	logger.Info("preparing to start server")

	r := mux.NewRouter()

	r.Path(handlers.URLRouteGameByRoomID).Methods(http.MethodGet).Handler(handlers.NewGameHandler(logger))

	// TODO: Add negroni packages middlewares: recovery... etc

	// Init game pool factory
	//gamePoolFactory, err := NewGamePoolFactory(uint16(connLimit), uint8(pgW), uint8(pgH))
	//if err != nil {
	//	logger.Fatal(&errStartingServer{err})
	//}
	//logrus.Info("game pool factory was created")

	// Init game pool manager which allocates connections on pools
	//gamePoolManager, err := NewGamePoolManager(gamePoolFactory, uint16(poolLimit))
	//if err != nil {
	//	logger.Fatal(&errStartingServer{err})
	//}
	//logger.Info("game pool manager was created")

	// Init game connection manager
	//gameConnManager := new(GameConnManager)
	//logger.Info("game connection manager was created")

	logger.Info("starting server")

	if err := http.ListenAndServe(":3001", r); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
