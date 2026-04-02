FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker


FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
EXPOSE 8080