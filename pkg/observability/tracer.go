// Package observability provides unified tracing and metrics capabilities.
package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps OpenTelemetry tracer functionality.
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new Tracer instance.
func NewTracer(name string) *Tracer {
	return &Tracer{
		tracer: otel.GetTracerProvider().Tracer(name),
	}
}

// Start starts a new span with the given name.
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// StartWithAttributes starts a new span with attributes.
func (t *Tracer) StartWithAttributes(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

// SpanFromContext returns the span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// SetSpanError records an error on the span.
func SetSpanError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// SetSpanOK sets the span status to OK.
func SetSpanOK(span trace.Span) {
	span.SetStatus(codes.Ok, "")
}

// AddSpanEvent adds an event to the span.
func AddSpanEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// WithSpanAttributes returns a span start option with attributes.
func WithSpanAttributes(attrs ...attribute.KeyValue) trace.SpanStartOption {
	return trace.WithAttributes(attrs...)
}

// Attribute helpers for common attributes.
var (
	// Component attribute key.
	AttrComponent = attribute.Key("component")
	// Operation attribute key.
	AttrOperation = attribute.Key("operation")
	// DBSystem attribute key.
	AttrDBSystem = attribute.Key("db.system")
	// DBName attribute key.
	AttrDBName = attribute.Key("db.name")
	// DBStatement attribute key.
	AttrDBStatement = attribute.Key("db.statement")
	// DBOperation attribute key.
	AttrDBOperation = attribute.Key("db.operation")
	// MessagingSystem attribute key.
	AttrMessagingSystem = attribute.Key("messaging.system")
	// MessagingDestination attribute key.
	AttrMessagingDestination = attribute.Key("messaging.destination")
	// MessagingOperation attribute key.
	AttrMessagingOperation = attribute.Key("messaging.operation")
)
