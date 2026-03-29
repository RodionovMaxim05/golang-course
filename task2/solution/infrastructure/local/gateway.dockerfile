FROM golang:alpine AS builder

WORKDIR /app

COPY api/ ./api/

COPY gateway/go.mod gateway/go.sum ./gateway/
WORKDIR /app/gateway
RUN go mod download

COPY gateway/ .
RUN go build -o gateway cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway/gateway .
CMD ["./gateway"]
