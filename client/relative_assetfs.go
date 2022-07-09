package client

import (
	"net/http"
	"path/filepath"
	"strings"
)

const pathToClient = "public/dist"

func adjustPath(path, prefix, root string) string {
	relative := strings.TrimPrefix(path, prefix)

	return filepath.Join(root, relative)
}

type relativeAssetFS struct {
	fs http.FileSystem

	// Place where files reside
	root string

	// URL prefix which should be removed from the path
	prefix string
}

func newRelativeAssetFS() *relativeAssetFS {
	return &relativeAssetFS{
		fs:     http.FS(clientFS),
		root:   pathToClient,
		prefix: URLRouteClient,
	}
}

func (fs *relativeAssetFS) Open(name string) (http.File, error) {
	return fs.fs.Open(adjustPath(name, fs.prefix, fs.root))
}
