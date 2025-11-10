FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download


COPY *.go ./

RUN go build -o ./dispatch

FROM alpine:3.22
WORKDIR /

COPY --from=builder /app/dispatch /dispatch


ENTRYPOINT ["/dispatch"]