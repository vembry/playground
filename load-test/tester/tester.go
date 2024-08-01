package tester

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tester[T any] struct {
	parameter T
	cfg       Config
}

// Config contain load-tester configuration
type Config struct {
	Logger   *slog.Logger
	Duration time.Duration
}

// New initiate load-tester instance
func New[T any](cfg Config, parameter T) *tester[T] {
	return &tester[T]{
		parameter: parameter,
		cfg:       cfg,
	}
}

// Do setup and execute load-tester
func (t *tester[T]) Do(execution func(ctx context.Context, l *slog.Logger, parameter T)) {

	// setup otel trace
	tracer := otel.Tracer("load-test-tracer")

	// tracer := otel.GetTracerProvider().Tracer("asd")

	// setup end time
	end := time.Now().Add(t.cfg.Duration)

	for time.Now().Before(end) {
		// start open-telemetry
		ctx, span := tracer.Start(
			context.Background(),
			"load-test-starter",
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// execute test script
		execution(ctx, t.cfg.Logger, t.parameter)
	}

}
