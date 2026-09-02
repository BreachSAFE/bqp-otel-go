<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# Producer CLI

## Contents

1. [Purpose](#purpose)
2. [Contract](#contract)
3. [Failure behavior](#failure-behavior)

## Purpose

`bqp-otel-event` is the process boundary for producers that cannot link the Go SDK. It
consumes one or more bounded `bqp.run.v1` JSON objects on stdin and emits them over OTLP/gRPC.
It never reads or rewrites producer stdout findings or stderr diagnostics.

## Contract

```text
bqp-otel-event --endpoint HOST:4317 [--insecure] [--service-name NAME]
               [--service-version VERSION] < events.jsonl
```

- `--endpoint` is required, or may be supplied as `BQP_OTEL_ENDPOINT`.
- `--insecure` is explicit opt-in for a private local Collector endpoint.
- Each input line is one JSON `bqp.run.v1` event; unknown fields are rejected.
- Input is bounded to 1 MiB per line and validated by the shared Go SDK.
- stdout is unused; diagnostics and malformed-input errors go to stderr.

## Failure behavior

The command exits `0` only when every input line validates and is accepted by the SDK. It
exits `2` for configuration or input errors. Collector delivery remains fail-open at the
SDK/exporter boundary; producer findings are not made dependent on telemetry success.
