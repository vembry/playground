package app

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewTracer() func() error {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("error on initiating otel exporter. err=%v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(resource.Environment()),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() error {
		if err := exporter.Shutdown(ctx); err != nil {
			return err
		}

		return tracerProvider.Shutdown(ctx)
	}
}
