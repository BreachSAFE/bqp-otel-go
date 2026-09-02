<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# ADR-0001: OCB-built BQP telemetry distribution

## Contents

1. [Status](#status)
2. [Context](#context)
3. [Decision](#decision)
4. [Data flow](#data-flow)
5. [Deployment boundary](#deployment-boundary)
6. [Repository shape](#repository-shape)
7. [Migration](#migration)
8. [Consequences](#consequences)

## Status

Proposed. This records the architecture; implementation and acceptance remain staged.

## Context

BQP products are Python and Go applications with different lifecycle contexts. QuReddy-App
currently writes bounded local JSONL traces in Python. The platform needs one telemetry
transport, bounded local evidence, optional OTLP export, and independent release cadence
without duplicating sink and redaction logic in every product.

## Decision

Create `bqp-otel-go` as an OpenTelemetry Collector Builder distribution. Start with pinned
stock OTLP receivers, batch/memory limiting, and debug/file exporters. The `bqp.run.v1`
contract is language-neutral. Producers emit allowlisted lifecycle metadata; the Collector
owns validation, central redaction, bounded persistence, batching, and export.

Custom BQP processors or exporters require a written gap against a stock component and a real
integration test. No Docker socket is permitted. Producer delivery is bounded and fail-open.

## Repository shape

```text
sdk/                 thin Go producer SDK (typed bqp.run.v1 events)
schema/              language-neutral contract and fixtures
config/              sandbox Collector configuration
otelcol-builder.yaml pinned OCB component manifest
dist/                generated Collector binary (ignored, never hand-edited)
```

The SDK and Collector are versioned together initially but remain separate ownership
boundaries. A future split is allowed only after a second producer proves the package
boundary is stable.

## Data flow

```text
QuReddy / QuReddy-App / evidence-go / PDF-go / future tools
                         │
                         │ OTLP carrying bqp.run.v1 metadata
                         ▼
                 bqp-otel-go Collector
                  │       │        │
                  ▼       ▼        ▼
              file sink  debug   remote OTLP
```

## Deployment boundary

The default deployment is a separate Collector container in the same Docker stack as the
producers. This gives simple single-stack deployment while preserving a process, network,
volume, and release boundary. Host-local and shared-service deployments remain compatible
through the same OTLP endpoint contract; see [`../deployment.md`](../deployment.md).

## Migration

1. Freeze current Python sink behavior.
2. Version `bqp.run.v1` and publish fixtures.
3. Build the stock-component OCB distribution.
4. Add a thin Python producer client.
5. Prove QuReddy-App → Collector → bounded file output with real containers.
6. Integrate QuReddy, evidence-go, PDF-go, and later ike-scan.
7. Add optional remote exporters.
8. Remove duplicated Python sink and central redaction code.

## Consequences

The Collector becomes independently releasable and can serve multiple BQP products. There is
an additional process and schema boundary, so version compatibility, health, bounded storage,
multi-architecture builds, SBOMs, signatures, and provenance become release requirements.
