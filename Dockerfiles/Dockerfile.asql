FROM golang:1.25 AS builder

WORKDIR /build

COPY go.mod go.sum ./
COPY internal/ ./internal/
COPY cmd/ ./cmd/

RUN go mod download
RUN CGO_ENABLED=1 go build \
    -ldflags="-extldflags '-Wl,-rpath,/usr/local/lib'" \
    -o app ./cmd/server/

FROM debian:bookworm-slim AS runtime

WORKDIR /app

COPY --from=builder /build/app /app/app

EXPOSE 8080
RUN chmod +x /app/app

CMD ["/app/app"]