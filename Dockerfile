# syntax=docker/dockerfile:1.4

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mockserver main.go

# Runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mockserver /app/
# COPY config.json /app/

EXPOSE 8080

ENTRYPOINT ["./mockserver"]
