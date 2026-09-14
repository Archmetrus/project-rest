#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p bin
for command in migrate user-service hr-service gateway; do
    go build -o "bin/$command" "./cmd/$command"
done
./bin/migrate user up
./bin/migrate hr up
pids=()
cleanup() {
    trap - EXIT INT TERM
    for pid in "${pids[@]}"; do
        kill -TERM "$pid" 2>/dev/null || true
    done
    for pid in "${pids[@]}"; do
        wait "$pid" 2>/dev/null || true
    done
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
for service in user-service hr-service gateway; do
    "./bin/$service" &
    pids+=("$!")
done
# If any service exits, clean up its siblings as well.
wait -n "${pids[@]}"
