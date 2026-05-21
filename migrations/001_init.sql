CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    query_string TEXT NOT NULL DEFAULT '',
    remote_addr TEXT NOT NULL DEFAULT '',
    request_host TEXT NOT NULL DEFAULT '',
    upstream_base TEXT NOT NULL DEFAULT '',
    status_code INTEGER,
    error_text TEXT,
    is_capture_path BOOLEAN NOT NULL DEFAULT FALSE,
    response_is_sse BOOLEAN NOT NULL DEFAULT FALSE,
    token_fingerprint TEXT,
    token_preview TEXT,
    model TEXT,
    stream BOOLEAN,
    request_headers JSONB NOT NULL DEFAULT '{}'::JSONB,
    response_headers JSONB NOT NULL DEFAULT '{}'::JSONB,
    request_content_type TEXT,
    response_content_type TEXT,
    request_body BYTEA,
    response_body BYTEA,
    request_json JSONB,
    response_json JSONB,
    user_text TEXT,
    assistant_text TEXT,
    usage_json JSONB,
    prompt_tokens BIGINT NOT NULL DEFAULT 0,
    completion_tokens BIGINT NOT NULL DEFAULT 0,
    total_tokens BIGINT NOT NULL DEFAULT 0,
    request_bytes BIGINT NOT NULL DEFAULT 0,
    response_bytes BIGINT NOT NULL DEFAULT 0,
    request_truncated BOOLEAN NOT NULL DEFAULT FALSE,
    response_truncated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_started_at
    ON audit_logs (started_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_token_fingerprint
    ON audit_logs (token_fingerprint);

CREATE INDEX IF NOT EXISTS idx_audit_logs_model
    ON audit_logs (model);

CREATE INDEX IF NOT EXISTS idx_audit_logs_status_code
    ON audit_logs (status_code);

CREATE INDEX IF NOT EXISTS idx_audit_logs_capture_started_at
    ON audit_logs (is_capture_path, started_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_total_tokens
    ON audit_logs (total_tokens DESC);

CREATE TABLE IF NOT EXISTS token_aliases (
    token_fingerprint TEXT PRIMARY KEY,
    token_alias TEXT NOT NULL UNIQUE,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_token_aliases_token_alias
    ON token_aliases (token_alias);
