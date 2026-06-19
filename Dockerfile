# ── Stage 1: Build ──────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o studentsync ./backend/main.go

# ── Stage 2: Runtime ─────────────────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/studentsync .
COPY --from=builder /app/frontend ./frontend  

EXPOSE 8080

CMD ["./studentsync"]