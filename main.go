package main

import (
	"flag"
	"net"
	"net/http"
	"runtime"
	"time"

	"bitbucket.org/pushkin_ivan/pool-websocket-handler"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

func main() {
	var host, port, shutdownPort, hashSalt, delay string
	flag.StringVar(&host, "host", "",
		"host on which game server handles requests")
	flag.StringVar(&port, "port", "8081",
		"port on which game server handles requests")
	flag.StringVar(&shutdownPort, "shutdown_port", "8082",
		"port on which server accepts for shutdown request")
	flag.StringVar(&hashSalt, "hash_salt", "",
		"salt for request verifier")
	flag.StringVar(&delay, "delay", "150ms", "game stream delay")

	var poolLimit, connLimit, pgW, pgH uint
	flag.UintVar(&poolLimit, "pool_limit", 10,
		"max pool number on server")
	flag.UintVar(&connLimit, "conn_limit", 4,
		"max connection number on pool")
	flag.UintVar(&pgW, "pg_w", 40, "playground width")
	flag.UintVar(&pgH, "pg_h", 28, "playground height")

	var wsReadBufferSize, wsWriteBufferSize int
	flag.IntVar(&wsReadBufferSize, "ws_read_buf", 4096,
		"websocket input buffer size")
	flag.IntVar(&wsWriteBufferSize, "ws_write_buf", 4096,
		"websocket output buffer size")

	flag.Parse()

	// Working listener is used for game servering
	workingListener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		glog.Exitln("Cannot create working listener:", err)
	}

	// Shutdown listener is used only for shutdown command. Listening
	// only for local requests
	shutdownListener, err := net.Listen("tcp", "127.0.0.1:"+
		shutdownPort)
	if err != nil {
		glog.Exitln("Cannot create shutdown listener:", err)
	}

	if glog.V(4) {
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
	if glog.V(4) {
		glog.Infoln("Pool factory was created")
	}

	// Init pool manager which allocates connections on pools
	poolManager, err := NewGamePoolManager(factory, uint8(poolLimit))
	if err != nil {
		glog.Exitln("Cannot create pool manager:", err)
	}
	if glog.V(4) {
		glog.Infoln("Pool manager was created")
	}

	streamDelay, err := time.ParseDuration(delay)
	if err != nil {
		glog.Exitln("Invalid delay:", err)
	}

	streamer, err := NewStreamer(cxt, streamDelay)
	if err != nil {
		glog.Exitln("Cannot create streamer:", err)
	}
	if glog.V(4) {
		glog.Infoln("Streamer was created")
	}

	// Init connection manager
	connManager, err := NewConnManager(streamer)
	if err != nil {
		glog.Exitln("Cannot create connection manager:", err)
	}
	if glog.V(4) {
		glog.Infoln("Connection manager was created")
	}

	// Init request verifier
	verifier := NewVerifier(hashSalt)
	if glog.V(4) {
		glog.Infoln("Request verifier was created")
	}

	// Configure websocket upgrader
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
		// Don't check origin
		CheckOrigin: func(*http.Request) bool { return true },
	}

	// Create pool handler
	handler := pwshandler.NewPoolHandler(
		poolManager, connManager, verifier, upgrader)
	if glog.V(4) {
		glog.Infoln("Game handler was init")
	}

	// Setup GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start goroutine looking for shutdown command
	go func() {
		// Waiting for shutdown command. We don't need of connection
		if _, err := shutdownListener.Accept(); err != nil {
			glog.Errorln("Accepting shutdown connection:", err)
		}
		if glog.V(3) {
			glog.Infoln("Accepted shutdown command")
		}

		// Closing shutdown listener
		if err := shutdownListener.Close(); err != nil {
			glog.Errorln("Closing shutdown listener:", err)
		}
		if glog.V(4) {
			glog.Infoln("Shutdown listener was closed")
		}

		// Finishing all goroutines
		cancel()
		if glog.V(3) {
			glog.Infoln("Root context was canceled")
		}
		if glog.V(4) {
			glog.Infoln("Wait...")
		}
		time.Sleep(time.Second * 2)

		// Closing working listener
		if err := workingListener.Close(); err != nil {
			glog.Errorln("Closing working listener:", err)
		}
		if glog.V(4) {
			glog.Infoln(
				"Working listener was closed.",
				"Server will shutdown with error:",
				"use of closed network connection",
			)
		}
	}()

	if glog.V(4) {
		glog.Infoln("Starting server")
	}

	// Start server
	if err = http.Serve(workingListener, handler); err != nil {
		glog.Errorln("Servering error:", err)
	}
}
