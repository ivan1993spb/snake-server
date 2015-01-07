// Copyright 2015 Pushkin Ivan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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

// TokenVerifierMux verifies each accepted connection by passed token
type TokenVerifierMux struct {
	Mux
	poolManager *GamePoolManager
	hashSalt    string
}

type errCannotCreateTokenVerifierMux struct {
	errStr string
}

func (e *errCannotCreateTokenVerifierMux) Error() string {
	return "cannot create token verifier mux: " + e.errStr
}

func NewTokenVerifierMux(m Mux, pm *GamePoolManager, hs string,
) (*TokenVerifierMux, error) {
	if m == nil {
		return nil, &errCannotCreateTokenVerifierMux{"mux is nil"}
	}
	if pm == nil {
		return nil, &errCannotCreateTokenVerifierMux{
			"pool manager is nil",
		}
	}
	if len(hs) == 0 {
		return nil, &errCannotCreateTokenVerifierMux{
			"empty hash salt",
		}
	}

	hashSalt := sha256.Sum256([]byte(hs))
	return &TokenVerifierMux{m, pm, hex.EncodeToString(hashSalt[:])},
		nil
}

type errConnNotTrusted struct {
	err error
}

func (e *errConnNotTrusted) Error() string {
	return "connection is not trusted: " + e.err.Error()
}

func (v *TokenVerifierMux) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("verifying token hash sum")
	}

	if tokenStr := r.FormValue(FORM_KEY_TOKEN); len(tokenStr) > 0 {
		data, err := base64.StdEncoding.DecodeString(tokenStr)
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

		validSum := sha256.Sum256([]byte(v.hashSalt + token.Part))

		for i := 0; i < sha256.Size; i++ {
			if validSum[i] != sum[i] {
				glog.Errorln(&errConnNotTrusted{
					errors.New("ivalid sum"),
				})
				goto forbidden
			}
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("token hash sum is valid")
			glog.Infoln("verifying token uniqueness")
		}

		for _, request := range v.poolManager.GetRequests() {
			// tokenStr is encoded token
			if tokenStr == request.FormValue(FORM_KEY_TOKEN) {
				if glog.V(INFOLOG_LEVEL_CONNS) {
					glog.Warningln("found equal token")
				}
				goto forbidden
			}
		}

		if glog.V(INFOLOG_LEVEL_CONNS) {
			glog.Infoln("token is unique")
		}
	} else if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Warningln("token was not received")
		goto forbidden
	}

	v.Mux.ServeHTTP(w, r)
	return

forbidden:

	glog.Warningln("forbidden")

	http.Error(w, http.StatusText(http.StatusForbidden),
		http.StatusForbidden)
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
		err := json.NewEncoder(w).Encode(map[string]uint{
			"pool_limit": poolLimit,
			"conn_limit": connLimit,
		})
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PlaygroundSizeHandler(pgW, pgH uint8) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]uint8{
			"playground_width":  pgW,
			"playground_height": pgH,
		})
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PoolCountHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]uint16{
			"pool_count": poolManager.PoolCount(),
		})
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func ConnCountHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]uint32{
			"conn_count": poolManager.ConnCount(),
		})
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PoolInfoListHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(poolManager.PoolInfoList())
		if err != nil {
			glog.Errorln(&errHandleRequest{err})
		}
	})
}

func PoolConnIdsHandler(poolManager *GamePoolManager) http.Handler {
	return JsonHandler(func(w http.ResponseWriter, r *http.Request) {
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

type ReportMux struct {
	Mux
}

func (rm *ReportMux) ServeHTTP(w http.ResponseWriter, r *http.Request,
) {
	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("received request:", r.URL.Path)
	}

	rm.Mux.ServeHTTP(w, r)

	if glog.V(INFOLOG_LEVEL_CONNS) {
		glog.Infoln("closing connection:", r.URL.Path)
	}
}
