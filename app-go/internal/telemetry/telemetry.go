package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func New() func() {
	var (
		shutdowns  []func()        = []func(){}
		ctx        context.Context = context.Background()
		shutdownFn func()          = func() {
			for _, shutdown := range shutdowns {
				shutdown()
			}
		}
	)

	// setup propagation
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// setup resource
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		log.Fatalf("error on initiating otel's resource. err=%v", err)
		return shutdownFn
	}

	// setup tracer
	traceProvider := newTraceProvider(ctx, res)
	otel.SetTracerProvider(traceProvider)
	shutdowns = append(shutdowns, func() {
		traceProvider.Shutdown(context.Background())
	})

	// setup meter
	meterProvider := newMeterProvider(ctx, res)
	otel.SetMeterProvider(meterProvider)
	shutdowns = append(shutdowns, func() {
		traceProvider.Shutdown(context.Background())
	})

	// return shutdown handler
	return shutdownFn
}

func newTraceProvider(ctx context.Context, res *resource.Resource) *trace.TracerProvider {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel's trace exporter. err=%v", err)
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(traceExporter)),
	)

	return traceProvider
}

func newMeterProvider(ctx context.Context, res *resource.Resource) *metric.MeterProvider {
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel's metric exporter. err=%v", err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	return meterProvider
}
