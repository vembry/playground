package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	customHandler := &traceHandler{Handler: handler}
	return slog.New(customHandler)
}
