FROM golang:1.23.3-alpine3.20 as builder

WORKDIR /tmp/build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o httpserver ./cmd/http

FROM gcr.io/distroless/static-debian12
COPY --from=builder /tmp/build/httpserver /

CMD ["/httpserver"]
