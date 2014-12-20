package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
)

const (
	INFOLOG_LEVEL_SERVER = iota + 1 // Messages about server
	INFOLOG_LEVEL_POOLS             // about pools
	INFOLOG_LEVEL_CONNS             // and connections
)

type errStartingServer struct {
	err error
}

func (e *errStartingServer) Error() string {
	return "Starting server error: " + e.err.Error()
}

func main() {
	var host, gamePort, sdPort, hashSalt string
	flag.StringVar(&host, "host", "",
		"host on which game server handles requests")
	flag.StringVar(&gamePort, "game_port", "8081",
		"port on which game server handles requests")
	flag.StringVar(&sdPort, "shutdown_port", "8082",
		"port on which server accepts for shutdown request")
	flag.StringVar(&hashSalt, "hash_salt", "",
		"salt for request verifier")

	var poolLimit, connLimit, pgW, pgH uint
	flag.UintVar(&poolLimit, "pool_limit", 10,
		"max pool number on server")
	flag.UintVar(&connLimit, "conn_limit", 4,
		"max connection number on pool")
	flag.UintVar(&pgW, "pg_w", 40, "playground width")
	flag.UintVar(&pgH, "pg_h", 28, "playground height")

	flag.Parse()

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Preparing to start server")
	}

	// Working listener is using for game servering
	workingListener, err := net.Listen("tcp", host+":"+gamePort)
	if err != nil {
		glog.Exitln(&errStartingServer{
			fmt.Errorf("Cannot create working listener: %s", err),
		})
	}

	// Shutdown listener is using only for shutdown command
	shutdownListener, err := net.Listen("tcp", "127.0.0.1:"+sdPort)
	if err != nil {
		glog.Exitln(&errStartingServer{
			fmt.Errorf("Cannot create shutdown listener: %s", err),
		})
	}

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Listeners was created")
	}

	// Get root context and cancel func for all goroutines on server
	cxt, cancel := context.WithCancel(context.Background())

	// Init pool factory
	factory, err := NewPGPoolFactory(cxt, uint8(connLimit),
		uint8(pgW), uint8(pgH))
	if err != nil {
		glog.Exitln(&errStartingServer{err})
	}
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Pool factory was created")
	}

	// Init pool manager which allocates connections on pools
	poolManager, err := NewGamePoolManager(factory, uint8(poolLimit))
	if err != nil {
		glog.Exitln(&errStartingServer{err})
	}
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Pool manager was created")
	}

	// Init connection manager
	connManager := NewConnManager()
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Connection manager was created")
	}

	// Init request verifier
	verifier := NewRequestVerifier(hashSalt)
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Request verifier was created")
	}

	// Setup GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start goroutine looking for shutdown command
	go func() {
		// Waiting for shutdown command
		if _, err := shutdownListener.Accept(); err != nil {
			glog.Errorln("Accepting shutdown connection error:", err)
		}
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("Accepted shutdown command")
		}

		// Closing shutdown listener
		if err := shutdownListener.Close(); err != nil {
			glog.Errorln("Closing shutdown listener error:", err)
		}
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("Shutdown listener was closed")
		}

		// Finishing all goroutines
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("Canceling root context")
		}
		go cancel()
		time.Sleep(time.Second)
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("Root context was canceled")
		}

		// Closing working listener
		if err := workingListener.Close(); err != nil {
			glog.Errorln("Closing working listener error:", err)
		}
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln(
				"Working listener was closed.",
				"Server will shutdown with error:",
				"use of closed network connection",
			)
		}
	}()

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("Starting server")
	}

	// Starting server
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
