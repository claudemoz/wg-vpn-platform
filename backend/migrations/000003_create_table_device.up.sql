CREATE TABLE IF NOT EXISTS "device" (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID         NOT NULL REFERENCES "user" (id),
    server_id       UUID         NOT NULL REFERENCES "server" (id),
    name            VARCHAR(100) NOT NULL,
    public_key      VARCHAR(64)  NOT NULL UNIQUE,
    assigned_ip     VARCHAR(45)  NOT NULL UNIQUE,
    last_handshake  TIMESTAMPTZ,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_device_user_id ON "device" (user_id);
CREATE INDEX IF NOT EXISTS idx_device_server_id ON "device" (server_id);
CREATE INDEX IF NOT EXISTS idx_device_deleted_at ON "device" (deleted_at);
CREATE INDEX IF NOT EXISTS idx_device_revoked_at ON "device" (revoked_at);
