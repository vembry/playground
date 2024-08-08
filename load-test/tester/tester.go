package tester

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tester struct {
	cfg Config
}

// Config contain load-tester configuration
type Config struct {
	Logger   *slog.Logger
	Duration time.Duration
}

// New initiate load-tester instance
func New(cfg Config) *tester {
	return &tester{
		cfg: cfg,
	}
}

// Do setup and execute load-tester
func (t *tester) Do(execution func(ctx context.Context, l *slog.Logger)) {

	// load otel trace
	tracer := otel.Tracer("load-test-tracer")

	// setup end time
	end := time.Now().Add(t.cfg.Duration)

	for time.Now().Before(end) {
		t.do(tracer, execution)
	}
}

func (t *tester) do(tracer trace.Tracer, execution func(ctx context.Context, l *slog.Logger)) {
	// start open-telemetry
	ctx, span := tracer.Start(
		context.Background(),
		"load-test-starter",
	)
	defer span.End()

	// execute test script
	execution(ctx, t.cfg.Logger)
}
