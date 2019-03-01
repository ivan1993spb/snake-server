package client

import "net/http"

const (
	URLRouteClient         = "/client"
	URLRouteServerEndpoint = "/"
)

func NewHandler() http.Handler {
	return http.FileServer(newRelateveAssetFS())
}
