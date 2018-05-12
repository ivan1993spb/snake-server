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
	"github.com/ivan1993spb/snake-server/middlewares"
)

const (
	defaultAddress     = ":8080"
	defaultGroupsLimit = 10
)

var (
	address     string
	groupsLimit int
	seed        int64
)

func init() {
	flag.StringVar(&address, "address", defaultAddress, "address to serve")
	flag.IntVar(&groupsLimit, "groups-limit", defaultGroupsLimit, "groups limit")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "random seed")
	flag.Parse()
}

func main() {
	logger := logrus.New()
	logger.Info("preparing to start server")

	logger.Infoln("address:", address)
	logger.Infoln("group limit:", groupsLimit)
	logger.Infoln("seed:", seed)

	rand.Seed(seed)

	groupManager, err := connections.NewConnectionGroupManager(groupsLimit)
	if err != nil {
		logger.Fatalln("cannot create connections group manager:", err)
	}

	r := mux.NewRouter()
	r.Path(handlers.URLRouteCreateGame).Methods(handlers.MethodCreateGame).Handler(handlers.NewCreateGameHandler(logger, groupManager))
	r.Path(handlers.URLRouteGetGameByID).Methods(handlers.MethodGetGame).Handler(handlers.NewGetGameHandler(logger, groupManager))
	r.Path(handlers.URLRouteDeleteGameByID).Methods(handlers.MethodDeleteGame).Handler(handlers.NewDeleteGameHandler(logger, groupManager))
	r.Path(handlers.URLRouteGetGames).Methods(handlers.MethodGetGames).Handler(handlers.NewGetGamesHandler(logger, groupManager))
	r.Path(handlers.URLRouteGameWebSocketByID).Methods(handlers.MethodGame).Handler(handlers.NewGameWebSocketHandler(logger, groupManager))

	n := negroni.New(middlewares.NewRecovery(logger), middlewares.NewLogger(logger))
	n.UseHandler(r)

	logger.Info("starting server")

	if err := http.ListenAndServe(address, n); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
