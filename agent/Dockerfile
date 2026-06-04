# ── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o netwatch-agent ./main.go

# ── Runtime stage ────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/netwatch-agent /netwatch-agent

ENTRYPOINT ["/netwatch-agent"]
