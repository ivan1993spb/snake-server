package client

import "net/http"

const URLRouteClient = "/client"

func NewHandler() http.Handler {
	return http.FileServer(newAssetFS())
}
