CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW(),
    username text NOT NULL UNIQUE,
    firstname text NOT NULL,
    lastname text NOT NULL,
    password_hash bytea NOT NULL,
    version integer NOT NULL DEFAULT 1
);