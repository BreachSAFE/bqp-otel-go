// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package sdk

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Client emits validated BQP events as short OTel spans. A zero Client uses the global
// provider, which is a no-op until a producer explicitly configures an exporter.
type Client struct {
	tracer trace.Tracer
}

// NewClient constructs a client from an explicit provider. A nil provider selects the
// global provider and preserves disabled-by-default behavior.
func NewClient(provider trace.TracerProvider) Client {
	if provider == nil {
		provider = otel.GetTracerProvider()
	}
	return Client{tracer: provider.Tracer("github.com/BreachSAFE/bqp-otel-go/sdk")}
}

// Emit validates the event, records it in one bounded span, and returns validation errors.
// Export failures are surfaced by the provider's exporter lifecycle, not hidden here.
func (c Client) Emit(ctx context.Context, event Event) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", ErrInvalidEvent)
	}
	if c.tracer == nil {
		return fmt.Errorf("%w: nil tracer", ErrInvalidEvent)
	}
	if err := event.validate(); err != nil {
		return err
	}
	_, span := c.tracer.Start(ctx, event.Event, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()
	span.SetAttributes(
		attribute.String("bqp.schema", Schema),
		attribute.String("bqp.component", event.Component),
		attribute.String("bqp.outcome", string(event.Outcome)),
		attribute.String("bqp.run_id", event.RunID),
		attribute.String("bqp.session_id", event.SessionID),
		attribute.Int64("bqp.duration_ms", event.DurationMS),
	)
	for key, value := range event.Attributes {
		span.SetAttributes(attribute.String(key, value))
	}
	if event.Outcome == OutcomeError || event.Outcome == OutcomeTimeout {
		span.SetStatus(codes.Error, string(event.Outcome))
	}
	return nil
}

// NewOTLPProvider creates an explicit OTLP/gRPC provider and its shutdown function. The
// caller owns shutdown and must invoke it with a bounded context. Insecure transport must
// be opted into explicitly for a private local network; remote endpoints stay TLS by default.
func NewOTLPProvider(ctx context.Context, endpoint, serviceName, serviceVersion string, insecure bool) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	if ctx == nil {
		return nil, nil, fmt.Errorf("%w: nil context", ErrInvalidEvent)
	}
	if endpoint == "" || serviceName == "" {
		return nil, nil, fmt.Errorf("%w: endpoint and service name are required", ErrInvalidEvent)
	}
	options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
	if insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}
	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.version", serviceVersion),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create OTel resource: %w", err)
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	return provider, provider.Shutdown, nil
}
