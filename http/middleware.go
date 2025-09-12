//||------------------------------------------------------------------------------------------------||
//|| HTTP Package
//|| middleware.go
//||------------------------------------------------------------------------------------------------||

package http

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//||------------------------------------------------------------------------------------------------||
//|| CORS Middleware (simple allow-list)
//||------------------------------------------------------------------------------------------------||

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, o := range allowedOrigins {
		allowed[strings.TrimSpace(o)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// wildcard
			if _, ok := allowed["*"]; ok && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

//||------------------------------------------------------------------------------------------------||
//|| loggingResponseWriter: capture status + size
//||------------------------------------------------------------------------------------------------||

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	if lrw.status == 0 {
		// implicit 200 if Write is called before WriteHeader
		lrw.status = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(p)
	lrw.size += n
	return n, err
}

//||------------------------------------------------------------------------------------------------||
//|| Logger Middleware (placeholder)
//||------------------------------------------------------------------------------------------------||

func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{ResponseWriter: w}

			// call downstream
			next.ServeHTTP(lrw, r)

			dur := time.Since(start)

			// route template if using mux (nice for grouping)
			var route string
			if rt := mux.CurrentRoute(r); rt != nil {
				if tmpl, err := rt.GetPathTemplate(); err == nil {
					route = tmpl
				}
			}
			if route == "" {
				route = r.URL.Path
			}

			// include querystring? switch to r.URL.String() if you prefer
			fmt.Printf("[HTTP][%3d] %-4s %-40s dur=%s\n",
				lrw.status, r.Method, route, dur)
		})
	}
}

//||------------------------------------------------------------------------------------------------||
//|| NotFound helper
//||------------------------------------------------------------------------------------------------||

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
}
