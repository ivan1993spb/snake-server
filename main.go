package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/ivan1993spb/pwshandler"
	"golang.org/x/net/context"
)

// Infolog leveles
const (
	INFOLOG_LEVEL_SERVER = iota + 1 // Server level
	INFOLOG_LEVEL_POOLS             // Pool level
	INFOLOG_LEVEL_CONNS             // Connection level
)

type errStartingServer struct {
	err error
}

func (e *errStartingServer) Error() string {
	return "starting server error: " + e.err.Error()
}

func main() {

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                  BEGIN PARSING PARAMETERS                   *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

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
		glog.Infoln("checking parameters")

		if len(host) == 0 {
			glog.Warningln("empty host")
		}
		if len(gamePort) == 0 {
			glog.Warningln("empty game port")
		} else if i, e := strconv.Atoi(gamePort); e != nil || i < 1 {
			glog.Warningln("invalid game port")
		}
		if len(sdPort) == 0 {
			glog.Warningln("empty shutdown port")
		} else if i, e := strconv.Atoi(sdPort); e != nil || i < 1 {
			glog.Warningln("invalid shutdown port")
		}
		if len(hashSalt) == 0 {
			glog.Warningln("empty hash salt")
		}
		if poolLimit == 0 {
			glog.Warningln("invalid pool limit")
		}
		if connLimit == 0 {
			glog.Warningln("invalid connection limit per pool")
		}
		if pgW*pgH == 0 {
			glog.Warningln("invalid playground proportions")
		}
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                   END PARSING PARAMETERS                    *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("preparing to start server")
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                  BEGIN CREATING LISTENERS                   *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Working listener is used for game servering
	workingListener, err := net.Listen("tcp", host+":"+gamePort)
	if err != nil {
		glog.Exitln(&errStartingServer{
			fmt.Errorf("cannot create working listener: %s", err),
		})
	}

	// Shutdown listener is used for shutdown command
	shutdownListener, err := net.Listen("tcp", "127.0.0.1:"+sdPort)
	if err != nil {
		glog.Exitln(&errStartingServer{
			fmt.Errorf("cannot create shutdown listener: %s", err),
		})
	}

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("listeners was created")
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                   END CREATING LISTENERS                    *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	cxt, cancel := context.WithCancel(context.Background())

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                     BEGIN INIT MODULES                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Init pool factory
	factory, err := NewPGPoolFactory(cxt, uint8(connLimit),
		uint8(pgW), uint8(pgH))
	if err != nil {
		glog.Exitln(&errStartingServer{err})
	}
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("pool factory was created")
	}

	// Init pool manager which allocates connections on pools
	poolManager, err := NewGamePoolManager(factory, uint8(poolLimit))
	if err != nil {
		glog.Exitln(&errStartingServer{err})
	}
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("pool manager was created")
	}

	// Init connection manager
	connManager := NewConnManager()
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("connection manager was created")
	}

	// Init request verifier
	verifier := NewRequestVerifier(hashSalt)
	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("request verifier was created")
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                      END INIT MODULES                       *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	// Setup GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("starting server")
	}

	// Start goroutine looking for shutdown command
	go func() {
		// Waiting for shutdown command
		if _, err := shutdownListener.Accept(); err != nil {
			glog.Errorln("accepting shutdown connection error:", err)
		}
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("accepted shutdown connection")
		}

		// Closing shutdown listener
		if err := shutdownListener.Close(); err != nil {
			glog.Errorln("closing shutdown listener error:", err)
		}
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("shutdown listener was closed")
		}

		// Finishing all goroutines
		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln("finishing all goroutines on server")
		}
		cancel()
		time.Sleep(time.Second)

		if glog.V(INFOLOG_LEVEL_SERVER) {
			glog.Infoln(
				"closing working listener;",
				"server will shutdown with error:",
				"use of closed network connection",
			)
		}
		// Closing working listener
		if err := workingListener.Close(); err != nil {
			glog.Errorln("closing working listener error:", err)
		}
	}()

	// Starting server
	err = http.Serve(
		workingListener,
		pwshandler.PoolHandler(poolManager, connManager, verifier),
	)
	if err != nil {
		glog.Errorln("servering error:", err)
	}

	// Flush log
	glog.Flush()

	time.Sleep(time.Second)

	if glog.V(INFOLOG_LEVEL_SERVER) {
		glog.Infoln("goodbye")
	}
}
