CREATE TABLE IF NOT EXISTS accounts (
    user_id bigserial PRIMARY KEY REFERENCES users ON DELETE CASCADE,
    balance integer DEFAULT 0 NOT NULL,
    created_at timestamp(0) WITH time zone DEFAULT NOW() NOT NULL
)