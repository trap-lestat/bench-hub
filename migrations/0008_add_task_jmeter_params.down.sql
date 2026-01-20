DROP INDEX IF EXISTS idx_locust_tasks_target_host;
ALTER TABLE locust_tasks
DROP COLUMN IF EXISTS target_host,
DROP COLUMN IF EXISTS jmeter_tpm;
