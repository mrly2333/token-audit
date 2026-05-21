# Compatibility and Validation

This project is currently published with a conservative compatibility statement.

## Upstream officially tested so far

- **NewAPI** as the upstream OpenAI-compatible gateway

## Request types covered by current implementation

- `/v1/chat/completions`
- `/v1/responses`
- `/v1/completions`
- `/v1/embeddings`
- JSON responses
- SSE streaming responses

## What is captured

- request timing and status
- model name when detectable
- token usage when provided by the upstream payload
- user text and assistant text when the payload format matches supported OpenAI-style schemas
- redacted request and response bodies up to the configured `max_body_bytes`

## Not yet formally validated

- gateways other than NewAPI
- custom or vendor-specific payload formats
- multimodal payloads beyond the currently parsed text and image patterns
- audio or realtime interfaces
- WebSocket-based transports
- large-scale multi-node deployments

## Practical expectation

If your upstream behaves like standard OpenAI or NewAPI JSON and SSE APIs, the proxy will likely work with little or no change. Even so, you should treat non-NewAPI deployments as unverified until you run your own staging tests.

## Recommended rollout approach

1. Start in a staging environment.
2. Send known requests through the proxy.
3. Confirm forwarding behavior, status codes, streaming output, and token accounting.
4. Review stored logs to verify redaction and payload parsing meet your expectations.
5. Only then place the proxy in front of production traffic.
