package backend

import (
	"net/http"

	"github.com/floj/loxy/config"
)

func newFileServer(c *config.BackendFileServer) http.Handler {
	root := http.Dir(c.Root)
	return http.FileServer(root)
}
