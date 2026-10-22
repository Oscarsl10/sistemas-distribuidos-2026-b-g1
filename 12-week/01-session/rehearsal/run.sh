#!/usr/bin/env bash
# Rehearsal of session 1: build once, promote the same tagged image, release blue-green without
# downtime, and roll back. Run from this folder with Docker running: ./run.sh
set -euo pipefail
cd "$(dirname "$0")"

REGISTRY=127.0.0.1:5000
IMAGE=$REGISTRY/deploy-demo
EDGE=http://127.0.0.1:8090
LOG=$(mktemp)

say() { printf '\n== %s\n' "$*"; }
switch() {  # switch blue|green: rewrite the one-line config and reload the edge (no restart)
  echo "set \$active_upstream $1:8080;" > edge/active.conf
  docker compose exec -T edge nginx -t >/dev/null 2>&1
  docker compose exec -T edge nginx -s reload >/dev/null 2>&1
}
now_ms() { echo $(( $(date +%s%N) / 1000000 )); }
cleanup() { docker compose down -v >/dev/null 2>&1 || true; docker rm -f rehearsal-registry >/dev/null 2>&1 || true; }
trap cleanup EXIT

say "1. Local registry"
docker rm -f rehearsal-registry >/dev/null 2>&1 || true
docker run -d --name rehearsal-registry -p 127.0.0.1:5000:5000 registry:2 >/dev/null

say "2. Build each version ONCE and push it"
for v in v1.0.0 v1.1.0; do
  docker build -q --build-arg VERSION=$v -t $IMAGE:$v . >/dev/null
  docker push -q $IMAGE:$v >/dev/null
done
BLUE_DIGEST=$(docker image inspect --format '{{index .RepoDigests 0}}' $IMAGE:v1.0.0)
GREEN_DIGEST=$(docker image inspect --format '{{index .RepoDigests 0}}' $IMAGE:v1.1.0)
echo "v1.0.0 -> $BLUE_DIGEST"
echo "v1.1.0 -> $GREEN_DIGEST"

say "3. Promotion: the same digests run in qa and in prod, only the environment changes"
for env in qa prod; do
  export ENVIRONMENT=$env BLUE_IMAGE=$BLUE_DIGEST GREEN_IMAGE=$GREEN_DIGEST
  echo "set \$active_upstream blue:8080;" > edge/active.conf
  docker compose up -d --wait >/dev/null 2>&1
  printf '%s: ' "$env"; curl -s $EDGE/version
  [ "$env" = prod ] || docker compose down -v >/dev/null 2>&1
done
echo "images used in prod:"; docker compose ps --format '  {{.Service}}  {{.Image}}  {{.Status}}'

say "4. Health checks: both colours healthy before any traffic moves"
docker compose ps --format '  {{.Service}}  {{.Status}}'

say "5. Blue-green release under load (a request every 50 ms)"
( while true; do curl -s -o /dev/null -w '%{http_code} ' $EDGE/version; sleep 0.05; done > "$LOG" ) &
LOAD=$!
sleep 2
T0=$(now_ms); switch green
until curl -s $EDGE/version | grep -q v1.1.0; do sleep 0.05; done
T1=$(now_ms)
echo "release to v1.1.0 done in $((T1 - T0)) ms: $(curl -s $EDGE/version)"
sleep 2

say "6. Rollback to the previous version"
R0=$(now_ms); switch blue
until curl -s $EDGE/version | grep -q v1.0.0; do sleep 0.05; done
R1=$(now_ms)
echo "rollback to v1.0.0 done in $((R1 - R0)) ms: $(curl -s $EDGE/version)"
sleep 2
kill $LOAD 2>/dev/null || true; wait $LOAD 2>/dev/null || true

say "7. Result of the load during the release and the rollback"
TOTAL=$(wc -w < "$LOG"); OK=$(tr ' ' '\n' < "$LOG" | grep -c '^200$' || true)
echo "requests: $TOTAL, answered 200: $OK, failed: $((TOTAL - OK))"
rm -f "$LOG"
