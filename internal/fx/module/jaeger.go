package module

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

var Jaeger = fx.Module(
	"jaeger",
	fx.Provide(
		NewJaeger,
	),
	// Instantiating constructor => DO NOT DELETE
	fx.Invoke(func(_ *sdktrace.TracerProvider) {}),
)

func NewJaeger(lc fx.Lifecycle, config *config.Config) (*sdktrace.TracerProvider, error) {
	tp, err := initTracer(config)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()

			// Shutdown tracer
			if err = tp.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("error shutting down tracer provider: %w", err)
			}

			// Close gRPC connection
			//if conn != nil {
			//	if err := conn.Close(); err != nil {
			//		return fmt.Errorf("error closing gRPC connection: %w", err)
			//	}
			//}

			return nil
		},
	})

	return tp, nil
}

func initTracer(config *config.Config) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(config.App.Name),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure sampling
	sampler := sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(config.Jaeger.SamplingRatio),
	)

	_, conn, err := newGRPCExporter(ctx, config)

	if err != nil {
		return nil, err
	}

	// Add up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Here we set global tracer
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	return tp, nil
}

func newGRPCExporter(ctx context.Context, cfg *config.Config) (sdktrace.SpanExporter, *grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if cfg.Jaeger.UseTLS {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// If credentials are provided, use per-RPC authentication
	if cfg.Jaeger.Username != "" && cfg.Jaeger.Password != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(basicAuth{
			username: cfg.Jaeger.Username,
			password: cfg.Jaeger.Password,
		}))
	}

	conn, err := grpc.NewClient(cfg.Jaeger.Addr, opts...)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to OTLP collector: %w", err)
	}

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("failed to initialize OTLP exporter: %w", err)
	}

	return exp, conn, nil
}

type basicAuth struct {
	username string
	password string
}

func (a basicAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + basicAuthHeader(a.username, a.password),
	}, nil
}

func (a basicAuth) RequireTransportSecurity() bool {
	return true
}

func basicAuthHeader(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
