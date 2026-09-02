<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# bqp-otel-go instructions

## Contents

1. [Scope](#scope)
2. [Boundaries](#boundaries)
3. [Required gates](#required-gates)

## Scope

This repository owns the BQP OpenTelemetry Collector distribution and its configuration.
It does not own scanner semantics, Gradio/WebSocket/PTY lifecycle, evidence interpretation,
or PDF rendering.

## Boundaries

- Use the official OpenTelemetry Collector Builder (OCB) and pinned stock components first.
- Treat `bqp.run.v1` as a language-neutral producer contract.
- Keep producers fail-open and bounded; collector failure must not break a scan or session.
- Never accept raw terminal bytes, credentials, OAuth URLs, or sensitive scan payloads as
  default telemetry fields.
- Add a custom processor/exporter only after a stock component is shown insufficient by a
  real requirement and test.
- Do not mount the Docker socket.

## Required gates

Run the golden-Go gate sequence: `gofmt`, `go vet`, `go test`, `go test -race`, `go mod verify`,
`staticcheck`, `govulncheck`, `osv-scanner`, Docker build/smoke, SBOM, and provenance checks.
Use real producer-to-collector integration tests; do not add mock-only acceptance tests.

Every Markdown file has a numbered table of contents. Every nontrivial change records the
10-step conformance status and performs the shared anti-pattern review.
