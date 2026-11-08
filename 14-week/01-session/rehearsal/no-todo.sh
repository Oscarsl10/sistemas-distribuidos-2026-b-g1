#!/usr/bin/env bash
# Fails when a Go file still has a placeholder marker. The gate blocks, it does not warn.
set -u
if grep -rn -E "TODO|FIXME|PLACEHOLDER" --include='*.go' . ; then
  echo "no-todo: FAIL, placeholder markers found above"
  exit 1
fi
echo "no-todo: OK, no placeholder markers in *.go files"
