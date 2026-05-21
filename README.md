# token-audit

`token-audit` 是一个透明的 OpenAI 兼容审计代理，内置可视化审计后台。它部署在现有上游 API 前面，原样转发请求，并将脱敏后的请求、响应和用量信息写入 PostgreSQL，方便后续排查、统计和审计。

当前版本主要针对 **NewAPI** 作为上游服务进行过验证。理论上也可能兼容其他 OpenAI 兼容网关，但这些集成目前还没有做正式测试。

## 项目功能

- 代理指定的 OpenAI 风格接口，例如 `/v1/chat/completions`、`/v1/responses`、`/v1/completions`、`/v1/embeddings`
- 记录请求耗时、状态码、Token 用量、模型名称、用户文本、助手文本以及部分原始载荷
- 使用基于 HMAC 的指纹替代真实 API Key，避免把原始令牌直接存进数据库
- 提供审计后台，可查看概览统计、请求检索、令牌别名和数据库维护操作
- 同时支持普通 JSON 响应和流式 SSE 响应

## 适用场景

如果你需要在 OpenAI 兼容网关前加一层轻量审计能力，这个项目适合用来回答下面这些问题：

- 哪个令牌消耗的用量最多？
- 哪些模型的调用流量最高？
- 哪些请求失败、超时或返回了异常内容？
- 某一次请求对应的提示词和响应文本是什么？

## 重要说明

- 这个代理会为审计目的存储经过脱敏处理的请求和响应内容。只有在你的隐私策略和数据留存策略允许记录提示词与回复内容时，才建议部署。
- 审计后台应放在可信网络边界之后，最好配合 VPN、反向代理鉴权或单点登录使用。
- 存储的请求头和响应头中会排除原始 `Authorization`、`Cookie` 和 `Set-Cookie`。
- 当前兼容性与测试边界请查看 [COMPATIBILITY.md](./COMPATIBILITY.md)。

## 快速开始

最快的启动方式是使用 Docker Compose：

1. 编辑 `config.docker.yaml`，把其中占位用的 `hmac_secret` 和 `admin_password_hash` 替换成真实值。
2. 把 `upstream_base` 指向你现有的 NewAPI 或其他 OpenAI 兼容网关。
3. 执行 `docker compose up -d --build`。
4. 打开 `http://127.0.0.1:3007/audit`。

完整部署方式请查看 [DEPLOYMENT.md](./DEPLOYMENT.md)。

## 仓库结构

- `cmd/newapi-audit-proxy`：程序入口
- `internal/proxy`：透明代理与响应流转逻辑
- `internal/audit`：请求解析、内容脱敏、数据存储与统计能力
- `internal/web`：审计后台和相关 API
- `migrations`：PostgreSQL 表结构与迁移脚本
- `scripts/hash_password.go`：生成后台管理员 bcrypt 密码哈希的辅助脚本
- `audit.sh`：Linux `systemd` 部署辅助脚本

## 发布说明

这个仓库当前是整理后的对外发布版本，保留了源码、部署文件和使用文档，移除了本地缓存、编译产物以及机器相关配置。
