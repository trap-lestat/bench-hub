ALTER TABLE locust_tasks
ADD COLUMN IF NOT EXISTS target_host varchar(255),
ADD COLUMN IF NOT EXISTS jmeter_tpm integer;

CREATE INDEX IF NOT EXISTS idx_locust_tasks_target_host ON locust_tasks (target_host);
