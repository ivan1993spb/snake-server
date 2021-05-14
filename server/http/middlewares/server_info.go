package middlewares

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

func NewServerInfo(name, version, build string) negroni.Handler {
	serverInfo := fmt.Sprintf("%s/%s (build %s)", name, version, build)

	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.Header().Set("Server", serverInfo)
		next(rw, r)
	})
}
