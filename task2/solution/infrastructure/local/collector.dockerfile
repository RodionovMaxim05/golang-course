FROM golang:alpine AS builder

WORKDIR /app

COPY api/ ./api/

COPY collector/go.mod collector/go.sum ./collector/
WORKDIR /app/collector
RUN go mod download

COPY collector/ .
RUN go build -o collector cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/collector/collector .
CMD ["./collector"]
