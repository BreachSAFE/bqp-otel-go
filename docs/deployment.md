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

## Hardened sandbox profile

Build and run the local profile with:

```sh
docker compose -f docker-compose.sandbox.yaml up --build
```

The profile initializes the named volume with the collector's non-root UID, then binds
receivers to loopback, drops all Linux capabilities, enables
`no-new-privileges`, uses a read-only root filesystem, and grants write access only to
the named `bqp-otel-events` volume. `/tmp` is an explicitly bounded memory filesystem.
The Collector health endpoint is available at `http://127.0.0.1:13133/`; its image
healthcheck probes this endpoint rather than merely parsing the configuration. File
persistence rotates at 10 MiB per active file, retains at most five backups, and removes
backups older than seven days (a nominal 60 MiB ceiling including the active file).
The image is built from the pinned golden-Go builder and pinned distroless runtime; the
Compose tag is development-only until the release workflow publishes a signed digest.
