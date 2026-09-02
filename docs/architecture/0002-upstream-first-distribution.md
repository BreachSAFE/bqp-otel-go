<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# ADR-0002: Upstream-first Collector distribution

## Contents

1. [Status](#status)
2. [Decision](#decision)
3. [Options](#options)
4. [Release boundary](#release-boundary)
5. [Exit criteria](#exit-criteria)

## Status

Accepted for the 0.1.0 implementation sequence.

## Decision

Use the official OpenTelemetry Collector Contrib distribution wherever its stock
components satisfy BQP requirements. Pin its version and image digest, supply a reviewed
BQP configuration, and publish our image/configuration with SBOM, provenance, and Cosign
verification. Keep `otelcol-builder` available as a reproducible fallback, not as a reason
to own a custom Collector fork.

Producer code uses the official OpenTelemetry Go API/SDK and OTLP exporters. BQP code owns
only the `bqp.run.v1` envelope, bounded allowlisted attributes, and product lifecycle
semantics.

## Options

| Option | Decision | Reason |
|---|---|---|
| Official `otelcol-contrib` binary | Select first | Maximum upstream code, fixes, and multi-arch coverage |
| OCB-generated stock binary | Fallback | Reproducible component selection when the official distribution is too broad |
| Custom BQP Collector component | Defer | Only justified by a tested stock-component gap |
| BQP-owned exporter/protocol | Reject | Duplicates mature OTLP behavior and creates brittle maintenance |

## Release boundary

The BQP release artifact is the pinned image/configuration bundle, not an unreviewed fork
of upstream Collector source. The release workflow must record the upstream version and
digest, emit an SBOM, attach provenance, sign by digest, and verify the signature identity
before promoting a mutable tag.

## Exit criteria

Move from upstream binary to OCB only when a real integration test demonstrates a required
component or build property unavailable in the selected official distribution. Add custom
code only after a written gap analysis, bounded-input design, and real producer-to-Collector
conformance test.
