# SPDX-FileCopyrightText: 2026 BreachSAFE <https://www.breachsafe.io>
# SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0

ARG GOLDEN_GO_IMAGE=ghcr.io/paul007ex/breachsafe-golden-go:1.27.0-cgo@sha256:542520836ac7f681826af610acb9fb5230e1e1b9e97ae2f2d2b744099d32f605
ARG DISTROLESS_IMAGE=gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

FROM ${GOLDEN_GO_IMAGE} AS builder
USER 0
WORKDIR /src
COPY . .
ARG OCB_VERSION=v0.159.0
RUN go install go.opentelemetry.io/collector/cmd/builder@${OCB_VERSION} \
    && GOFLAGS="${GOFLAGS} -buildvcs=false" builder --config otelcol-builder.yaml \
    && CGO_ENABLED=0 GOFLAGS="${GOFLAGS} -buildvcs=false" go build -trimpath -o dist/bqp-otel-healthcheck ./cmd/bqp-otel-healthcheck \
    && test -x dist/bqp-otel \
    && test -x dist/bqp-otel-healthcheck

FROM ${DISTROLESS_IMAGE} AS runtime
COPY --from=builder /src/dist/bqp-otel /bqp-otel
COPY --from=builder /src/dist/bqp-otel-healthcheck /bqp-otel-healthcheck
COPY config/collector.sandbox.yaml /etc/bqp-otel/config.yaml
VOLUME ["/var/lib/bqp-otel"]
EXPOSE 4317 4318 13133 55679
USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD ["/bqp-otel-healthcheck", "http://127.0.0.1:13133/"]
ENTRYPOINT ["/bqp-otel"]
CMD ["--config=/etc/bqp-otel/config.yaml"]
