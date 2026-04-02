CREATE TABLE IF NOT EXISTS tags (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS member_tags (
    member_id uuid NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    assigned_at timestamptz NOT NULL DEFAULT now(),
    assigned_by uuid REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (member_id, tag_id)
);

CREATE TABLE IF NOT EXISTS segment_definitions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    filter_expression jsonb NOT NULL DEFAULT '{}',
    schedule text NOT NULL DEFAULT 'manual' CHECK (schedule IN ('manual', 'nightly')),
    created_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS segment_runs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    segment_id uuid NOT NULL REFERENCES segment_definitions(id) ON DELETE CASCADE,
    ran_at timestamptz NOT NULL DEFAULT now(),
    member_count integer NOT NULL DEFAULT 0,
    triggered_by text NOT NULL DEFAULT 'manual' CHECK (triggered_by IN ('manual', 'scheduler'))
);

CREATE TABLE IF NOT EXISTS tag_versions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exported_by uuid REFERENCES users(id) ON DELETE SET NULL,
    exported_at timestamptz NOT NULL DEFAULT now(),
    snapshot jsonb NOT NULL DEFAULT '{}',
    imported_at timestamptz,
    imported_by uuid REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_member_tags_member ON member_tags(member_id);
CREATE INDEX IF NOT EXISTS idx_member_tags_tag ON member_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_segment_runs_segment ON segment_runs(segment_id, ran_at DESC);
