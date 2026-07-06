#!/usr/bin/env bash
# test/stop_hook_test.sh — verify Stop hook commands resolve the repo without
# relying on Claude-only runtime variables.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmp_root="$(mktemp -d)"
tmp_out="$(mktemp)"
trap 'rm -rf "$tmp_root"; rm -f "$tmp_out"' EXIT

mkdir -p "$tmp_root/.my-day" "$tmp_root/runs" "$tmp_root/config"
cp "$repo_root/config/pools.json" "$tmp_root/config/pools.json"
codex_home="$tmp_root/codex-home"
codex_thread_id="stop-hook-test-thread"

pass=0
fail=0

ok() { echo "PASS $1"; pass=$((pass + 1)); }
no() { echo "FAIL $1"; fail=$((fail + 1)); }

stop_command() {
    local cfg="$1"
    jq -r '.hooks.Stop[0].hooks[0].command' "$repo_root/$cfg"
}

assert_hook_runs_without_claude_project_dir() {
    local cfg="$1" cmd rc
    cmd="$(stop_command "$cfg")"
    (
        cd "$repo_root" || exit 99
        env -u CLAUDE_PROJECT_DIR \
            JARVIS_ROOT="$tmp_root" \
            JARVIS_RUNS_DIR="$tmp_root/runs" \
            bash -c "$cmd"
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ]; then
        ok "$cfg Stop hook runs when CLAUDE_PROJECT_DIR is unset"
    else
        no "$cfg Stop hook runs when CLAUDE_PROJECT_DIR is unset (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_hook_has_readable_root_error() {
    local cfg="$1" cmd rc outside
    cmd="$(stop_command "$cfg")"
    outside="$(mktemp -d)"
    (
        cd "$outside" || exit 99
        env -u CLAUDE_PROJECT_DIR \
            CODEX_HOME="$tmp_root/missing-codex-home" \
            CODEX_THREAD_ID="missing-stop-hook-thread" \
            JARVIS_ROOT="$tmp_root" \
            JARVIS_RUNS_DIR="$tmp_root/runs" \
            bash -c "$cmd"
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -ne 0 ] && grep -q "cannot resolve repo root" "$tmp_out"; then
        ok "$cfg Stop hook reports readable repo-root error outside a git repo"
    else
        no "$cfg Stop hook reports readable repo-root error outside a git repo (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_preflight_records_codex_stop_runner() {
    rm -rf "$codex_home"
    (
        cd "$repo_root" || exit 99
        CODEX_HOME="$codex_home" \
            CODEX_THREAD_ID="$codex_thread_id" \
            bash bootstrap/preflight.sh
    ) >"$tmp_out" 2>&1

    local runner="$codex_home/jarvis/run-stop-hook.sh"
    local root_file="$codex_home/jarvis/session-roots/$codex_thread_id"
    if [ -x "$runner" ] && [ "$(cat "$root_file" 2>/dev/null)" = "$repo_root" ]; then
        ok "preflight records Codex Stop hook runner and session root"
    else
        no "preflight records Codex Stop hook runner and session root (output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_codex_hook_runs_outside_git_with_recorded_root() {
    local cmd rc outside
    cmd="$(stop_command ".codex/hooks.json")"
    outside="$(mktemp -d)"
    (
        cd "$outside" || exit 99
        env -u CLAUDE_PROJECT_DIR \
            CODEX_HOME="$codex_home" \
            CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_ROOT="$tmp_root" \
            JARVIS_RUNS_DIR="$tmp_root/runs" \
            bash -c "$cmd"
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -eq 0 ]; then
        ok ".codex/hooks.json Stop hook runs outside a git repo with recorded root"
    else
        no ".codex/hooks.json Stop hook runs outside a git repo with recorded root (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_preflight_records_codex_stop_runner
assert_codex_hook_runs_outside_git_with_recorded_root

for cfg in .codex/hooks.json .claude/settings.json; do
    assert_hook_runs_without_claude_project_dir "$cfg"
    assert_hook_has_readable_root_error "$cfg"
done

echo ""
echo "stop_hook_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ] && { echo "ALLPASS"; exit 0; } || { echo "FAILED"; exit 1; }
