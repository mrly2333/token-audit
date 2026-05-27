<div align="center">

# token-audit

**Transparent AI API Audit Proxy with Built-in Dashboard**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](docker-compose.yml)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey?style=flat-square)]()

A lightweight, transparent proxy that sits in front of your OpenAI-compatible API gateway, silently logging every request and response to PostgreSQL — with a clean web dashboard for searching, filtering, and analyzing your AI traffic.

[English](#features) | [中文](#中文说明)

</div>

---

## Features

- **Multi-Protocol Support** — OpenAI (`/v1/chat/completions`, `/v1/responses`, `/v1/completions`, `/v1/embeddings`) and Claude (`/v1/messages`)
- **Full Request Capture** — Headers, body, user prompts, assistant responses, model name, token usage, latency
- **Streaming Aware** — Handles both JSON and SSE streaming responses, reconstructing text from delta chunks
- **Token Privacy** — HMAC-based fingerprinting replaces raw API keys; real tokens are never stored
- **Built-in Dashboard** — Search logs, filter by token/model/status, view per-token and per-model statistics
- **Auto Migration** — Database schema updates automatically on startup
- **One-Command Deploy** — Docker Compose ready, or use the included `systemd` helper script

## Architecture

```
                          ┌─────────────────────┐
  Client ──────────────►  │     token-audit      │
  (OpenAI / Claude API)   │   (this project)     │
                          └──────────┬───────────┘
                                     │
                          ┌──────────▼───────────┐
                          │   Upstream Gateway    │
                          │  (NewAPI / OpenAI /   │
                          │   Claude / etc.)      │
                          └──────────────────────┘
                                     │
                          ┌──────────▼───────────┐
                          │     PostgreSQL        │
                          │   (audit storage)     │
                          └──────────────────────┘
```

## Quick Start

### Docker Compose (recommended)

```bash
# 1. Clone the repo
git clone https://github.com/mrly2333/token-audit.git
cd token-audit

# 2. Edit config
cp config.example.yaml config.yaml
# Set: upstream_base, hmac_secret, admin_password_hash

# 3. Generate admin password hash
docker run --rm -v "$PWD:/src" -w /src golang:1.23 \
  go run ./scripts/hash_password.go "your-strong-password"

# 4. Start
docker compose up -d --build

# 5. Open dashboard
# http://127.0.0.1:3007/audit
```

### Binary + systemd

```bash
go build -o bin/newapi-audit-proxy ./cmd/newapi-audit-proxy
cp config.example.yaml config.yaml  # edit config
bash audit.sh install
```

## Screenshots

> _Add screenshots of the dashboard here to increase visibility._

| Dashboard Overview | Request Logs | Token Statistics |
|---|---|---|
| ![overview](docs/screenshots/overview.png) | ![logs](docs/screenshots/logs.png) | ![tokens](docs/screenshots/tokens.png) |

## Use Cases

- **Cost Attribution** — Which token/model consumes the most?
- **Debugging** — What prompt caused a 500 error? What did the model actually return?
- **Compliance** — Keep an audit trail of all AI API interactions
- **Monitoring** — Track error rates, latency, and token usage trends

## Documentation

- [Deployment Guide](DEPLOYMENT.md) — Full setup instructions
- [Compatibility](COMPATIBILITY.md) — Tested upstream gateways and supported APIs
- [Contributing](CONTRIBUTING.md) — How to contribute

---

## 中文说明

`token-audit` 是一个透明的 OpenAI 兼容审计代理，内置可视化审计后台。它部署在现有上游 API 前面，原样转发请求，并将脱敏后的请求、响应和用量信息写入 PostgreSQL，方便后续排查、统计和审计。

### 支持的协议

| 协议 | 端点 | 认证方式 |
|------|------|----------|
| OpenAI | `/v1/chat/completions`, `/v1/responses`, `/v1/completions`, `/v1/embeddings` | `Authorization: Bearer` |
| Claude | `/v1/messages` | `x-api-key` |

### 适用场景

- 哪个令牌消耗的用量最多？
- 哪些模型的调用流量最高？
- 哪些请求失败、超时或返回了异常内容？
- 某一次请求对应的提示词和响应文本是什么？

### 重要说明

- 审计后台应放在可信网络边界之后，配合 VPN 或反向代理鉴权使用
- 存储的请求头中会排除 `Authorization`、`Cookie` 和 `Set-Cookie`
- Token 指纹由 `hmac_secret` 决定，轮换密钥后指纹会变化

---

## Tech Stack

- **Language:** Go 1.23
- **Database:** PostgreSQL 16+
- **Dependencies:** pgx/v5 (PostgreSQL driver), golang.org/x/crypto (bcrypt), gopkg.in/yaml.v3
- **Container:** Docker + Docker Compose

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=mrly2333/token-audit&type=Date)](https://star-history.com/#mrly2333/token-audit&Date)

## License

[MIT](LICENSE)
