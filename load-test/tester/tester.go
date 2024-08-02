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

	// load otel trace
	tracer := otel.Tracer("load-test-tracer")

	// setup end time
	end := time.Now().Add(t.cfg.Duration)

	for time.Now().Before(end) {
		t.do(tracer, execution)
	}

}

func (t *tester[T]) do(tracer trace.Tracer, execution func(ctx context.Context, l *slog.Logger, parameter T)) {
	// start open-telemetry
	ctx, span := tracer.Start(
		context.Background(),
		"load-test-starter",
	)
	defer span.End()

	// execute test script
	execution(ctx, t.cfg.Logger, t.parameter)
}
