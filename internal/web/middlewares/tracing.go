package middlewares

import (
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig holds configuration for the tracing middleware.
type TracingConfig struct {
	Tracer     *observability.Tracer
	Propagator propagation.TextMapPropagator
	SkipPaths  map[string]bool
}

// Tracing returns a middleware that traces HTTP requests.
func Tracing(cfg TracingConfig) echo.MiddlewareFunc {
	if cfg.Tracer == nil {
		cfg.Tracer = observability.NewTracer("http.server")
	}
	if cfg.Propagator == nil {
		cfg.Propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip tracing for certain paths
			if cfg.SkipPaths != nil && cfg.SkipPaths[c.Path()] {
				return next(c)
			}

			req := c.Request()

			// Extract trace context from incoming request
			ctx := cfg.Propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

			// Start a new span
			spanName := c.Path()
			if spanName == "" {
				spanName = req.URL.Path
			}

			ctx, span := cfg.Tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(req.Method),
					semconv.HTTPTarget(req.URL.Path),
					semconv.HTTPScheme(req.URL.Scheme),
					semconv.NetHostName(req.Host),
					attribute.String("http.user_agent", req.UserAgent()),
					attribute.String("http.client_ip", c.RealIP()),
				),
			)
			defer span.End()

			// Update request context
			c.SetRequest(req.WithContext(ctx))

			// Process request
			err := next(c)

			// Record response attributes
			statusCode := c.Response().Status
			span.SetAttributes(
				semconv.HTTPStatusCode(statusCode),
				attribute.Int64("http.response_size", c.Response().Size),
			)

			if err != nil {
				observability.SetSpanError(span, err)
			} else if statusCode >= 400 {
				span.SetStatus(codes.Error, "HTTP error")
			} else {
				observability.SetSpanOK(span)
			}

			return err
		}
	}
}

// MetricsConfig holds configuration for the metrics middleware.
type MetricsConfig struct {
	Meter     *observability.Meter
	Namespace string
	SkipPaths map[string]bool
}

// Metrics returns a middleware that records HTTP metrics.
func Metrics(cfg MetricsConfig) echo.MiddlewareFunc {
	if cfg.Meter == nil {
		cfg.Meter = observability.NewMeter("http.server")
	}
	if cfg.Namespace == "" {
		cfg.Namespace = "http"
	}

	// Create metrics
	requestsTotal, _ := cfg.Meter.Counter(cfg.Namespace+"_requests_total", "Total number of HTTP requests")
	requestDuration, _ := cfg.Meter.Histogram(cfg.Namespace+"_request_duration_seconds", "HTTP request duration in seconds", "s")
	requestsInFlight, _ := cfg.Meter.UpDownCounter(cfg.Namespace+"_requests_in_flight", "Number of HTTP requests currently being processed")
	responseSize, _ := cfg.Meter.Histogram(cfg.Namespace+"_response_size_bytes", "HTTP response size in bytes", "By")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip metrics for certain paths
			if cfg.SkipPaths != nil && cfg.SkipPaths[c.Path()] {
				return next(c)
			}

			ctx := c.Request().Context()
			startTime := time.Now()

			// Track in-flight requests
			attrs := metric.WithAttributes(
				attribute.String("method", c.Request().Method),
				attribute.String("path", c.Path()),
			)

			requestsInFlight.Add(ctx, 1, attrs)
			defer requestsInFlight.Add(ctx, -1, attrs)

			// Process request
			err := next(c)

			// Record metrics
			duration := time.Since(startTime).Seconds()
			statusCode := c.Response().Status

			responseAttrs := metric.WithAttributes(
				attribute.String("method", c.Request().Method),
				attribute.String("path", c.Path()),
				attribute.Int("status_code", statusCode),
			)

			requestsTotal.Add(ctx, 1, responseAttrs)
			requestDuration.Record(ctx, duration, responseAttrs)
			responseSize.Record(ctx, float64(c.Response().Size), responseAttrs)

			return err
		}
	}
}
