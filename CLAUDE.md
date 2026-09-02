<!-- SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0 -->

# bqp-otel-go

## Contents

1. [Working rules](#working-rules)

## Working rules

Read [`AGENTS.md`](AGENTS.md) first. This is a standalone Go consumer of the BreachSAFE
golden toolchain image. Follow the common Go engineering, quality, CI/CD, container, and
release skills. Start with the official stock Collector distribution; use OCB only when a
measured BQP requirement proves the upstream binary insufficient, and add custom components
only after a real integration test.
