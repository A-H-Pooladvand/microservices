package metric

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func Provider() metric.MeterProvider {
	return otel.GetMeterProvider()
}

func Meter() metric.Meter {
	return Provider().Meter("tms")
}
