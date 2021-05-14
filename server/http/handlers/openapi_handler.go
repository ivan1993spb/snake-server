//go:generate go-bindata-assetfs -nometadata -prefix ../ -pkg handlers ../openapi.yaml

package handlers

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

const URLRouteOpenAPI = "/openapi.yaml"

func NewOpenAPIHandler() http.Handler {
	return http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
	})
}
