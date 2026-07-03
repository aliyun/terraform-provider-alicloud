#!/usr/bin/env bash
# bootstrap/triage-one.sh — per-item orchestration bookend: claim → done → release
#
# Called AFTER the subagent has finished actual triage work, with the finished
# summary and target status. Codifies the claim→work→done→release sequence so the
# bookend is enforced by the script rather than left to the model.
#
# Usage:
#   triage-one.sh <id> <pool> <project> <summary> <status>
#
#   id       — Aone work item identifier (e.g. WI-1234)
#   pool     — pool name (informational; passed through for future use)
#   project  — Aone project id for claim/release scoping
#   summary  — finished triage summary (passed to wrap.sh done)
#   status   — target Aone status to set (mandatory; passed to wrap.sh done)
#
# Sequence:
#   1. claim.sh claim <id> <project>
#      → exit 1 (lost race): print SKIP and exit 0
#   2. wrap.sh done <id> <summary> <status>
#      → failure: log.sh escalate <id> "done failed"; do NOT release; exit 1
#   3. claim.sh release <id> <project>
#   4. print "DONE: <id>"; exit 0
#
# Override paths for testing:
#   TRIAGE_CLAIM_CMD   — defaults to co-located claim.sh
#   TRIAGE_WRAP_CMD    — defaults to co-located wrap.sh
#   TRIAGE_LOG_CMD     — defaults to co-located log.sh
#
# Respects JARVIS_ROOT env override.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------------------------------------------------------------------------
# Arg validation
# ---------------------------------------------------------------------------
if [ "${#}" -ne 5 ]; then
    echo "Usage: triage-one.sh <id> <pool> <project> <summary> <status>" >&2
    exit 1
fi

id="$1"
pool="$2"
project="$3"
summary="$4"
status="$5"

# ---------------------------------------------------------------------------
# Resolve helper commands (co-located defaults, overridable for tests)
# ---------------------------------------------------------------------------
CLAIM_CMD="${TRIAGE_CLAIM_CMD:-$script_dir/claim.sh}"
WRAP_CMD="${TRIAGE_WRAP_CMD:-$script_dir/wrap.sh}"
LOG_CMD="${TRIAGE_LOG_CMD:-$script_dir/log.sh}"
COORD_CMD="${TRIAGE_COORD_CMD:-$script_dir/coord.sh}"

# ---------------------------------------------------------------------------
# Register this triage instance for coordination tracking so checkpoints
# have a real owner_instance and crashes leave a resumable record.
# ---------------------------------------------------------------------------
# Pass $$ (this bookend's real pid) so coord.sh dead can check liveness via kill -0;
# without it register embeds coord.sh's own short-lived pid and the instance reads dead
# immediately (see coord.sh register comment).
COORD_ID=$(bash "$COORD_CMD" register triage "$$" 2>/dev/null || true)
export COORD_ID

# ---------------------------------------------------------------------------
# Step 1: Claim — lost race → SKIP (not an error)
# ---------------------------------------------------------------------------
if ! "$CLAIM_CMD" claim "$id" "$project"; then
    echo "SKIP: $id"
    exit 0
fi

# Coord: mark instance as having claimed this item, with worktree/branch/repo
# context so a crash leaves a resumable record.
bash "$COORD_CMD" checkpoint "$id" claimed "$(pwd)" "$(git rev-parse --abbrev-ref HEAD 2>/dev/null)" "$(basename "$(pwd)")" || true

# ---------------------------------------------------------------------------
# Step 2: wrap.sh done — failure → escalate + exit 1 (no release)
# ---------------------------------------------------------------------------
if ! "$WRAP_CMD" done "$id" "$summary" "$status"; then
    "$LOG_CMD" escalate "$id" "done failed"
    exit 1
fi

# ---------------------------------------------------------------------------
# Step 3: Release
# ---------------------------------------------------------------------------
# Coord: mark item as done before releasing the claim (non-blocking)
bash "$COORD_CMD" checkpoint "$id" done || true
"$CLAIM_CMD" release "$id" "$project"

# ---------------------------------------------------------------------------
# Step 4: Done
# ---------------------------------------------------------------------------
echo "DONE: $id"
exit 0
