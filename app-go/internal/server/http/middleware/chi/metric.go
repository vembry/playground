package chi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
)

const (
	metricNameRequestDurationMs = "http.server.request.duration"
	metricUnitRequestDurationMs = "ms"
	metricDescRequestDurationMs = "Measures the latency of HTTP requests processed by the server, in milliseconds."
)

// NewRequestDurationMillis is a copy of riandyrn/otelchi's metric with adjustment to follow otel conventions
func NewRequestDurationMillis(serverName string) func(next http.Handler) http.Handler {
	// init metric, here we are using histogram for capturing request duration
	histogram, err := otel.GetMeterProvider().Meter(metricNameRequestDurationMs).Int64Histogram(
		metricNameRequestDurationMs,
		otelmetric.WithDescription(metricDescRequestDurationMs),
		otelmetric.WithUnit(metricUnitRequestDurationMs),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to create %s histogram: %v", metricNameRequestDurationMs, err))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// capture the start time of the request
			startTime := time.Now()

			// execute next http handler
			next.ServeHTTP(w, r)

			// try to extract path from context
			ctx := chi.RouteContext(r.Context())
			path := ctx.RoutePattern()

			// record the request duration
			duration := time.Since(startTime)

			// construct attributes
			attrs := make([]attribute.KeyValue, 0)
			attrs = append(attrs, attribute.String("http.url.path", path)) // adding path into the metric
			attrs = append(attrs, httpconv.ServerRequest(serverName, r)...)

			histogram.Record(
				r.Context(),
				int64(duration.Milliseconds()),
				otelmetric.WithAttributes(attrs...),
			)
		})
	}
}
