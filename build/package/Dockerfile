FROM golang:1.16-buster AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./. ./

RUN go build -v ./cmd/autotest

FROM debian:buster-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/autotest /app/autotest

CMD ["/app/autotest"]