#!/usr/bin/env bash
# Ensure the alicloud provider repo is synced. Path is resolved via
# bootstrap/workspace.sh dir terraform_provider(base 配置不存绝对路径,
# 本机覆盖走 workspaces.local.json / JARVIS_WORKSPACE_ROOT). No repo -> clone;
# existing repo -> fetch + reset --hard FETCH_HEAD (强制对齐 upstream HEAD).
#
# workspace 定位:只读查证镜像。开发/评审改动必须走 worktree(CLAUDE.md 纪律 1);
# 主目录任何 in-progress 状态都会被本脚本的 reset --hard 清掉。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
REPO_DIR="$(bash "$ROOT/bootstrap/workspace.sh" dir terraform_provider)"
if [ -z "$REPO_DIR" ] || [ "$REPO_DIR" = "null" ]; then
  echo "[sync-provider] cannot resolve terraform_provider workspace dir (bootstrap/workspace.sh dir terraform_provider)" >&2
  exit 1
fi
REMOTE="https://github.com/aliyun/terraform-provider-alicloud.git"

if [ ! -d "$REPO_DIR/.git" ]; then
  echo "[sync-provider] cloning into $REPO_DIR ..."
  mkdir -p "$(dirname "$REPO_DIR")"
  git clone --depth 1 "$REMOTE" "$REPO_DIR"
else
  echo "[sync-provider] updating $REPO_DIR ..."
  DEF=$(git -C "$REPO_DIR" remote show origin | sed -n 's/.*HEAD branch: //p')
  DEF="${DEF:-master}"
  git -C "$REPO_DIR" fetch --depth 1 origin "$DEF"
  # 强制对齐 upstream HEAD——避免 workspace 的本地 branch 停留在旧 HEAD 而 grep 到 stale 代码
  git -C "$REPO_DIR" reset --hard FETCH_HEAD
fi
echo "[sync-provider] ready: $REPO_DIR @ $(git -C "$REPO_DIR" rev-parse --short HEAD)"
