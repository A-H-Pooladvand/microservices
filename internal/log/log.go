package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Error logs an error message with optional fields.
func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

// ErrorCtx logs an error message with trace context.
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Error(msg, appendTraceFields(ctx, fields)...)
}

// Info logs an info message with optional fields.
func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

// InfoCtx logs an info message with trace context.
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Info(msg, appendTraceFields(ctx, fields)...)
}

// Panic logs a panic message with optional fields.
func Panic(msg string, fields ...zap.Field) {
	zap.L().Panic(msg, fields...)
}

// PanicCtx logs a panic message with trace context.
func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Panic(msg, appendTraceFields(ctx, fields)...)
}

// Warn logs a warning message with optional fields.
func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

// WarnCtx logs a warning message with trace context.
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Warn(msg, appendTraceFields(ctx, fields)...)
}

// Fatal logs a fatal message with optional fields.
func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}

// FatalCtx logs a fatal message with trace context.
func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, appendTraceFields(ctx, fields)...)
}

// Debug logs a debug message with optional fields.
func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

// DebugCtx logs a debug message with trace context.
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	zap.L().Debug(msg, appendTraceFields(ctx, fields)...)
}

// Logger returns a new logger with trace context.
func Logger(ctx context.Context) *zap.Logger {
	return zap.L().With(traceFields(ctx)...)
}

// appendTraceFields appends trace fields to the given fields.
func appendTraceFields(ctx context.Context, fields []zap.Field) []zap.Field {
	return append(fields, traceFields(ctx)...)
}

// traceFields extracts trace fields from context.
func traceFields(ctx context.Context) []zap.Field {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return nil
	}

	return []zap.Field{
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("span_id", span.SpanContext().SpanID().String()),
	}
}

// TraceID returns a zap field with the trace ID.
func TraceID(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}

// SpanID returns a zap field with the span ID.
func SpanID(spanID string) zap.Field {
	return zap.String("span_id", spanID)
}
