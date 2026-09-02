<!-- SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io> -->
<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# Deployment modes

## Contents

1. [Default: same stack](#default-same-stack)
2. [Host-local mode](#host-local-mode)
3. [Shared service mode](#shared-service-mode)
4. [Endpoint and failure contract](#endpoint-and-failure-contract)

## Default: same stack

The supported default places the Collector in the same Docker Compose/application stack as
the producers. It is still a separate container and process, with a private network and a
named writable telemetry volume.

```text
┌─────────────────────────────────────────────────────────────┐
│ BQP application stack                                        │
│                                                             │
│  ┌───────────────┐   OTLP/private network   ┌─────────────┐ │
│  │ QuReddy-App   │────────────────────────▶│ bqp-otel    │ │
│  │ QuReddy/tools │                          │ Collector   │ │
│  └───────────────┘                          └──────┬──────┘ │
│                                                   │        │
│                                      named volume │        │
│                                                   ▼        │
│                                      bqp events JSONL       │
└─────────────────────────────────────────────────────────────┘
```

No Docker socket is mounted. Producers address the service by its Compose DNS name, not by
`localhost`. The Collector receiver is not published to the host unless a local operator
explicitly needs it.

## Host-local mode

For development, the Collector may run on the host while producers run in Docker. The
producer endpoint is explicitly configured to the host gateway; no code changes are needed.

```text
Docker producers ── OTLP ──▶ host bqp-otel ──▶ local volume/backends
```

## Shared service mode

Multiple stacks may send OTLP to a dedicated Collector service. This mode requires TLS or
mTLS, authentication at the network boundary, and per-producer resource identity.

```text
QuReddy-App ─┐
QuReddy     ─┼── mTLS OTLP ──▶ shared bqp-otel ──▶ backends
evidence-go ─┘
```

## Endpoint and failure contract

- OTLP/gRPC is the default internal transport; OTLP/HTTP is available for constrained hosts.
- Producers are disabled by default and use bounded, fail-open delivery when enabled.
- Collector persistence is through a named writable volume, never the image layer.
- A Collector outage must not turn a scan, session, or evidence build into a product failure.
- The Collector image is independently versioned, scanned, signed, and verified.
