package app

import (
	"api/common"
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// prometheus is an instance containing app's prometheus related lifecycle
type prometheus struct {
	server *http.Server
}

// NewPrometheus is to initiate prometheus instance
func NewPrometheus(cfg *EnvConfig) *prometheus {
	return &prometheus{
		server: constructPrometheusServer(cfg),
	}
}

// Start is to initiate metric provider to be scraped by prometheus
func (p *prometheus) Start() {
	go func() {
		if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start. err=%v", err)
		}
	}()
	common.WatchForExitSignal()
	p.server.Shutdown(context.Background())
}

// constructPrometheusServer is to construct a server to be scraped by prometheus.
// Deliberately setting up different server just for prometheus handler for isolation
func constructPrometheusServer(cfg *EnvConfig) *http.Server {
	// setup handler
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// construct server
	return &http.Server{
		Addr:    cfg.PrometheusHttpAddress,
		Handler: mux,
	}
}
