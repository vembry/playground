package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metric struct {
	server             *http.Server
	httpRequestLatency *prometheus.HistogramVec
}

func (m *metric) Name() string {
	return "httpserver"
}

func NewMetric(cfg *EnvConfig) *metric {
	return &metric{
		server: constructPrometheusServer(cfg),
		httpRequestLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_latency_milliseconds",
			Help:    "Histogram for Latency of requests in milliseconds.",
			Buckets: []float64{100, 200, 300, 500, 700, 1000, 1500, 2000, 5000, 10000},
		}, []string{"route", "method", "status_code"}),
	}
}

func (m *metric) RecordInbound(route string, method string, statusCode string, duration time.Duration) {
	m.httpRequestLatency.WithLabelValues(route, method, statusCode).Observe(float64(duration.Milliseconds()))
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

// Start is to initiate metric provider to be scraped by prometheus
func (m *metric) Start() {
	go func() {
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start. err=%v", err)
		}
	}()
}

// Stop is to shutdown metric provider
func (m *metric) Stop() {
	m.server.Shutdown(context.Background())
}
