CREATE TABLE IF NOT EXISTS exports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    requested_by uuid REFERENCES users(id) ON DELETE SET NULL,
    format text NOT NULL CHECK (format IN ('csv', 'excel', 'pdf')),
    scope text NOT NULL CHECK (scope IN ('occupancy', 'bookings', 'exceptions')),
    segment_id uuid REFERENCES segment_definitions(id) ON DELETE SET NULL,
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'ready', 'failed')),
    file_path text,
    query_from timestamptz,
    query_to timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_exports_requested_by ON exports(requested_by, created_at DESC);
