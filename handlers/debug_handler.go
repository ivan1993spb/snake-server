package handlers

import (
	"net/http"
	"net/http/pprof"
)

const URLRouteDebug = "/debug"

func NewDebugHandler() http.Handler {
	h := http.NewServeMux()
	h.HandleFunc("/debug/pprof/", pprof.Index)
	h.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	h.HandleFunc("/debug/pprof/profile", pprof.Profile)
	h.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	h.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return h
}
