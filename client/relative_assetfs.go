//go:generate go-bindata-assetfs -nometadata -pkg client dist/...

package client

import (
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
)

const pathToClient = "dist"

func adjustPath(path string) string {
	return strings.TrimPrefix(path, URLRouteClient)
}

type relativeAssetFS struct {
	assetFS *assetfs.AssetFS
}

func newRelativeAssetFS() *relativeAssetFS {
	return &relativeAssetFS{
		assetFS: &assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    pathToClient,
		},
	}
}

func (fs *relativeAssetFS) Open(name string) (http.File, error) {
	return fs.assetFS.Open(adjustPath(name))
}
