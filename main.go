package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/handlers"
)

const (
	defaultAddress = ":8080"

	defaultGroupsLimit      = 10
	defaultConnectionsLimit = 4 // ?

	defaultWidth  = 100 // ?
	defaultHeight = 100 // ?
)

var (
	address string

	groupsLimit int

	// TODO: ?
	//emptyRoomExpire time.Duration // Create if users will be able to add rooms

	// Room properties
	connectionsLimit uint

	pgW, pgH uint

	seed int64
)

func init() {
	flag.StringVar(&address, "address", defaultAddress, "address to serve")
	flag.IntVar(&groupsLimit, "max-groups", defaultGroupsLimit, "max groups count on server")

	flag.UintVar(&connectionsLimit, "default-conn-limit", defaultConnectionsLimit, "default connection count for group")
	flag.UintVar(&pgW, "width", defaultWidth, "default map width")
	flag.UintVar(&pgH, "height", defaultHeight, "default map height")

	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "random seed")
	flag.Parse()
}

func main() {
	logger := logrus.New()
	logger.Info("preparing to start server")

	logger.Infoln("address:", address)

	logger.Infoln("group limit:", groupsLimit)

	rand.Seed(seed)
	logger.Infoln("seed:", seed)

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit)
	if err != nil {
		logger.Fatalln("cannot create connections group manager:", err)
	}

	r := mux.NewRouter()
	r.Path(handlers.URLRouteGameWebSocket).Methods(handlers.MethodGame).Handler(handlers.NewGameWebSocketHandler(logger, groupManager))
	r.Path(handlers.URLRouteCreateGame).Methods(handlers.MethodCreateGame).Handler(handlers.NewCreateGameHandler(logger, groupManager))
	r.Path(handlers.URLRouteDeleteGameByID).Methods(handlers.MethodDeleteGame).Handler(handlers.NewDeleteGameHandler(logger, groupManager))

	// TODO: Check is it necessary to use recovery middleware.
	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(r)

	logger.Info("starting server")

	if err := http.ListenAndServe(address, n); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
