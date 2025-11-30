package grpc_test

import (
	"context"
	"errors"
	"testing"

	grpcpkg "github.com/a-h-pooladvand/microservices/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewInterceptorConfig(t *testing.T) {
	cfg, err := grpcpkg.NewInterceptorConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.Tracer == nil {
		t.Error("expected tracer to be set")
	}

	if cfg.Meter == nil {
		t.Error("expected meter to be set")
	}

	if cfg.Metrics == nil {
		t.Error("expected metrics to be set")
	}
}

func TestUnaryMetricsInterceptor(t *testing.T) {
	cfg, err := grpcpkg.NewInterceptorConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	interceptor := grpcpkg.UnaryMetricsInterceptor(cfg)

	// Create mock handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp != "response" {
		t.Errorf("expected 'response', got %v", resp)
	}
}

func TestUnaryMetricsInterceptor_Error(t *testing.T) {
	cfg, err := grpcpkg.NewInterceptorConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	interceptor := grpcpkg.UnaryMetricsInterceptor(cfg)

	expectedErr := status.Error(codes.NotFound, "not found")

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedErr
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	_, err = interceptor(context.Background(), "request", info, handler)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestUnaryLoggingInterceptor(t *testing.T) {
	interceptor := grpcpkg.UnaryLoggingInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp != "response" {
		t.Errorf("expected 'response', got %v", resp)
	}
}

func TestUnaryRecoverInterceptor(t *testing.T) {
	interceptor := grpcpkg.UnaryRecoverInterceptor()

	// Handler that panics
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	// Should recover from panic and return error
	resp, err := interceptor(context.Background(), "request", info, handler)
	if err == nil {
		t.Error("expected error after panic recovery, got nil")
	}

	if resp != nil {
		t.Errorf("expected nil response after panic, got %v", resp)
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}

	if st.Code() != codes.Internal {
		t.Errorf("expected Internal code, got %v", st.Code())
	}
}

func TestUnaryRecoverInterceptor_NoPanic(t *testing.T) {
	interceptor := grpcpkg.UnaryRecoverInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	resp, err := interceptor(context.Background(), "request", info, handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp != "response" {
		t.Errorf("expected 'response', got %v", resp)
	}
}
