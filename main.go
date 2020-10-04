package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/ivan1993spb/snake-server/config"
	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/server/http"
)

const ServerName = "Snake-Server"

var (
	Version = "dev"
	Build   = "dev"
	Author  = "Ivan Pushkin"
	License = "MIT"
)

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

func main() {
	cfg, err := configurate()
	logger := logger(cfg.Server.Log)
	if err != nil {
		logger.Fatalln("cannot load config:", err)
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
		"broadcast":    cfg.Server.EnableBroadcast,
		"web":          cfg.Server.EnableWeb,
		"cors":         !cfg.Server.ForbidCORS,
	}).Info("preparing to start server")

	if cfg.Server.EnableBroadcast {
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

	server := http.NewServer(cfg, groupManager, logger, Author, License, Version, Build)

	logger.WithFields(logrus.Fields{
		"address": cfg.Server.Address,
		"tls":     cfg.Server.TLS.Enable,
	}).Info("starting server")

	if err := server.Run(); err != nil {
		logger.Fatalf("server error: %s", err)
	}
}
