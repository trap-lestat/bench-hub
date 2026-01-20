DROP INDEX IF EXISTS idx_locust_scripts_script_type;
ALTER TABLE locust_scripts DROP COLUMN IF EXISTS script_type;
