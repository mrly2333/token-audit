# 部署指南

本文说明 `token-audit` 当前支持的部署方式，以及每种方式背后的运行假设。

## 架构说明

程序只监听一个端口，同时提供两类能力：

- 挂载在 `web_base_path` 下的审计后台，例如 `/audit`
- 其他路径上的透明代理转发能力

典型调用链路如下：

1. 你的客户端把 OpenAI 兼容请求发送到 `token-audit`。
2. `token-audit` 将请求转发到配置好的 `upstream_base`。
3. 命中采集路径的请求会被写入 PostgreSQL。
4. 运维或审计人员通过内置后台查看日志和统计信息。

## 环境要求

- PostgreSQL 16 或更高版本
- 一个已经运行的 OpenAI 兼容上游网关
- 一台能够同时访问 PostgreSQL 和上游网关的部署主机
- 如果要本地编译源码，需要 Go 1.23 及以上版本
- 如果使用容器部署，需要 Docker 和 Docker Compose

## 配置项说明

`config.example.yaml` 和 `config.docker.yaml` 中包含当前支持的配置项：

- `listen_addr`：代理和审计后台监听的地址
- `upstream_base`：要转发到的上游网关基础地址
- `web_base_path`：审计后台路径前缀，通常使用 `/audit`
- `postgres_dsn`：用于迁移和运行时存储的 PostgreSQL DSN
- `hmac_secret`：用于对 Bearer Token 生成指纹的密钥，不直接保存原始令牌
- `admin_username`：内置后台管理员用户名
- `admin_password_hash`：内置后台管理员密码对应的 bcrypt 哈希
- `max_body_bytes`：单次请求或响应最多保留多少字节的内容
- `capture_paths`：需要完整审计的上游路径列表

## 生成管理员密码哈希

如果本机有 Go：

```bash
go run ./scripts/hash_password.go "<你的强密码>"
```

如果你想用 Docker 临时生成：

```bash
docker run --rm -v "$PWD:/src" -w /src golang:1.23 \
  go run ./scripts/hash_password.go "<你的强密码>"
```

把生成出的 bcrypt 哈希填入 `admin_password_hash`。

## 方式一：Docker Compose

这是单机部署时最省事的方案。

1. 编辑 `config.docker.yaml`。
2. 设置真实的 `hmac_secret`。
3. 生成并填入有效的 `admin_password_hash`。
4. 把 `upstream_base` 指向你已经运行的上游网关。
5. 启动服务：

```bash
docker compose up -d --build
```

该编排默认会启动：

- `postgres`：本地 PostgreSQL 容器
- `audit-proxy`：编译后的审计代理容器

默认访问地址：

- 代理和审计服务：`http://127.0.0.1:3007`
- 审计后台：`http://127.0.0.1:3007/audit`
- PostgreSQL：`127.0.0.1:5432`

## 方式二：二进制 + systemd

如果你准备把它部署到 Linux 服务器，并且更偏向原生服务方式，可以使用这一套。

1. 安装 Go 1.23 或更高版本。
2. 编译程序：

```bash
go build -o bin/newapi-audit-proxy ./cmd/newapi-audit-proxy
```

3. 将 `config.example.yaml` 复制为 `config.yaml`，并填入真实配置。
4. 安装并启动服务：

```bash
bash audit.sh install
```

常用服务命令：

```bash
bash audit.sh status
bash audit.sh restart
bash audit.sh log
```

`audit.sh` 会在当前项目目录基础上安装一个 `systemd` 服务单元并启动它。

## 反向代理与对外暴露建议

如果你要把审计后台暴露给更多人访问，建议前面再加一层反向代理，并限制访问来源。推荐至少做到：

- 只允许可信 IP 访问
- 为 `/audit` 增加 SSO、VPN 或反向代理鉴权
- 在反向代理层终止 TLS
- 不要把 PostgreSQL 直接暴露到公网

## 运维注意事项

- 程序启动时会自动执行数据库迁移。
- 对命中采集路径的请求，会保存脱敏后的提示词和响应内容。
- Token 指纹由 `hmac_secret` 决定，同一个密钥下结果稳定；如果你轮换密钥，后续新请求的指纹会发生变化。
- 对流式响应的采集属于尽力而为模式，会根据 SSE 数据重建可展示文本和用量摘要。

## 故障排查

如果服务能启动，但后台无法登录：

- 检查 `admin_password_hash` 是否是真实有效的 bcrypt 哈希
- 检查浏览器访问路径是否与 YAML 中配置的 `web_base_path` 一致

如果请求能转发，但后台没有日志：

- 检查请求路径是否包含在 `capture_paths` 中
- 检查应用进程是否能访问 PostgreSQL
- 查看服务日志中是否有迁移失败或插入失败的报错

如果客户端收到上游错误：

- 检查 `upstream_base` 是否填写正确
- 检查上游网关是否真的支持你正在转发的路径
- 检查部署主机是否能够访问上游网关
