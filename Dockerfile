FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./...
RUN go build -o /app/bin/gommitizen

# Path: Dockerfile
FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/bin/ /app/bin/

ENTRYPOINT ["/app/bin/gommitizen"]
