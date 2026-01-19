CREATE TABLE IF NOT EXISTS locust_scripts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(128) UNIQUE NOT NULL,
    description text,
    content text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_locust_scripts_created_at ON locust_scripts (created_at);
