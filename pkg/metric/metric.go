package metric

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

const exportInterval = 15 * time.Second

func New(ctx context.Context, service string) (func(context.Context) error, error) {
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
	)

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(exportInterval))),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return func(ctx context.Context) error {
		err = mp.ForceFlush(ctx)
		if err != nil {
			err = mp.Shutdown(ctx)
		}

		return err
	}, nil
}
