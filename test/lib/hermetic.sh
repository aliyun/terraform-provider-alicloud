#!/usr/bin/env bash
# test/lib/hermetic.sh — isolate a shell test from the live Jarvis control plane.
#
# WHY THIS EXISTS (Aone 85192197)
#
# bootstrap/claim.sh and bootstrap/wrap.sh decide "am I inside an interactive
# Claude/Codex session?" from four environment variables, and when the answer is
# yes they claim/write through the REAL control plane:
#
#   bootstrap/lib.sh:90   jarvis_interactive_context()
#   bootstrap/claim.sh:516 _is_interactive_context()   (a second copy of the same rule)
#
# Claude Code and Codex export those variables into every tool subprocess, so a
# test inherits them and its `claim.sh claim <fake-id>` becomes a real claim.
# test/claim_test.sh used to unset only two of the four, which minted real
# control-plane Tasks for ids the test made up:
#
#   aone:1086837:9001   task 1130, generation 5   (claim_test.sh WORKITEM_ID)
#   aone:1086837:771…7712  tasks 6304/6007/6041/6335/6042/4451/1135/6989
#                          (claim_test.sh's `claim "77${i}"` concurrency case)
#
# Five of those reached FAILED_FINAL after the reaper burned retry 4/3, and each
# retry sent a real `tag jarvis-claimed` write at Aone. One stole the developer
# session's only claim slot, which armed the PreToolUse fence and blocked
# unrelated tool calls until it was suspended by hand.
#
# USAGE — near the top of any test that executes claim.sh / wrap.sh /
# triage-one.sh, after its temp dirs exist:
#
#   source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/hermetic.sh"
#   jarvis_test_hermetic_init "$my_tmp_dir"
#
# Pass a directory the test already cleans up; the tripwire stub is written
# there. Omitting it falls back to a fresh mktemp -d that the test must clean.
#
# For a test that DELIBERATELY exercises the interactive path with its own stub
# runner (test/interactive_claim_test.sh, test/receipted_wrap_test.sh), call
# jarvis_test_hermetic_pin_control_plane on its own and install the stub after.
#
# test/hermetic_isolation_test.sh enforces both halves: that every test executing
# those scripts is wired up here, and that this file's signal list still covers
# every variable the two context helpers read.

# Every variable that makes jarvis_interactive_context()/_is_interactive_context()
# report "interactive". Keep in sync with bootstrap/lib.sh and bootstrap/claim.sh —
# test/hermetic_isolation_test.sh fails if a new one appears there but not here.
JARVIS_TEST_INTERACTIVE_SIGNALS="CLAUDE_CODE_SESSION_ID CODEX_THREAD_ID JARVIS_INTERACTIVE_CLIENT JARVIS_INTERACTIVE_SESSION_ID"

# Discard port: connections fail immediately instead of reaching production.
JARVIS_TEST_CONTROL_PLANE_SENTINEL="http://127.0.0.1:9"

jarvis_test_hermetic_signals() {
    printf '%s\n' $JARVIS_TEST_INTERACTIVE_SIGNALS
}

# Make the interactive helpers report "not interactive".
jarvis_test_hermetic_isolate_signals() {
    local name
    for name in $JARVIS_TEST_INTERACTIVE_SIGNALS; do
        unset "$name"
    done
}

# Pin the control-plane endpoint at the discard port. bootstrap/runtime-config.sh
# replays the caller's JARVIS_* environment last (see its saved_file handling), so
# this survives the gitignored env files and wins over the real endpoint. Only
# non-empty values are preserved, hence a non-empty sentinel.
jarvis_test_hermetic_pin_control_plane() {
    export JARVIS_CONTROL_PLANE_BASE_URL="$JARVIS_TEST_CONTROL_PLANE_SENTINEL"
}

# A runner that fails loudly instead of quietly reaching the real worker. Tests
# wired up here should never invoke it: with the signals isolated, claim.sh and
# wrap.sh skip the interactive path entirely. If it ever fires, a code path went
# interactive behind the test's back — exactly the regression to catch.
jarvis_test_hermetic_tripwire_runner() {
    local dir="${1:-}"
    if [ -z "$dir" ]; then
        dir="$(mktemp -d)" || return 1
    fi
    local stub="$dir/hermetic-tripwire-runner.sh"
    cat > "$stub" <<'TRIPWIRE'
#!/usr/bin/env bash
echo "hermetic tripwire: this test reached the interactive worker CLI (args: $*)." >&2
echo "hermetic tripwire: an interactive context leaked in, so claim.sh/wrap.sh" >&2
echo "hermetic tripwire: would have written to the live control plane." >&2
echo "hermetic tripwire: see test/lib/hermetic.sh and Aone 85192197." >&2
exit 97
TRIPWIRE
    chmod +x "$stub" || return 1
    export JARVIS_INTERACTIVE_WORKER_RUNNER="$stub"
    printf '%s\n' "$stub"
}

jarvis_test_hermetic_init() {
    jarvis_test_hermetic_isolate_signals
    jarvis_test_hermetic_pin_control_plane
    jarvis_test_hermetic_tripwire_runner "${1:-}" >/dev/null
}
