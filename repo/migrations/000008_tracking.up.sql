CREATE TABLE IF NOT EXISTS vehicle_positions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id uuid NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    received_at timestamptz NOT NULL DEFAULT now(),
    device_time timestamptz,
    device_time_trusted boolean NOT NULL DEFAULT false,
    suspect boolean NOT NULL DEFAULT false,
    confirmed boolean NOT NULL DEFAULT true,
    discarded_at timestamptz
);

CREATE TABLE IF NOT EXISTS stop_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id uuid NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    started_at timestamptz NOT NULL,
    detected_at timestamptz NOT NULL,
    latitude double precision,
    longitude double precision
);

CREATE INDEX IF NOT EXISTS idx_vehicle_positions_vehicle_received ON vehicle_positions(vehicle_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_vehicle_positions_pending_suspect ON vehicle_positions(vehicle_id, suspect, confirmed, discarded_at);
CREATE INDEX IF NOT EXISTS idx_stop_events_vehicle_detected ON stop_events(vehicle_id, detected_at DESC);
