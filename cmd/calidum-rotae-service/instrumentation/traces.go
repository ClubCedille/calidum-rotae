package instrumentation

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	CALIDUM_ROTAE_TRACER     = "calidum-rotae-tracer"
	CALIDUM_ROTAE_SERVICE    = "calidum_rotae_service"
	DISCORD_PROVIDER_SERVICE = "discord_provider"
	EMAIL_PROVIDER_SERVICE   = "email_provider"
)

type Traces struct {
	CalidumRotaeTracer trace.Tracer
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

func newTracerProvider(otlpExporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
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
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(otlpExporter)),
		sdktrace.WithResource(r),
	), nil
}

func SetupOpenTelemetry(ctx context.Context, host, port string) (*sdktrace.TracerProvider, error) {
	otlpExporter, err := newOTLPExporter(ctx, host, port)
	if err != nil {
		return nil, err
	}

	// consoleExporter, err := newConsoleExporter(ctx)
	// if err != nil {
	// 	return nil, err
	// }
    //    tp, err := newTracerProvider(consoleExporter, otlpExporter)
	tp, err := newTracerProvider(otlpExporter)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return tp, nil
}

func (tracer *Traces) GrpcSpan(ctx context.Context, spanName, funcName, service string) (context.Context, trace.Span) {
	return tracer.CalidumRotaeTracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.method", funcName),
			attribute.String("rpc.service", service),
			attribute.Int("rpc.grpc.status_code", 200),
		),
	)
}

func (tracer *Traces) HttpPostSpan(ctx context.Context, g *gin.Context, spanName string) (context.Context, trace.Span) {
	return tracer.CalidumRotaeTracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("http.method", "POST"),
			attribute.String("http.scheme", "http"),
			attribute.String("http.target", g.Request.URL.Path),
			attribute.String("http.host", g.Request.Host),
			attribute.Int("http.status_code", g.Writer.Status()),
			attribute.String("http.user_agent", g.Request.UserAgent()),
		),
	)
}
