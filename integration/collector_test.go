// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/BreachSAFE/bqp-otel-go/sdk"
	"go.opentelemetry.io/otel/attribute"
)

func TestSDKEmitsToRunningCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpoint := os.Getenv("BQP_OTEL_ENDPOINT")
	if endpoint == "" {
		endpoint = "127.0.0.1:4317"
	}
	provider, shutdown, err := sdk.NewOTLPProvider(ctx, endpoint, "bqp-integration", "test", true)
	if err != nil {
		t.Fatalf("NewOTLPProvider() error = %v", err)
	}
	client := sdk.NewClient(provider)
	event, err := sdk.NewEvent("integration.probe", "bqp-otel-go", sdk.OutcomeOK, "run-integration", "session-integration", 1, map[string]string{"probe": "real"})
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}
	if err := client.Emit(ctx, event); err != nil {
		t.Fatalf("Emit() error = %v", err)
	}
	tracer := provider.Tracer("integration/raw")
	_, span := tracer.Start(ctx, "integration.redaction")
	span.SetAttributes(
		attribute.String("token", "must-not-persist"),
		attribute.String("api_token", "must-not-persist-api-token"),
		attribute.String("Authorization", "must-not-persist-authorization"),
		attribute.String("db_password", "must-not-persist-password"),
		attribute.String("private-key", "must-not-persist-private-key"),
	)
	span.End()
	if err := shutdown(ctx); err != nil {
		t.Fatalf("provider shutdown: %v", err)
	}
}
