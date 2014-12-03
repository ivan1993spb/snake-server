package main

import (
	"flag"
	"time"
)

var Config = &struct {
	// Host and port on which game server handles requests
	Host, Port string
	// Port on which server accepts for shutdown request
	ShutdownPort string
	// Max pool number on server
	PoolLimit uint8
	// Max connection number on pool
	ConnLimit uint8
	// Playground size
	PgW, PgH uint8
	// HashSalt for request verifier
	HashSalt string
	// Websocket input and output buffsizes
	WsReadBufferSize, WsWriteBufferSize int
	// Output stram delay
	Delay time.Duration
}{}

func init() {
	flag.StringVar(&Config.Host, "host", "",
		"host on which game server handles requests")
	flag.StringVar(&Config.Port, "port", "8081",
		"port on which game server handles requests")
	flag.StringVar(&Config.ShutdownPort, "shutdown_port", "8082",
		"port on which server accepts for shutdown request")

	var tmp uint
	flag.UintVar(&tmp, "pool_limit", 10, "max pool number on server")
	Config.PoolLimit = uint8(tmp)
	flag.UintVar(&tmp, "conn_limit", 4,
		"max connection number on pool")
	Config.ConnLimit = uint8(tmp)
	flag.UintVar(&tmp, "pg_w", 40, "playground width")
	Config.PgW = uint8(tmp)
	flag.UintVar(&tmp, "pg_h", 28, "playground height")
	Config.PgH = uint8(tmp)

	flag.StringVar(&Config.HashSalt, "hash_salt", "",
		"salt for request verifier")
	flag.IntVar(&Config.WsReadBufferSize, "ws_read_buf", 4096,
		"websocket input buffer size")
	flag.IntVar(&Config.WsWriteBufferSize, "ws_write_buf", 4096,
		"websocket output buffer size")

	var s string
	flag.StringVar(&s, "delay", "150ms", "game stream delay")
	delay, err := time.ParseDuration(s)
	if err == nil {
		Config.Delay = delay
	}
}
