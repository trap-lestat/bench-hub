#!/usr/bin/env bash
set -euo pipefail

HOST=${HOST:-http://localhost:8080}
USERS=${USERS:-50}
SPAWN_RATE=${SPAWN_RATE:-5}
DURATION=${DURATION:-5m}
OUT_DIR=${OUT_DIR:-reports}

mkdir -p "$OUT_DIR"
TS=$(date +%Y%m%d%H%M%S)

locust -f locust/locustfile.py \
  --headless \
  -u "$USERS" \
  -r "$SPAWN_RATE" \
  --run-time "$DURATION" \
  --host "$HOST" \
  --csv "$OUT_DIR/locust_$TS" \
  --html "$OUT_DIR/locust_$TS.html"

echo "reports generated in $OUT_DIR"
