DROP INDEX IF EXISTS idx_device_server_assigned_ip;
ALTER TABLE "device" ADD CONSTRAINT device_assigned_ip_key UNIQUE (assigned_ip);

ALTER TABLE "server" DROP COLUMN IF EXISTS subnet;
