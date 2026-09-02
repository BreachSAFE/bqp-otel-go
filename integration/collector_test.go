// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/BreachSAFE/bqp-otel-go/sdk"
)

func TestSDKEmitsToRunningCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	provider, shutdown, err := sdk.NewOTLPProvider(ctx, "127.0.0.1:4317", "bqp-integration", "test", true)
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
	if err := shutdown(ctx); err != nil {
		t.Fatalf("provider shutdown: %v", err)
	}
}
