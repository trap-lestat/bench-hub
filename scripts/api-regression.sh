#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${BASE_URL:-http://localhost:8080}
USERNAME=${USERNAME:-admin}
PASSWORD=${PASSWORD:-admin123}

login_payload=$(printf '{"username":"%s","password":"%s"}' "$USERNAME" "$PASSWORD")

auth=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d "$login_payload")

access_token=$(echo "$auth" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
if [ -z "$access_token" ]; then
  echo "login failed"
  echo "$auth"
  exit 1
fi

auth_header="Authorization: Bearer $access_token"

curl -s -X GET "$BASE_URL/api/v1/users" -H "$auth_header" >/dev/null
curl -s -X GET "$BASE_URL/api/v1/items" -H "$auth_header" >/dev/null

echo "api regression ok"
