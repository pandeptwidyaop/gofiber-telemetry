package services

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

func setupOTLPExporter(ctx context.Context) (*otlptrace.Exporter, error) {

	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithHeaders(map[string]string{"Content-Type": "application/json"}),
			otlptracehttp.WithInsecure(),
		),
	)
}

func setupOTLPTraceProvider(exporter *otlptrace.Exporter, name string) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultExportTimeout*time.Millisecond)),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name),
			),
		),
	)
}

func InitTracing(ctx context.Context, serviceName string) (*trace.TracerProvider, error) {
	exporter, err := setupOTLPExporter(ctx)

	if err != nil {
		return nil, err
	}

	tp := setupOTLPTraceProvider(exporter, serviceName)

	otel.SetTracerProvider(tp)

	return tp, nil
}
