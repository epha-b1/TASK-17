CREATE TABLE IF NOT EXISTS exceptions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id uuid NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    exception_type text NOT NULL CHECK (exception_type IN ('gate_stuck', 'sensor_offline', 'camera_error')),
    status text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'acknowledged')),
    acknowledged_by uuid REFERENCES users(id) ON DELETE SET NULL,
    acknowledged_at timestamptz,
    note text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_exceptions_status_created_at ON exceptions(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_exceptions_device_id ON exceptions(device_id);
