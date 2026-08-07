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

prompt_command() {
    jq -r '.hooks.UserPromptSubmit[0].hooks[0].command' \
        "$repo_root/.codex/hooks.json"
}

tool_activity_command() {
    local event="$1"
    jq -r --arg event "$event" '.hooks[$event][0].hooks[0].command' \
        "$repo_root/.codex/hooks.json"
}

assert_hook_runs_without_claude_project_dir() {
    local cfg="$1" cmd rc
    cmd="$(stop_command "$cfg")"
    (
        cd "$repo_root" || exit 99
        env -u CLAUDE_PROJECT_DIR \
            JARVIS_ROOT="$tmp_root" \
            JARVIS_RUNS_DIR="$tmp_root/runs" \
            bash -c "$cmd" <<'EOF'
{"hook_event_name":"Stop","session_id":"stop-hook-test-thread","turn_id":"stop-turn"}
EOF
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
    if [ "$rc" -eq 2 ] \
            && grep -Eq "cannot resolve (a current )?repo" "$tmp_out"; then
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
            JARVIS_PREFLIGHT_SKIP_RECONCILE=1 \
            bash bootstrap/preflight.sh
    ) >"$tmp_out" 2>&1

    local runner="$codex_home/jarvis/run-stop-hook.sh"
    local root_file="$codex_home/jarvis/session-roots/$codex_thread_id"
    if [ -x "$runner" ] \
            && [ "$(cat "$root_file" 2>/dev/null)" = "$repo_root" ] \
            && grep -q 'run-interactive-worker-hook.sh.*codex' "$runner"; then
        ok "preflight records Codex Stop hook runner and session root"
    else
        no "preflight records Codex Stop hook runner and session root (output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_preflight_fails_when_root_mapping_cannot_be_written() {
    local bad_home="$tmp_root/bad-codex-home" rc
    mkdir -p "$bad_home/jarvis/repo-root"
    (
        cd "$repo_root" || exit 99
        env -u CODEX_THREAD_ID CODEX_HOME="$bad_home" \
            JARVIS_PREFLIGHT_SKIP_RECONCILE=1 \
            bash bootstrap/preflight.sh
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -ne 0 ] \
            && grep -q 'failed to install Codex hook runner' "$tmp_out"; then
        ok "preflight fails closed when Codex root mapping cannot be installed"
    else
        no "preflight install failure is not masked (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
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
            bash -c "$cmd" <<EOF
{"hook_event_name":"Stop","session_id":"$codex_thread_id","turn_id":"stop-turn"}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -eq 0 ]; then
        ok ".codex/hooks.json Stop hook runs outside a git repo with recorded root"
    else
        no ".codex/hooks.json Stop hook runs outside a git repo with recorded root (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_codex_prompt_runs_outside_git_with_recorded_root() {
    local cmd rc outside
    cmd="$(prompt_command)"
    outside="$(mktemp -d)"
    (
        cd "$outside" || exit 99
        env -u CLAUDE_PROJECT_DIR \
            CODEX_HOME="$codex_home" \
            CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_INTERACTIVE_STATE_DIR="$tmp_root/interactive-state" \
            bash -c "$cmd" <<EOF
{"hook_event_name":"UserPromptSubmit","session_id":"$codex_thread_id","turn_id":"outside-turn"}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -eq 0 ] && jq -e . "$tmp_out" >/dev/null 2>&1; then
        ok ".codex UserPromptSubmit runs outside git with recorded root"
    else
        no ".codex UserPromptSubmit outside git (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_codex_tool_activity_hook_runs_outside_git() {
    local event cmd rc outside
    for event in PreToolUse PostToolUse; do
        cmd="$(tool_activity_command "$event")"
        outside="$(mktemp -d)"
        (
            cd "$outside" || exit 99
            CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
                JARVIS_INTERACTIVE_STATE_DIR="$tmp_root/interactive-state" \
                bash -c "$cmd" <<EOF
{"hook_event_name":"$event","session_id":"$codex_thread_id","turn_id":"long-turn","tool_name":"exec_command"}
EOF
        ) >"$tmp_out" 2>&1
        rc=$?
        rm -rf "$outside"
        if [ "$event" = "PreToolUse" ] && [ "$rc" -eq 2 ] \
                && grep -q 'SessionStart' "$tmp_out"; then
            ok ".codex PreToolUse blocks an unregistered Worker outside git"
        elif [ "$event" = "PostToolUse" ] && [ "$rc" -eq 0 ] \
                && jq -e . "$tmp_out" >/dev/null 2>&1; then
            ok ".codex PostToolUse activity hook invokes Worker outside git"
        else
            no ".codex $event activity hook (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
        fi
    done
}

assert_codex_pretool_root_resolution_fails_closed() {
    local cmd rc outside
    cmd="$(tool_activity_command PreToolUse)"
    outside="$(mktemp -d)"
    (
        cd "$outside" || exit 99
        CODEX_HOME="$tmp_root/missing-pretool-home" \
            CODEX_THREAD_ID="missing-pretool-thread" \
            JARVIS_ROOT="$tmp_root" bash -c "$cmd" <<'EOF'
{"hook_event_name":"PreToolUse","session_id":"missing-pretool-thread","turn_id":"t","tool_name":"Bash","tool_input":{"command":"a1 update"}}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -eq 2 ] && grep -q 'worker-fence-hook: cannot resolve' "$tmp_out"; then
        ok "Codex PreToolUse fails closed when no current repo can be resolved"
    else
        no "Codex PreToolUse root failure blocks (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_pretool_wrapper_normalizes_all_failures_to_block() {
    local cmd rc fake_manager fake_marker wrapper_only
    cmd="$(tool_activity_command PreToolUse)"
    fake_manager="$tmp_root/fake-pretool-manager"
    fake_marker="$tmp_root/fake-pretool-manager-invoked"
    cat > "$fake_manager" <<EOF
#!/usr/bin/env bash
touch "$fake_marker"
printf '{}\n'
EOF
    chmod +x "$fake_manager"
    (
        cd "$repo_root" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_INTERACTIVE_WORKER_PYTHON="/definitely/missing/python" \
            JARVIS_INTERACTIVE_WORKER_MANAGER="$fake_manager" \
            bash -c "$cmd" <<'EOF'
{malformed-json
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 2 ] && [ ! -e "$fake_marker" ] \
            && grep -q 'fail-closed' "$tmp_out"; then
        ok "PreToolUse ignores replaceable Python/manager overrides"
    else
        no "PreToolUse trusted manager cannot be replaced (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    wrapper_only="$tmp_root/pretool-wrapper-only"
    mkdir -p "$wrapper_only"
    cp "$repo_root/bootstrap/run-interactive-worker-hook.sh" \
        "$wrapper_only/run-interactive-worker-hook.sh"
    cp "$repo_root/bootstrap/runtime-config.sh" \
        "$wrapper_only/runtime-config.sh"
    bash "$wrapper_only/run-interactive-worker-hook.sh" codex PreToolUse \
        >"$tmp_out" 2>&1 <<EOF
{"hook_event_name":"PreToolUse","session_id":"$codex_thread_id","turn_id":"t","tool_name":"Bash","tool_input":{"command":"a1 update"}}
EOF
    rc=$?
    if [ "$rc" -eq 2 ] && grep -q 'trusted PreToolUse manager is unavailable' "$tmp_out"; then
        ok "PreToolUse missing trusted manager is normalized to exit 2"
    else
        no "PreToolUse missing trusted manager blocks (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    (
        cd "$repo_root" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            bash -c "$cmd" <<'EOF'
{malformed-json
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 2 ] && grep -q 'fail-closed' "$tmp_out"; then
        ok "PreToolUse malformed payload is blocked fail-closed"
    else
        no "PreToolUse malformed payload blocks (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    (
        cd "$repo_root" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            bash -c "$cmd" <<EOF
{"hook_event_name":"PreToolUse","turn_id":"t","tool_name":"Bash","tool_input":{"command":"a1 update"}}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 2 ] && grep -q 'fail-closed' "$tmp_out"; then
        ok "PreToolUse missing session id is blocked fail-closed"
    else
        no "PreToolUse missing session id blocks (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_lifecycle_prefers_system_python() {
    local wrapper="$repo_root/bootstrap/run-interactive-worker-hook.sh"
    if grep -q 'python_bin="/usr/bin/python3"' "$wrapper" \
            && grep -q 'JARVIS_INTERACTIVE_WORKER_PYTHON' "$wrapper" \
            && grep -q 'python_bin="python3"' "$wrapper"; then
        ok "interactive lifecycle prefers system Python with override and fallback"
    else
        no "interactive lifecycle Python trust-store selection is incomplete"
    fi
}

assert_interactive_wrapper_strips_admin_token() {
    local fake_python="$tmp_root/fake-worker-python"
    local empty_env="$tmp_root/empty-worker.env"
    local rc
    : >"$empty_env"
    chmod 600 "$empty_env"
    cat >"$fake_python" <<'EOF'
#!/usr/bin/env bash
if [ -n "${JARVIS_CONTROL_PLANE_ADMIN_TOKEN:-}" ]; then
    echo "admin token reached interactive worker" >&2
    exit 91
fi
exit 0
EOF
    chmod +x "$fake_python"
    JARVIS_CONTROL_PLANE_ADMIN_TOKEN="operator-only" \
        JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$empty_env" \
        JARVIS_INTERACTIVE_BRIDGE_ENV="$empty_env" \
        XDG_CONFIG_HOME="$tmp_root/no-runtime-config" \
        HOME="$tmp_root" \
        JARVIS_INTERACTIVE_WORKER_PYTHON="$fake_python" \
        JARVIS_INTERACTIVE_WORKER_MANAGER="/unused/worker.py" \
        bash "$repo_root/bootstrap/run-interactive-worker-hook.sh" cli status \
        >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ]; then
        ok "interactive wrapper strips operator admin token"
    else
        no "interactive wrapper leaked operator admin token (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_claude_has_all_tool_worker_fence_handler() {
    local command matcher
    command="$(jq -r '.hooks.PreToolUse[0].hooks[0].command' \
        "$repo_root/.claude/settings.json")"
    matcher="$(jq -r '.hooks.PreToolUse[0].matcher // ""' \
        "$repo_root/.claude/settings.json")"
    if [ -z "$matcher" ] \
            && printf '%s' "$command" | grep -q 'run-interactive-worker-hook.sh.*claude' \
            && printf '%s' "$command" | grep -q 'exit 2'; then
        ok "Claude PreToolUse installs an all-tool fail-closed Worker fence"
    else
        no "Claude PreToolUse all-tool Worker fence is missing"
    fi
}

assert_stale_global_runner_cannot_bypass_current_hooks() {
    local runner="$codex_home/jarvis/run-stop-hook.sh"
    local runner_backup="$tmp_root/current-runner-backup"
    local invoked="$tmp_root/stale-runner-invoked" outside cmd rc
    cp "$runner" "$runner_backup"
    cat > "$runner" <<EOF
#!/usr/bin/env bash
touch "$invoked"
exit 0
EOF
    chmod +x "$runner"
    outside="$(mktemp -d)"
    cmd="$(prompt_command)"
    (
        cd "$outside" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_INTERACTIVE_STATE_DIR="$tmp_root/interactive-state" \
            bash -c "$cmd" <<EOF
{"hook_event_name":"UserPromptSubmit","session_id":"$codex_thread_id","turn_id":"upgrade-turn"}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ] && [ ! -e "$invoked" ] && jq -e . "$tmp_out" >/dev/null 2>&1; then
        ok "current repo hook bypasses a stale installed prompt runner"
    else
        no "stale prompt runner is bypassed (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    cmd="$(stop_command ".codex/hooks.json")"
    (
        cd "$outside" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_ROOT="$tmp_root" JARVIS_RUNS_DIR="$tmp_root/runs" \
            bash -c "$cmd" <<EOF
{"hook_event_name":"Stop","session_id":"$codex_thread_id","turn_id":"stale-runner-stop"}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    if [ "$rc" -eq 0 ] && [ ! -e "$invoked" ]; then
        ok "current repo hook bypasses a stale installed Stop runner"
    else
        no "stale Stop runner is bypassed (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    rm -f "$codex_home/jarvis/repo-root" \
        "$codex_home/jarvis/session-roots/$codex_thread_id" "$invoked"
    cmd="$(prompt_command)"
    (
        cd "$outside" || exit 99
        CODEX_HOME="$codex_home" CODEX_THREAD_ID="$codex_thread_id" \
            JARVIS_ROOT="$tmp_root" bash -c "$cmd" <<EOF
{"hook_event_name":"UserPromptSubmit","session_id":"$codex_thread_id","turn_id":"blocked-upgrade"}
EOF
    ) >"$tmp_out" 2>&1
    rc=$?
    rm -rf "$outside"
    if [ "$rc" -eq 2 ] && [ ! -e "$invoked" ] \
            && grep -q 'run bootstrap/preflight.sh' "$tmp_out"; then
        ok "stale prompt runner fails closed when no current repo can be resolved"
    else
        no "stale prompt runner cannot fail open (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi

    cp "$runner_backup" "$runner"
    chmod +x "$runner"
    mkdir -p "$codex_home/jarvis/session-roots"
    printf '%s\n' "$repo_root" > "$codex_home/jarvis/repo-root"
    printf '%s\n' "$repo_root" > \
        "$codex_home/jarvis/session-roots/$codex_thread_id"
}

assert_codex_stop_has_single_ordered_handler() {
    local count command
    count="$(jq '.hooks.Stop[0].hooks | length' "$repo_root/.codex/hooks.json")"
    command="$(stop_command ".codex/hooks.json")"
    if [ "$count" -eq 1 ] && printf '%s' "$command" | grep -q 'run-stop-hook.sh.*codex'; then
        ok "Codex Stop uses one ordered wrapper handler"
    else
        no "Codex Stop uses one ordered wrapper handler (count=$count command=$command)"
    fi
}

make_fake_stop_scripts() {
    local dir="$1" wrap_rc="$2" worker_rc="$3"
    mkdir -p "$dir"
    cp "$repo_root/bootstrap/run-stop-hook.sh" "$dir/run-stop-hook.sh"
    cp "$repo_root/bootstrap/runtime-config.sh" "$dir/runtime-config.sh"
    # Force the primary control-plane check into its documented unavailable
    # fallback. Do not depend on a machine-local manager path or runtime env.
    cat > "$dir/jarvis-interactive-worker.py" <<'PY'
raise SystemExit(1)
PY
    cat > "$dir/wrap-check.sh" <<'EOF'
#!/usr/bin/env bash
cat >/dev/null
echo wrap >> "$STOP_TEST_ORDER"
echo wrap-stdout
echo wrap-stderr >&2
exit "$STOP_TEST_WRAP_RC"
EOF
    cat > "$dir/run-interactive-worker-hook.sh" <<'EOF'
#!/usr/bin/env bash
cat > "$STOP_TEST_PAYLOAD"
echo worker >> "$STOP_TEST_ORDER"
echo worker-stdout
echo worker-stderr >&2
exit "$STOP_TEST_WORKER_RC"
EOF
    chmod +x "$dir/"*.sh
    export STOP_TEST_WRAP_RC="$wrap_rc"
    export STOP_TEST_WORKER_RC="$worker_rc"
}

assert_ordered_wrapper_success_is_single_json() {
    local dir="$tmp_root/fake-stop-success" order="$tmp_root/order-success"
    local received="$tmp_root/payload-success" stdout="$tmp_root/stdout-success"
    local stderr="$tmp_root/stderr-success" payload rc
    payload='{"hook_event_name":"Stop","session_id":"s1","turn_id":"t1"}'
    make_fake_stop_scripts "$dir" 0 0
    STOP_TEST_ORDER="$order" STOP_TEST_PAYLOAD="$received" \
        bash "$dir/run-stop-hook.sh" codex >"$stdout" 2>"$stderr" <<<"$payload"
    rc=$?
    if [ "$rc" -eq 0 ] \
            && [ "$(cat "$stdout")" = '{}' ] \
            && [ "$(tr '\n' ' ' < "$order")" = 'wrap worker ' ] \
            && [ "$(cat "$received")" = "$payload" ] \
            && grep -q 'wrap-stdout' "$stderr"; then
        ok "ordered Stop passes exact payload and emits one canonical JSON object"
    else
        no "ordered Stop success contract (rc=$rc stdout=$(tr '\n' ' ' <"$stdout") order=$(tr '\n' ' ' <"$order" 2>/dev/null))"
    fi
}

assert_wrap_block_never_signals_worker() {
    local dir="$tmp_root/fake-stop-block" order="$tmp_root/order-block"
    local received="$tmp_root/payload-block" rc
    make_fake_stop_scripts "$dir" 2 0
    STOP_TEST_ORDER="$order" STOP_TEST_PAYLOAD="$received" \
        bash "$dir/run-stop-hook.sh" codex >"$tmp_out" 2>&1 <<<'{"hook_event_name":"Stop"}'
    rc=$?
    if [ "$rc" -eq 2 ] && [ "$(cat "$order")" = 'wrap' ] && [ ! -e "$received" ]; then
        ok "wrap-check block leaves Worker turn active"
    else
        no "wrap-check block leaves Worker turn active (rc=$rc order=$(tr '\n' ' ' <"$order" 2>/dev/null))"
    fi
}

assert_worker_failure_blocks_stop() {
    local dir="$tmp_root/fake-stop-worker-fail" order="$tmp_root/order-worker-fail"
    local received="$tmp_root/payload-worker-fail" rc
    make_fake_stop_scripts "$dir" 0 9
    STOP_TEST_ORDER="$order" STOP_TEST_PAYLOAD="$received" \
        bash "$dir/run-stop-hook.sh" codex >"$tmp_out" 2>&1 <<<'{"hook_event_name":"Stop"}'
    rc=$?
    if [ "$rc" -eq 2 ] && grep -q 'interactive worker Stop signal failed' "$tmp_out"; then
        ok "Worker Stop accounting failure blocks Stop fail-closed"
    else
        no "Worker Stop accounting failure blocks Stop fail-closed (rc=$rc output=$(tr '\n' ' ' <"$tmp_out"))"
    fi
}

assert_preflight_records_codex_stop_runner
assert_preflight_fails_when_root_mapping_cannot_be_written
assert_codex_hook_runs_outside_git_with_recorded_root
assert_codex_prompt_runs_outside_git_with_recorded_root
assert_codex_tool_activity_hook_runs_outside_git
assert_codex_pretool_root_resolution_fails_closed
assert_pretool_wrapper_normalizes_all_failures_to_block
assert_lifecycle_prefers_system_python
assert_interactive_wrapper_strips_admin_token
assert_claude_has_all_tool_worker_fence_handler
assert_stale_global_runner_cannot_bypass_current_hooks
assert_codex_stop_has_single_ordered_handler

for cfg in .codex/hooks.json .claude/settings.json; do
    assert_hook_runs_without_claude_project_dir "$cfg"
    assert_hook_has_readable_root_error "$cfg"
done

assert_ordered_wrapper_success_is_single_json
assert_wrap_block_never_signals_worker
assert_worker_failure_blocks_stop

echo ""
echo "stop_hook_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ] && { echo "ALLPASS"; exit 0; } || { echo "FAILED"; exit 1; }
