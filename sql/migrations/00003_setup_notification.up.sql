begin;
CREATE TABLE notifications (
    id uuid primary key not null default gen_random_uuid(),
    user_id UUID not null,
    message text not null,
    read boolean DEFAULT FALSE,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);
-- Index for faster lookup by user_id
CREATE INDEX idx_notifications_user_id ON notifications (user_id);
commit;