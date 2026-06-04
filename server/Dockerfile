# ── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
COPY . .

ENV GONOSUMDB="*"
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS=-mod=mod go build -ldflags="-s -w" -o netwatch-server ./main.go

# ── Runtime stage ────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/netwatch-server /netwatch-server

EXPOSE 8080
ENTRYPOINT ["/netwatch-server"]
