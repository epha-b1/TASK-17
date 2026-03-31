ALTER TABLE reconciliation_runs
    ADD COLUMN IF NOT EXISTS run_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS started_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS completed_at timestamptz,
    ADD COLUMN IF NOT EXISTS triggered_by text NOT NULL DEFAULT 'scheduler',
    ADD COLUMN IF NOT EXISTS zones_checked integer NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS discrepancies_found integer NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS created_at timestamptz NOT NULL DEFAULT now();

ALTER TABLE compensating_events
    ADD COLUMN IF NOT EXISTS reconciliation_run_id uuid REFERENCES reconciliation_runs(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS zone_id uuid REFERENCES zones(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS event_type text,
    ADD COLUMN IF NOT EXISTS quantity integer,
    ADD COLUMN IF NOT EXISTS event_derived_stalls integer,
    ADD COLUMN IF NOT EXISTS authoritative_stalls integer,
    ADD COLUMN IF NOT EXISTS delta integer,
    ADD COLUMN IF NOT EXISTS reason text,
    ADD COLUMN IF NOT EXISTS stall_count integer,
    ADD COLUMN IF NOT EXISTS detail jsonb NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS created_at timestamptz NOT NULL DEFAULT now();

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='compensating_events' AND column_name='event_type'
    ) THEN
        BEGIN
            ALTER TABLE compensating_events
                DROP CONSTRAINT IF EXISTS compensating_events_event_type_check,
                ADD CONSTRAINT compensating_events_event_type_check
                    CHECK (event_type IN ('hold', 'release'));
        EXCEPTION WHEN undefined_column THEN
            NULL;
        END;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='reconciliation_runs' AND column_name='triggered_by'
    ) THEN
        BEGIN
            ALTER TABLE reconciliation_runs
                DROP CONSTRAINT IF EXISTS reconciliation_runs_triggered_by_check,
                ADD CONSTRAINT reconciliation_runs_triggered_by_check
                    CHECK (triggered_by IN ('manual', 'scheduler', 'test'));
        EXCEPTION WHEN undefined_column THEN
            NULL;
        END;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_reconciliation_runs_run_at ON reconciliation_runs(run_at DESC);
CREATE INDEX IF NOT EXISTS idx_compensating_events_run ON compensating_events(reconciliation_run_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_compensating_events_zone ON compensating_events(zone_id, created_at DESC);
