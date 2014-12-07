package main

import (
	"flag"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
)

const (
	INFOLOG_LEVEL_ABOUT_SERVER  = iota + 1 // Messages about server
	INFOLOG_LEVEL_ABOUT_STREAMS            // Messages about streams
	INFOLOG_LEVEL_ABOUT_POOLS              // Messages about pools
	INFOLOG_LEVEL_ABOUT_CONNS              // About connections
)

func main() {
	var host, port, sdPort, hashSalt string
	flag.StringVar(&host, "host", "",
		"host on which game server handles requests")
	flag.StringVar(&port, "port", "8081",
		"port on which game server handles requests")
	flag.StringVar(&sdPort, "shutdown_port", "8082",
		"port on which server accepts for shutdown request")
	flag.StringVar(&hashSalt, "hash_salt", "",
		"salt for request verifier")

	var delay time.Duration
	flag.DurationVar(&delay, "delay", time.Millisecond*150,
		"stream delay")

	var poolLimit, connLimit, pgW, pgH uint
	flag.UintVar(&poolLimit, "pool_limit", 10,
		"max pool number on server")
	flag.UintVar(&connLimit, "conn_limit", 4,
		"max connection number on pool")
	flag.UintVar(&pgW, "pg_w", 40, "playground width")
	flag.UintVar(&pgH, "pg_h", 28, "playground height")

	flag.Parse()

	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Preparing to start server")
	}

	// Working listener is used for game servering
	workingListener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		glog.Exitln("Cannot create working listener:", err)
	}

	// Shutdown listener is used only for shutdown command. Listening
	// only local requests
	shutdownListener, err := net.Listen("tcp", "127.0.0.1:"+sdPort)
	if err != nil {
		glog.Exitln("Cannot create shutdown listener:", err)
	}

	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Listeners was created")
	}

	// Gets root context and cancel func for all goroutines on server
	cxt, cancel := context.WithCancel(context.Background())

	// Init pool factory
	factory, err := NewPGPoolFactory(cxt, uint8(connLimit),
		uint8(pgW), uint8(pgH))
	if err != nil {
		glog.Exitln("Cannot create pool factory:", err)
	}
	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Pool factory was created")
	}

	// Init pool manager which allocates connections on pools
	poolManager, err := NewGamePoolManager(factory, uint8(poolLimit))
	if err != nil {
		glog.Exitln("Cannot create pool manager:", err)
	}
	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Pool manager was created")
	}

	streamer, err := NewStreamer(cxt, delay)
	if err != nil {
		glog.Exitln("Cannot create streamer:", err)
	}
	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Streamer was created")
	}

	// Init connection manager
	connManager, err := NewConnManager(streamer)
	if err != nil {
		glog.Exitln("Cannot create connection manager:", err)
	}
	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Connection manager was created")
	}

	// Init request verifier
	verifier := NewVerifier(hashSalt)
	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Request verifier was created")
	}

	// Setup GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start goroutine looking for shutdown command
	go func() {
		// Waiting for shutdown command. We don't need of connection
		if _, err := shutdownListener.Accept(); err != nil {
			glog.Errorln("Accepting shutdown connection:", err)
		}
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln("Accepted shutdown command")
		}

		// Closing shutdown listener
		if err := shutdownListener.Close(); err != nil {
			glog.Errorln("Closing shutdown listener:", err)
		}
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln("Shutdown listener was closed")
		}

		// Finishing all goroutines
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln("Canceling root context")
		}
		go cancel()
		time.Sleep(time.Second)
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln("Root context was canceled")
		}

		// Closing working listener
		if err := workingListener.Close(); err != nil {
			glog.Errorln("Closing working listener:", err)
		}
		if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
			glog.Infoln(
				"Working listener was closed.",
				"Server will shutdown with error:",
				"use of closed network connection",
			)
		}
	}()

	if glog.V(INFOLOG_LEVEL_ABOUT_SERVER) {
		glog.Infoln("Starting server")
	}
	// Start server
	err = http.Serve(
		workingListener,
		pwshandler.PoolHandler(poolManager, connManager, verifier),
	)
	if err != nil {
		glog.Errorln("Servering error:", err)
	}

	// Flush log
	glog.Flush()
}
