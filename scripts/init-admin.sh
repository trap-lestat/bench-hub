#!/bin/sh
set -e

ADMIN_USERNAME=${ADMIN_USERNAME:-admin}
# bcrypt hash for "admin123"
ADMIN_PASSWORD_HASH=${ADMIN_PASSWORD_HASH:-'$2a$10$7oHAl0cDzsl2RkTE3RzF3.2gRPt1G8mpliQ5xpt..pDyROilG4BIW'}

DB_NAME=${DB_NAME:-bench_hub}
DB_USER=${DB_USER:-postgres}
DB_PASS=${DB_PASS:-postgres}
DB_HOST=${DB_HOST:-}

SQL=$(cat <<'SQL'
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
INSERT INTO users (id, username, password_hash)
VALUES (gen_random_uuid(), :'admin_user', :'password_hash')
ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash;
SQL
)

if [ -n "$DB_HOST" ]; then
  if ! command -v psql >/dev/null 2>&1; then
    echo "psql not found. Install psql or run without DB_HOST to use docker compose." >&2
    exit 1
  fi
  printf "%s\n" "$SQL" | PGPASSWORD="$DB_PASS" psql -v ON_ERROR_STOP=1 \
    -v admin_user="$ADMIN_USERNAME" \
    -v password_hash="$ADMIN_PASSWORD_HASH" \
    -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME"
else
  printf "%s\n" "$SQL" | docker compose exec -T db psql -v ON_ERROR_STOP=1 \
    -v admin_user="$ADMIN_USERNAME" \
    -v password_hash="$ADMIN_PASSWORD_HASH" \
    -U "$DB_USER" -d "$DB_NAME"
fi

echo "admin user ready: $ADMIN_USERNAME"
