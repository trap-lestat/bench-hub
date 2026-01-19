CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO users (id, username, password_hash)
VALUES (gen_random_uuid(), 'admin', '$2a$10$7oHAl0cDzsl2RkTE3RzF3.2gRPt1G8mpliQ5xpt..pDyROilG4BIW')
ON CONFLICT (username) DO NOTHING;
