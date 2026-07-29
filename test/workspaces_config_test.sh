#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

jq -e '.workspaces.terraform_generator_v4.git_url=="git@gitlab.alibaba-inc.com:opensource-tools/terraform-generator-v4.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.terraform_generator_v4.default_branch=="main"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.jarvis.git_url=="git@gitlab.alibaba-inc.com:terraflow/jarvis-preview.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.jarvis.default_branch=="master"' \
  "$repo_root/config/workspaces.json" >/dev/null

# AutoWonder 项目事实坐标及默认分支。
jq -e '
  .workspaces.auto_wonder.repo == "auto-wonder" and
  .workspaces.auto_wonder.git_url == "git@gitlab.alibaba-inc.com:sdlc-autopilot/auto-wonder.git" and
  .workspaces.auto_wonder.app == 341827 and
  .workspaces.auto_wonder.project == 2087214 and
  .workspaces.auto_wonder.default_branch == "master" and
  .workspaces.auto_wonder.ops == {} and
  .workspaces.auto_wonder.desc == "AutoWonder 自驱研发平台"
' "$repo_root/config/workspaces.json" >/dev/null

mkdir -p "$tmpdir/auto-wonder"
resolved_auto_wonder="$(JARVIS_WORKSPACES_LOCAL=none JARVIS_WORKSPACE_ROOT="$tmpdir" \
  bash "$repo_root/bootstrap/workspace.sh" dir auto_wonder)"
if [ "$resolved_auto_wonder" != "$tmpdir/auto-wonder" ]; then
  echo "auto_wonder: expected $tmpdir/auto-wonder, got $resolved_auto_wonder" >&2
  exit 1
fi

# AutoWonder 客户端运行时仓库事实坐标及默认分支。
jq -e '
  .workspaces.auto_wonder_client_runtime.repo == "auto-wonder-client-runtime" and
  .workspaces.auto_wonder_client_runtime.git_url == "git@gitlab.alibaba-inc.com:sdlc-autopilot/auto-wonder-client-runtime.git" and
  .workspaces.auto_wonder_client_runtime.default_branch == "master" and
  .workspaces.auto_wonder_client_runtime.ops == {} and
  .workspaces.auto_wonder_client_runtime.desc == "AutoWonder 客户端运行时"
' "$repo_root/config/workspaces.json" >/dev/null

mkdir -p "$tmpdir/auto-wonder-client-runtime"
resolved_auto_wonder_client_runtime="$(JARVIS_WORKSPACES_LOCAL=none JARVIS_WORKSPACE_ROOT="$tmpdir" \
  bash "$repo_root/bootstrap/workspace.sh" dir auto_wonder_client_runtime)"
if [ "$resolved_auto_wonder_client_runtime" != "$tmpdir/auto-wonder-client-runtime" ]; then
  echo "auto_wonder_client_runtime: expected $tmpdir/auto-wonder-client-runtime, got $resolved_auto_wonder_client_runtime" >&2
  exit 1
fi

# 自动化服务台六个交付仓库：事实坐标、默认分支和池归属必须完整登记。
jq -e '
  def platform_workspace($key; $repo; $url):
    .workspaces[$key] as $w |
    $w.repo == $repo and
    $w.git_url == $url and
    $w.default_branch == "master" and
    (($w.pools // []) | index("automation_platform") != null);

  platform_workspace("automation_platform"; "aliyun-automation-platform"; "git@gitlab.alibaba-inc.com:aliyun-automation-platform/aliyun-automation-platform.git") and
  platform_workspace("automation_platform_frontend"; "iac-service"; "git@gitlab.alibaba-inc.com:aliyun-api/iac-service.git") and
  platform_workspace("automation_platform_runtime"; "iac-service-runtime"; "git@gitlab.alibaba-inc.com:opensource-tools/iac-service-runtime.git") and
  platform_workspace("automation_platform_function_test"; "automation-function-test"; "git@gitlab.alibaba-inc.com:aliyun-automation-platform/automation-function-test.git") and
  platform_workspace("automation_platform_api"; "IaCService_pop_IaCService_2021-08-06"; "git@gitlab.alibaba-inc.com:cloudspec-model/IaCService_pop_IaCService_2021-08-06.git") and
  platform_workspace("automation_platform_api_inner"; "IaCService-inner_pop_IaCService-inner_2021-09-01"; "git@gitlab.alibaba-inc.com:cloudspec-model/IaCService-inner_pop_IaCService-inner_2021-09-01.git") and

  .workspaces.automation_platform.project == 1091779 and
  .workspaces.automation_platform.app == 172823 and
  .workspaces.automation_platform.repo_id == 2156624 and
  .workspaces.automation_platform.pipelines == {"prestage":66,"prod":67} and
  .workspaces.automation_platform.delivery == "delivery-aliyun-automation-platform.md" and

  .workspaces.automation_platform.ops == {"build":"./mvnw -q -DskipTests package","test":"./mvnw -q test"} and
  .workspaces.automation_platform_frontend.ops == {"build":"npm run build","test":"npm test -- --runInBand","lint":"npm run lint"} and
  .workspaces.automation_platform_runtime.ops == {"build":"go build ./...","test":"go test ./..."} and
  .workspaces.automation_platform_function_test.ops == {"test":"mvn -q test"} and
  .workspaces.automation_platform_api.ops == {"build":"cloudspec build"} and
  .workspaces.automation_platform_api_inner.ops == {"build":"cloudspec build"}
' "$repo_root/config/workspaces.json" >/dev/null

# 池边界必须恰好是六个产品仓库，Agent 仓库不能被误纳入；交付坐标只属于后端。
jq -e '
  ([.workspaces | to_entries[] |
      select(((.value.pools // []) | index("automation_platform")) != null) |
      .key] | sort) == [
        "automation_platform",
        "automation_platform_api",
        "automation_platform_api_inner",
        "automation_platform_frontend",
        "automation_platform_function_test",
        "automation_platform_runtime"
      ] and
  ([.workspaces | to_entries[] |
      select(.key != "automation_platform") |
      select(((.value.pools // []) | index("automation_platform")) != null) |
      .value |
      has("project") or has("app") or has("repo_id") or has("pipelines") or has("delivery")]
    | any | not)
' "$repo_root/config/workspaces.json" >/dev/null

while read -r workspace_key repo_dir; do
  mkdir -p "$tmpdir/$repo_dir"
  resolved_platform="$(JARVIS_WORKSPACE_ROOT="$tmpdir" bash "$repo_root/bootstrap/workspace.sh" dir "$workspace_key")"
  if [ "$resolved_platform" != "$tmpdir/$repo_dir" ]; then
    echo "$workspace_key: expected $tmpdir/$repo_dir, got $resolved_platform" >&2
    exit 1
  fi
done <<'WORKSPACES'
automation_platform aliyun-automation-platform
automation_platform_frontend iac-service
automation_platform_runtime iac-service-runtime
automation_platform_function_test automation-function-test
automation_platform_api IaCService_pop_IaCService_2021-08-06
automation_platform_api_inner IaCService-inner_pop_IaCService-inner_2021-09-01
WORKSPACES

# terraform_provider remote registration must reflect the real machine layout:
# origin = upstream aliyun, fork = api-tool-agent (F1 fix — the old ChenHanZhang /
# upstream_remote=alicloud registration was stale and contradicted sync-provider.sh).
jq -e '.workspaces.terraform_provider.git_url=="https://github.com/aliyun/terraform-provider-alicloud.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.terraform_provider.upstream_remote=="origin"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.terraform_provider.fork_remote=="fork"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.terraform_provider.jarvis_github_login=="api-tool-agent"' \
  "$repo_root/config/workspaces.json" >/dev/null

# tf_playground 数据仓登记(git 化场景语料库,直推 master 模型;probe.sh probe_playground_dir 会解析)
jq -e '.workspaces.tf_playground.git_url=="git@gitlab.alibaba-inc.com:terraflow/tf_playground.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.tf_playground.default_branch=="master"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.tf_playground.repo=="tf_playground"' \
  "$repo_root/config/workspaces.json" >/dev/null

# fixture 段:关闭 workspaces.local.json 合并(用 JARVIS_WORKSPACES_LOCAL=none)。
# 本机若有 local.json(gitignored)把 tf_playground/terraform_provider 等指向别处,
# 会绕过 JARVIS_WORKSPACE_ROOT 的 fixture 目录导致断言飘。工作纪律见 workspace.sh 头注释。
mkdir -p "$tmpdir/tf_playground"
resolved_pg="$(JARVIS_WORKSPACES_LOCAL=none JARVIS_WORKSPACE_ROOT="$tmpdir" \
  bash "$repo_root/bootstrap/workspace.sh" dir tf_playground)"
if [ "$resolved_pg" != "$tmpdir/tf_playground" ]; then
  echo "expected $tmpdir/tf_playground, got $resolved_pg" >&2
  exit 1
fi

mkdir -p "$tmpdir/terraform-generator-v4"
resolved="$(JARVIS_WORKSPACES_LOCAL=none JARVIS_WORKSPACE_ROOT="$tmpdir" \
  bash "$repo_root/bootstrap/workspace.sh" dir terraform_generator_v4)"

if [ "$resolved" != "$tmpdir/terraform-generator-v4" ]; then
  echo "expected $tmpdir/terraform-generator-v4, got $resolved" >&2
  exit 1
fi

mkdir -p "$tmpdir/jarvis"
resolved_jarvis="$(JARVIS_WORKSPACES_LOCAL=none JARVIS_WORKSPACE_ROOT="$tmpdir" \
  bash "$repo_root/bootstrap/workspace.sh" dir jarvis)"

if [ "$resolved_jarvis" != "$tmpdir/jarvis" ]; then
  echo "expected $tmpdir/jarvis, got $resolved_jarvis" >&2
  exit 1
fi

# --- worktree fallback for gitignored workspaces.local.json ------------------
# config/workspaces.local.json is gitignored, so it does NOT propagate into a
# `git worktree add` checkout — only the tracked workspaces.json comes along.
# workspace.sh must therefore fall back to the PRIMARY repo's local.json
# (resolved via `git rev-parse --git-common-dir`) when run from a worktree, or the
# machine-local path override is silently lost inside every worktree (bit us twice).
# Simulate with a throwaway git repo + a real `git worktree add`.
wt_main="$tmpdir/wtmain"
mkdir -p "$wt_main/config" "$wt_main/bootstrap"
git -C "$wt_main" init -q
git -C "$wt_main" config core.hooksPath /dev/null   # hermetic: ignore the machine's global git hooks
git -C "$wt_main" config user.email t@example.com
git -C "$wt_main" config user.name tester
# base (tracked): 'demo' has NO path/git_url — resolution needs the local override.
printf '%s\n' '{"workspaces":{"demo":{"repo":"demo-repo"}}}' > "$wt_main/config/workspaces.json"
cp "$repo_root/bootstrap/workspace.sh" "$wt_main/bootstrap/workspace.sh"
printf '%s\n' 'config/workspaces.local.json' > "$wt_main/.gitignore"
git -C "$wt_main" add -A
git -C "$wt_main" commit -q -m init
# local override (UNtracked / gitignored): points demo at a real on-disk dir.
demo_target="$tmpdir/demo-on-disk"
mkdir -p "$demo_target"
printf '{"workspaces":{"demo":{"path":"%s"}}}\n' "$demo_target" > "$wt_main/config/workspaces.local.json"
# add a worktree — tracked files come along, the gitignored local.json does not.
wt_child="$tmpdir/wtchild"
git -C "$wt_main" worktree add -q "$wt_child" -b wt-test
if [ -f "$wt_child/config/workspaces.local.json" ]; then
  echo "sandbox invariant broken: local.json should be absent inside the worktree" >&2
  exit 1
fi
# From the worktree, workspace.sh must fall back to main's local.json and resolve
# demo → $demo_target. Without the fallback it sees base alone (demo has no path,
# no git_url) → exits 4 (missing_capability) → empty stdout. Point JARVIS_WORKSPACE_ROOT
# at a nonexistent dir so branch (b) ROOT/repo can never accidentally match.
resolved_wt="$(JARVIS_ROOT="$wt_child" JARVIS_WORKSPACE_ROOT="$tmpdir/noexist" \
  bash "$wt_child/bootstrap/workspace.sh" dir demo 2>/dev/null || true)"
if [ "$resolved_wt" != "$demo_target" ]; then
  echo "worktree fallback: expected $demo_target, got '$resolved_wt'" >&2
  exit 1
fi

# Control: from the MAIN checkout (local.json present) resolution is unchanged and
# the fallback must NOT alter behavior.
resolved_main="$(JARVIS_ROOT="$wt_main" JARVIS_WORKSPACE_ROOT="$tmpdir/noexist" \
  bash "$wt_main/bootstrap/workspace.sh" dir demo 2>/dev/null || true)"
if [ "$resolved_main" != "$demo_target" ]; then
  echo "main-checkout control: expected $demo_target, got '$resolved_main'" >&2
  exit 1
fi

echo "workspaces_config_test: PASS"
