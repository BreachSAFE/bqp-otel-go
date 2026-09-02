<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# bqp-otel-go

## Contents

1. [Purpose](#purpose)
2. [Architecture](#architecture)
3. [Current status](#current-status)
4. [Delivery plan](#delivery-plan)
5. [Deployment](#deployment)
6. [Repository rules](#repository-rules)
7. [Producer CLI](#producer-cli)

## Purpose

`bqp-otel-go` is the BreachSAFE Quantum Platform OpenTelemetry Collector distribution.
It packages the official OpenTelemetry Collector Contrib distribution first, with a pinned
configuration and release evidence. OCB is retained as a reproducible fallback if a real
requirement cannot be met by the upstream binary; no BQP-specific Collector component is
assumed.

## Architecture

```text
BQP producers (QuReddy, QuReddy-App, evidence-go, PDF-go)
                     │ bqp.run.v1 over OTLP
                     ▼
              bqp-otel-go
       OCB-built Collector binary
         validate → bound → persist → export
```

Products own lifecycle context. This repository owns telemetry transport and export.
The authoritative plan is [`docs/architecture/0001-collector-distribution.md`](docs/architecture/0001-collector-distribution.md).

## Current status

The initial architecture checkpoint is on `main` (`45aa9f5`). No custom processor or
exporter is shipped. The first implementation must use pinned stock OTel components and
prove a real producer-to-collector path before central redaction or backend-specific code
is added.

## Delivery plan

The [0.1.0 milestone](https://github.com/BreachSAFE/bqp-otel-go/milestone/1) tracks the
work in dependency order: reproducible distribution, schema fixtures, hardened sandbox,
producer integrations, backend profiles, then real conformance and release gates. The
issue tracker is the execution record; the ADR is the architecture record.

## Deployment

The default deployment is a separate `bqp-otel` container in the same Docker stack as the
producers. Host-local and shared-service modes use the same OTLP contract and are documented
in [`docs/deployment.md`](docs/deployment.md).

## Repository rules

Read [`AGENTS.md`](AGENTS.md) and [`CLAUDE.md`](CLAUDE.md) before changing this repository.
Use the BreachSAFE golden Go image; do not copy its Dockerfile into this repository.

## Producer CLI

Producers that cannot link the SDK can use [`docs/producer-cli.md`](docs/producer-cli.md)
and the `bqp-otel-event` command to send bounded `bqp.run.v1` JSONL over OTLP without
changing their stdout findings or stderr diagnostics.
