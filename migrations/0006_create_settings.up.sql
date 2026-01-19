CREATE TABLE IF NOT EXISTS app_settings (
    key varchar(64) PRIMARY KEY,
    value text NOT NULL,
    updated_at timestamp NOT NULL DEFAULT now()
);
