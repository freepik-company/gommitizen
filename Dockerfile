FROM golang:1.21.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    gcc \
    musl-dev \
    git

RUN git config --global user.email "gommitizen@localhost" && \
    git config --global user.name "Gommitizen"

COPY . .

RUN go mod download

RUN go test -v ./...
RUN go build -o /app/bin/gommitizen

# Path: Dockerfile
FROM alpine:3.13

WORKDIR /code

RUN apk add --no-cache \
    git

RUN git config --global user.email "gommitizen@localhost" && \
    git config --global user.name "Gommitizen"

COPY --from=builder /app/bin/ /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/gommitizen"]
