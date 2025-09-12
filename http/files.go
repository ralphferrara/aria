//||------------------------------------------------------------------------------------------------||
//|| HTTP Package
//|| init.go
//||------------------------------------------------------------------------------------------------||

package http

import (
	"fmt"
	"net/http"
)

//||------------------------------------------------------------------------------------------------||
//|| MountStatic: Mounts a static file server to a given route prefix
//||------------------------------------------------------------------------------------------------||

func (h *HTTPWrapper) Mount(routePrefix string, dir string) error {
	if h.Router == nil {
		return fmt.Errorf("HTTPWrapper [%s] does not use mux backend", h.Name)
	}

	fs := http.FileServer(http.Dir(dir))
	h.Router.PathPrefix(routePrefix).Handler(http.StripPrefix(routePrefix, fs))
	return nil
}
