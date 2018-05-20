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
	defaultGroupsLimit = 100
	defaultConnsLimit  = 1000
)

var (
	address     string
	groupsLimit int
	connsLimit  int
	seed        int64
	flagJSONLog bool
	logLevel    string
)

func init() {
	flag.StringVar(&address, "address", defaultAddress, "address to serve")
	flag.IntVar(&groupsLimit, "groups-limit", defaultGroupsLimit, "groups limit")
	flag.IntVar(&connsLimit, "conns-limit", defaultConnsLimit, "web-socket connections limit")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "random seed")
	flag.BoolVar(&flagJSONLog, "log-json", false, "use json format for logger")
	flag.StringVar(&logLevel, "log-level", "info", "set log level: panic, fatal, error, warning (warn), info or debug")
	flag.Parse()
}

func logger() logrus.FieldLogger {
	logger := logrus.New()
	if flagJSONLog {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	return logger
}

func main() {
	logger := logger()

	logger.WithFields(logrus.Fields{
		"conns_limit":  connsLimit,
		"groups_limit": groupsLimit,
		"seed":         seed,
		"log_level":    logLevel,
	}).Info("preparing to start server")

	rand.Seed(seed)

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
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

	logger.WithField("address", address).Info("starting server")

	if err := http.ListenAndServe(address, n); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
