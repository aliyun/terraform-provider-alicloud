#!/usr/bin/env bash
# Resolve and run the project Stop hook. This wrapper is intentionally tiny so
# Claude and Codex hook configs can call the same repo-local entrypoint.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wrap_check="$script_dir/wrap-check.sh"

if [ ! -f "$wrap_check" ]; then
    echo "stop-hook: wrap-check.sh not found at $wrap_check" >&2
    exit 127
fi

exec bash "$wrap_check"
