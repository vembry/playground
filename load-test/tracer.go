package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func newTelemetry() func() {
	os.Setenv("OTEL_SDK_DISABLED", "false")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://host.docker.internal:10002")
	os.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "5000")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_SERVICE_NAME", "load-tester")
	os.Setenv("OTEL_TRACES_SAMPLER", "always_on")

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel exporter. err=%v", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		log.Fatalf("error on initiating otel resource. err=%v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}
}
