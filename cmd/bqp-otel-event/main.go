// SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

// Command bqp-otel-event bridges bounded bqp.run.v1 JSONL to an OTLP collector.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/BreachSAFE/bqp-otel-go/sdk"
)

const maxLineBytes = 1 << 20

func main() {
	endpoint := flag.String("endpoint", os.Getenv("BQP_OTEL_ENDPOINT"), "OTLP/gRPC endpoint")
	service := flag.String("service-name", "bqp-producer", "OTel service.name")
	version := flag.String("service-version", "unknown", "OTel service.version")
	insecure := flag.Bool("insecure", false, "use insecure transport for a private local endpoint")
	flag.Parse()
	if strings.TrimSpace(*endpoint) == "" {
		fail("-endpoint or BQP_OTEL_ENDPOINT is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	provider, shutdown, err := sdk.NewOTLPProvider(ctx, *endpoint, *service, *version, *insecure)
	if err != nil {
		fail("configure OTLP provider: %v", err)
	}
	defer func() {
		if shutdownErr := shutdown(context.Background()); shutdownErr != nil {
			fmt.Fprintf(os.Stderr, "bqp-otel-event: shutdown: %v\n", shutdownErr)
		}
	}()

	if err := emitJSONL(ctx, os.Stdin, sdk.NewClient(provider)); err != nil {
		fail("emit JSONL: %v", err)
	}
}

func emitJSONL(ctx context.Context, input io.Reader, client sdk.Client) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), maxLineBytes)
	line := 0
	for scanner.Scan() {
		line++
		var event sdk.Event
		decoder := json.NewDecoder(strings.NewReader(scanner.Text()))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&event); err != nil {
			return fmt.Errorf("line %d: decode: %w", line, err)
		}
		if event.Schema != sdk.Schema {
			return fmt.Errorf("line %d: schema must be %q", line, sdk.Schema)
		}
		if err := client.Emit(ctx, event); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "bqp-otel-event: "+format+"\n", args...)
	os.Exit(2)
}
