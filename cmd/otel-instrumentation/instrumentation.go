package instrumentation

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	EMAIL_SERVICE            = "email-service"
	DISCORD_SERVICE          = "discord-service"
	DISCORD_PROVIDER_SERVICE = "discord-provider-service"
	EMAIL_PROVIDER_SERVICE   = "email-provider-service"
	CALIDUM_ROTAE_SERVICE    = "calidum-rotate-services"
)

type Tracer struct {
	EmailTracer   trace.Tracer
	DiscordTracer trace.Tracer
	Enabled       bool
}

// OTLP exporter
func newOTLPExporter(ctx context.Context, host, port string) (sdktrace.SpanExporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%s", host, port)),
		otlptracehttp.WithInsecure(),
	)
	return otlptrace.New(ctx, client)
}

// Console exporter
func newConsoleExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
}

func newTracerProvider(consoleExporter, otlpExporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(CALIDUM_ROTAE_SERVICE),
		),
	)

	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(consoleExporter)),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(otlpExporter)),
		sdktrace.WithResource(r),
	), nil
}

func SetupOpenTelemetry(ctx context.Context, host, port string) (*sdktrace.TracerProvider, error) {
	otlpExporter, err := newOTLPExporter(ctx, host, port)
	if err != nil {
		return nil, err
	}

	consoleExporter, err := newConsoleExporter(ctx)
	if err != nil {
		return nil, err
	}

	tp, err := newTracerProvider(consoleExporter, otlpExporter)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return tp, nil
}

func (t *Tracer) StartSpanAndInitTracers(ctx context.Context, name string) (context.Context, trace.Span) {
	if t.Enabled {
		if t.EmailTracer == nil {
			t.EmailTracer = otel.Tracer(EMAIL_SERVICE)
		}

		if t.DiscordTracer == nil {
			t.DiscordTracer = otel.Tracer(DISCORD_SERVICE)
		}

		switch name {
		case EMAIL_PROVIDER_SERVICE:
			return t.EmailTracer.Start(ctx, name)
		case DISCORD_PROVIDER_SERVICE:
			return t.DiscordTracer.Start(ctx, name)
		}
	}

	return ctx, nil
}

func (t *Tracer) AddErrorEvent(ctx context.Context, span trace.Span, errorMessage string) trace.Span {
	if t.Enabled {
		span.AddEvent("error", trace.WithAttributes(
			attribute.String("error.message", errorMessage),
		))
		span.SetStatus(codes.Error, errorMessage)
		return span
	}

	return nil
}

func (t *Tracer) EndSpan(span trace.Span, message string) {
	if t.Enabled {
		span.SetStatus(codes.Ok, message)
		span.End()
	}
}
