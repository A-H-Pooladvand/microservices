package trace

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{}

func newTracer() Tracer {
	return Tracer{}
}

func (t Tracer) FromContext(ctx context.Context) trace.Tracer {
	span := trace.SpanFromContext(ctx)

	// Todo:: Change app to config
	return span.TracerProvider().Tracer("app")
}
