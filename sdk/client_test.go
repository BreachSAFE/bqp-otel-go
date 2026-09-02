// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package sdk

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestClientEmitsValidatedSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := trace.NewTracerProvider(trace.WithSyncer(exporter))
	client := NewClient(provider)
	event, err := NewEvent("scan.start", "qureddy", OutcomeOK, "run-1", "session-1", 4, map[string]string{"target.kind": "tls"})
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}
	if err := client.Emit(context.Background(), event); err != nil {
		t.Fatalf("Emit() error = %v", err)
	}
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("exported spans = %d, want 1", len(spans))
	}
	if got := spans[0].Name; got != "scan.start" {
		t.Fatalf("span name = %q, want scan.start", got)
	}
	if got := len(spans[0].Attributes); got != 7 {
		t.Fatalf("span attribute count = %d, want 7", got)
	}
}

func TestClientRejectsDirectlyConstructedInvalidEvent(t *testing.T) {
	client := NewClient(nil)
	err := client.Emit(context.Background(), Event{Event: "", Component: "qureddy", Outcome: OutcomeOK})
	if err == nil {
		t.Fatal("Emit() accepted invalid event")
	}
}
