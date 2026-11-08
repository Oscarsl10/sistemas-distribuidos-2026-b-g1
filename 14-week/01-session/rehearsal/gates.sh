#!/usr/bin/env bash
# The quality gates of the pipeline, in the order of the pipeline. The first red gate stops the run.
set -eu
cd "$(dirname "$0")"

echo "== lint: formatting"
unformatted="$(gofmt -l .)"
if [ -n "$unformatted" ]; then echo "lint: FAIL, not formatted: $unformatted"; exit 1; fi
echo "lint: OK, gofmt reports 0 files"

echo "== typecheck: vet and compile"
go vet ./...
go build ./...
echo "typecheck: OK"

echo "== test: unit tests"
out="$(go test -count=1 -v ./... 2>&1)" || { echo "$out"; echo "test: FAIL"; exit 1; }
echo "$out" | grep -E '^(--- PASS|--- FAIL)' || true
echo "test: OK, $(echo "$out" | grep -c -- '--- PASS') PASS lines, $(echo "$out" | grep -c -- '--- FAIL') FAIL lines"

echo "== coverage-check: measured must be >= declared floor"
cov="$(go test -count=1 -cover ./... | awk '/coverage:/ {gsub("%","",$5); print $5}' | sort -n | head -1)"
floor="$(cat coverage.floor)"
awk -v m="$cov" -v f="$floor" 'BEGIN { if (m + 0 < f + 0) { printf "coverage-check: FAIL, measured %s%% < floor %s%%\n", m, f; exit 1 } printf "coverage-check: OK, measured %s%% >= floor %s%%\n", m, f }'

echo "== no-todo: no placeholders left in the code (the dependency, secret and SAST scans are not part of this rehearsal)"
./no-todo.sh

echo "== all gates green: the change is releasable"
