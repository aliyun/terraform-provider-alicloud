#!/usr/bin/env bash
# Hermetic coverage for the a1 CR/branch-exit execution gate.

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

python3 - "$guard" <<'PY'
import importlib.util
import pathlib
import sys

guard_path = pathlib.Path(sys.argv[1])
spec = importlib.util.spec_from_file_location("a1_command_guard", guard_path)
g = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = g
assert spec.loader is not None
spec.loader.exec_module(g)

allowed_bindings = g._configured_nonprod_pipeline_bindings()
assert {("260634", 65), ("260634", 420), ("cloudspec", 420)}.issubset(
    allowed_bindings)
assert ("260634", 67) not in allowed_bindings

malformed_policies = (
    None,
    {},
    {"pools": []},
    {"pools": {"x": "bad"}},
    {"pools": {"x": {"apps": {}}}},
    {"pools": {"x": {"apps": [{"app": 1, "name": "x"}]}}},
    {"pools": {"x": {"apps": [{
        "app": 1, "name": "x",
        "pipelines": {"prestage": "420", "prod": 67},
    }]}}},
    {"pools": {"x": {"apps": [{
        "app": 1, "name": "x", "pipelines": {"prestage": 420},
    }]}}},
    {"pools": {"x": {"apps": [{
        "app": 1, "name": "x",
        "pipelines": {"daily": 67, "prestage": 420, "prod": 67},
    }]}}},
)
for policy in malformed_policies:
    assert g._parse_nonprod_pipeline_bindings(policy) == frozenset()

renumbered = g._parse_nonprod_pipeline_bindings({"pools": {"x": {"apps": [{
    "app": 1,
    "name": "x",
    "pipelines": {"daily": 65, "prestage": 420, "prod": 999},
}]}}})
assert {("1", 65), ("1", 420), ("x", 420)}.issubset(renumbered)
assert ("1", 999) not in renumbered

blocked = {
    "a1 app pipeline exit-cr --pipeline-id=66": "app pipeline exit-cr",
    "a1 --format=json app pipeline quit --pipeline-id 66": "app pipeline quit",
    "a1 --no-update-check app cr quit 123 --pipeline-id=66": "app cr quit",
    "a1 app --quiet cr --format json quit 123 --pipeline-id 66": "app cr quit",
    "bin/a1id -- app cr quit 123 --pipeline-id 66": "app cr quit",
    "bash bin/a1id as jarvis -- --quiet app cr quit 123 --pipeline-id=66": "app cr quit",
    "env X=1 /usr/local/bin/a1 app --verbose pipeline --format json exit-cr --pipeline-id=66": "app pipeline exit-cr",
    "/bin/bash -lc 'a1 app cr quit 123 --pipeline-id 66'": "app cr quit",
    "/bin/zsh -fc 'bin/a1id -- app pipeline quit --pipeline-id=66'": "app pipeline quit",
    "exec a1 app pipeline quit --pipeline-id=66": "app pipeline quit",
    "(a1 app cr quit 123 --pipeline-id 66)": "app cr quit",
    "if a1 app pipeline exit-cr --pipeline-id 66; then :; fi": "app pipeline exit-cr",
    "! a1 app pipeline quit --pipeline-id 66": "app pipeline quit",
    "time -p a1 app cr quit 123 --pipeline-id 66": "app cr quit",
    "nohup a1 app pipeline quit --pipeline-id 66": "app pipeline quit",
    "a1 app cr submit 123 --app 260634 --pipeline-id 67":
        "non-production allowlist",
    "a1 app cr submit 123 --app cloudspec --pipeline-id +67":
        "non-production allowlist",
    "a1 app cr submit 123 --app=260634 --pipeline-id=067":
        "non-production allowlist",
    "a1 app cr submit 123 --app 260634 --pipeline-id 420 --pipeline-id 67":
        "non-production allowlist",
    "a1 app cr submit 123 --app 260634 --pipeline-id 999":
        "non-production allowlist",
    "a1 app cr submit 123 --pipeline-id 420":
        "cannot be proven non-production",
    "a1 --format=json app cr submit 123 --app 260634 --pipeline-id=67":
        "non-production allowlist",
    "a1 app --quiet cr --format json submit 123 --app 260634 --pipeline-id 67":
        "non-production allowlist",
    "bin/a1id -- app cr submit 123 --app=260634 --pipeline-id=67":
        "non-production allowlist",
    "bash bin/a1id as jarvis -- app cr submit 123 --app 260634 --pipeline-id 67":
        "non-production allowlist",
    "pid=67; a1 app cr submit 123 --app 260634 --pipeline-id $pid":
        "non-production allowlist",
    "pipeline=--pipeline-id=67; bin/a1id -- app cr submit 123 --app 260634 $pipeline":
        "non-production allowlist",
    "/bin/bash -lc 'a1 app cr submit 123 --app 260634 --pipeline-id 67'":
        "non-production allowlist",
    "a1 cd-pipeline run 67 --app 260634 --cr-id 123": "non-production allowlist",
    "a1 cd-pipeline run +67 --app 260634 --cr-id 123": "non-production allowlist",
    "a1 cd-pipeline run 067 --app 260634 --cr-id 123": "non-production allowlist",
    "a1 cd-pipeline run --app 260634 67 --cr-id 123": "non-production allowlist",
    "bin/a1id -- cd-pipeline run 67 --app=260634 --cr-id=123": "non-production allowlist",
    "pid=67; bin/a1id -- cd-pipeline run $pid --app 260634 --cr-id 123":
        "non-production allowlist",
    "a1 cd-pipeline run 999 --app 260634 --cr-id 123": "non-production allowlist",
    "a1 cd-pipeline run 65 --app 999999 --cr-id 123":
        "non-production allowlist",
    "a1 cd-pipeline run 65 --mvn 43819 --cr-id 123":
        "cannot be proven non-production",
    "a1 cd-pipeline run rerun --pipeline-id 67 --app 260634":
        "non-production allowlist",
    "a1 cd-pipeline run rerun --pipeline-id 420 --pipeline-id 67 --app 260634":
        "non-production allowlist",
    "a1 cd-pipeline run rerun --pipeline-id $pid --app 260634":
        "cannot be proven non-production",
    "a1 cd-pipeline run rerun 3087676553": "cannot be proven non-production",
    "bin/a1id -- cd-pipeline run task status deploy-order-list --task-id 3344115914":
        "task actions",
    "a1 cd-pipeline run cancel 3087676553": "app pipeline exit-cr",
    "a1 app pipeline reenter --pipeline-id 67 --app 260634":
        "non-production allowlist",
    "a1 app pipeline run --pipeline-id 999 --app 260634":
        "non-production allowlist",
    "a1 app pipeline retry --pipeline-id 67 --app 260634":
        "non-production allowlist",
    "a1 app pipeline stage job task status resume --task-id 3344115914":
        "task actions",
    "a1 app pipeline stage job task CR_PUBLISH_RELEASE resume": "task actions",
    "a1 app pipeline lib-deploy start --app-id 260634 --cr-ids 123":
        "explicit numeric non-prod pipeline ID",
    "a1 app pipeline lib-deploy submit --app-id 260634 --flow-inst-id 1":
        "explicit numeric non-prod pipeline ID",
    "a1 cd-pipeline bind 999 --app 260634": "explicit numeric non-prod pipeline ID",
    "a1 app pipeline future-mutation --pipeline-id 420":
        "explicit numeric non-prod pipeline ID",
    "curl -fsS 'https://acube.example/api/createBuildTaskV2' -d '{}'":
        "createBuildTaskV2",
    "url='https://acube.example/api/createBuildTaskV2'; wget -qO- \"$url\"":
        "createBuildTaskV2",
    "/bin/bash -lc 'endpoint=https://acube/api/createBuildTaskV2; curl \"$endpoint\"'":
        "createBuildTaskV2",
    "op=createBuildTask;suffix=V2; bash -lc \"curl https://acube/api/$op$suffix\"":
        "createBuildTaskV2",
    "op=createBuildTask;suffix=V2; sh -c \"wget -qO- https://acube/api/$op$suffix\"":
        "createBuildTaskV2",
    "op=createBuildTask;suffix=V2; bash -lc \"python3 -c 'import requests; requests.post(\\\"https://acube/api/$op$suffix\\\")'\"":
        "createBuildTaskV2",
    "export op=createBuildTask suffix=V2; bash -lc 'curl https://acube/api/$op$suffix'":
        "createBuildTaskV2",
    "export op=createBuildTask suffix=V2; bash -lc 'python3 -c \"import requests; requests.post(\\\"https://acube/api/$op$suffix\\\")\"'":
        "createBuildTaskV2",
    "env op=createBuildTask suffix=V2 sh -c 'wget -qO- https://acube/api/$op$suffix'":
        "createBuildTaskV2",
    "env op=createBuildTask suffix=V2 sh -c 'node -e \"fetch(\\\"https://acube/api/$op$suffix\\\")\"'":
        "createBuildTaskV2",
    "env op=createBuildTask suffix=V2 python3 -c 'import requests; requests.post(\"https://acube/api/$op$suffix\")'":
        "createBuildTaskV2",
    "env op=createBuildTask suffix=V2 node -e 'fetch(\"https://acube/api/$op$suffix\")'":
        "createBuildTaskV2",
    "op=createBuildTask;suffix=V2;curl https://acube/api/$op$suffix":
        "createBuildTaskV2",
    "client=curl;op=createBuildTask;ver=V2;$client https://acube/api/${op}${ver}":
        "createBuildTaskV2",
    "client=wget;prefix=createBuild;middle=Task;suffix=V2;$client -qO- https://acube/api/$prefix$middle$suffix":
        "createBuildTaskV2",
    "python -c 'import requests; requests.post(\"https://acube/api/createBuildTaskV2\")'":
        "createBuildTaskV2",
    "op=createBuildTask;ver=V2;python -c \"import requests; requests.post('https://acube/api/$op$ver')\"":
        "createBuildTaskV2",
    "python3 -c 'import urllib.request; urllib.request.urlopen(\"https://acube/api/createBuildTaskV2\")'":
        "createBuildTaskV2",
    "node -e 'fetch(\"https://acube/api/createBuildTaskV2\", {method:\"POST\"})'":
        "createBuildTaskV2",
    "op=createBuildTask;ver=V2;node -e \"fetch('https://acube/api/$op$ver')\"":
        "createBuildTaskV2",
    "python3 -c 'import subprocess; subprocess.run([\"curl\",\"https://acube/api/createBuildTaskV2\"])'":
        "createBuildTaskV2",
}
for command, marker in blocked.items():
    for tool_name in ("Bash", "exec_command", "mcp__shell__exec_command"):
        event = {"tool_name": tool_name, "tool_input": {"command": command}}
        reason = g.pretool_a1_block_reason(event)
        if not reason or marker not in reason:
            raise SystemExit(
                "expected block for %r/%s, got %r" % (command, tool_name, reason))

allowed = (
    "a1 app cr get 123 --format json",
    "bin/a1id -- app cr get 123 --format json",
    "a1 app cr submit 123 --app 260634 --pipeline-id 420",
    "bin/a1id -- app cr submit 123 --app=260634 --pipeline-id=420",
    "a1 app pipeline status --pipeline-id 67",
    "a1 app cr submit --help",
    "a1 app cr submit 123",
    "a1 cd-pipeline run 420 --app 260634 --cr-id 123",
    "bin/a1id -- cd-pipeline run --app 260634 420 --cr-id 123",
    "a1 cd-pipeline run rerun --pipeline-id 420 --app 260634",
    "a1 cd-pipeline run task status --task-id 3344115914",
    "a1 cd-pipeline run status 67 --app 260634",
    "a1 cd-pipeline run --help",
    "a1 app pipeline reenter --pipeline-id 420 --app 260634",
    "a1 app pipeline run --pipeline-id 66 --app 283346",
    "a1 app pipeline retry --pipeline-id 420 --app 260634",
    "a1 app pipeline stage job task status --task-id 3344115914",
    "a1 app pipeline stage job task CR_PUBLISH_RELEASE status",
    "a1 app pipeline stage job task BUILD log",
    "a1 app pipeline list --app 260634",
    "a1 app pipeline lib-deploy preview --app-id 260634 --flow-inst-id 1",
    "a1 cd-pipeline get 420",
    "a1 cd-pipeline list --app 260634",
    "rg -n exit-cr docs",
    "echo a1 app pipeline exit-cr",
    "printf '%s' 'a1 app cr quit'",
    "command -v a1",
    "rg -n createBuildTaskV2 .",
    "grep -R createBuildTaskV2 docs",
    "printf '%s' 'curl https://acube/api/createBuildTaskV2'",
    "curl https://example.com; printf createBuildTaskV2",
    "python3 -c 'print(\"createBuildTaskV2\")'",
    "export op=createBuildTask suffix=V2; rg -n \"$op$suffix\" .",
    "env op=createBuildTask suffix=V2 grep -R \"$op$suffix\" docs",
    "a1 project workitem create --project 528766 --title ordinary",
    "a1 project workitem relation create --project 528766 --workitem 1 --target 2",
)
for command in allowed:
    event = {"tool_name": "exec_command", "tool_input": {"cmd": command}}
    reason = g.pretool_a1_block_reason(event)
    if reason:
        raise SystemExit("unexpected block for %r: %s" % (command, reason))

calls = []
class ExecCalled(Exception):
    pass
def fake_exec(binary, argv, env):
    calls.append((binary, argv, env))
    raise ExecCalled(argv)
g.os.execvpe = fake_exec

for argv in (
    ["app", "pipeline", "exit-cr", "--pipeline-id", "66"],
    ["--format=json", "app", "pipeline", "quit", "--pipeline-id=66"],
    ["app", "cr", "quit", "123", "--pipeline-id", "66"],
    ["app", "cr", "submit", "123", "--app", "260634", "--pipeline-id", "67"],
    ["--format=json", "app", "cr", "submit", "123", "--app", "260634", "--pipeline-id=67"],
    ["cd-pipeline", "run", "67", "--app", "260634", "--cr-id", "123"],
    ["cd-pipeline", "run", "rerun", "--pipeline-id", "67", "--app", "260634"],
    ["cd-pipeline", "run", "rerun", "3087676553"],
    ["cd-pipeline", "run", "task", "status", "resume", "--task-id", "3344115914"],
):
    try:
        g.run_guarded("a1", argv)
    except g.GuardError as exc:
        assert ("disabled" in str(exc)
                or "non-production allowlist" in str(exc))
    else:
        raise SystemExit("blocked mutation reached the real a1: %r" % argv)
assert not calls

try:
    g.run_guarded("a1", ["app", "cr", "get", "123", "--format", "json"])
except ExecCalled as exc:
    assert exc.args[0] == ["a1", "app", "cr", "get", "123", "--format", "json"]
else:
    raise SystemExit("ordinary a1 command did not pass through")
PY
if [ $? -eq 0 ]; then
    ok "PreTool and runtime guards reject every CR/branch exit"
else
    no "PreTool/runtime guard coverage"
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

# The early wrapper gate must reject before identity migration, whoami, or any
# other real-a1 call, including the formerly allowed targeted quit.
blocked_argv=(
  '-- app pipeline exit-cr --pipeline-id=66'
  '-- app pipeline quit --pipeline-id 66'
  '-- app cr quit 123 --pipeline-id 66'
  '-- app --quiet cr --format json quit 123 --pipeline-id=66'
  'as jarvis -- --quiet app cr quit 123 --pipeline-id=66'
  '-- app cr submit 123 --app 260634 --pipeline-id 67'
  '-- --format=json app cr submit 123 --app 260634 --pipeline-id=67'
  'as jarvis -- app --quiet cr --format json submit 123 --app 260634 --pipeline-id 67'
  '-- cd-pipeline run 67 --app 260634 --cr-id 123'
  '-- cd-pipeline run rerun --pipeline-id 67 --app 260634'
  '-- cd-pipeline run rerun 3087676553'
  '-- cd-pipeline run task status resume --task-id 3344115914'
)
for command in "${blocked_argv[@]}"; do
    rm -rf "$identity_root/identities"
    mkdir -p "$identity_root"
    printf 'live\n' > "$identity_root/auth.yaml"
    : > "$capture"
    # shellcheck disable=SC2086
    run_a1id $command >/dev/null 2>"$tmp/blocked.err"
    rc=$?
    if [ "$rc" -ne 0 ] && [ ! -s "$capture" ] \
        && grep -Eq 'stop and ask a human|human must approve|cannot be proven non-production|task actions are disabled|non-production allowlist' "$tmp/blocked.err"; then
        ok "a1id early-denies with zero real-a1 calls: $command"
    else
        no "a1id deny failed rc=$rc command=$command log=$(cat "$capture")"
    fi
done

mkdir -p "$identity_root/identities/jarvis"
printf 'seeded\n' > "$identity_root/identities/jarvis/auth.yaml"
: > "$capture"
run_a1id -- app cr get 123 --format json >/dev/null 2>"$tmp/ordinary.err"
rc=$?
if [ "$rc" -eq 0 ] && grep -Fq 'ARGS=app cr get 123 --format json' "$capture"; then
    ok "ordinary a1id commands still pass through"
else
    no "ordinary a1 command regression rc=$rc err=$(cat "$tmp/ordinary.err")"
fi

: > "$capture"
run_a1id -- app cr submit 123 --app 260634 --pipeline-id 420 \
    >/dev/null 2>"$tmp/pre-release.err"
rc=$?
if [ "$rc" -eq 0 ] \
    && grep -Fq 'ARGS=app cr submit 123 --app 260634 --pipeline-id 420' "$capture"; then
    ok "CloudSpec pre-release pipeline 420 remains allowed"
else
    no "pre-release pipeline regression rc=$rc err=$(cat "$tmp/pre-release.err")"
fi

: > "$capture"
run_a1id -- cd-pipeline run 420 --app 260634 --cr-id 123 \
    >/dev/null 2>"$tmp/canonical-pre-release.err"
rc=$?
if [ "$rc" -eq 0 ] \
    && grep -Fq 'ARGS=cd-pipeline run 420 --app 260634 --cr-id 123' "$capture"; then
    ok "canonical CloudSpec pre-release pipeline 420 remains allowed"
else
    no "canonical pre-release regression rc=$rc err=$(cat "$tmp/canonical-pre-release.err")"
fi

: > "$capture"
run_a1id -- project workitem create --project 528766 --title ordinary \
    >/dev/null 2>"$tmp/aone-create.err"
rc=$?
if [ "$rc" -eq 0 ] \
    && grep -Fq 'ARGS=project workitem create --project 528766 --title ordinary' "$capture"; then
    ok "ordinary non-post-PR Aone create remains allowed"
else
    no "ordinary Aone create regression rc=$rc err=$(cat "$tmp/aone-create.err")"
fi

: > "$capture"
run_a1id -- project workitem relation add 84846271 relate:84881882 \
    >/dev/null 2>"$tmp/aone-relation.err"
rc=$?
if [ "$rc" -eq 0 ] \
    && grep -Fq 'ARGS=project workitem relation add 84846271 relate:84881882' "$capture"; then
    ok "ordinary non-post-PR Aone relation remains allowed"
else
    no "ordinary Aone relation regression rc=$rc err=$(cat "$tmp/aone-relation.err")"
fi

/usr/bin/python3 -I "$manager" --help >/dev/null 2>&1
if [ $? -eq 0 ]; then
    ok "trusted -I interactive-worker entrypoint imports the guard"
else
    no "interactive-worker -I import smoke"
fi

# ── PreToolUse: raw assignee writes must route through aone-assign.sh ────────
# Regression for tf_customer assignee overwrite (#84486902 / #84955165 /
# #85115148): the route phase told the model to "幂等同步源单 assignee" and it
# called a1 directly, so a human takeover was silently overwritten every
# dispatch. Nothing in the repo wrote --assignee, so a prompt rule alone had no
# enforcement point; this gate is that point.
pretool() {
    python3 "$guard" --check-pretool-command "$1" >/dev/null 2>&1
    echo $?
}
blocked_cmds=(
    "bin/a1id -- project workitem update 85020657 --assignee 484483"
    "bin/a1id -- project workitem update 85020657 --assignee=484483"
    "bin/a1id as terraform-rd -- project workitem update 85020657 --assignee 521957"
    "a1 project workitem update 85020657 --assignee 484483"
)
for cmd in "${blocked_cmds[@]}"; do
    if [ "$(pretool "$cmd")" = "2" ]; then
        ok "raw assignee write blocked: ${cmd:0:52}"
    else
        no "raw assignee write NOT blocked: $cmd"
    fi
done

allowed_cmds=(
    "bash bootstrap/aone-assign.sh 85020657 484483"
    "bin/a1id -- project workitem update 85020657 --status 问题解决中"
    "bin/a1id -- project workitem update 85020657 --tag jarvis-idle"
    "bin/a1id -- project workitem get 85020657"
    "grep -rn -- --assignee bootstrap/"
)
for cmd in "${allowed_cmds[@]}"; do
    if [ "$(pretool "$cmd")" != "2" ]; then
        ok "not over-blocked: ${cmd:0:52}"
    else
        no "over-blocked: $cmd"
    fi
done

echo "a1_command_guard_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
