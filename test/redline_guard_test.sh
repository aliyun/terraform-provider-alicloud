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

expect_at() { # <desc> <want-exit> <cwd> <command>
    local desc="$1" want="$2" run_cwd="$3" cmd="$4" got
    got="$(cd "$run_cwd" && run_guard "$cmd")"
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
expect "outer split variables survive bash -lc curl" 2 \
    "op=createBuildTask;suffix=V2; bash -lc \"curl https://acube.example/api/\$op\$suffix\""
expect "outer split variables survive sh -c wget" 2 \
    "op=createBuildTask;suffix=V2; sh -c \"wget -qO- https://acube.example/api/\$op\$suffix\""
expect "outer split variables survive bash -lc python" 2 \
    "op=createBuildTask;suffix=V2; bash -lc \"python3 -c 'import requests; requests.post(\\\"https://acube.example/api/\$op\$suffix\\\")'\""
expect "exported split variables survive bash -lc curl" 2 \
    "export op=createBuildTask suffix=V2; bash -lc 'curl https://acube.example/api/\$op\$suffix'"
expect "exported split variables survive bash -lc python" 2 \
    "export op=createBuildTask suffix=V2; bash -lc 'python3 -c \"import requests; requests.post(\\\"https://acube.example/api/\$op\$suffix\\\")\"'"
expect "env split variables survive sh -c wget" 2 \
    "env op=createBuildTask suffix=V2 sh -c 'wget -qO- https://acube.example/api/\$op\$suffix'"
expect "env split variables survive sh -c node" 2 \
    "env op=createBuildTask suffix=V2 sh -c 'node -e \"fetch(\\\"https://acube.example/api/\$op\$suffix\\\")\"'"
expect "env split variables reach direct python" 2 \
    "env op=createBuildTask suffix=V2 python3 -c 'import requests; requests.post(\"https://acube.example/api/\$op\$suffix\")'"
expect "env split variables reach direct node" 2 \
    "env op=createBuildTask suffix=V2 node -e 'fetch(\"https://acube.example/api/\$op\$suffix\")'"
expect "split variable createBuildTaskV2 curl blocked" 2 \
    "op=createBuildTask;suffix=V2;curl https://acube.example/api/\$op\$suffix"
expect "variable client and split API name blocked" 2 \
    "client=curl;op=createBuildTask;suffix=V2;\$client https://acube.example/api/\${op}\${suffix}"
expect "split API name in python request blocked" 2 \
    "op=createBuildTask;suffix=V2;python3 -c \"import requests; requests.post('https://acube.example/api/\$op\$suffix')\""
expect "split API name in node fetch blocked" 2 \
    "op=createBuildTask;suffix=V2;node -e \"fetch('https://acube.example/api/\$op\$suffix')\""
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
expect "exported split marker rg audit allowed" 0 \
    "export op=createBuildTask suffix=V2; rg -n \"\$op\$suffix\" ."
expect "env split marker grep audit allowed" 0 \
    "env op=createBuildTask suffix=V2 grep -R \"\$op\$suffix\" docs"
expect "ordinary Aone create remains allowed" 0 \
    "a1 project workitem create --project 528766 --title ordinary"
expect "ordinary Aone relation remains allowed" 0 \
    "a1 project workitem relation create --project 528766 --workitem 1 --target 2"

# --- CloudSpec local-resource-test permanent red line -----------------------
expect "direct aliyun cspec test blocked" 2 \
    "aliyun cspec test --name Example"
expect "absolute aliyun cspec test blocked" 2 \
    "/usr/local/bin/aliyun cspec test"
expect "env/nohup wrapper blocked" 2 \
    "env PROFILE=test nohup aliyun cspec test --name Example &"
expect "time wrapper blocked" 2 \
    "time -p aliyun cspec test"
expect "command wrapper blocked" 2 \
    "command aliyun cspec test"
expect "aliyun spaced profile flag blocked" 2 \
    "aliyun --profile test cspec test"
expect "aliyun equals profile flag blocked" 2 \
    "aliyun --profile=test cspec test"
expect "aliyun short profile flag blocked" 2 \
    "aliyun -p test cspec test"
expect "aliyun source profile and sts flags blocked" 2 \
    "aliyun --source-profile base --sts-region cn-hangzhou --sts-endpoint sts.aliyuncs.com cspec test"
expect "aliyun endpoint and timeout flags blocked" 2 \
    "aliyun -e cspec.example --connect-timeout 5 --read-timeout=30 cspec test"
expect "bash -lc wrapper blocked" 2 \
    "bash -lc 'aliyun cspec test --name Example'"
expect "simple variable wrapper blocked" 2 \
    "cli=aliyun; section=cspec; action=test; \$cli \$section \$action"
expect "bare dollar command substitution blocked" 2 \
    'result=$(aliyun cspec test --name Example)'
expect "double-quoted dollar command substitution blocked" 2 \
    'printf "%s\n" "$(aliyun --profile test cspec test)"'
expect "bare backtick command substitution blocked" 2 \
    'result=`aliyun cspec test --name Example`'
expect "nested command substitution blocked" 2 \
    'result=$(printf "%s" "$(aliyun cspec test)")'
expect "sudo wrapper blocked" 2 \
    "sudo -u nobody aliyun cspec test"
expect "timeout wrapper blocked" 2 \
    "timeout --signal TERM 30s aliyun cspec test"
expect "xargs wrapper blocked" 2 \
    "printf x | xargs -n1 aliyun cspec test"
expect "parallel wrapper blocked" 2 \
    "parallel --jobs 2 aliyun cspec test ::: one two"
expect "env split-string wrapper blocked" 2 \
    "env -S 'aliyun --profile test cspec test'"
expect "JARVIS_MASTER_OK cannot bypass aliyun cspec test" 2 \
    "aliyun cspec test" JARVIS_MASTER_OK=1
expect "rg aliyun cspec test audit allowed" 0 \
    "rg -n 'aliyun cspec test' ."
expect "grep aliyun cspec test audit allowed" 0 \
    "grep -R 'aliyun cspec test' docs"
expect "printf aliyun cspec test audit allowed" 0 \
    "printf '%s' 'aliyun cspec test'"
expect "single-quoted dollar substitution audit allowed" 0 \
    "printf '%s' '\$(aliyun cspec test)'"
expect "single-quoted backtick audit allowed" 0 \
    "printf '%s' '\`aliyun cspec test\`'"
expect "escaped dollar substitution audit allowed" 0 \
    'printf "%s" "\$(aliyun cspec test)"'
expect "escaped backtick substitution audit allowed" 0 \
    'printf "%s" "\`aliyun cspec test\`"'
expect "xargs echo audit allowed" 0 \
    "printf 'aliyun cspec test' | xargs echo"
expect "parallel printf audit allowed" 0 \
    "parallel printf ::: 'aliyun cspec test'"
expect "env split-string printf audit allowed" 0 \
    "env -S 'printf aliyun-cspec-test'"
expect "sudo rg audit allowed" 0 \
    "sudo rg -n 'aliyun cspec test' docs"
expect "go test remains allowed" 0 \
    "go test ./alicloud -run TestAccExample"
expect "aliyun cspec build remains allowed" 0 \
    "aliyun cspec build"
expect "aliyun cspec check remains allowed" 0 \
    "aliyun cspec check --name Example"
expect "missing classifier still blocks direct aliyun cspec test" 2 \
    "aliyun --profile test cspec test" \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"
expect "missing classifier blocks dollar command substitution" 2 \
    'result=$(aliyun cspec test)' \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"
expect "missing classifier blocks backtick command substitution" 2 \
    'result=`aliyun cspec test`' \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"
expect "missing classifier still allows printf audit" 0 \
    "printf '%s' 'aliyun cspec test'" \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"

# --- direct AMP permanent red line ------------------------------------------
expect "direct amp blocked" 2 \
    "amp doctor -o json"
expect "absolute amp blocked" 2 \
    "/Users/example/.local/bin/amp publish pre"
expect "direct amp_safe shebang blocked" 2 \
    "$repo_root/bootstrap/amp_safe.py doctor"
expect "non-isolated python amp_safe blocked" 2 \
    "/usr/bin/python3 $repo_root/bootstrap/amp_safe.py doctor"
expect "PATH python amp_safe blocked" 2 \
    "python3 -I $repo_root/bootstrap/amp_safe.py doctor"
expect "same-name copied amp_safe blocked" 2 \
    "/usr/bin/python3 -I /tmp/amp_safe.py doctor"
expect "env command amp blocked" 2 \
    "env PROFILE=test command amp branch get"
expect "variable-built amp blocked" 2 \
    "cli=amp; action=doctor; \$cli \$action"
expect "nested shell amp blocked" 2 \
    "bash -lc 'amp publish daily --dry-run'"
expect "command substitution amp blocked" 2 \
    'result=$(amp api list -o json)'
expect "xargs amp blocked" 2 \
    "printf x | xargs amp doctor"
expect "python subprocess amp blocked" 2 \
    "python3 -c 'import subprocess; subprocess.run([\"amp\",\"doctor\"])'"
expect "node child process amp blocked" 2 \
    "node -e 'require(\"child_process\").spawn(\"amp\",[\"doctor\"])'"
expect "eval amp blocked" 2 \
    "eval 'amp doctor'"
expect "find exec amp blocked" 2 \
    "find . -maxdepth 0 -exec amp doctor ';'"
expect "awk system amp blocked" 2 \
    "awk 'BEGIN { system(\"amp doctor\") }'"
printf '%s\n' '#!/usr/bin/env bash' 'amp doctor' > "$tmp/run-amp.sh"
chmod +x "$tmp/run-amp.sh"
expect "shell script amp blocked" 2 \
    "bash $tmp/run-amp.sh"
expect "sourced shell script amp blocked" 2 \
    "source $tmp/run-amp.sh"
expect_at "relative shell script amp blocked" 2 "$tmp" \
    "bash ./run-amp.sh"
expect_at "relative sourced script amp blocked" 2 "$tmp" \
    "source ./run-amp.sh"
expect_at "relative direct shebang script amp blocked" 2 "$tmp" \
    "./run-amp.sh"
expect "absolute direct shebang script amp blocked" 2 \
    "$tmp/run-amp.sh"
expect "python concatenated subprocess amp blocked" 2 \
    "python3 -c 'import subprocess; subprocess.run([\"am\"+\"p\",\"doctor\"])'"
expect "trusted amp wrapper command allowed" 0 \
    "/usr/bin/python3 -I $repo_root/bootstrap/amp_safe.py doctor"
expect "trusted wrapper cannot mask chained raw amp" 2 \
    "/usr/bin/python3 -I $repo_root/bootstrap/amp_safe.py doctor; amp doctor"
expect "amp text audit allowed" 0 \
    "rg -n 'amp publish' ."
expect "missing classifier still blocks direct amp" 2 \
    "amp doctor" JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"
expect "missing classifier allows quoted substitution audit" 0 \
    "printf '%s' '\$(aliyun cspec test)'" \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"

printf '%s\n' '#!/usr/bin/env python3' 'raise SystemExit(70)' > "$tmp/crashing-guard.py"
chmod +x "$tmp/crashing-guard.py"
expect "crashing classifier still blocks direct cspec test" 2 \
    "sudo aliyun --source-profile base cspec test" \
    JARVIS_A1_COMMAND_GUARD="$tmp/crashing-guard.py"
expect "crashing classifier still blocks command substitution" 2 \
    'result=$(aliyun cspec test)' \
    JARVIS_A1_COMMAND_GUARD="$tmp/crashing-guard.py"
expect "crashing classifier still allows rg audit" 0 \
    "rg -n 'aliyun cspec test' docs" \
    JARVIS_A1_COMMAND_GUARD="$tmp/crashing-guard.py"

# --- production delivery hard gate -----------------------------------------
expect "legacy pipeline 67 submit blocked" 2 \
    "a1 app cr submit 123 --pipeline-id 67"
expect "canonical CloudSpec pipeline 67 run blocked" 2 \
    "bin/a1id -- cd-pipeline run 67 --app 260634 --cr-id 123"
expect "JARVIS_MASTER_OK cannot bypass pipeline 67" 2 \
    "bin/a1id -- cd-pipeline run 67 --app 260634 --cr-id 123" JARVIS_MASTER_OK=1
expect "CloudSpec pre-release pipeline 420 allowed" 0 \
    "bin/a1id -- cd-pipeline run 420 --app 260634 --cr-id 123"
expect "unknown future pipeline fails closed" 2 \
    "bin/a1id -- cd-pipeline run 999 --app 260634 --cr-id 123"
expect "missing classifier blocks a1 fail-closed" 2 \
    "bin/a1id -- cd-pipeline run 67 --app 260634 --cr-id 123" \
    JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"
expect "missing classifier does not block unrelated shell" 0 \
    "git status --short" JARVIS_A1_COMMAND_GUARD="$tmp/missing-guard.py"

echo
echo "pass=$pass fail=$fail"
[ "$fail" -eq 0 ]
