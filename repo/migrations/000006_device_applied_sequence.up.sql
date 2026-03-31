ALTER TABLE devices
    ADD COLUMN IF NOT EXISTS last_applied_sequence_number bigint NOT NULL DEFAULT 0;
