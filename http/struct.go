//||------------------------------------------------------------------------------------------------||
//|| HTTP Package
//|| struct.go
//||------------------------------------------------------------------------------------------------||

package http

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"net/http"

	"github.com/gorilla/mux"
)

//||------------------------------------------------------------------------------------------------||
//|| Backend Enum
//||------------------------------------------------------------------------------------------------||

type HTTPBackend string

const (
	BackendMux  HTTPBackend = "MUX"
	BackendHTTP HTTPBackend = "HTTP"
)

//||------------------------------------------------------------------------------------------------||
//|| HTTPWrapper
//||------------------------------------------------------------------------------------------------||

type HTTPWrapper struct {
	Name       string
	HTTPConfig HTTPConfig
	Server     *http.Server
	Handler    http.Handler
	Router     *mux.Router
	ServeMux   *http.ServeMux
}

//||------------------------------------------------------------------------------------------------||
//|| HTTPConfig
//||------------------------------------------------------------------------------------------------||

type HTTPConfig struct {
	Backend      HTTPBackend
	Port         string
	Middleware   func(next http.Handler) http.Handler
	ErrorHandler func(next http.Handler) http.Handler
	CorsHandler  func(next http.Handler) http.Handler
}

//||------------------------------------------------------------------------------------------------||
//|| Defaults
//||------------------------------------------------------------------------------------------------||

type Defaults struct {
	Cors       func(http.Handler) http.Handler
	Middleware func(http.Handler) http.Handler
	Error      func(http.Handler) http.Handler
	NotFound   http.Handler
}

//||------------------------------------------------------------------------------------------------||
//|| Package Defaults (can be overridden)
//||------------------------------------------------------------------------------------------------||

var httpDefaults = Defaults{
	Cors:       nil,               // disabled
	Middleware: nil,               // disabled
	Error:      nil,               // disabled
	NotFound:   NotFoundHandler(), // simple 404
}

//||------------------------------------------------------------------------------------------------||
//|| SetDefaults
//||------------------------------------------------------------------------------------------------||

func SetDefaults(d Defaults) {
	if d.Cors != nil {
		httpDefaults.Cors = d.Cors
	}
	if d.Middleware != nil {
		httpDefaults.Middleware = d.Middleware
	}
	if d.Error != nil {
		httpDefaults.Error = d.Error
	}
	if d.NotFound != nil {
		httpDefaults.NotFound = d.NotFound
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Wrapper SetDefaults
//||------------------------------------------------------------------------------------------------||

func (h *HTTPWrapper) SetDefaults(d Defaults) {
	if d.Cors != nil {
		h.HTTPConfig.CorsHandler = d.Cors
	}
	if d.Middleware != nil {
		h.HTTPConfig.Middleware = d.Middleware
	}
	if d.Error != nil {
		h.HTTPConfig.ErrorHandler = d.Error
	}
}

//||------------------------------------------------------------------------------------------------||
//|| ResetDefaults
//||------------------------------------------------------------------------------------------------||

func ResetDefaults() {
	httpDefaults = Defaults{
		Cors:       nil,
		Middleware: nil,
		Error:      nil,
		NotFound:   NotFoundHandler(),
	}
}

//||------------------------------------------------------------------------------------------------||
//|| DefaultConfig
//||------------------------------------------------------------------------------------------------||

func DefaultConfig(port string) HTTPConfig {
	return HTTPConfig{
		Backend: BackendMux,
		Port:    port,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| ApplyDefaults
//||------------------------------------------------------------------------------------------------||

func ApplyDefaults(cfg HTTPConfig) HTTPConfig {
	out := cfg
	if out.CorsHandler == nil && httpDefaults.Cors != nil {
		out.CorsHandler = httpDefaults.Cors
	}
	if out.Middleware == nil && httpDefaults.Middleware != nil {
		out.Middleware = httpDefaults.Middleware
	}
	if out.ErrorHandler == nil && httpDefaults.Error != nil {
		out.ErrorHandler = httpDefaults.Error
	}
	return out
}
