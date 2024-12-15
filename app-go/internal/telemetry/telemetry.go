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
		shutdowns []func()        = []func(){}
		ctx       context.Context = context.Background()
	)

	// setup propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// setup resource
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		log.Fatalf("error on initiating otel's resource. err=%v", err)
		return func() {}
	}

	// setup tracer
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel's trace exporter. err=%v", err)
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(traceExporter)),
	)
	otel.SetTracerProvider(traceProvider)
	shutdowns = append(shutdowns, func() {
		traceProvider.Shutdown(context.Background())
	})

	// setup meter
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel's metric exporter. err=%v", err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider)
	shutdowns = append(shutdowns, func() {
		traceProvider.Shutdown(context.Background())
	})

	// return shutdown handler
	return func() {
		for _, shutdown := range shutdowns {
			shutdown()
		}
	}
}
