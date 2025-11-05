FROM golang AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /worker

FROM alpine:latest
WORKDIR /

COPY --from=builder /worker /worker

ENTRYPOINT ["/worker"]