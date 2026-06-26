#!/usr/bin/env bash
# bootstrap/watchdog.sh — escalate tasks whose owner instance is dead
#
# For each task reported by `coord.sh list-orphans`, writes
# escalation/<aid>.md via log.sh so a human (or adopt flow) can decide.
# Does NOT trigger triage or any Aone write operations.
#
# Usage: bash bootstrap/watchdog.sh
#
# Environment overrides (for testing):
#   JARVIS_ROOT           — repo root (default: git rev-parse --show-toplevel)
#   JARVIS_ESCALATION_DIR — escalation dir (default: <repo_root>/escalation)

set -uo pipefail

# ---------------------------------------------------------------------------
# Resolve JARVIS_ROOT
# ---------------------------------------------------------------------------
if [ -z "${JARVIS_ROOT:-}" ]; then
    JARVIS_ROOT="$(git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel 2>/dev/null \
        || git rev-parse --show-toplevel)"
fi
export JARVIS_ROOT

# ---------------------------------------------------------------------------
# Source log.sh (provides escalate function) — same pattern as sweep.sh
# ---------------------------------------------------------------------------
_watchdog_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/log.sh
source "$_watchdog_dir/log.sh"

# ---------------------------------------------------------------------------
# Escalate each orphaned task
# ---------------------------------------------------------------------------
while IFS= read -r aid; do
    [ -z "$aid" ] && continue
    escalate "$aid" "owner dead, awaiting adopt"
done < <(bash "$_watchdog_dir/coord.sh" list-orphans)
