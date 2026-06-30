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
_watchdog_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$_watchdog_dir/lib.sh"
JARVIS_ROOT="$(jarvis_root)"
export JARVIS_ROOT
# shellcheck source=bootstrap/log.sh
source "$_watchdog_dir/log.sh"

# ---------------------------------------------------------------------------
# Escalate each orphaned task — skip if already escalated (dedup)
# ---------------------------------------------------------------------------
_esc_dir="${JARVIS_ESCALATION_DIR:-$JARVIS_ROOT/escalation}"
while IFS= read -r aid; do
    [ -z "$aid" ] && continue
    # Only escalate once per aid; if file already exists the task is known.
    [ -f "$_esc_dir/$aid.md" ] && continue
    escalate "$aid" "owner dead, awaiting adopt"
done < <(bash "$_watchdog_dir/coord.sh" list-orphans)
