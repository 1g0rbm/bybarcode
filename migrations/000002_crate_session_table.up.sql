CREATE TABLE IF NOT EXISTS sessions
(
    id            UUID PRIMARY KEY,
    token         UUID NOT NULL,
    refresh_token UUID NOT NULL,
    account_id    BIGINT REFERENCES account (id) ON DELETE CASCADE,
    expire_at     TIMESTAMP WITH TIME ZONE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'::text),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'::text)
);
