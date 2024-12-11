package app

import (
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func NewTelemetry(ctx context.Context) (shutdown func() error, err error) {
	var (
		shutdownFuncs []func() error
		// shutdown      func() error
		// err           error
	)

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func() error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn())
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) { err = errors.Join(inErr, shutdown()) }

	disabled := os.Getenv("OTEL_SDK_DISABLED") == "true"

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		handleErr(err)
		return
	}

	tracerShutdown, err := newTracerProvider(ctx, disabled, res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerShutdown)

	meterShutdown, err := newMeterProvider(ctx, disabled, res)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterShutdown)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		handleErr(err)
		return
	}

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context, disabled bool, res *resource.Resource) (func() error, error) {
	if disabled {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		return func() error { return nil }, nil
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	provider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	return func() error {
		return provider.Shutdown(context.Background())
	}, nil
}

func newMeterProvider(ctx context.Context, disabled bool, res *resource.Resource) (func() error, error) {
	if disabled {
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
		return func() error { return nil }, nil
	}

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)

	return func() error {
		return provider.Shutdown(context.Background())
	}, nil
}
