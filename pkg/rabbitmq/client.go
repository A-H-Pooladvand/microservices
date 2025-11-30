package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

// Client wraps AMQP connection with observability support.
type Client struct {
	conn       *amqp.Connection
	config     Config
	tracer     *observability.Tracer
	metrics    *observability.MessagingMetrics
	logger     *zap.Logger
	propagator propagation.TextMapPropagator
	mu         sync.RWMutex
	closed     bool
}

// New creates a new RabbitMQ client with observability.
func New(cfg Config, logger *zap.Logger) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	conn, err := amqp.DialConfig(cfg.DSN(), amqp.Config{
		Heartbeat: cfg.Heartbeat,
		Properties: amqp.Table{
			"connection_name": cfg.ConnectionName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	// Initialize observability
	tracer := observability.NewTracer("rabbitmq")
	meter := observability.NewMeter("rabbitmq")
	metrics, err := observability.NewMessagingMetrics(meter, "rabbitmq")
	if err != nil {
		return nil, fmt.Errorf("failed to create messaging metrics: %w", err)
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Client{
		conn:       conn,
		config:     cfg,
		tracer:     tracer,
		metrics:    metrics,
		logger:     logger,
		propagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	}, nil
}

// Channel creates a new channel with tracing.
func (c *Client) Channel(ctx context.Context) (*Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, ErrConnectionFailed
	}

	ctx, span := c.tracer.StartWithAttributes(ctx, "rabbitmq.Channel",
		attribute.String("messaging.system", "rabbitmq"),
	)
	defer span.End()

	ch, err := c.conn.Channel()
	if err != nil {
		observability.SetSpanError(span, err)
		return nil, fmt.Errorf("create channel: %w", err)
	}

	if err := ch.Qos(c.config.PrefetchCount, 0, c.config.PrefetchGlobal); err != nil {
		_ = ch.Close()
		observability.SetSpanError(span, err)
		return nil, fmt.Errorf("set qos: %w", err)
	}

	observability.SetSpanOK(span)
	return &Channel{
		Channel:    ch,
		tracer:     c.tracer,
		metrics:    c.metrics,
		logger:     c.logger,
		propagator: c.propagator,
	}, nil
}

// Close closes the connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return c.conn.Close()
}

// IsClosed returns true if the connection is closed.
func (c *Client) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed || c.conn.IsClosed()
}

// Channel wraps AMQP channel with observability support.
type Channel struct {
	*amqp.Channel
	tracer     *observability.Tracer
	metrics    *observability.MessagingMetrics
	logger     *zap.Logger
	propagator propagation.TextMapPropagator
}

// Publish publishes a message with tracing.
func (ch *Channel) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ctx, span := ch.tracer.StartWithAttributes(ctx, "rabbitmq.Publish",
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.destination", exchange),
		attribute.String("messaging.destination_kind", "exchange"),
		attribute.String("messaging.rabbitmq.routing_key", key),
	)
	defer span.End()

	startTime := time.Now()

	// Inject trace context into headers
	if msg.Headers == nil {
		msg.Headers = amqp.Table{}
	}
	ch.propagator.Inject(ctx, AMQPHeaderCarrier(msg.Headers))

	err := ch.Channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)

	ch.recordPublishMetrics(ctx, exchange, key, startTime, err)

	if err != nil {
		observability.SetSpanError(span, err)
		ch.logger.Error("failed to publish message",
			zap.String("exchange", exchange),
			zap.String("routing_key", key),
			zap.Error(err),
		)
		return fmt.Errorf("publish: %w", err)
	}

	observability.SetSpanOK(span)
	return nil
}

// Consume starts consuming messages with tracing.
func (ch *Channel) Consume(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	ctx, span := ch.tracer.StartWithAttributes(ctx, "rabbitmq.Consume",
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.destination", queue),
		attribute.String("messaging.destination_kind", "queue"),
		attribute.String("messaging.consumer_id", consumer),
	)
	defer span.End()

	deliveries, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	if err != nil {
		observability.SetSpanError(span, err)
		return nil, fmt.Errorf("consume: %w", err)
	}

	// Wrap deliveries with tracing
	tracedDeliveries := make(chan amqp.Delivery)
	go func() {
		defer close(tracedDeliveries)
		for d := range deliveries {
			ch.metrics.MessagesConsumed.Add(ctx, 1,
				metric.WithAttributes(attribute.String("queue", queue)),
			)
			tracedDeliveries <- d
		}
	}()

	observability.SetSpanOK(span)
	return tracedDeliveries, nil
}

// ExchangeDeclare declares an exchange with tracing.
func (ch *Channel) ExchangeDeclare(ctx context.Context, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	ctx, span := ch.tracer.StartWithAttributes(ctx, "rabbitmq.ExchangeDeclare",
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.rabbitmq.exchange", name),
		attribute.String("messaging.rabbitmq.exchange_type", kind),
	)
	defer span.End()

	err := ch.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		observability.SetSpanError(span, err)
		return fmt.Errorf("exchange declare: %w", err)
	}

	observability.SetSpanOK(span)
	return nil
}

// QueueDeclare declares a queue with tracing.
func (ch *Channel) QueueDeclare(ctx context.Context, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	ctx, span := ch.tracer.StartWithAttributes(ctx, "rabbitmq.QueueDeclare",
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.rabbitmq.queue", name),
	)
	defer span.End()

	q, err := ch.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	if err != nil {
		observability.SetSpanError(span, err)
		return amqp.Queue{}, fmt.Errorf("queue declare: %w", err)
	}

	observability.SetSpanOK(span)
	return q, nil
}

// QueueBind binds a queue to an exchange with tracing.
func (ch *Channel) QueueBind(ctx context.Context, name, key, exchange string, noWait bool, args amqp.Table) error {
	ctx, span := ch.tracer.StartWithAttributes(ctx, "rabbitmq.QueueBind",
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.rabbitmq.queue", name),
		attribute.String("messaging.rabbitmq.exchange", exchange),
		attribute.String("messaging.rabbitmq.routing_key", key),
	)
	defer span.End()

	err := ch.Channel.QueueBind(name, key, exchange, noWait, args)
	if err != nil {
		observability.SetSpanError(span, err)
		return fmt.Errorf("queue bind: %w", err)
	}

	observability.SetSpanOK(span)
	return nil
}

func (ch *Channel) recordPublishMetrics(ctx context.Context, exchange, key string, startTime time.Time, err error) {
	duration := time.Since(startTime).Seconds()
	attrs := metric.WithAttributes(
		attribute.String("exchange", exchange),
		attribute.String("routing_key", key),
	)

	ch.metrics.ProcessingTime.Record(ctx, duration, attrs)

	if err != nil {
		ch.metrics.ErrorsTotal.Add(ctx, 1, attrs)
	} else {
		ch.metrics.MessagesPublished.Add(ctx, 1, attrs)
	}
}

// AMQPHeaderCarrier implements propagation.TextMapCarrier for AMQP headers.
type AMQPHeaderCarrier amqp.Table

// Get returns the value for a key.
func (c AMQPHeaderCarrier) Get(key string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Set sets a key-value pair.
func (c AMQPHeaderCarrier) Set(key, val string) {
	c[key] = val
}

// Keys returns all keys.
func (c AMQPHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// ExtractTraceContext extracts trace context from delivery headers.
func ExtractTraceContext(ctx context.Context, d amqp.Delivery) context.Context {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	return propagator.Extract(ctx, AMQPHeaderCarrier(d.Headers))
}
