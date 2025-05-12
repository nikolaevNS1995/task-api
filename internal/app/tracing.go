package app

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"task-api/pkg/config"
	"time"
)

type Tracer struct {
	provider *trace.TracerProvider
}

func NewTracerProvider() *Tracer {
	return &Tracer{}
}

func InitTracerProvider(lc fx.Lifecycle, cfg *config.AppConfig, tracer *Tracer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var exporter trace.SpanExporter
			var err error

			if cfg.Telemetry.Local {
				exporter, err = stdouttrace.New()
			} else {
				exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(cfg.Telemetry.Host+":"+cfg.Telemetry.Port), otlptracegrpc.WithInsecure())
			}
			if err != nil {
				return fmt.Errorf("failed to create exporter: %w", err)
			}

			tp := trace.NewTracerProvider(
				trace.WithBatcher(exporter),
				trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(cfg.AppName))),
			)
			otel.SetTracerProvider(tp)
			tracer.provider = tp
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if tracer.provider == nil {
				return nil
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return tracer.provider.Shutdown(ctx)
		},
	})
}
