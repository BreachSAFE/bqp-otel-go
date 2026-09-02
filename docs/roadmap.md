<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# bqp-otel-go roadmap

## Contents

1. [Milestone](#milestone)
2. [Execution order](#execution-order)
3. [Done versus remaining](#done-versus-remaining)

## Milestone

The active milestone is [0.1.0 — BQP OTel Collector distribution](https://github.com/BreachSAFE/bqp-otel-go/milestone/1).

## Execution order

1. [#1](https://github.com/BreachSAFE/bqp-otel-go/issues/1) — consumer CI on the pinned golden-Go image.
2. [#2](https://github.com/BreachSAFE/bqp-otel-go/issues/2) — generated distribution from pinned stock OTel components.
3. [#3](https://github.com/BreachSAFE/bqp-otel-go/issues/3) — close `bqp.run.v1` and add fixtures.
4. [#4](https://github.com/BreachSAFE/bqp-otel-go/issues/4) — hardened sandbox profile.
5. [#5](https://github.com/BreachSAFE/bqp-otel-go/issues/5) — validation and redaction decision.
6. [#6](https://github.com/BreachSAFE/bqp-otel-go/issues/6) — QuReddy-App producer integration.
7. [#7](https://github.com/BreachSAFE/bqp-otel-go/issues/7) — QuReddy scanner producer integration.
8. [#8](https://github.com/BreachSAFE/bqp-otel-go/issues/8) — evidence and PDF producer adapters.
9. [#9](https://github.com/BreachSAFE/bqp-otel-go/issues/9) — Grafana/Tempo/Loki profile.
10. [#10](https://github.com/BreachSAFE/bqp-otel-go/issues/10) — real conformance and release gates.

## Done versus remaining

| State | Scope |
|---|---|
| ✅ | Repository scope, ADR, draft config, schema boundary, pinned CI contract, Go SDK foundation |
| 🟡 | Generated collector binary, Go SDK foundation, and SDK→Collector integration gate wired |
| ⏳ | Producer integrations and backend profiles |
| ⏳ | Real multi-process conformance, SBOM, signing, and release |

The local OCB compile was attempted with the pinned image and OCB version. It reached
dependency resolution but stopped when the disposable container filesystem exhausted its
`/go/pkg` space; this is recorded as an environment limitation, not a passing build.
