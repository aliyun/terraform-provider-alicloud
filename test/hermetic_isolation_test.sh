#!/usr/bin/env bash
# test/hermetic_isolation_test.sh — keep the control plane out of the test suite.
#
# Two failure modes put real Tasks into the live control plane (Aone 85192197),
# and both are invisible when they happen: the test still passes locally, the
# junk shows up days later as FAILED_FINAL rows and stolen claim slots. So both
# are checked mechanically here rather than left to reviewer memory.
#
#   1. A test executes claim.sh / wrap.sh / triage-one.sh without wiring itself to
#      test/lib/hermetic.sh, and inherits the host session's interactive signals.
#   2. bootstrap gains a NEW interactive signal, and hermetic.sh keeps unsetting
#      only the old ones — which is exactly how JARVIS_INTERACTIVE_CLIENT /
#      JARVIS_INTERACTIVE_SESSION_ID slipped past claim_test.sh's two-name unset.
set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
harness="$test_dir/lib/hermetic.sh"

fail=0
no(){ echo "FAIL: $1" >&2; fail=1; }
ok(){ echo "PASS: $1"; }

[ -f "$harness" ] || { echo "FAIL: missing $harness" >&2; exit 1; }

# ── 1. signal-list parity with the two context helpers ───────────────────────
# Extract the variables each helper reads, then require the harness to cover them.
context_signals() {
    local file="$1" fn="$2"
    awk -v fn="$fn" '
        $0 ~ "^" fn "\\(\\)" { inside = 1 }
        inside { print }
        inside && /^}/ { exit }
    ' "$file" | grep -oE '\$\{[A-Z_][A-Z0-9_]*' | sed 's/^\${//' | sort -u
}

harness_signals="$(
    grep -E '^JARVIS_TEST_INTERACTIVE_SIGNALS=' "$harness" \
        | sed 's/^[^=]*=//; s/"//g' | tr ' ' '\n' | grep -v '^$' | sort -u
)"

# The rule lives in two places (bootstrap/lib.sh for wrap.sh, a private copy in
# bootstrap/claim.sh). Identical today; if one gains a signal the other misses,
# one script starts claiming through the control plane while the other does not.
lib_set="$(context_signals "$repo_root/bootstrap/lib.sh" jarvis_interactive_context)"
claim_set="$(context_signals "$repo_root/bootstrap/claim.sh" _is_interactive_context)"
if [ "$lib_set" = "$claim_set" ]; then
    ok "the two interactive-context helpers read the same signal set"
else
    no "lib.sh::jarvis_interactive_context and claim.sh::_is_interactive_context disagree:
$(diff <(printf '%s\n' "$lib_set") <(printf '%s\n' "$claim_set") | sed 's/^/    /')"
fi
if [ -z "$harness_signals" ]; then
    no "cannot read JARVIS_TEST_INTERACTIVE_SIGNALS from $harness"
else
    for spec in "bootstrap/lib.sh:jarvis_interactive_context" \
                "bootstrap/claim.sh:_is_interactive_context"; do
        file="$repo_root/${spec%%:*}"
        fn="${spec##*:}"
        [ -f "$file" ] || { no "missing $file"; continue; }
        found="$(context_signals "$file" "$fn")"
        if [ -z "$found" ]; then
            no "no signals parsed from ${spec%%:*}::$fn (did it get renamed?)"
            continue
        fi
        missing=""
        while read -r sig; do
            [ -n "$sig" ] || continue
            printf '%s\n' "$harness_signals" | grep -qx "$sig" || missing="$missing $sig"
        done <<< "$found"
        if [ -n "$missing" ]; then
            no "${spec%%:*}::$fn reads$missing but test/lib/hermetic.sh does not isolate it"
        else
            ok "${spec%%:*}::$fn signals all covered by the harness"
        fi
    done
fi

# ── 2. every test that runs a fenced script is wired to the harness ──────────
# Executing form only: `bash <path>/claim.sh …`, or a variable holding that path
# which the test later runs. Doc/lint tests that merely mention the filename do
# not match, so they are not forced to adopt the harness.
executors=""
for f in "$repo_root"/test/*.sh; do
    base="$(basename "$f")"
    [ "$base" = "hermetic_isolation_test.sh" ] && continue
    grep -qE '(bash|env)[^|;&]*(bootstrap/)?(claim|wrap|triage-one)\.sh|^[A-Za-z_]+=.*(bootstrap|BOOTSTRAP_DIR)[^ ]*/(claim|wrap|triage-one)\.sh' "$f" \
        || continue
    executors="$executors $base"
done
[ -n "$executors" ] || no "found no tests executing claim.sh/wrap.sh/triage-one.sh — has the detection broken?"

for base in $executors; do
    f="$repo_root/test/$base"
    # A test that builds its own stub copy of the script never runs the real one.
    if grep -qE 'cat > "?\$[A-Za-z_]+(/[a-z]+)*/(claim|wrap|triage-one)\.sh' "$f"; then
        ok "$base stubs the script itself (no real execution)"
        continue
    fi
    if ! grep -q 'lib/hermetic.sh' "$f"; then
        no "$base executes a fenced script but does not source test/lib/hermetic.sh"
        continue
    fi
    # Either full isolation, or (for suites that deliberately simulate a session)
    # explicit signal isolation plus a pinned endpoint.
    if grep -q 'jarvis_test_hermetic_init' "$f"; then
        ok "$base uses jarvis_test_hermetic_init"
    elif grep -q 'jarvis_test_hermetic_isolate_signals' "$f" \
        && grep -q 'jarvis_test_hermetic_pin_control_plane' "$f"; then
        if grep -q 'JARVIS_INTERACTIVE_WORKER_RUNNER' "$f"; then
            ok "$base simulates a session with its own stub runner + pinned endpoint"
        else
            no "$base opens an interactive signal without installing a stub runner"
        fi
    else
        no "$base sources the harness but calls neither init nor isolate+pin"
    fi
done

# ── 3. the harness itself must actually do the three things ──────────────────
for term in 'unset "$name"' 'JARVIS_CONTROL_PLANE_BASE_URL' 'JARVIS_INTERACTIVE_WORKER_RUNNER' 'exit 97'; do
    grep -Fq -- "$term" "$harness" || no "harness lost '$term'"
done

# The sentinel must not be a reachable endpoint.
sentinel="$(grep -E '^JARVIS_TEST_CONTROL_PLANE_SENTINEL=' "$harness" | sed 's/^[^=]*=//; s/"//g')"
case "$sentinel" in
    http://127.0.0.1:*|http://localhost:*) ok "control-plane sentinel is loopback ($sentinel)" ;;
    *) no "control-plane sentinel must be loopback, got '$sentinel'" ;;
esac

if [ "$fail" = 0 ]; then
    echo "hermetic_isolation_test: PASS"
else
    echo "hermetic_isolation_test: FAIL" >&2
fi
exit "$fail"
