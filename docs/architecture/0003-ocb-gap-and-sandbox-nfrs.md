<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# ADR-0003: OCB gap analysis and sandbox NFRs

## Contents

1. [Status](#status)
2. [Context](#context)
3. [Gap analysis](#gap-analysis)
4. [Decision](#decision)
5. [Sandbox NFR targets](#sandbox-nfr-targets)
6. [Consequences](#consequences)

## Status

Accepted for the 0.1.0 sandbox distribution; targets are requirements, not benchmark
results. Re-measure before production sizing.

## Context

ADR-0002 requires the official Collector Contrib distribution whenever its stock
components satisfy the BQP requirement, with OCB as a documented fallback. ADR-0001
selected OCB before this gap was recorded. The sandbox also claimed bounded persistence
without stating measurable storage, throughput, latency, or health targets.

## Gap analysis

The BQP artifact requires exactly the stock component closure below:

```text
receiver:  otlp
processors: memory_limiter, redaction, batch
exporters:  debug, file
extensions: health_check, zpages
```

The official Contrib distribution is a broad multi-component binary and cannot produce
this minimal allowlist as a separate artifact. OCB is therefore selected as the fallback
for this sandbox: it reproducibly builds only the reviewed closure while retaining stock
component implementations. No BQP-specific Collector code is present. Re-evaluate this
choice if the official distribution adds a supported minimal-build mechanism or if a
production requirement needs a component outside this closure.

## Decision

Keep the pinned OCB build for the sandbox and record the component manifest in
`otelcol-builder.yaml`. ADR-0002 remains the governing upstream-first policy; this ADR is
the required gap record, not a replacement for it. The redaction and health behavior use
stock Contrib components rather than custom processors.

## Sandbox NFR targets

These are acceptance targets for the sandbox profile, not claims of measured capacity:

| NFR | Target | Configuration/evidence |
|---|---:|---|
| sustained intake | 500 events/s | benchmark required before production sizing |
| burst intake | 1,000 events/s for 30 s | benchmark required before production sizing |
| producer-to-file p99 | ≤ 2 s | batch flush and integration measurement |
| collector memory | 128 MiB limit + 32 MiB spike | `memory_limiter` |
| active file | ≤ 10 MiB | file rotation |
| retained backups | ≤ 5 files / 7 days | file rotation |
| nominal local-trace ceiling | ≤ 60 MiB | active file plus backups |
| health probe response | HTTP 200 within 2 s | `health_check` extension and image probe |

Abrupt termination may lose up to one batch/flush interval of local telemetry; producers
remain fail-open and product results are not coupled to collector durability.

## Consequences

The image is smaller and has a narrower component attack surface than a general Contrib
image, at the cost of owning an OCB build and its reproducibility evidence. The explicit
targets make future load testing falsifiable; until those benchmarks run, they remain
release-planning constraints rather than capacity claims.
