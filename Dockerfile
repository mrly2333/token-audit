FROM golang:1.23 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/newapi-audit-proxy ./cmd/newapi-audit-proxy

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /out/newapi-audit-proxy /app/newapi-audit-proxy
COPY config.docker.yaml /app/config.yaml

EXPOSE 3007

CMD ["/app/newapi-audit-proxy", "-config", "/app/config.yaml"]
