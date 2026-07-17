#!/usr/bin/env bash
# Hermetic coverage for the a1 destructive-CR execution gate.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
guard="$repo_root/bootstrap/a1_command_guard.py"
manager="$repo_root/bootstrap/jarvis-interactive-worker.py"
a1id="$repo_root/bin/a1id"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

pass=0
fail=0
ok(){ echo "PASS: $1"; pass=$((pass + 1)); }
no(){ echo "FAIL: $1" >&2; fail=$((fail + 1)); }

python3 - "$guard" "$manager" <<'PY'
import contextlib
import importlib.util
import pathlib
import sys
import time

guard_path = pathlib.Path(sys.argv[1])
manager_path = pathlib.Path(sys.argv[2])

def load(name, path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    sys.modules[name] = module
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module

g = load("a1_command_guard", guard_path)

blocked = {
    "a1 app pipeline exit-cr --pipeline-id=66": "exit-cr",
    "a1 --format=json app pipeline quit --pipeline-id 66": "pipeline quit",
    "a1 --quiet=true app pipeline exit-cr --pipeline-id=66": "exit-cr",
    "a1 -q=false app pipeline quit --pipeline-id 66": "pipeline quit",
    "env X=1 /usr/local/bin/a1 app --verbose pipeline --format json exit-cr --pipeline-id=66": "exit-cr",
    "bin/a1id -- app pipeline --format=json quit --pipeline-id 66": "pipeline quit",
    "bash bin/a1id as jarvis -- --quiet app pipeline exit-cr --pipeline-id=66": "exit-cr",
    "a1 --no-update-check app cr quit 123 --pipeline-id=66": "direct",
    "/bin/bash -lc 'a1 app pipeline exit-cr --pipeline-id 66'": "exit-cr",
    "/bin/zsh -fc 'a1 app pipeline quit --pipeline-id=66'": "pipeline quit",
    "exec a1 app pipeline quit --pipeline-id=66": "pipeline quit",
    "(a1 app pipeline exit-cr --pipeline-id 66)": "exit-cr",
    "if a1 app pipeline exit-cr --pipeline-id 66; then :; fi": "exit-cr",
    "! a1 app pipeline quit --pipeline-id 66": "pipeline quit",
    "time -p a1 app pipeline exit-cr --pipeline-id 66": "exit-cr",
    "nohup a1 app pipeline quit --pipeline-id 66": "pipeline quit",
    "time env X=1 a1 app pipeline quit --pipeline-id 66": "pipeline quit",
}
for command, marker in blocked.items():
    for tool_name in ("Bash", "exec_command", "mcp__shell__exec_command"):
        event = {"tool_name": tool_name, "tool_input": {"command": command}}
        reason = g.pretool_a1_block_reason(event)
        if not reason or marker not in reason:
            raise SystemExit("expected block for %r/%s, got %r" % (
                command, tool_name, reason))

allowed = (
    "bin/a1id -- app cr quit 123 --pipeline-id 66",
    "a1 app cr get 123 --format json",
    "rg -n exit-cr docs",
    "echo a1 app pipeline exit-cr",
    "printf '%s' 'a1 app pipeline quit'",
    "command -v a1",
)
for command in allowed:
    event = {"tool_name": "exec_command", "tool_input": {"cmd": command}}
    reason = g.pretool_a1_block_reason(event)
    if reason:
        raise SystemExit("unexpected block for %r: %s" % (command, reason))
target_event = {
    "tool_name": "exec_command",
    "tool_input": {"cmd": "bin/a1id -- app cr quit 123 --pipeline-id=66"},
}
assert g.pretool_targeted_a1id_quit(target_event) == ("123", "66")

for malformed in (
    ["app", "cr", "quit", "--pipeline-id", "66"],
    ["app", "cr", "quit", "123"],
    ["app", "cr", "quit", "123", "456", "--pipeline-id", "66"],
    ["app", "cr", "quit", "not-a-number", "--pipeline-id", "66"],
    ["app", "cr", "quit", "0", "--pipeline-id", "66"],
    ["app", "cr", "quit", "123", "--pipeline-id", "0"],
    ["app", "cr", "quit", "123", "--pipeline-id", "66", "--pipeline-id", "67"],
    ["--config", "/tmp/other", "app", "cr", "quit", "123", "--pipeline-id", "66"],
    ["app", "cr", "quit", "123", "--pipeline-id", "66", "--force"],
):
    try:
        g._parse_targeted_quit(malformed)
    except g.GuardError:
        pass
    else:
        raise SystemExit("malformed targeted quit was accepted: %r" % malformed)

# CR ownership requires all three dimensions: claimed work item, branch, repo.
cr = {
    "appId": 260634,
    "appName": "safe-app",
    "crId": 12345678,
    "workItemIds": [123456],
    "crItems": [{"crItemContent": {
        "branchName": "feature/safe",
        "trunkUrl": "git@gitlab.alibaba-inc.com:demo/safe-app.git",
    }}],
}
app_id = g._verify_cr_ownership(
    cr, cr_id="12345678", requested_app=None,
    claimed_aone_id="123456",
    origin="gitlab.alibaba-inc.com/demo/safe-app", branch="feature/safe")
assert app_id == "260634"
for mutation, marker in (
    ({"workItemIds": [654321]}, "claimed Aone"),
    ({"crItems": [{"crItemContent": {
        "branchName": "feature/other",
        "trunkUrl": "git@gitlab.alibaba-inc.com:demo/safe-app.git"}}]},
     "repository/branch"),
):
    changed = dict(cr)
    changed.update(mutation)
    try:
        g._verify_cr_ownership(
            changed, cr_id="12345678", requested_app=None,
            claimed_aone_id="123456",
            origin="gitlab.alibaba-inc.com/demo/safe-app", branch="feature/safe")
    except g.GuardError as exc:
        assert marker in str(exc)
    else:
        raise SystemExit("ownership mismatch was accepted: %s" % marker)

# Pipeline status -> latest instance -> branch membership is strict.
responses = iter((
    {"appId": 260634, "pipelineId": 66, "pipelineInstanceId": 9001},
    {"appId": 260634, "pipelineId": 66, "pipelineInstanceId": 9001,
     "changeRequests": [{"crId": 12345678}]},
))
original_read_json = g._read_a1_json
g._read_a1_json = lambda *_args, **_kwargs: next(responses)
g._verify_pipeline_membership(
    "a1", app_id="260634", pipeline_id="66", cr_id="12345678")
g._read_a1_json = lambda *_args, **_kwargs: {
    "appId": 260634, "pipelineId": 66, "pipelineInstanceId": 9001,
    "changeRequests": [],
}
try:
    # First response is intentionally valid status; second is empty branch set.
    seq = iter((
        {"appId": 260634, "pipelineId": 66, "pipelineInstanceId": 9001},
        {"appId": 260634, "pipelineId": 66, "pipelineInstanceId": 9001,
         "changeRequests": []},
    ))
    g._read_a1_json = lambda *_args, **_kwargs: next(seq)
    g._verify_pipeline_membership(
        "a1", app_id="260634", pipeline_id="66", cr_id="12345678")
except g.GuardError as exc:
    assert "not attached" in str(exc)
else:
    raise SystemExit("pipeline membership mismatch was accepted")
g._read_a1_json = original_read_json

# Full guarded path checks exact authority twice and canonicalizes the mutation.
fingerprint = {
    "aoneId": "123456", "assignmentEpoch": "session:s:f:g",
    "workerKeyDigest": "0123456789abcdef",
}
authority_calls = []
g._worker_authority = lambda phase, cr_id, pipeline_id, token="": (
    authority_calls.append((phase, cr_id, pipeline_id, token)) or
    (fingerprint, "permit-token"))
g._current_repo_branch = lambda: (
    "gitlab.alibaba-inc.com/demo/safe-app", "feature/safe")
g._read_cr = lambda *_args, **_kwargs: cr
g._verify_pipeline_membership = lambda *_args, **_kwargs: None

class ExecCalled(Exception):
    pass
def fake_exec(_bin, argv, _env):
    raise ExecCalled(argv)
g.os.execvpe = fake_exec
try:
    g.run_guarded("a1", [
        "--format=json", "app", "cr", "quit", "12345678",
        "--pipeline-id=66",
    ])
except ExecCalled as exc:
    assert exc.args[0] == [
        "a1", "app", "cr", "quit", "12345678",
        "--pipeline-id", "66", "--app", "260634", "--format=json",
    ]
else:
    raise SystemExit("guarded path did not exec canonical targeted quit")
assert authority_calls == [
    ("begin", "12345678", "66", ""),
    ("confirm", "12345678", "66", "permit-token"),
]

# Manager begin/confirm reuses locked current/lease authority and consumes a
# PreTool-issued permit.  Patch only host/lease primitives, not permit logic.
m = load("jarvis_interactive_worker_test", manager_path)
class FakeStore:
    def __init__(self, state):
        self.state = state
    @contextlib.contextmanager
    def locked(self):
        yield
    def load_unlocked(self):
        return self.state
    def save_unlocked(self, state):
        self.state = state

current = {
    "aoneId": "123456", "taskId": "t", "sessionId": "s",
    "fenceToken": "f", "generation": "g", "runtimeSessionId": "r",
    "heartbeatEnabled": True,
}
state = {
    "client": "codex", "clientSessionId": "thread", "workerKey": "worker",
    "current": current, "turnActive": True, "activeTurnId": "turn",
    "a1CrQuitPermit": {
        "token": "permit-token", "phase": "READY", "crId": "12345678",
        "pipelineId": "66", "aoneId": "123456",
        "assignmentEpoch": "session:s:f:g", "rootTurnId": "turn",
        "expiresAt": time.time() + 60,
    },
}
store = FakeStore(state)
m._runtime_context = lambda: ("codex", "thread")
m._current_store = lambda: store
m._calling_process_matches = lambda *_args: True
m._local_tool_block_reason = lambda *_args: None
m._codex_turn_live = lambda *_args: True
m._session_permit_block_reason = lambda *_args, **_kwargs: None
begin = m.current_authority("begin", "12345678", "66")
assert begin["permitToken"] == "permit-token"
assert store.state["a1CrQuitPermit"]["phase"] == "VERIFYING"
confirm = m.current_authority("confirm", "12345678", "66", "permit-token")
assert confirm["assignmentEpoch"] == begin["assignmentEpoch"]
assert "a1CrQuitPermit" not in store.state

# Missing/mismatched/expired permits are all fail closed.
for permit in (None, {
    "token": "permit-token", "phase": "READY", "crId": "99999999",
    "pipelineId": "66", "aoneId": "123456",
    "assignmentEpoch": "session:s:f:g", "rootTurnId": "turn",
    "expiresAt": time.time() + 60,
}, {
    "token": "permit-token", "phase": "READY", "crId": "12345678",
    "pipelineId": "66", "aoneId": "123456",
    "assignmentEpoch": "session:s:f:g", "rootTurnId": "turn",
    "expiresAt": time.time() - 1,
}):
    if permit is None:
        store.state.pop("a1CrQuitPermit", None)
    else:
        store.state["a1CrQuitPermit"] = permit
    try:
        m.current_authority("begin", "12345678", "66")
    except RuntimeError:
        pass
    else:
        raise SystemExit("invalid PreTool permit was accepted: %r" % permit)
PY
if [ $? -eq 0 ]; then
    ok "parser, ownership, pipeline, canonical quit, and authority permit invariants"
else
    no "Python guard/authority invariants"
fi

bin_dir="$tmp/bin"
identity_root="$tmp/a1"
capture="$tmp/a1.log"
mkdir -p "$bin_dir" "$identity_root"
cat > "$bin_dir/a1" <<'STUB'
#!/usr/bin/env bash
printf 'A1_CONFIG_DIR=%s ARGS=%s\n' "${A1_CONFIG_DIR:-}" "$*" >> "$STUB_CAPTURE"
exit 0
STUB
chmod +x "$bin_dir/a1"

run_a1id(){
    A1ID_ROOT="$identity_root" A1_BIN="$bin_dir/a1" STUB_CAPTURE="$capture" \
        bash "$a1id" "$@"
}

# No identity is seeded: a first-run wrapper would normally probe live whoami.
# Early bulk denial must happen before even that read reaches real a1.
for variant in exit-cr quit; do
    rm -rf "$identity_root/identities"
    mkdir -p "$identity_root"
    printf 'live\n' > "$identity_root/auth.yaml"
    : > "$capture"
    run_a1id -- --format=json app pipeline "$variant" --pipeline-id=66 \
        >/dev/null 2>"$tmp/blocked.err"
    rc=$?
    if [ "$rc" -ne 0 ] && [ ! -s "$capture" ] \
        && grep -Fq 'app cr quit' "$tmp/blocked.err"; then
        ok "a1id early-denies pipeline $variant with zero real-a1 calls"
    else
        no "a1id early deny $variant rc=$rc log=$(cat "$capture")"
    fi
done

# Interspersed global flags and the explicit identity wrapper use the same
# parser; seed auth so no unrelated identity probe can affect the assertion.
mkdir -p "$identity_root/identities/jarvis"
printf 'seeded\n' > "$identity_root/identities/jarvis/auth.yaml"
: > "$capture"
run_a1id as jarvis -- app --quiet pipeline --format json exit-cr --pipeline-id 66 \
    >/dev/null 2>"$tmp/as-blocked.err"
rc=$?
if [ "$rc" -ne 0 ] && [ ! -s "$capture" ]; then
    ok "a1id as/global-flag variant is denied before real a1"
else
    no "a1id as/global variant bypass rc=$rc log=$(cat "$capture")"
fi

# Regression for the old broad token-presence implementation: an app value
# literally named "pipeline" is not the app/pipeline/quit command path.
: > "$capture"
run_a1id -- app cr quit 1 --pipeline-id 2 --app pipeline \
    >/dev/null 2>"$tmp/not-bulk.err"
rc=$?
if [ "$rc" -ne 0 ] && ! grep -Fq 'permanently disabled' "$tmp/not-bulk.err" \
    && [ "$(grep -Fc 'ARGS=app cr quit' "$capture" || true)" -eq 0 ]; then
    ok "targeted quit with --app pipeline is not misclassified as bulk quit"
else
    no "early parser false positive rc=$rc err=$(cat "$tmp/not-bulk.err")"
fi

: > "$capture"
run_a1id -- app cr get 123 --format json >/dev/null 2>"$tmp/ordinary.err"
rc=$?
if [ "$rc" -eq 0 ] && grep -Fq 'ARGS=app cr get 123 --format json' "$capture"; then
    ok "ordinary a1id commands still pass through"
else
    no "ordinary a1 command regression rc=$rc err=$(cat "$tmp/ordinary.err")"
fi

/usr/bin/python3 -I "$manager" --help >/dev/null 2>&1
if [ $? -eq 0 ]; then
    ok "trusted -I interactive-worker entrypoint imports the guard"
else
    no "interactive-worker -I import smoke"
fi

echo "a1_command_guard_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
