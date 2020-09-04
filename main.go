package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/client"
	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/handlers"
	"github.com/ivan1993spb/snake-server/middlewares"
)

const ServerName = "Snake-Server"

const (
	defaultAddress     = ":8080"
	defaultGroupsLimit = 100
	defaultConnsLimit  = 1000
)

var (
	Version = "dev"
	Build   = "dev"
	Author  = "Ivan Pushkin"
	License = "MIT"
)

var (
	address string

	flagEnableTLS bool
	certFile      string
	keyFile       string

	groupsLimit int
	connsLimit  int
	seed        int64

	flagJSONLog bool
	logLevel    string

	enableBroadcast bool

	enableWeb bool

	forbidCORS bool
)

const logName = "api"

func usage() {
	fmt.Fprint(os.Stderr, "Welcome to snake-server!\n\n")
	fmt.Fprintf(os.Stderr, "Server version %s, build %s\n\n", Version, Build)
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&address, "address", defaultAddress, "address to serve")
	flag.BoolVar(&flagEnableTLS, "tls-enable", false, "enable TLS")
	flag.StringVar(&certFile, "tls-cert", "", "path to certificate file")
	flag.StringVar(&keyFile, "tls-key", "", "path to key file")
	flag.IntVar(&groupsLimit, "groups-limit", defaultGroupsLimit, "game groups limit")
	flag.IntVar(&connsLimit, "conns-limit", defaultConnsLimit, "web-socket connections limit")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "random seed")
	flag.BoolVar(&flagJSONLog, "log-json", false, "use json format for logger")
	flag.StringVar(&logLevel, "log-level", "info", "set log level: panic, fatal, error, warning (warn), info or debug")
	flag.BoolVar(&enableBroadcast, "enable-broadcast", false, "enable broadcasting API method")
	flag.BoolVar(&enableWeb, "enable-web", false, "enable web client")
	flag.BoolVar(&forbidCORS, "forbid-cors", false, "forbid cross-origin resource sharing")
	flag.Usage = usage
	flag.Parse()
}

func logger() *logrus.Logger {
	logger := logrus.New()
	if flagJSONLog {
		logger.Formatter = &logrus.JSONFormatter{}
	} else if runtime.GOOS == "windows" {
		// Log Output on Windows shows Bash format
		// See: https://gitlab.com/gitlab-org/gitlab-runner/issues/6
		// See: https://github.com/sirupsen/logrus/issues/172
		logger.Formatter = &logrus.TextFormatter{
			DisableColors: true,
		}
	}
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	return logger
}

func serve(h http.Handler) error {
	if flagEnableTLS {
		return http.ListenAndServeTLS(address, certFile, keyFile, h)
	}
	return http.ListenAndServe(address, h)
}

func main() {
	logger := logger()

	logger.WithFields(logrus.Fields{
		"author":  Author,
		"license": License,
		"version": Version,
		"build":   Build,
	}).Info("welcome to snake-server!")

	logger.WithFields(logrus.Fields{
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
	}).Info("golang info")

	logger.WithFields(logrus.Fields{
		"conns_limit":  connsLimit,
		"groups_limit": groupsLimit,
		"seed":         seed,
		"log_level":    logLevel,
		"broadcast":    enableBroadcast,
		"web":          enableWeb,
		"cors":         !forbidCORS,
	}).Info("preparing to start server")

	if enableBroadcast {
		logger.Warning("broadcasting API method is enabled!")
	}

	rand.Seed(seed)

	groupManager, err := connections.NewConnectionGroupManager(logger, groupsLimit, connsLimit)
	if err != nil {
		logger.Fatalln("cannot create connections group manager:", err)
	}
	if err := prometheus.Register(groupManager); err != nil {
		logger.Fatalln("cannot register connection group manager as a metric collector:", err)
	}

	rootRouter := mux.NewRouter().StrictSlash(true)
	rootRouter.Path("/metrics").Handler(promhttp.Handler())
	rootRouter.Path(handlers.URLRouteOpenAPI).Handler(handlers.NewOpenAPIHandler())
	if enableWeb {
		rootRouter.Path(client.URLRouteServerEndpoint).Handler(http.RedirectHandler(client.URLRouteClient, http.StatusFound))
		rootRouter.PathPrefix(client.URLRouteClient).Handler(negroni.New(gzip.Gzip(gzip.DefaultCompression), negroni.Wrap(client.NewHandler())))
	} else {
		rootRouter.Path(handlers.URLRouteWelcome).Methods(handlers.MethodWelcome).Handler(handlers.NewWelcomeHandler(logger))
	}
	rootRouter.NotFoundHandler = handlers.NewNotFoundHandler(logger)

	// Web-Socket routes
	wsRouter := rootRouter.PathPrefix("/ws").Subrouter()
	wsRouter.Path(handlers.URLRouteGameWebSocketByID).Methods(handlers.MethodGame).Handler(handlers.NewGameWebSocketHandler(logger, groupManager))

	// API routes
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiRouter.Path(handlers.URLRouteGetInfo).Methods(handlers.MethodGetInfo).Handler(handlers.NewGetInfoHandler(logger, Author, License, Version, Build))
	apiRouter.Path(handlers.URLRouteGetCapacity).Methods(handlers.MethodGetCapacity).Handler(handlers.NewGetCapacityHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRouteCreateGame).Methods(handlers.MethodCreateGame).Handler(handlers.NewCreateGameHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRouteGetGameByID).Methods(handlers.MethodGetGame).Handler(handlers.NewGetGameHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRouteDeleteGameByID).Methods(handlers.MethodDeleteGame).Handler(handlers.NewDeleteGameHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRouteGetGames).Methods(handlers.MethodGetGames).Handler(handlers.NewGetGamesHandler(logger, groupManager))
	if enableBroadcast {
		apiRouter.Path(handlers.URLRouteBroadcast).Methods(handlers.MethodBroadcast).Handler(handlers.NewBroadcastHandler(logger, groupManager))
	}
	apiRouter.Path(handlers.URLRouteGetObjects).Methods(handlers.MethodGetObjects).Handler(handlers.NewGetObjectsHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRoutePing).Methods(handlers.MethodPing).Handler(handlers.NewPingHandler(logger))

	n := negroni.New(
		middlewares.NewRecovery(logger),
		middlewares.NewServerInfo(ServerName, Version, Build),
		middlewares.NewLogger(logger, logName),
	)

	if !forbidCORS {
		n.Use(middlewares.NewCORS())
	}

	n.UseHandler(rootRouter)

	logger.WithFields(logrus.Fields{
		"address": address,
		"tls":     flagEnableTLS,
	}).Info("starting server")

	if err := serve(n); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
