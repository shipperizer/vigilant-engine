FROM --platform=$BUILDPLATFORM golang:1.16 AS builder

LABEL org.opencontainers.image.source=https://github.com/shipperizer/vigilant-engine

ARG SKAFFOLD_GO_GCFLAGS
ARG TARGETOS
ARG TARGETARCH
ARG app_name=app

ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GO_BIN=/go/bin/app
ENV GRPC_HEALTH_PROBE_VERSION=v0.3.6
RUN apt-get update
RUN apt-get install -y awscli
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-${GOOS}-${GOARCH} && \
  chmod +x /bin/grpc_health_probe

RUN go env

WORKDIR /var/app

COPY . .

RUN make build

FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/shipperizer/vigilant-engine

COPY --from=builder /go/bin/app /app
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/app"]
