# Build stage — Debian 13 (trixie) with Go 1.26.2
FROM debian:trixie AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION=1.26.2
RUN curl -fsSL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
    | tar -C /usr/local -xzf -
ENV PATH=/usr/local/go/bin:$PATH
ENV GOTOOLCHAIN=local

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gohub main.go

# Runtime stage — Debian 13 (trixie-slim)
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/gohub .
COPY --from=builder /app/docker-entrypoint.sh /usr/local/bin/
COPY --from=builder /app/.env.docker .env
COPY --from=builder /app/database/migrations ./database/migrations

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["docker-entrypoint.sh"]
