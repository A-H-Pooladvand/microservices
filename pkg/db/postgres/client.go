package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Client wraps gorm.DB with observability support.
type Client struct {
	*gorm.DB
	tracer  *observability.Tracer
	metrics *observability.DBMetrics
	config  Config
}

// New creates a new postgres client and returns a pointer to the Client instance.
func New(config Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	dialect := postgres.Open(config.DSN())

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	connection, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	// Configure connection pool
	connection.SetMaxOpenConns(config.MaxOpenConns)
	connection.SetMaxIdleConns(config.MaxIdleConns)
	connection.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Wait for connection
	for {
		if err = connection.Ping(); err == nil {
			break
		}

		select {
		case <-time.After(500 * time.Millisecond):
			continue
		case <-ctx.Done():
			_ = connection.Close()
			return nil, errors.New("unable to connect to postgres client, context deadline exceeded")
		}
	}

	// Initialize observability
	tracer := observability.NewTracer("postgres")
	meter := observability.NewMeter("postgres")
	metrics, err := observability.NewDBMetrics(meter, "postgres")
	if err != nil {
		return nil, fmt.Errorf("failed to create db metrics: %w", err)
	}

	client := &Client{
		DB:      db,
		tracer:  tracer,
		metrics: metrics,
		config:  config,
	}

	// Register GORM callbacks for tracing
	client.registerCallbacks()

	return client, nil
}

// registerCallbacks registers GORM callbacks for tracing and metrics.
func (c *Client) registerCallbacks() {
	// Create callback
	_ = c.DB.Callback().Create().Before("gorm:create").Register("observability:before_create", c.beforeCallback("CREATE"))
	_ = c.DB.Callback().Create().After("gorm:create").Register("observability:after_create", c.afterCallback("CREATE"))

	// Query callback
	_ = c.DB.Callback().Query().Before("gorm:query").Register("observability:before_query", c.beforeCallback("SELECT"))
	_ = c.DB.Callback().Query().After("gorm:query").Register("observability:after_query", c.afterCallback("SELECT"))

	// Update callback
	_ = c.DB.Callback().Update().Before("gorm:update").Register("observability:before_update", c.beforeCallback("UPDATE"))
	_ = c.DB.Callback().Update().After("gorm:update").Register("observability:after_update", c.afterCallback("UPDATE"))

	// Delete callback
	_ = c.DB.Callback().Delete().Before("gorm:delete").Register("observability:before_delete", c.beforeCallback("DELETE"))
	_ = c.DB.Callback().Delete().After("gorm:delete").Register("observability:after_delete", c.afterCallback("DELETE"))

	// Raw callback
	_ = c.DB.Callback().Raw().Before("gorm:raw").Register("observability:before_raw", c.beforeCallback("RAW"))
	_ = c.DB.Callback().Raw().After("gorm:raw").Register("observability:after_raw", c.afterCallback("RAW"))

	// Row callback
	_ = c.DB.Callback().Row().Before("gorm:row").Register("observability:before_row", c.beforeCallback("ROW"))
	_ = c.DB.Callback().Row().After("gorm:row").Register("observability:after_row", c.afterCallback("ROW"))
}

type spanKey struct{}

type spanData struct {
	startTime time.Time
}

func (c *Client) beforeCallback(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context == nil {
			return
		}

		ctx, _ := c.tracer.StartWithAttributes(db.Statement.Context, "postgres."+operation,
			observability.AttrDBSystem.String("postgresql"),
			observability.AttrDBName.String(c.config.DB),
			observability.AttrDBOperation.String(operation),
		)

		// Store start time and update context (span is already in ctx)
		db.Statement.Context = context.WithValue(ctx, spanKey{}, spanData{
			startTime: time.Now(),
		})

		c.metrics.ConnectionsActive.Add(ctx, 1)
	}
}

func (c *Client) afterCallback(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context == nil {
			return
		}

		ctx := db.Statement.Context
		span := observability.SpanFromContext(ctx)

		// Record query
		if db.Statement.SQL.String() != "" {
			span.SetAttributes(observability.AttrDBStatement.String(db.Statement.SQL.String()))
		}

		// Record metrics
		attrs := metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("table", db.Statement.Table),
		)

		if data, ok := ctx.Value(spanKey{}).(spanData); ok {
			duration := time.Since(data.startTime).Seconds()
			c.metrics.QueryDuration.Record(ctx, duration, attrs)
		}

		c.metrics.QueriesTotal.Add(ctx, 1, attrs)

		c.metrics.ConnectionsActive.Add(ctx, -1)

		// Record error if any
		if db.Error != nil {
			observability.SetSpanError(span, db.Error)
			c.metrics.ErrorsTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("operation", operation),
				attribute.String("error_type", db.Error.Error()),
			))
		} else {
			observability.SetSpanOK(span)
		}

		span.End()
	}
}

// WithContext returns a new client with the given context.
func (c *Client) WithContext(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

// Close closes the database connection.
func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping checks the database connection.
func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns database statistics.
func (c *Client) Stats() (Stats, error) {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return Stats{}, err
	}

	stats := sqlDB.Stats()
	return Stats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUse:              stats.InUse,
		Idle:               stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}, nil
}

// Stats represents database connection pool statistics.
type Stats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}
