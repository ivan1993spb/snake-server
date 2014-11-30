package main

import (
	"flag"
	"net"
	"net/http"
	"runtime"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

func main() {
	flag.Parse()

	// Working listener is used for game servering
	workingListener, err :=
		net.Listen("tcp", Config.Host+":"+Config.Port)
	if err != nil {
		glog.Exitln("Cannot create working listener", err)
	}

	// Shutdown listener is used only for shutdown command
	shutdownListener, err :=
		net.Listen("tcp", "127.0.0.1:"+Config.ShutdownPort)
	if err != nil {
		glog.Exitln("Cannot create shutdown listener", err)
	}

	// Gets root context and cancel func for all goroutines on server
	cxt, cancel := context.WithCancel(context.Background())

	// Init pool factory
	factory := NewPGPoolFactory(cxt, Config.ConnLimit,
		Config.PgW, Config.PgH)

	// Init pool manager which allocates connections on pools
	poolManager := NewGamePoolManager(factory, Config.PoolLimit)

	// Init connection manager
	connManager := NewConnManager()

	// Init verifier
	verifier := NewVerifier(Config.HashSalt)

	// Configure websocket upgrader
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  Config.WsReadBufferSize,
		WriteBufferSize: Config.WsWriteBufferSize,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	// Create pool handler
	handler := pwshandler.NewPoolHandler(
		poolManager, connManager, verifier, upgrader)

	// Shutdown goroutine
	go func() {
		// Waiting for shutdown command
		if _, err := shutdownListener.Accept(); err != nil {
			glog.Errorln("Accepting shutdown connection:", err)
		}

		// Closing shutdown listener
		if err := shutdownListener.Close(); err != nil {
			glog.Errorln("Closing shutdown listener:", err)
		}

		// Finishing all server goroutines
		cancel()

		// Closing working listener
		if err := workingListener.Close(); err != nil {
			glog.Errorln("Closing working listener:", err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())
	if err = http.Serve(workingListener, handler); err != nil {
		glog.Errorln("Game servering error:", err)
	}
}

// const LOG_DEBUG_LEVEL glog.Level = 1

// func DebugInfoln(args ...interface{}) {
// 	if glog.V(LOG_DEBUG_LEVEL) {
// 		glog.Infoln(args...)
// 	}
// }
