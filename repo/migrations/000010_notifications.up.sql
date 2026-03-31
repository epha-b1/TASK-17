CREATE TABLE IF NOT EXISTS notification_topics (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notification_subscriptions (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id uuid NOT NULL REFERENCES notification_topics(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, topic_id)
);

CREATE TABLE IF NOT EXISTS user_dnd_settings (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    start_time time NOT NULL DEFAULT '22:00',
    end_time time NOT NULL DEFAULT '06:00',
    enabled boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id uuid REFERENCES notification_topics(id) ON DELETE SET NULL,
    title text NOT NULL,
    body text NOT NULL,
    read boolean NOT NULL DEFAULT false,
    dismissed boolean NOT NULL DEFAULT false,
    booking_id uuid REFERENCES reservations(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    read_at timestamptz,
    dismissed_at timestamptz
);

CREATE TABLE IF NOT EXISTS notification_jobs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id uuid NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    booking_id uuid REFERENCES reservations(id) ON DELETE SET NULL,
    topic_id uuid REFERENCES notification_topics(id) ON DELETE SET NULL,
    channel text NOT NULL DEFAULT 'in_app' CHECK (channel IN ('in_app', 'sms', 'email')),
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'delivered', 'failed', 'suppressed', 'deferred')),
    attempt_count integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz,
    last_error text,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    delivered_at timestamptz,
    downloaded_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notification_jobs_status_next_attempt ON notification_jobs(status, next_attempt_at);
CREATE INDEX IF NOT EXISTS idx_notification_jobs_booking_user_created ON notification_jobs(booking_id, user_id, created_at DESC);

INSERT INTO notification_topics(name)
VALUES ('booking_success'), ('booking_changed'), ('expiry_approaching'), ('arrears_reminder')
ON CONFLICT (name) DO NOTHING;
