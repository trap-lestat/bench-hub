CREATE TABLE IF NOT EXISTS locust_tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(128) NOT NULL,
    script_id uuid NOT NULL REFERENCES locust_scripts(id) ON DELETE CASCADE,
    users_count integer NOT NULL,
    spawn_rate integer NOT NULL,
    duration_seconds integer NOT NULL,
    status varchar(32) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    started_at timestamp,
    finished_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_locust_tasks_created_at ON locust_tasks (created_at);
CREATE INDEX IF NOT EXISTS idx_locust_tasks_status ON locust_tasks (status);
