CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar(64) UNIQUE NOT NULL,
    password_hash varchar(128) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);
