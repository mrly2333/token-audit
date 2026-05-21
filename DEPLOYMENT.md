# Deployment Guide

This document describes the supported deployment patterns for `token-audit` and the assumptions behind each one.

## Architecture

The proxy listens on a single port and serves two things at once:

- Audit web UI under `web_base_path` such as `/audit`
- Transparent proxy traffic for all other paths

Typical flow:

1. Your clients send OpenAI-compatible traffic to `token-audit`.
2. `token-audit` forwards the request to the configured `upstream_base`.
3. Matching capture paths are recorded in PostgreSQL.
4. Operators review logs in the built-in web console.

## Requirements

- PostgreSQL 16 or newer
- One existing upstream OpenAI-compatible gateway
- A deployment host that can reach both PostgreSQL and the upstream gateway
- Go 1.23+ if you want to build locally from source
- Docker and Docker Compose if you prefer container deployment

## Configuration fields

`config.example.yaml` and `config.docker.yaml` document the supported fields:

- `listen_addr`: bind address for the proxy and audit web UI
- `upstream_base`: base URL of the upstream gateway to forward traffic to
- `web_base_path`: path prefix for the audit UI, usually `/audit`
- `postgres_dsn`: PostgreSQL DSN used for migrations and runtime storage
- `hmac_secret`: secret used to fingerprint bearer tokens without storing them directly
- `admin_username`: username for the built-in admin login
- `admin_password_hash`: bcrypt hash for the built-in admin login
- `max_body_bytes`: max bytes to retain per request or response payload
- `capture_paths`: upstream paths that should be fully audited

## Generating the admin password hash

If Go is available locally:

```bash
go run ./scripts/hash_password.go "<strong-password>"
```

If you prefer Docker:

```bash
docker run --rm -v "$PWD:/src" -w /src golang:1.23 \
  go run ./scripts/hash_password.go "<strong-password>"
```

Paste the generated bcrypt hash into `admin_password_hash`.

## Option 1: Docker Compose

This is the easiest deployment model for a single host.

1. Edit `config.docker.yaml`.
2. Set a real `hmac_secret`.
3. Generate and paste a valid `admin_password_hash`.
4. Point `upstream_base` at your running upstream gateway.
5. Start the stack:

```bash
docker compose up -d --build
```

What this stack provides:

- `postgres`: local PostgreSQL container
- `audit-proxy`: the compiled proxy container

Default access points:

- Proxy and audit service: `http://127.0.0.1:3007`
- Audit UI: `http://127.0.0.1:3007/audit`
- PostgreSQL: `127.0.0.1:5432`

## Option 2: Binary + systemd

This option is suitable for a Linux server where you want a native service instead of Docker.

1. Install Go 1.23+.
2. Build the binary:

```bash
go build -o bin/newapi-audit-proxy ./cmd/newapi-audit-proxy
```

3. Copy `config.example.yaml` to `config.yaml` and fill in real values.
4. Install and start the service:

```bash
bash audit.sh install
```

Useful service commands:

```bash
bash audit.sh status
bash audit.sh restart
bash audit.sh log
```

`audit.sh` installs a `systemd` unit that runs the binary from the current project directory.

## Reverse proxy and exposure guidance

If you expose the audit UI to a wider audience, place a reverse proxy in front of it and restrict access. Recommended controls:

- allow-list trusted source IPs
- protect `/audit` with SSO, VPN, or reverse proxy authentication
- terminate TLS at the reverse proxy
- keep PostgreSQL private and non-public

## Operational notes

- Database migrations run automatically on startup.
- The application stores redacted prompt and response content for captured routes.
- Token fingerprints are deterministic per `hmac_secret`; rotating the secret changes future fingerprints.
- Streaming responses are captured in a best-effort manner by reconstructing SSE output into stored text and usage summaries.

## Troubleshooting

If the service starts but login fails:

- verify that `admin_password_hash` is a real bcrypt hash
- verify that the browser is opening the same `web_base_path` configured in YAML

If requests pass through but no logs appear:

- verify that the request path is included in `capture_paths`
- verify that PostgreSQL is reachable from the application process
- check service logs for migration or insert errors

If clients receive upstream errors:

- verify `upstream_base`
- verify that the upstream gateway accepts the same paths you are forwarding
- confirm that the upstream gateway is reachable from the deployment host
