package instrumentation

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	EMAIL_SERVICE            = "email-service"
	DISCORD_SERVICE          = "discord-service"
	DISCORD_PROVIDER_SERVICE = "discord-provider-service"
	EMAIL_PROVIDER_SERVICE   = "email-provider-service"
)

type Tracer struct {
	EmailTracer   trace.Tracer
	DiscordTracer trace.Tracer
	Enabled       bool
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
