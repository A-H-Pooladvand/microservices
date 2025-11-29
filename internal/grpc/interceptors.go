package grpc

import (
	"context"
	"time"

	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorConfig holds configuration for gRPC interceptors.
type InterceptorConfig struct {
	Tracer  *observability.Tracer
	Meter   *observability.Meter
	Metrics *GRPCMetrics
}

// GRPCMetrics holds metrics for gRPC operations.
type GRPCMetrics struct {
	RequestsTotal    metric.Int64Counter
	RequestDuration  metric.Float64Histogram
	RequestsInFlight metric.Int64UpDownCounter
	ErrorsTotal      metric.Int64Counter
}

// NewInterceptorConfig creates a new InterceptorConfig with default values.
func NewInterceptorConfig() (*InterceptorConfig, error) {
	tracer := observability.NewTracer("grpc.server")
	meter := observability.NewMeter("grpc.server")

	requestsTotal, err := meter.Counter("grpc_server_requests_total", "Total number of gRPC requests")
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Histogram("grpc_server_request_duration_seconds", "gRPC request duration in seconds", "s")
	if err != nil {
		return nil, err
	}

	requestsInFlight, err := meter.UpDownCounter("grpc_server_requests_in_flight", "Number of gRPC requests currently being processed")
	if err != nil {
		return nil, err
	}

	errorsTotal, err := meter.Counter("grpc_server_errors_total", "Total number of gRPC errors")
	if err != nil {
		return nil, err
	}

	return &InterceptorConfig{
		Tracer: tracer,
		Meter:  meter,
		Metrics: &GRPCMetrics{
			RequestsTotal:    requestsTotal,
			RequestDuration:  requestDuration,
			RequestsInFlight: requestsInFlight,
			ErrorsTotal:      errorsTotal,
		},
	}, nil
}

// UnaryMetricsInterceptor returns a unary server interceptor that records metrics.
func UnaryMetricsInterceptor(cfg *InterceptorConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		attrs := metric.WithAttributes(
			attribute.String("grpc.method", info.FullMethod),
		)

		cfg.Metrics.RequestsInFlight.Add(ctx, 1, attrs)
		defer cfg.Metrics.RequestsInFlight.Add(ctx, -1, attrs)

		resp, err := handler(ctx, req)

		duration := time.Since(startTime).Seconds()
		code := status.Code(err)

		responseAttrs := metric.WithAttributes(
			attribute.String("grpc.method", info.FullMethod),
			attribute.String("grpc.code", code.String()),
		)

		cfg.Metrics.RequestsTotal.Add(ctx, 1, responseAttrs)
		cfg.Metrics.RequestDuration.Record(ctx, duration, responseAttrs)

		if err != nil {
			cfg.Metrics.ErrorsTotal.Add(ctx, 1, responseAttrs)
		}

		return resp, err
	}
}

// UnaryLoggingInterceptor returns a unary server interceptor that logs requests.
func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		code := status.Code(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("code", code.String()),
			zap.Duration("duration", duration),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			if code == codes.Internal || code == codes.Unknown {
				log.ErrorCtx(ctx, "gRPC request failed", fields...)
			} else {
				log.WarnCtx(ctx, "gRPC request error", fields...)
			}
		} else {
			log.InfoCtx(ctx, "gRPC request completed", fields...)
		}

		return resp, err
	}
}

// UnaryRecoverInterceptor returns a unary server interceptor that recovers from panics.
func UnaryRecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorCtx(ctx, "gRPC panic recovered",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// StreamMetricsInterceptor returns a stream server interceptor that records metrics.
func StreamMetricsInterceptor(cfg *InterceptorConfig) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		ctx := ss.Context()

		attrs := metric.WithAttributes(
			attribute.String("grpc.method", info.FullMethod),
			attribute.Bool("grpc.stream", true),
		)

		cfg.Metrics.RequestsInFlight.Add(ctx, 1, attrs)
		defer cfg.Metrics.RequestsInFlight.Add(ctx, -1, attrs)

		err := handler(srv, ss)

		duration := time.Since(startTime).Seconds()
		code := status.Code(err)

		responseAttrs := metric.WithAttributes(
			attribute.String("grpc.method", info.FullMethod),
			attribute.String("grpc.code", code.String()),
			attribute.Bool("grpc.stream", true),
		)

		cfg.Metrics.RequestsTotal.Add(ctx, 1, responseAttrs)
		cfg.Metrics.RequestDuration.Record(ctx, duration, responseAttrs)

		if err != nil {
			cfg.Metrics.ErrorsTotal.Add(ctx, 1, responseAttrs)
		}

		return err
	}
}

// StreamLoggingInterceptor returns a stream server interceptor that logs requests.
func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		ctx := ss.Context()

		err := handler(srv, ss)

		duration := time.Since(startTime)
		code := status.Code(err)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("code", code.String()),
			zap.Duration("duration", duration),
			zap.Bool("stream", true),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			log.ErrorCtx(ctx, "gRPC stream failed", fields...)
		} else {
			log.InfoCtx(ctx, "gRPC stream completed", fields...)
		}

		return err
	}
}

// StreamRecoverInterceptor returns a stream server interceptor that recovers from panics.
func StreamRecoverInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorCtx(ss.Context(), "gRPC stream panic recovered",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, ss)
	}
}
