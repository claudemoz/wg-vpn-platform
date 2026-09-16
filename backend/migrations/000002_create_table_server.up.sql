CREATE TABLE IF NOT EXISTS "server" (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    region          VARCHAR(10)  NOT NULL,
    name            VARCHAR(100) NOT NULL,
    public_endpoint VARCHAR(255) NOT NULL,
    wg_public_key   VARCHAR(64)  NOT NULL UNIQUE,
    grpc_endpoint   VARCHAR(255) NOT NULL,
    max_peers       INT          NOT NULL DEFAULT 100,
    current_peers   INT          NOT NULL DEFAULT 0,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_server_region ON "server" (region);
CREATE INDEX IF NOT EXISTS idx_server_is_active ON "server" (is_active);
CREATE INDEX IF NOT EXISTS idx_server_deleted_at ON "server" (deleted_at);
