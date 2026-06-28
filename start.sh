#!/usr/bin/env bash
# Start the Mosaic App frontend and backend together.

set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

PORT="${SERVER_PORT:-8080}"
PIDS=()

cleanup() {
  trap - INT TERM
  echo ""
  echo "Stopping services..."
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
  exit 0
}

trap cleanup INT TERM

if [ ! -d frontend/node_modules ]; then
  echo "Installing frontend dependencies..."
  (cd frontend && npm install)
fi

echo "Starting backend on http://localhost:${PORT}..."
go run . &
PIDS+=($!)

echo "Starting frontend on http://localhost:3000..."
(cd frontend && BROWSER=none npm start) &
PIDS+=($!)

echo ""
echo "Mosaic App is running:"
echo "  App:      http://localhost:3000"
echo "  API:      http://localhost:${PORT}/api/health"
echo ""
echo "Press Ctrl+C to stop both services."

wait
