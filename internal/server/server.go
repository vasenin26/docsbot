package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vasenin26/docsbot/internal/handlers"
)

// NewHandler configures HTTP routes for health and metrics.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.HealthHandler)
	mux.HandleFunc("/readyz", handlers.HealthHandler)
	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

// Run starts HTTP server on given address with reasonable timeouts.
func Run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      NewHandler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv.ListenAndServe()
}
