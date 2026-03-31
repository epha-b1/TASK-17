CREATE TABLE IF NOT EXISTS capacity_snapshots (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id uuid NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    snapshot_at timestamptz NOT NULL DEFAULT now(),
    authoritative_stalls integer NOT NULL CHECK (authoritative_stalls >= 0)
);

CREATE TABLE IF NOT EXISTS reconciliation_runs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    run_at timestamptz NOT NULL DEFAULT now(),
    zones_checked integer NOT NULL DEFAULT 0,
    discrepancies_found integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS compensating_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reconciliation_run_id uuid NOT NULL REFERENCES reconciliation_runs(id) ON DELETE CASCADE,
    zone_id uuid NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    event_type text NOT NULL CHECK (event_type IN ('hold', 'release')),
    stall_count integer NOT NULL CHECK (stall_count > 0),
    detail jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reconciliation_runs_run_at ON reconciliation_runs(run_at DESC);
CREATE INDEX IF NOT EXISTS idx_compensating_events_run ON compensating_events(reconciliation_run_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_compensating_events_zone ON compensating_events(zone_id, created_at DESC);
