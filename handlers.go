package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
)

type Mux interface {
	http.Handler
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(
		http.ResponseWriter,
		*http.Request,
	))
	Handler(r *http.Request) (h http.Handler, pattern string)
}

const FORM_KEY_TOKEN = "token"

type Token struct {
	Sum  string `json:"sum"`
	Part string `json:"part"`
}

// SecurityMux verifies each accepted connection by passed token
type SecurityMux struct {
	*http.ServeMux
	hashSalt string
}

func NewSecurityMux(hashSalt string) (Mux, error) {
	if len(hashSalt) > 0 {
		return &SecurityMux{http.NewServeMux(), hashSalt}, nil
	}

	return nil,
		errors.New("cannot create security mux: empty hash salt")
}

type errConnNotTrusted struct {
	err error
}

func (e *errConnNotTrusted) Error() string {
	return "connection is not trusted: " + e.err.Error()
}

func (sm *SecurityMux) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("verifying connection token")
	}

	if data := r.FormValue(FORM_KEY_TOKEN); len(data) > 0 {
		data, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			glog.Errorln(&errConnNotTrusted{err})
			goto forbidden
		}

		var token *Token
		if err = json.Unmarshal(data, &token); err != nil {
			glog.Errorln(&errConnNotTrusted{err})
			goto forbidden
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("token was received and parsed")
			glog.Infoln("checking sum")
		}

		// Checking token

		sum, err := hex.DecodeString(token.Sum)
		if err != nil {
			glog.Errorln(&errConnNotTrusted{err})
			goto forbidden
		}

		if len(sum) != sha256.Size {
			glog.Errorln(&errConnNotTrusted{
				errors.New("ivalid sum size")},
			)
			goto forbidden
		}

		validSum := sha256.Sum256([]byte(sm.hashSalt + token.Part))

		for i := 0; i < sha256.Size; i++ {
			if validSum[i] != sum[i] {
				glog.Errorln(&errConnNotTrusted{
					errors.New("ivalid sum"),
				})
				goto forbidden
			}
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("token is valid")
		}
	} else if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Warningln("token was not received")
		goto forbidden
	}

	sm.ServeMux.ServeHTTP(w, r)
	return

forbidden:

	http.Error(w, http.StatusText(http.StatusForbidden),
		http.StatusForbidden)
}

// UniqueRequestsHandler verifies connection uniqueness by token
func UniqueRequestsHandler(h http.Handler,
	poolManager *GamePoolManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("verifying connection uniqueness")
		}

		if token := r.FormValue(FORM_KEY_TOKEN); len(token) > 0 {
			for _, request := range poolManager.GetRequests() {
				if token == request.FormValue(FORM_KEY_TOKEN) {
					if glog.V(INFOLOG_LEVEL_CONNS) {
						glog.Warningln("found equal token")
					}
					goto forbidden
				}
			}

			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Infoln("connection is unique")
			}
		} else {
			if glog.V(INFOLOG_LEVEL_CONNS) {
				glog.Warningln("token was not received")
			}
			goto forbidden
		}

		h.ServeHTTP(w, r)
		return

	forbidden:

		http.Error(w, http.StatusText(http.StatusForbidden),
			http.StatusForbidden)
	}
}

// JsonHandler inserts json content-type in response
func JsonHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(
			"Content-Type",
			"application/json; charset=utf-8",
		)

		h.ServeHTTP(w, r)
	}
}

type errHandleRequest struct {
	err error
}

func (e *errHandleRequest) Error() string {
	return "cannot handle request: " + e.err.Error()
}

func ServerLimitsHandler(poolLimit, connLimit uint) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received request for limits")
		}

		_, err := fmt.Fprintf(w, `{"pool_limit":%d,"conn_limit":%d}`,
			poolLimit, connLimit)
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PlaygroundSizeHandler(pgW, pgH uint8) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received request for playground size")
		}

		_, err := fmt.Fprintf(
			w, `{"playground_width":%d,"playground_height":%d}`,
			pgW, pgH,
		)
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func ConnCountHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received request for count of opened conns")
		}

		_, err := fmt.Fprintf(w, `{"conn_count":%d}`,
			poolManager.ConnCount())
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PoolInfoListHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received request for pool info list")
		}

		err := json.NewEncoder(w).Encode(poolManager.PoolInfoList())
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PoolConnIdsHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, r *http.Request) {
		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("received request for pool connection ids")
		}

		// Connection ids
		var ids []uint16

		if id, err := strconv.Atoi(
			r.FormValue(FORM_KEY_POOL_ID)); err != nil {
			glog.Errorln("cannot get pool id:", err)
		} else {
			id := uint16(id)
			pool, err := poolManager.GetPool(id)
			if err != nil {
				glog.Errorln("cannot get pool:", err)
			} else {
				ids = pool.ConnIds()
			}
		}

		if err := json.NewEncoder(w).Encode(ids); err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}
