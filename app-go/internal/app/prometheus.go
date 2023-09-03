package app

import (
	"app-go/common"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metric is an instance containing app's prometheus related lifecycle
type metric struct {
	server *http.Server

	// httpRequestLatency is a histogram to record http request-response latency
	httpRequestLatency *prometheus.HistogramVec

	// taskLatency is a histogram to record worker-task latency
	workerLatency *prometheus.HistogramVec
}

// NewMetric is to initiate prometheus instance
func NewMetric(cfg *EnvConfig) *metric {
	return &metric{
		server: constructPrometheusServer(cfg),
		httpRequestLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_latency_milliseconds",
			Help:    "Histogram for Latency of requests in milliseconds.",
			Buckets: []float64{100, 200, 300, 500, 700, 1000, 1500, 2000, 5000, 10000},
		}, []string{"route", "method", "status_code"}),
		workerLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "worker_latency_milliseconds",
			Help:    "Histogram for Latency of requests in milliseconds.",
			Buckets: []float64{100, 200, 300, 500, 700, 1000, 1500, 2000, 5000, 10000},
		}, []string{"task_type", "is_error"}),
	}
}

// Start is to initiate metric provider to be scraped by prometheus
func (m *metric) Start() {
	go func() {
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start. err=%v", err)
		}
	}()
	common.WatchForExitSignal()
	m.server.Shutdown(context.Background())
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

// GinRequest is a gin-gonic/gin middleware to record http request activities
func (m *metric) GinRequest(c *gin.Context) {
	// initiate time
	start := time.Now()

	// continue request
	c.Next()

	// construct values to be passed onto histogram observation for latency
	duration := time.Since(start)
	route := c.FullPath()
	method := c.Request.Method
	statusCode := strconv.Itoa(c.Writer.Status())

	// save latency observation
	m.httpRequestLatency.WithLabelValues(route, method, statusCode).Observe(float64(duration.Milliseconds()))
}

// AsynqTask is a hibiken/asynq middleware to record worker task activities
func (m *metric) AsynqTask(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		// initiate time
		start := time.Now()

		// continue process
		err := h.ProcessTask(ctx, t)

		// construct values to be passed onto histogram observation for latency
		duration := time.Since(start)
		taskType := t.Type()
		isError := "false"
		if err != nil {
			isError = "true"
		}

		// save latency observation
		m.workerLatency.WithLabelValues(taskType, isError).Observe(float64(duration.Milliseconds()))

		return err
	})
}
