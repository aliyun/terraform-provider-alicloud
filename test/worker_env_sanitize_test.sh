#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
sanitizer="$repo_root/bootstrap/worker-env-sanitize.sh"
temporary="$(mktemp -d 2>/dev/null || mktemp -d -t worker-env-test)"
trap 'rm -rf "$temporary"' EXIT

source_file="$temporary/scheduler.env"
worker_file="$temporary/worker.env"
cat >"$source_file" <<'EOF'
JARVIS_CONTROL_PLANE_TOKEN=ordinary-worker-token
JARVIS_CONTROL_PLANE_ADMIN_TOKEN=operator-only-placeholder
export JARVIS_CONTROL_PLANE_ADMIN_TOKEN
JARVIS_CONTROL_PLANE_BASE_URL=https://control.example
EOF

bash "$sanitizer" "$source_file" "$worker_file"

grep -q '^JARVIS_CONTROL_PLANE_TOKEN=ordinary-worker-token$' "$worker_file"
grep -q '^JARVIS_CONTROL_PLANE_BASE_URL=https://control.example$' "$worker_file"
! grep -q 'JARVIS_CONTROL_PLANE_ADMIN_TOKEN' "$worker_file"
! grep -q 'operator-only-placeholder' "$worker_file"

mode="$(stat -f '%Lp' "$worker_file" 2>/dev/null || stat -c '%a' "$worker_file")"
[ "$mode" = "600" ]

echo "worker_env_sanitize_test: PASS"
