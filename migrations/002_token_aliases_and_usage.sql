ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS prompt_tokens BIGINT NOT NULL DEFAULT 0;

ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS completion_tokens BIGINT NOT NULL DEFAULT 0;

ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS total_tokens BIGINT NOT NULL DEFAULT 0;

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
