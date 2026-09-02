// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package main

import (
	"strings"
	"testing"

	"github.com/BreachSAFE/bqp-otel-go/sdk"
)

func TestEmitJSONLValidatesAndConsumesEvents(t *testing.T) {
	input := `{"event":"scan.start","component":"qureddy","outcome":"ok","attributes":{"target":"example.test"}}
{"event":"scan.finish","component":"qureddy","outcome":"error"}
`

	if err := emitJSONL(t.Context(), strings.NewReader(input), sdk.NewClient(nil)); err != nil {
		t.Fatalf("emitJSONL() error = %v", err)
	}
}

func TestEmitJSONLRejectsUnknownFields(t *testing.T) {
	input := `{"event":"scan.start","component":"qureddy","outcome":"ok","secret":"no"}`

	if err := emitJSONL(t.Context(), strings.NewReader(input), sdk.NewClient(nil)); err == nil {
		t.Fatal("emitJSONL() accepted an unknown field")
	}
}

func TestEmitJSONLRejectsInvalidEvent(t *testing.T) {
	input := `{"event":"scan.start","component":"qureddy","outcome":"bad"}`

	if err := emitJSONL(t.Context(), strings.NewReader(input), sdk.NewClient(nil)); err == nil {
		t.Fatal("emitJSONL() accepted an invalid outcome")
	}
}
