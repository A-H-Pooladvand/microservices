package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func FromContext(ctx context.Context) trace.Tracer {
	return trace.SpanFromContext(ctx).TracerProvider().Tracer("app")
}

func Tracer() trace.Tracer {
	return Provider().Tracer("app")
}

func Provider() trace.TracerProvider {
	return otel.GetTracerProvider()
}
