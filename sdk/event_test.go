// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package sdk

import (
	"errors"
	"testing"
)

func TestNewEventCopiesAndSortsAttributes(t *testing.T) {
	attributes := map[string]string{"zeta": "last", "alpha": "first"}
	event, err := NewEvent("scan.start", "qureddy", OutcomeOK, "run-1", "session-1", 3, attributes)
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}
	attributes["alpha"] = "mutated"
	if event.Attributes["alpha"] != "first" {
		t.Fatalf("event retained caller-owned attribute map")
	}
	want := []string{"alpha", "zeta"}
	got := event.AttributeKeys()
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AttributeKeys() = %v, want %v", got, want)
		}
	}
}

func TestNewEventRejectsSensitiveAndUnboundedInputs(t *testing.T) {
	if _, err := NewEvent("scan.start", "qureddy", OutcomeOK, "run", "session", 0, map[string]string{"oauth_token": "x"}); !errors.Is(err, ErrSensitiveKey) {
		t.Fatalf("sensitive key error = %v", err)
	}
	if _, err := NewEvent("scan\nstart", "qureddy", OutcomeOK, "run", "session", 0, nil); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("newline error = %v", err)
	}
	if _, err := NewEvent("scan.start", "qureddy", OutcomeOK, "run", "session", -1, nil); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("negative duration error = %v", err)
	}
}
