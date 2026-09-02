<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# Producer SDK boundary

## Contents

1. [Purpose](#purpose)
2. [Ownership](#ownership)
3. [Emission contract](#emission-contract)
4. [Redaction boundary](#redaction-boundary)
5. [Non-goals](#non-goals)

## Purpose

The `bqp-otel-go/sdk` package is the thin Go producer boundary for tools that need
to report `bqp.run.v1` lifecycle events. It is not a second Collector and does not persist
or interpret findings.

## Ownership

| Layer | Owns |
|---|---|
| Producer SDK | typed event construction, context/run correlation, bounded attributes, fail-open delivery |
| Collector binary | OTLP reception, central validation/redaction, batching, persistence, export |
| Product | scanner semantics, stdout results, evidence files, user-facing errors |

The SDK must not read the Docker socket, capture raw terminal streams, or place credentials,
OAuth URLs, or full findings payloads in attributes.

## Emission contract

The SDK will accept an explicit `context.Context`, copy caller-owned attributes at the
boundary, enforce event and attribute limits before encoding, and expose delivery failures
through an observable diagnostic path without making telemetry availability a product
success criterion. OTLP endpoint, timeout, and enablement are explicit configuration;
disabled-by-default remains the safe product default until a producer opts in.

## Redaction boundary

The SDK rejects sensitive attribute keys before export. The sandbox Collector adds a
defense-in-depth stock `redaction` processor that masks credential, OAuth, terminal, and
stream key variants from telemetry received from any producer. A real integration test sends
literal and variant sensitive keys through the generated Collector and asserts that the
persisted JSONL artifact does not contain their values. This is a deny-list of field names,
not a claim that arbitrary payload text can be safely scrubbed; producers must never put
secrets or raw output in telemetry.

## Non-goals

The SDK will not duplicate the Collector's file exporter, redaction policy, retry queue, or
backend adapters. A custom SDK abstraction requires a concrete compatibility or ownership
invariant and a real producer-to-Collector test.
