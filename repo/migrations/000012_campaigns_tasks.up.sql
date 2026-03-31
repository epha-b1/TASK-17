CREATE TABLE IF NOT EXISTS campaigns (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    target_role text,
    created_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id uuid NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    description text NOT NULL,
    deadline timestamptz,
    reminder_interval_minutes integer NOT NULL DEFAULT 60 CHECK (reminder_interval_minutes > 0),
    last_reminder_at timestamptz,
    completed_at timestamptz,
    completed_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_by uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_campaigns_created_at ON campaigns(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_campaign ON tasks(campaign_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_due ON tasks(deadline, completed_at, last_reminder_at);

INSERT INTO notification_topics(name)
VALUES ('task_reminder')
ON CONFLICT (name) DO NOTHING;
