// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

package schema_test

import (
	"encoding/json"
	"os"
	"testing"

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
