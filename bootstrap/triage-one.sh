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
#      → exit 3 (missing required field): print legal candidates, exit 1
#      → other non-zero: report claim failure, exit 1
#   2. wrap.sh done <id> <summary> <status>
#      → failure: do NOT release; exit 1
#   3. claim.sh release <id> <project>
#   4. print "DONE: <id>"; exit 0
#
# Override paths for testing:
#   TRIAGE_CLAIM_CMD   — defaults to co-located claim.sh
#   TRIAGE_WRAP_CMD    — defaults to co-located wrap.sh
#   TRIAGE_FIELDS_CMD  — defaults to co-located aone-fields.sh
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
FIELDS_CMD="${TRIAGE_FIELDS_CMD:-$script_dir/aone-fields.sh}"

# ---------------------------------------------------------------------------
# Step 1: Claim — lost race → SKIP; missing required fields → print candidates
# ---------------------------------------------------------------------------
"$CLAIM_CMD" claim "$id" "$project"
claim_rc=$?
if [ "$claim_rc" -eq 3 ]; then
    if ! missing_json="$(bash "$FIELDS_CMD" missing "$id")"; then
        missing_json='[]'
    fi
    echo "MISSING_REQUIRED_FIELDS: $id $missing_json" >&2
    exit 1
elif [ "$claim_rc" -eq 1 ]; then
    echo "SKIP: $id"
    exit 0
elif [ "$claim_rc" -ne 0 ]; then
    echo "ERROR: claim failed for $id (rc=$claim_rc)" >&2
    exit 1
fi

# ---------------------------------------------------------------------------
# Step 2: wrap.sh done — failure → exit 1 (no release)
# ---------------------------------------------------------------------------
if ! "$WRAP_CMD" done "$id" "$summary" "$status"; then
    echo "ERROR: done failed for $id; claim remains held" >&2
    exit 1
fi

# ---------------------------------------------------------------------------
# Step 3: Release
# ---------------------------------------------------------------------------
"$CLAIM_CMD" release "$id" "$project"

# ---------------------------------------------------------------------------
# Step 4: Done
# ---------------------------------------------------------------------------
echo "DONE: $id"
exit 0
