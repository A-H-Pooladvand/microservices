package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"po/configs"
)

func Provide(lc fx.Lifecycle, config *configs.Jaeger, app *configs.App) (tracer trace.Tracer) {
	var tracerProvider *sdktrace.TracerProvider

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			res, err := resource.New(ctx,
				resource.WithAttributes(
					// the service name used to display traces in backends
					semconv.ServiceName(app.Name),
				),
			)

			if err != nil {
				return err
			}

			// It connects the OpenTelemetry Collector through local gRPC connection.
			// You may replace `localhost:4317` with your endpoint.
			conn, err := grpc.NewClient(config.Addr,
				// Note the use of insecure transport here. TLS is recommended in production.
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			if err != nil {
				return err
			}

			// Set up a trace exporter
			traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
			if err != nil {
				return err
			}

			// Register the trace exporter with a TracerProvider, using a batch
			// span processor to aggregate spans before export.
			bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
			tracerProvider = sdktrace.NewTracerProvider(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithResource(res),
				sdktrace.WithSpanProcessor(bsp),
			)

			otel.SetTracerProvider(tracerProvider)

			otel.SetTextMapPropagator(propagation.TraceContext{})

			tracer = otel.Tracer(app.Name)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return tracerProvider.Shutdown(ctx)
		},
	})

	return tracer
}
