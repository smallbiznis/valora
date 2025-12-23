CREATE TABLE IF NOT EXISTS sessions (
    id                  BIGINT PRIMARY KEY,
    user_id             BIGINT NOT NULL,
    session_token_hash  TEXT NOT NULL,
    user_agent          TEXT,
    ip_address          TEXT,
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_token_hash
    ON sessions (session_token_hash);

CREATE INDEX IF NOT EXISTS idx_sessions_user
    ON sessions (user_id);
