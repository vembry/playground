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

func init() {
	// TODO: remove this,
	// this is a quick way to define otel prerequisite
	os.Setenv("OTEL_SDK_DISABLED", "false")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://host.docker.internal:10002")
	os.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "5000")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_SERVICE_NAME", "broker")
	os.Setenv("OTEL_TRACES_SAMPLER", "always_on")
}

func NewTracer() func() {
	ctx := context.Background()

	// setup resource
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		log.Fatalf("error on initiating otel's resource. err=%v", err)
	}

	// initiate tracer
	tracerShutdownHandler, err := setTracerProvider(ctx, res)
	if err != nil {
		log.Fatalf("error on setting tracer provider. err=%v", err)
	}

	return func() {
		if err := tracerShutdownHandler(); err != nil {
			log.Printf("error on shutting down tracer. err=%v", err)
		}
	}
}

func setTracerProvider(ctx context.Context, res *resource.Resource) (func() error, error) {
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
		return provider.Shutdown(ctx)
	}, nil
}
