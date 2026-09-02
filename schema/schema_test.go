// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package schema_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/BreachSAFE/bqp-otel-go/sdk"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestFixturesConformToBQPEventSchema(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	var schema any
	decodeJSON(t, "bqp.run.v1.json", &schema)
	if err := compiler.AddResource("bqp.run.v1.json", schema); err != nil {
		t.Fatalf("AddResource() error = %v", err)
	}
	compiled, err := compiler.Compile("bqp.run.v1.json")
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}
	var valid any
	decodeJSON(t, "fixtures/valid-event.json", &valid)
	if err := compiled.Validate(valid); err != nil {
		t.Fatalf("valid fixture rejected: %v", err)
	}
	var invalid any
	decodeJSON(t, "fixtures/invalid-missing-outcome.json", &invalid)
	if err := compiled.Validate(invalid); err == nil {
		t.Fatal("invalid fixture accepted")
	}
	var invalidNewline any
	decodeJSON(t, "fixtures/invalid-newline-event.json", &invalidNewline)
	if err := compiled.Validate(invalidNewline); err == nil {
		t.Fatal("newline event fixture accepted")
	}
}

func TestSDKAndJSONSchemaAgreeOnCanonicalCorpus(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	var rawSchema any
	decodeJSON(t, "bqp.run.v1.json", &rawSchema)
	if err := compiler.AddResource("bqp.run.v1.json", rawSchema); err != nil {
		t.Fatalf("AddResource() error = %v", err)
	}
	compiled, err := compiler.Compile("bqp.run.v1.json")
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	attributes := map[string]string{"target.kind": "tls", "probe": "real"}
	event, err := sdk.NewEvent("scan.start", "qureddy", sdk.OutcomeOK, "run-1", "session-1", 4, attributes)
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var document any
	if err := json.Unmarshal(encoded, &document); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if err := compiled.Validate(document); err != nil {
		t.Fatalf("JSON schema rejected SDK event: %v", err)
	}

	invalidCases := []struct {
		name      string
		event     string
		component string
		outcome   sdk.Outcome
		attrs     map[string]string
	}{
		{name: "newline", event: "scan\nstart", component: "qureddy", outcome: sdk.OutcomeOK},
		{name: "negative duration", event: "scan.start", component: "qureddy", outcome: sdk.OutcomeOK},
		{name: "too many attributes", event: "scan.start", component: "qureddy", outcome: sdk.OutcomeOK, attrs: oversizedAttributes()},
		{name: "oversized attribute key", event: "scan.start", component: "qureddy", outcome: sdk.OutcomeOK, attrs: map[string]string{strings.Repeat("k", sdk.MaxEventLength+1): "value"}},
	}
	for _, testCase := range invalidCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.attrs == nil {
				testCase.attrs = map[string]string{"probe": "real"}
			}
			duration := int64(0)
			if testCase.name == "negative duration" {
				duration = -1
			}
			if _, err := sdk.NewEvent(testCase.event, testCase.component, testCase.outcome, "run", "session", duration, testCase.attrs); err == nil {
				t.Fatal("SDK accepted invalid corpus entry")
			}
			document := map[string]any{
				"schema": "bqp.run.v1", "event": testCase.event, "component": testCase.component,
				"outcome": testCase.outcome, "run_id": "run", "session_id": "session",
				"duration_ms": duration, "attributes": testCase.attrs,
			}
			if err := compiled.Validate(document); err == nil {
				t.Fatal("JSON schema accepted invalid corpus entry")
			}
		})
	}
}

func oversizedAttributes() map[string]string {
	attributes := make(map[string]string, sdk.MaxAttributes+1)
	for index := 0; index <= sdk.MaxAttributes; index++ {
		attributes["key"] = "value"
		attributes[fmt.Sprintf("key-%d", index)] = "value"
	}
	return attributes
}

func decodeJSON(t *testing.T, name string, target *any) {
	t.Helper()
	file, err := os.Open(name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(target); err != nil {
		t.Fatalf("decode %s: %v", name, err)
	}
}
