#!/usr/bin/env bash
# Hermetic coverage for bootstrap/redline-guard.sh: master-push red line,
# push-master-allowlist data-repo exemption (fail-closed), upstream slug
# block, catastrophic rm block, JARVIS_MASTER_OK escape hatch.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
guard="$repo_root/bootstrap/redline-guard.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

pass=0
fail=0
ok(){ echo "PASS: $1"; pass=$((pass + 1)); }
no(){ echo "FAIL: $1" >&2; fail=$((fail + 1)); }

mkrepo() { # <dir> <origin-url>
    git init -q "$1" && git -C "$1" remote add origin "$2"
}
mkrepo "$tmp/playground" "git@gitlab.alibaba-inc.com:terraflow/tf_playground.git"
mkrepo "$tmp/jarvis"     "git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git"
mkrepo "$tmp/evil"       "git@gitlab.alibaba-inc.com:terraflow/tf_playground_evil.git"

run_guard() { # <command> [env KEY=VAL...]; echoes exit code
    local cmd="$1"; shift
    printf '%s' "$(jq -n --arg c "$cmd" '{tool_name:"Bash",tool_input:{command:$c}}')" \
        | env -u JARVIS_MASTER_OK "$@" bash "$guard" >/dev/null 2>&1
    echo $?
}

expect() { # <desc> <want-exit> <command> [env...]
    local desc="$1" want="$2"; shift 2
    local got; got="$(run_guard "$@")"
    if [ "$got" = "$want" ]; then ok "$desc"; else no "$desc (want exit $want, got $got)"; fi
}

# --- master push red line ----------------------------------------------------
expect "bare master push blocked (cwd unknown, fail-closed)" 2 \
    "git push origin master"
expect "-C non-allowlisted repo master push blocked" 2 \
    "git -C $tmp/jarvis push origin master"
expect "slug over-match still blocked (tf_playground_evil)" 2 \
    "git -C $tmp/evil push origin master"
expect "cd with intermediate command stays blocked (fail-closed)" 2 \
    "cd $tmp/playground; echo hi && git push origin master"
expect "-C nonexistent dir stays blocked" 2 \
    "git -C $tmp/nope push origin master"

# --- allowlisted data repo ---------------------------------------------------
expect "-C allowlisted playground master push allowed" 0 \
    "git -C $tmp/playground push origin master"
expect "quoted -C allowlisted playground allowed" 0 \
    "git -C \"$tmp/playground\" push origin master"
expect "cd && push allowlisted playground allowed" 0 \
    "cd $tmp/playground && git push origin master"
expect "HEAD:master refspec on allowlisted repo allowed" 0 \
    "git -C $tmp/playground push origin HEAD:master"
expect "unknown remote name on allowlisted repo blocked" 2 \
    "git -C $tmp/playground push upstream master"

# --- unaffected paths --------------------------------------------------------
expect "feature branch push passes" 0 \
    "git push origin worktree-some-feature"
expect "upstream aliyun slug still blocked" 2 \
    "git push git@github.com:aliyun/terraform-provider-alicloud.git feature"
expect "JARVIS_MASTER_OK=1 escape hatch allows" 0 \
    "git push origin master" JARVIS_MASTER_OK=1
expect "rm -rf / still blocked" 2 \
    "rm -rf /"
expect "repo-scoped rm -rf passes" 0 \
    "rm -rf $tmp/playground/scratch"

# --- Acube downstream task red line -----------------------------------------
expect "direct createBuildTaskV2 curl blocked" 2 \
    "curl -fsS https://acube.example/api/createBuildTaskV2 -d '{}'"
expect "variable-wrapped createBuildTaskV2 wget blocked" 2 \
    "endpoint=https://acube.example/api/createBuildTaskV2; wget -qO- \"\$endpoint\""
expect "bash -lc createBuildTaskV2 curl blocked" 2 \
    "/bin/bash -lc 'endpoint=https://acube.example/api/createBuildTaskV2; curl \"\$endpoint\"'"
expect "split variable createBuildTaskV2 curl blocked" 2 \
    "op=createBuildTask;suffix=V2;curl https://acube.example/api/\$op\$suffix"
expect "python requests createBuildTaskV2 blocked" 2 \
    "python3 -c 'import requests; requests.post(\"https://acube.example/api/createBuildTaskV2\")'"
expect "node fetch createBuildTaskV2 blocked" 2 \
    "node -e 'fetch(\"https://acube.example/api/createBuildTaskV2\", {method:\"POST\"})'"
expect "JARVIS_MASTER_OK cannot bypass createBuildTaskV2" 2 \
    "curl -fsS https://acube.example/api/createBuildTaskV2" JARVIS_MASTER_OK=1
expect "rg createBuildTaskV2 audit allowed" 0 \
    "rg -n createBuildTaskV2 ."
expect "grep createBuildTaskV2 audit allowed" 0 \
    "grep -R createBuildTaskV2 docs"
expect "unrelated curl then marker printf allowed" 0 \
    "curl https://example.com; printf createBuildTaskV2"
expect "python marker print audit allowed" 0 \
    "python3 -c 'print(\"createBuildTaskV2\")'"

echo
echo "pass=$pass fail=$fail"
[ "$fail" -eq 0 ]
