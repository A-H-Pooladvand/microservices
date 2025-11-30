package observability_test

import (
	"context"
	"testing"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracer() (*observability.Tracer, *tracetest.SpanRecorder) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	otel.SetTracerProvider(tp)
	return observability.NewTracer("test"), recorder
}

func TestTracer_Start(t *testing.T) {
	tracer, recorder := setupTestTracer()

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	if span == nil {
		t.Fatal("expected span to be created")
	}

	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name() != "test-span" {
		t.Errorf("expected span name 'test-span', got '%s'", spans[0].Name())
	}
}

func TestTracer_StartWithAttributes(t *testing.T) {
	tracer, recorder := setupTestTracer()

	ctx := context.Background()
	ctx, span := tracer.StartWithAttributes(ctx, "test-span-attrs",
		attribute.String("key1", "value1"),
		attribute.Int("key2", 42),
	)
	defer span.End()

	if span == nil {
		t.Fatal("expected span to be created")
	}

	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	attrs := spans[0].Attributes()
	if len(attrs) < 2 {
		t.Errorf("expected at least 2 attributes, got %d", len(attrs))
	}
}

func TestSpanFromContext(t *testing.T) {
	tracer, _ := setupTestTracer()

	ctx := context.Background()
	ctx, expectedSpan := tracer.Start(ctx, "test-span")
	defer expectedSpan.End()

	actualSpan := observability.SpanFromContext(ctx)

	if actualSpan.SpanContext().SpanID() != expectedSpan.SpanContext().SpanID() {
		t.Error("expected same span from context")
	}
}

func TestSetSpanError(t *testing.T) {
	tracer, recorder := setupTestTracer()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "error-span")

	testErr := context.DeadlineExceeded
	observability.SetSpanError(span, testErr)
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	events := spans[0].Events()
	hasException := false
	for _, event := range events {
		if event.Name == "exception" {
			hasException = true
			break
		}
	}

	if !hasException {
		t.Error("expected exception event on span")
	}
}

func TestSetSpanOK(t *testing.T) {
	tracer, recorder := setupTestTracer()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "ok-span")

	observability.SetSpanOK(span)
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
}

func TestAddSpanEvent(t *testing.T) {
	tracer, recorder := setupTestTracer()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "event-span")

	observability.AddSpanEvent(span, "test-event", attribute.String("detail", "value"))
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	events := spans[0].Events()
	hasEvent := false
	for _, event := range events {
		if event.Name == "test-event" {
			hasEvent = true
			break
		}
	}

	if !hasEvent {
		t.Error("expected test-event on span")
	}
}
