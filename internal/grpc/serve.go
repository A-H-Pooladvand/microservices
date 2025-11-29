package grpc

import (
	"context"
	"net"
	"time"

	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/a-h-pooladvand/microservices/routes"
	"github.com/fatih/color"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds gRPC server configuration.
type ServerConfig struct {
	Address               string
	MaxRecvMsgSize        int
	MaxSendMsgSize        int
	MaxConcurrentStreams  uint32
	ConnectionTimeout     time.Duration
	GracefulShutdownTime  time.Duration
	EnableReflection      bool
}

// DefaultServerConfig returns default server configuration.
func DefaultServerConfig(addr string) ServerConfig {
	return ServerConfig{
		Address:              addr,
		MaxRecvMsgSize:       4 * 1024 * 1024,  // 4MB
		MaxSendMsgSize:       4 * 1024 * 1024,  // 4MB
		MaxConcurrentStreams: 100,
		ConnectionTimeout:    30 * time.Second,
		GracefulShutdownTime: 30 * time.Second,
		EnableReflection:     true,
	}
}

// Invoke creates and starts the gRPC server.
func Invoke(lc fx.Lifecycle, cfg *config.Config, h *handler.Grpc) *grpc.Server {
	serverCfg := DefaultServerConfig(cfg.GRPC.Addr)

	// Create interceptor config
	interceptorCfg, err := NewInterceptorConfig()
	if err != nil {
		log.Error("failed to create interceptor config", zap.Error(err))
	}

	// Build server options
	opts := buildServerOptions(serverCfg, interceptorCfg)

	server := grpc.NewServer(opts...)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Skip gRPC server if address is not configured
			if serverCfg.Address == "" {
				log.Info("gRPC server disabled (no address configured)")
				return nil
			}

			lis, err := net.Listen("tcp", ":"+serverCfg.Address)
			if err != nil {
				return err
			}

			// Register reflection service for debugging
			if serverCfg.EnableReflection {
				reflection.Register(server)
			}

			routes.RegisterGrpcRoutes(server, h)

			go func() {
				color.Green("gRPC server started on [::]:" + serverCfg.Address)

				if err = server.Serve(lis); err != nil {
					log.Error("gRPC serve error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Use GracefulStop with timeout for production
			done := make(chan struct{})
			go func() {
				server.GracefulStop()
				close(done)
			}()

			select {
			case <-done:
				log.Info("gRPC server stopped gracefully")
			case <-time.After(serverCfg.GracefulShutdownTime):
				log.Warn("gRPC server graceful shutdown timed out, forcing stop")
				server.Stop()
			}

			return nil
		},
	})

	return server
}

func buildServerOptions(cfg ServerConfig, interceptorCfg *InterceptorConfig) []grpc.ServerOption {
	// Unary interceptors (order matters: recovery -> logging -> metrics -> tracing)
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		UnaryRecoverInterceptor(),
		UnaryLoggingInterceptor(),
	}
	if interceptorCfg != nil {
		unaryInterceptors = append(unaryInterceptors, UnaryMetricsInterceptor(interceptorCfg))
	}

	// Stream interceptors
	streamInterceptors := []grpc.StreamServerInterceptor{
		StreamRecoverInterceptor(),
		StreamLoggingInterceptor(),
	}
	if interceptorCfg != nil {
		streamInterceptors = append(streamInterceptors, StreamMetricsInterceptor(interceptorCfg))
	}

	return []grpc.ServerOption{
		// OpenTelemetry tracing
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// Server limits
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams),
		grpc.ConnectionTimeout(cfg.ConnectionTimeout),

		// Interceptors
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
}
