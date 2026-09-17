-- Each VPN server gets its own address pool.
ALTER TABLE "server"
    ADD COLUMN IF NOT EXISTS subnet VARCHAR(43) NOT NULL DEFAULT '10.8.0.0/24';

-- An assigned IP only has to be unique within a server, not across the platform.
ALTER TABLE "device" DROP CONSTRAINT IF EXISTS device_assigned_ip_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_device_server_assigned_ip ON "device" (server_id, assigned_ip);
