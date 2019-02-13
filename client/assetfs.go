// Use command for bindata generation:
// go-bindata-assetfs -nometadata -pkg client -o ./client/bindata.go client/dist/...

package client

import (
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
)

const pathToClient = "client/dist"

func adjustPath(path string) string {
	return strings.TrimPrefix(path, URLRouteClient)
}

type assetFS struct {
	assetFS *assetfs.AssetFS
}

func newAssetFS() *assetFS {
	return &assetFS{
		assetFS: &assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    pathToClient,
		},
	}
}

func (fs *assetFS) Open(name string) (http.File, error) {
	return fs.assetFS.Open(adjustPath(name))
}
