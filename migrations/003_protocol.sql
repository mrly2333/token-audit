ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS protocol TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_audit_logs_protocol
    ON audit_logs (protocol);
