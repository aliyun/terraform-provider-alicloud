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
jq -e '.workspaces.jarvis.git_url=="git@gitlab.alibaba-inc.com:terraflow/jarvis.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.jarvis.default_branch=="master"' \
  "$repo_root/config/workspaces.json" >/dev/null

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
