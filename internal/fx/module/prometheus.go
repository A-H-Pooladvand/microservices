package module

import (
	"context"
	"github.com/a-h-pooladvand/microservices/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
	"time"
)

var Prometheus = fx.Module("prometheus",
	fx.Provide(
		NewPrometheus,
	),
	fx.Invoke(func(_ *metric.MeterProvider) {}),
)

func NewPrometheus(lc fx.Lifecycle, config *config.Config) (provider *metric.MeterProvider, err error) {
	exporter, err := prometheus.New()

	if err != nil {
		return
	}

	provider = metric.NewMeterProvider(metric.WithReader(exporter))
	provider.Meter(config.App.Name)
	otel.SetMeterProvider(provider)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()

			return exporter.Shutdown(shutdownCtx)
		},
	})

	return
}
