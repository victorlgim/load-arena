#!/usr/bin/env bash
set -euo pipefail

URL="${1:-http://localhost:8080/cpu?n=60000}"
N="${N:-10000}"
C="${C:-200}"

echo "Running hey: n=$N c=$C url=$URL"
hey -n "$N" -c "$C" "$URL"
