//||------------------------------------------------------------------------------------------------||
//|| HTTP Package
//|| init.go
//||------------------------------------------------------------------------------------------------||

package http

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Init
//||------------------------------------------------------------------------------------------------||

func Init(cfg *config.Config) (map[string]*HTTPWrapper, error) {

	//||------------------------------------------------------------------------------------------------||
	//|| Empty Output
	//||------------------------------------------------------------------------------------------------||

	out := map[string]*HTTPWrapper{}
	if cfg == nil || cfg.HTTP == nil {
		return out, nil
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Loop Configured HTTP Entries
	//||------------------------------------------------------------------------------------------------||

	for name, h := range cfg.HTTP {

		//||------------------------------------------------------------------------------------------------||
		//|| Normalize Backend
		//||------------------------------------------------------------------------------------------------||

		backend := BackendHTTP
		switch h.Backend {
		case "mux", "MUX":
			backend = BackendMux
		}

		//||------------------------------------------------------------------------------------------------||
		//|| Build Config (defaults will fill handlers)
		//||------------------------------------------------------------------------------------------------||

		hcfg := HTTPConfig{
			Backend:      backend,
			Port:         strconv.Itoa(h.Port),
			Middleware:   nil,
			ErrorHandler: nil,
			CorsHandler:  nil,
		}
		hcfg = ApplyDefaults(hcfg)

		//||------------------------------------------------------------------------------------------------||
		//|| Start Server
		//||------------------------------------------------------------------------------------------------||

		out[name] = InitHTTP(name, hcfg)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Return Map
	//||------------------------------------------------------------------------------------------------||

	return out, nil
}

//||------------------------------------------------------------------------------------------------||
//|| InitHTTP
//||------------------------------------------------------------------------------------------------||

func InitHTTP(name string, cfg HTTPConfig) *HTTPWrapper {

	//||------------------------------------------------------------------------------------------------||
	//|| Fill handler chain from package defaults where cfg has nils
	//||------------------------------------------------------------------------------------------------||

	cfg = ApplyDefaults(cfg)

	//||------------------------------------------------------------------------------------------------||
	//|| Select Router
	//||------------------------------------------------------------------------------------------------||

	var baseHandler http.Handler
	var mrouter *mux.Router
	var smux *http.ServeMux

	switch cfg.Backend {
	case BackendMux:
		mrouter = mux.NewRouter()
		if httpDefaults.NotFound != nil {
			mrouter.NotFoundHandler = httpDefaults.NotFound
		}
		baseHandler = mrouter
	case BackendHTTP:
		smux = http.NewServeMux()
		if httpDefaults.NotFound != nil {
			smux.Handle("/", httpDefaults.NotFound)
		}
		baseHandler = smux
	default:
		smux = http.NewServeMux()
		if httpDefaults.NotFound != nil {
			smux.Handle("/", httpDefaults.NotFound)
		}
		baseHandler = smux
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Chain: CORS -> Middleware -> ErrorHandler
	//||------------------------------------------------------------------------------------------------||

	handler := baseHandler
	if cfg.CorsHandler != nil {
		handler = cfg.CorsHandler(handler)
	}
	if cfg.Middleware != nil {
		handler = cfg.Middleware(handler)
	}
	if cfg.ErrorHandler != nil {
		handler = cfg.ErrorHandler(handler)
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Server
	//||------------------------------------------------------------------------------------------------||

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Wrapper
	//||------------------------------------------------------------------------------------------------||

	wrapper := &HTTPWrapper{
		Name:       name,
		HTTPConfig: cfg,
		Server:     srv,
		Handler:    handler,
		Router:     mrouter,
		ServeMux:   smux,
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Wrapper
	//||------------------------------------------------------------------------------------------------||

	return wrapper
}

//||------------------------------------------------------------------------------------------------||
//|| Start: Starts the HTTP server
//||------------------------------------------------------------------------------------------------||

func (h *HTTPWrapper) Start() error {
	if h.Server == nil {
		return fmt.Errorf("HTTP server [%s] is not initialized", h.Name)
	}

	// Rebuild handler chain using latest config values
	handler := h.RouterOrMux()
	if h.HTTPConfig.CorsHandler != nil {
		handler = h.HTTPConfig.CorsHandler(handler)
	}
	if h.HTTPConfig.Middleware != nil {
		handler = h.HTTPConfig.Middleware(handler)
	}
	if h.HTTPConfig.ErrorHandler != nil {
		handler = h.HTTPConfig.ErrorHandler(handler)
	}

	h.Server.Handler = handler // ← apply rebuilt chain

	fmt.Printf("\n[HTTP] - Starting HTTP server [%s] on port %s (backend: %s)\n", h.Name, h.HTTPConfig.Port, h.HTTPConfig.Backend)

	go func() {
		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[HTTP] - Server [%s] failed: %v\n", h.Name, err)
			os.Exit(1)
		}
	}()

	return nil
}

//||------------------------------------------------------------------------------------------------||
//|| Router or Mux
//||------------------------------------------------------------------------------------------------||

func (h *HTTPWrapper) RouterOrMux() http.Handler {
	if h.Router != nil {
		return h.Router
	}
	if h.ServeMux != nil {
		return h.ServeMux
	}
	return http.DefaultServeMux
}
