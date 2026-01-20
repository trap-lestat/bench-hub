ALTER TABLE locust_scripts
ADD COLUMN IF NOT EXISTS script_type varchar(20) NOT NULL DEFAULT 'locust';

CREATE INDEX IF NOT EXISTS idx_locust_scripts_script_type ON locust_scripts (script_type);
