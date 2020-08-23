package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/server/http"
)

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

func RunServer(server *http.Server) error {
	// TODO: Refactor this function.
	if flagEnableTLS {
		return server.ListenAndServeTLS(certFile, keyFile)
	}
	return server.ListenAndServe()
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

	server := &http.Server{
		Addr:         address,
		Logger:       logger,
		GroupManager: groupManager,
	}

	server.InitRoutes(enableWeb, enableBroadcast, forbidCORS, Author, License, Version, Build)

	logger.WithFields(logrus.Fields{
		"address": address,
		"tls":     flagEnableTLS,
	}).Info("starting server")

	if err := RunServer(server); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
