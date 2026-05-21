# token-audit

`token-audit` is a transparent OpenAI-compatible audit proxy with a built-in review console. It sits in front of an existing upstream API, forwards requests unchanged, and stores redacted request and response metadata in PostgreSQL for later inspection.

The current implementation is primarily validated against **NewAPI** as the upstream service. It may work with other OpenAI-compatible gateways, but those integrations have not been formally tested yet.

## What it does

- Proxies selected OpenAI-style endpoints such as `/v1/chat/completions`, `/v1/responses`, `/v1/completions`, and `/v1/embeddings`
- Stores request timing, status codes, token usage, model names, user text, assistant text, and selected raw payloads
- Replaces raw API keys with HMAC-based fingerprints so the original token is not stored in the database
- Provides a web console for dashboards, request search, token aliasing, and database maintenance actions
- Supports both standard JSON responses and streaming SSE responses

## When to use it

This project is useful when you need a lightweight audit layer in front of an OpenAI-compatible gateway and want to answer questions like:

- Which token consumed the most usage?
- Which models are generating the most traffic?
- Which requests failed, timed out, or returned unusual payloads?
- What prompt and response text was associated with a specific request?

## Important notes

- This proxy intentionally stores redacted request and response content for auditing. Do not deploy it unless your privacy and retention policy allows prompt and completion logging.
- The admin console should be protected behind a trusted network boundary, VPN, or reverse proxy authentication layer.
- Raw `Authorization`, `Cookie`, and `Set-Cookie` headers are excluded from stored headers.
- Current compatibility and testing notes are documented in [COMPATIBILITY.md](./COMPATIBILITY.md).

## Quick start

The fastest path is Docker Compose:

1. Edit `config.docker.yaml` and replace the placeholder `hmac_secret` and `admin_password_hash`.
2. Point `upstream_base` at your existing NewAPI or OpenAI-compatible gateway.
3. Run `docker compose up -d --build`.
4. Open `http://127.0.0.1:3007/audit`.

For a full production-oriented walkthrough, see [DEPLOYMENT.md](./DEPLOYMENT.md).

## Repository contents

- `cmd/newapi-audit-proxy`: application entrypoint
- `internal/proxy`: transparent proxy logic and response streaming
- `internal/audit`: request parsing, redaction, storage, and reporting
- `internal/web`: audit console and API endpoints
- `migrations`: PostgreSQL schema and migrations
- `scripts/hash_password.go`: helper to generate the admin bcrypt password hash
- `audit.sh`: Linux systemd deployment helper

## Validation status

This repository is being published as a cleaned release snapshot. Source code, deployment files, and documentation are included; local binaries, caches, and machine-specific config have been intentionally excluded.
