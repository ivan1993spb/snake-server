package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ivan1993spb/snake-server/pkg/openapi/server"
)

func Serve() {
	r := chi.NewRouter()
	r.Mount("/", server.Handler(&Handler{}))
	//e := echo.New()
	//openapi.RegisterHandlers(e, &Handler{})
	//e.Start("0.0.0.0:9219")
	if err := http.ListenAndServe("0.0.0.0:9219", r); err != nil {
		fmt.Println(err)
	}
}
