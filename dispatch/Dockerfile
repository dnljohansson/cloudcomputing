FROM golang AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dispatch

FROM alpine:latest
WORKDIR /

COPY --from=builder /dispatch /dispatch

ENTRYPOINT ["/dispatch"]