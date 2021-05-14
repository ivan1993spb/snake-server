package middlewares

import (
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func NewCORS() negroni.Handler {
	return cors.AllowAll()
}
