CREATE TABLE IF NOT EXISTS locust_reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id uuid REFERENCES locust_tasks(id) ON DELETE SET NULL,
    name varchar(128) NOT NULL,
    report_type varchar(32) NOT NULL,
    file_path text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_locust_reports_created_at ON locust_reports (created_at);
