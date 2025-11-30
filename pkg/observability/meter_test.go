package observability_test

import (
	"context"
	"testing"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func setupTestMeter(t *testing.T) (*observability.Meter, *sdkmetric.ManualReader) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(provider)
	return observability.NewMeter("test"), reader
}

func TestMeter_Counter(t *testing.T) {
	meter, reader := setupTestMeter(t)

	counter, err := meter.Counter("test_counter", "Test counter description")
	if err != nil {
		t.Fatalf("failed to create counter: %v", err)
	}

	ctx := context.Background()
	counter.Add(ctx, 5)
	counter.Add(ctx, 3)

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	// Verify metrics were recorded
	if len(rm.ScopeMetrics) == 0 {
		t.Fatal("expected scope metrics")
	}
}

func TestMeter_Histogram(t *testing.T) {
	meter, reader := setupTestMeter(t)

	histogram, err := meter.Histogram("test_histogram", "Test histogram description", "s")
	if err != nil {
		t.Fatalf("failed to create histogram: %v", err)
	}

	ctx := context.Background()
	histogram.Record(ctx, 0.5)
	histogram.Record(ctx, 1.0)
	histogram.Record(ctx, 1.5)

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	// Verify metrics were recorded
	if len(rm.ScopeMetrics) == 0 {
		t.Fatal("expected scope metrics")
	}
}

func TestMeter_UpDownCounter(t *testing.T) {
	meter, reader := setupTestMeter(t)

	counter, err := meter.UpDownCounter("test_updown_counter", "Test up-down counter description")
	if err != nil {
		t.Fatalf("failed to create up-down counter: %v", err)
	}

	ctx := context.Background()
	counter.Add(ctx, 5)
	counter.Add(ctx, -2)
	counter.Add(ctx, 3)

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	// Verify metrics were recorded
	if len(rm.ScopeMetrics) == 0 {
		t.Fatal("expected scope metrics")
	}
}

func TestRecordDuration(t *testing.T) {
	meter, reader := setupTestMeter(t)

	histogram, err := meter.Histogram("duration_histogram", "Test duration histogram", "s")
	if err != nil {
		t.Fatalf("failed to create histogram: %v", err)
	}

	ctx := context.Background()
	startTime := time.Now().Add(-100 * time.Millisecond)
	
	observability.RecordDuration(ctx, histogram, startTime, attribute.String("operation", "test"))

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	if len(rm.ScopeMetrics) == 0 {
		t.Fatal("expected scope metrics")
	}
}

func TestNewOperationMetrics(t *testing.T) {
	meter, _ := setupTestMeter(t)

	metrics, err := observability.NewOperationMetrics(meter, "test_operation")
	if err != nil {
		t.Fatalf("failed to create operation metrics: %v", err)
	}

	if metrics.RequestsTotal == nil {
		t.Error("expected RequestsTotal counter")
	}
	if metrics.RequestDuration == nil {
		t.Error("expected RequestDuration histogram")
	}
	if metrics.RequestsInFlight == nil {
		t.Error("expected RequestsInFlight counter")
	}
	if metrics.ErrorsTotal == nil {
		t.Error("expected ErrorsTotal counter")
	}
}

func TestNewDBMetrics(t *testing.T) {
	meter, _ := setupTestMeter(t)

	metrics, err := observability.NewDBMetrics(meter, "test_db")
	if err != nil {
		t.Fatalf("failed to create DB metrics: %v", err)
	}

	if metrics.QueriesTotal == nil {
		t.Error("expected QueriesTotal counter")
	}
	if metrics.QueryDuration == nil {
		t.Error("expected QueryDuration histogram")
	}
	if metrics.ConnectionsActive == nil {
		t.Error("expected ConnectionsActive counter")
	}
	if metrics.ErrorsTotal == nil {
		t.Error("expected ErrorsTotal counter")
	}
}

func TestNewCacheMetrics(t *testing.T) {
	meter, _ := setupTestMeter(t)

	metrics, err := observability.NewCacheMetrics(meter, "test_cache")
	if err != nil {
		t.Fatalf("failed to create cache metrics: %v", err)
	}

	if metrics.HitsTotal == nil {
		t.Error("expected HitsTotal counter")
	}
	if metrics.MissesTotal == nil {
		t.Error("expected MissesTotal counter")
	}
	if metrics.OperationTime == nil {
		t.Error("expected OperationTime histogram")
	}
	if metrics.ErrorsTotal == nil {
		t.Error("expected ErrorsTotal counter")
	}
}

func TestNewMessagingMetrics(t *testing.T) {
	meter, _ := setupTestMeter(t)

	metrics, err := observability.NewMessagingMetrics(meter, "test_messaging")
	if err != nil {
		t.Fatalf("failed to create messaging metrics: %v", err)
	}

	if metrics.MessagesPublished == nil {
		t.Error("expected MessagesPublished counter")
	}
	if metrics.MessagesConsumed == nil {
		t.Error("expected MessagesConsumed counter")
	}
	if metrics.ProcessingTime == nil {
		t.Error("expected ProcessingTime histogram")
	}
	if metrics.ErrorsTotal == nil {
		t.Error("expected ErrorsTotal counter")
	}
}

func TestIncrementCounter(t *testing.T) {
	meter, reader := setupTestMeter(t)

	counter, err := meter.Counter("increment_test_counter", "Test counter")
	if err != nil {
		t.Fatalf("failed to create counter: %v", err)
	}

	ctx := context.Background()
	
	// Use helper function
	observability.IncrementCounter(ctx, counter, attribute.String("label", "test"))

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	if len(rm.ScopeMetrics) == 0 {
		t.Fatal("expected scope metrics")
	}
}
