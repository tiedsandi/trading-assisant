package app

import (
	"context"
	"net/http"
	"time"
)

type RouteRegistrar interface {
	RegisterRoutes(*http.ServeMux)
}

// NewHandler composes health/readiness and optional feature routes.
func NewHandler(checkDatabase func(context.Context) error, features ...RouteRegistrar) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if checkDatabase == nil || checkDatabase(ctx) != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("{\"status\":\"not_ready\"}\n"))
			return
		}
		_, _ = w.Write([]byte("{\"status\":\"ready\"}\n"))
	})
	for _, feature := range features {
		feature.RegisterRoutes(mux)
	}
	return mux
}
