package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"

	"github.com/evalphobia/logrus_sentry"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/client"
	"github.com/ivan1993spb/snake-server/config"
	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/handlers"
	"github.com/ivan1993spb/snake-server/middlewares"
)

const ServerName = "Snake-Server"

var (
	Version = "dev"
	Build   = "dev"
	Author  = "Ivan Pushkin"
	License = "MIT"
)

const logName = "api"

func usage() {
	fmt.Fprint(os.Stderr, "Welcome to snake-server!\n\n")
	fmt.Fprintf(os.Stderr, "Server version %s, build %s\n\n", Version, Build)
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	flag.PrintDefaults()
}

func configurate() (config.Config, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.Usage = usage
	cfg, err := config.Configurate(afero.NewOsFs(), f, os.Args[1:])
	return cfg, err
}

func logger(configLog config.Log) *logrus.Logger {
	logger := logrus.New()
	if configLog.EnableJSON {
		logger.Formatter = &logrus.JSONFormatter{}
	} else if runtime.GOOS == "windows" {
		// Log Output on Windows shows Bash format
		// See: https://gitlab.com/gitlab-org/gitlab-runner/issues/6
		// See: https://github.com/sirupsen/logrus/issues/172
		logger.Formatter = &logrus.TextFormatter{
			DisableColors: true,
		}
	}
	if level, err := logrus.ParseLevel(configLog.Level); err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
	return logger
}

func serve(h http.Handler, address string, configTLS config.TLS) error {
	if configTLS.Enable {
		return http.ListenAndServeTLS(address, configTLS.Cert, configTLS.Key, h)
	}
	return http.ListenAndServe(address, h)
}

func main() {
	cfg, err := configurate()
	logger := logger(cfg.Server.Log)
	if err != nil {
		logger.Fatalln("cannot load config:", err)
	}

	if cfg.Server.Sentry.Enable {
		hook, err := logrus_sentry.NewAsyncSentryHook(cfg.Server.Sentry.DSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err == nil {
			logger.Hooks.Add(hook)
		}
	}

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
		"conns_limit":  cfg.Server.Limits.Conns,
		"groups_limit": cfg.Server.Limits.Groups,
		"seed":         cfg.Server.Seed,
		"log_level":    cfg.Server.Log.Level,
		"broadcast":    cfg.Server.Flags.EnableBroadcast,
		"web":          cfg.Server.Flags.EnableWeb,
		"cors":         !cfg.Server.Flags.ForbidCORS,
	}).Info("preparing to start server")

	if cfg.Server.Flags.EnableBroadcast {
		logger.Warning("broadcasting API method is enabled!")
	}

	rand.Seed(cfg.Server.Seed)

	groupManager, err := connections.NewConnectionGroupManager(logger, cfg.Server.Limits.Groups, cfg.Server.Limits.Conns)
	if err != nil {
		logger.Fatalln("cannot create connections group manager:", err)
	}
	if err := prometheus.Register(groupManager); err != nil {
		logger.Fatalln("cannot register connection group manager as a metric collector:", err)
	}

	rootRouter := mux.NewRouter().StrictSlash(true)
	rootRouter.Path("/metrics").Handler(promhttp.Handler())
	rootRouter.Path(handlers.URLRouteOpenAPI).Handler(handlers.NewOpenAPIHandler())
	if cfg.Server.Flags.EnableWeb {
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
	if cfg.Server.Flags.EnableBroadcast {
		apiRouter.Path(handlers.URLRouteBroadcast).Methods(handlers.MethodBroadcast).Handler(handlers.NewBroadcastHandler(logger, groupManager))
	}
	apiRouter.Path(handlers.URLRouteGetObjects).Methods(handlers.MethodGetObjects).Handler(handlers.NewGetObjectsHandler(logger, groupManager))
	apiRouter.Path(handlers.URLRoutePing).Methods(handlers.MethodPing).Handler(handlers.NewPingHandler(logger))

	n := negroni.New(
		middlewares.NewRecovery(logger),
		middlewares.NewServerInfo(ServerName, Version, Build),
		middlewares.NewLogger(logger, logName),
	)

	if !cfg.Server.Flags.ForbidCORS {
		n.Use(middlewares.NewCORS())
	}

	n.UseHandler(rootRouter)

	logger.WithFields(logrus.Fields{
		"address": cfg.Server.Address,
		"tls":     cfg.Server.TLS.Enable,
	}).Info("starting server")

	if err := serve(n, cfg.Server.Address, cfg.Server.TLS); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
