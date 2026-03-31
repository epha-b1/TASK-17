CREATE TABLE IF NOT EXISTS devices (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid REFERENCES organizations(id) ON DELETE SET NULL,
    device_key text UNIQUE NOT NULL,
    device_type text NOT NULL CHECK (device_type IN ('camera', 'gate', 'geomagnetic')),
    zone_id uuid REFERENCES zones(id) ON DELETE SET NULL,
    status text NOT NULL DEFAULT 'online' CHECK (status IN ('online', 'offline')),
    registered_at timestamptz NOT NULL DEFAULT now(),
    last_applied_sequence_number bigint NOT NULL DEFAULT 0,
    last_sequence_number bigint NOT NULL DEFAULT 0,
    last_event_received_at timestamptz
);

CREATE TABLE IF NOT EXISTS device_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id uuid NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    event_key text UNIQUE NOT NULL,
    sequence_number bigint NOT NULL,
    event_type text NOT NULL,
    payload jsonb,
    received_at timestamptz NOT NULL DEFAULT now(),
    device_time timestamptz,
    device_time_trusted boolean NOT NULL DEFAULT false,
    late boolean NOT NULL DEFAULT false,
    processed boolean NOT NULL DEFAULT false,
    replay_count integer NOT NULL DEFAULT 0,
    replayed_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_devices_zone ON devices(zone_id);
CREATE INDEX IF NOT EXISTS idx_device_events_device_time ON device_events(device_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_device_events_late ON device_events(late);
