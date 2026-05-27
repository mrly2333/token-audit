# 兼容性与验证范围

## 已验证的上游网关

| 上游 | 状态 |
|------|------|
| NewAPI | 已验证 |
| OpenAI API | 兼容（理论支持） |
| Claude API (Anthropic) | 已验证 |

## 支持的 API 端点

### OpenAI 兼容协议

| 端点 | 说明 |
|------|------|
| `/v1/chat/completions` | Chat Completions API |
| `/v1/responses` | Responses API |
| `/v1/completions` | Completions API |
| `/v1/embeddings` | Embeddings API |

### Claude 协议

| 端点 | 说明 |
|------|------|
| `/v1/messages` | Claude Messages API |

## 响应格式支持

- 普通 JSON 响应
- SSE 流式响应（OpenAI `choices[].delta` 和 Claude `content_block_delta`）

## 当前会采集的内容

- 请求耗时和状态码
- 能够识别时的模型名称
- 上游返回中自带的 Token 用量信息
- 符合当前 OpenAI / Claude 风格结构时提取出的用户文本和助手文本
- 在 `max_body_bytes` 限制范围内保存的脱敏请求体与响应体

## 尚未做正式验证的部分

- NewAPI 和 Anthropic 以外的其他网关
- 自定义或厂商私有的载荷格式
- 目前解析逻辑之外的多模态数据格式
- 音频接口和实时接口
- 基于 WebSocket 的传输方式
- 大规模多节点部署场景

## 实际使用预期

如果你的上游行为与标准 OpenAI 或 Claude 的 JSON / SSE 接口足够接近，这个代理大概率可以较小改动甚至直接工作。但只要不是已验证的上游，现阶段都建议先按"未正式验证"处理，再自行完成预发布测试。

## 推荐上线方式

1. 先在预发布或测试环境部署。
2. 通过代理发送几组已知请求。
3. 检查转发结果、状态码、流式输出和 Token 统计是否符合预期。
4. 回到审计后台确认脱敏与解析结果是否满足你的要求。
5. 确认无误后再接入正式流量。
