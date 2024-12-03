FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /tmp/build
COPY . .

RUN go mod download && \
    apk add dumb-init
RUN CGO_ENABLED=0 go build -o httpserver ./cmd/http

FROM gcr.io/distroless/static-debian12
COPY --from=builder /tmp/build/httpserver /usr/bin/dumb-init /

ENTRYPOINT ["/dumb-init", "--"]
CMD ["/httpserver"]
