package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Meter wraps OpenTelemetry meter functionality.
type Meter struct {
	meter metric.Meter
}

// NewMeter creates a new Meter instance.
func NewMeter(name string) *Meter {
	return &Meter{
		meter: otel.GetMeterProvider().Meter(name),
	}
}

// Counter creates a new Int64Counter.
func (m *Meter) Counter(name, description string) (metric.Int64Counter, error) {
	return m.meter.Int64Counter(name, metric.WithDescription(description))
}

// Histogram creates a new Float64Histogram.
func (m *Meter) Histogram(name, description, unit string) (metric.Float64Histogram, error) {
	return m.meter.Float64Histogram(name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// Gauge creates a new Float64Gauge.
func (m *Meter) Gauge(name, description, unit string) (metric.Float64Gauge, error) {
	return m.meter.Float64Gauge(name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// UpDownCounter creates a new Int64UpDownCounter.
func (m *Meter) UpDownCounter(name, description string) (metric.Int64UpDownCounter, error) {
	return m.meter.Int64UpDownCounter(name, metric.WithDescription(description))
}

// RecordDuration is a helper to record duration metrics.
func RecordDuration(ctx context.Context, histogram metric.Float64Histogram, startTime time.Time, attrs ...attribute.KeyValue) {
	duration := time.Since(startTime).Seconds()
	histogram.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// IncrementCounter is a helper to increment a counter.
func IncrementCounter(ctx context.Context, counter metric.Int64Counter, attrs ...attribute.KeyValue) {
	counter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// Metrics helper types for common metrics patterns.
type (
	// OperationMetrics holds common metrics for operations.
	OperationMetrics struct {
		RequestsTotal    metric.Int64Counter
		RequestDuration  metric.Float64Histogram
		RequestsInFlight metric.Int64UpDownCounter
		ErrorsTotal      metric.Int64Counter
	}

	// DBMetrics holds common metrics for database operations.
	DBMetrics struct {
		QueriesTotal      metric.Int64Counter
		QueryDuration     metric.Float64Histogram
		ConnectionsActive metric.Int64UpDownCounter
		ErrorsTotal       metric.Int64Counter
	}

	// CacheMetrics holds common metrics for cache operations.
	CacheMetrics struct {
		HitsTotal     metric.Int64Counter
		MissesTotal   metric.Int64Counter
		OperationTime metric.Float64Histogram
		ErrorsTotal   metric.Int64Counter
	}

	// MessagingMetrics holds common metrics for messaging operations.
	MessagingMetrics struct {
		MessagesPublished metric.Int64Counter
		MessagesConsumed  metric.Int64Counter
		ProcessingTime    metric.Float64Histogram
		ErrorsTotal       metric.Int64Counter
	}
)

// NewOperationMetrics creates a new OperationMetrics instance.
func NewOperationMetrics(m *Meter, prefix string) (*OperationMetrics, error) {
	requestsTotal, err := m.Counter(prefix+"_requests_total", "Total number of requests")
	if err != nil {
		return nil, err
	}

	requestDuration, err := m.Histogram(prefix+"_request_duration_seconds", "Request duration in seconds", "s")
	if err != nil {
		return nil, err
	}

	requestsInFlight, err := m.UpDownCounter(prefix+"_requests_in_flight", "Number of requests currently being processed")
	if err != nil {
		return nil, err
	}

	errorsTotal, err := m.Counter(prefix+"_errors_total", "Total number of errors")
	if err != nil {
		return nil, err
	}

	return &OperationMetrics{
		RequestsTotal:    requestsTotal,
		RequestDuration:  requestDuration,
		RequestsInFlight: requestsInFlight,
		ErrorsTotal:      errorsTotal,
	}, nil
}

// NewDBMetrics creates a new DBMetrics instance.
func NewDBMetrics(m *Meter, prefix string) (*DBMetrics, error) {
	queriesTotal, err := m.Counter(prefix+"_queries_total", "Total number of database queries")
	if err != nil {
		return nil, err
	}

	queryDuration, err := m.Histogram(prefix+"_query_duration_seconds", "Query duration in seconds", "s")
	if err != nil {
		return nil, err
	}

	connectionsActive, err := m.UpDownCounter(prefix+"_connections_active", "Number of active database connections")
	if err != nil {
		return nil, err
	}

	errorsTotal, err := m.Counter(prefix+"_errors_total", "Total number of database errors")
	if err != nil {
		return nil, err
	}

	return &DBMetrics{
		QueriesTotal:      queriesTotal,
		QueryDuration:     queryDuration,
		ConnectionsActive: connectionsActive,
		ErrorsTotal:       errorsTotal,
	}, nil
}

// NewCacheMetrics creates a new CacheMetrics instance.
func NewCacheMetrics(m *Meter, prefix string) (*CacheMetrics, error) {
	hitsTotal, err := m.Counter(prefix+"_hits_total", "Total number of cache hits")
	if err != nil {
		return nil, err
	}

	missesTotal, err := m.Counter(prefix+"_misses_total", "Total number of cache misses")
	if err != nil {
		return nil, err
	}

	operationTime, err := m.Histogram(prefix+"_operation_duration_seconds", "Cache operation duration in seconds", "s")
	if err != nil {
		return nil, err
	}

	errorsTotal, err := m.Counter(prefix+"_errors_total", "Total number of cache errors")
	if err != nil {
		return nil, err
	}

	return &CacheMetrics{
		HitsTotal:     hitsTotal,
		MissesTotal:   missesTotal,
		OperationTime: operationTime,
		ErrorsTotal:   errorsTotal,
	}, nil
}

// NewMessagingMetrics creates a new MessagingMetrics instance.
func NewMessagingMetrics(m *Meter, prefix string) (*MessagingMetrics, error) {
	messagesPublished, err := m.Counter(prefix+"_messages_published_total", "Total number of messages published")
	if err != nil {
		return nil, err
	}

	messagesConsumed, err := m.Counter(prefix+"_messages_consumed_total", "Total number of messages consumed")
	if err != nil {
		return nil, err
	}

	processingTime, err := m.Histogram(prefix+"_processing_duration_seconds", "Message processing duration in seconds", "s")
	if err != nil {
		return nil, err
	}

	errorsTotal, err := m.Counter(prefix+"_errors_total", "Total number of messaging errors")
	if err != nil {
		return nil, err
	}

	return &MessagingMetrics{
		MessagesPublished: messagesPublished,
		MessagesConsumed:  messagesConsumed,
		ProcessingTime:    processingTime,
		ErrorsTotal:       errorsTotal,
	}, nil
}
