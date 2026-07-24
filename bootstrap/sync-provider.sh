#!/usr/bin/env bash
# bootstrap/sync-provider.sh — alicloud provider 只读查证镜像同步(真单点)。
# 消费方:.claude/skills/{aone-triage,terraform-pr-review}/scripts/sync-provider.sh
# 薄 wrapper(及其 .agents 镜像),真实逻辑只在本文件维护。
#
# 路径经 bootstrap/workspace.sh dir terraform_provider 解析(base 配置不存绝对路径,
# 本机覆盖走 workspaces.local.json / JARVIS_WORKSPACE_ROOT)。
# 无库 -> clone;有库 -> fetch + reset --hard FETCH_HEAD(强制对齐 upstream HEAD)。
#
# workspace 定位:只读查证镜像。开发/评审改动必须走 worktree(CLAUDE.md 纪律 1);
# 主目录任何 in-progress 状态都会被本脚本的 reset --hard 清掉。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPO_DIR="$(bash "$ROOT/bootstrap/workspace.sh" dir terraform_provider)"
if [ -z "$REPO_DIR" ] || [ "$REPO_DIR" = "null" ]; then
  echo "[sync-provider] cannot resolve terraform_provider workspace dir (bootstrap/workspace.sh dir terraform_provider)" >&2
  exit 1
fi
# 同步源解析(按机器角色分流):
#   1) JARVIS_PROVIDER_GIT_URL 显式覆写(单机逃生阀)
#   2) role=worker → workspaces.json .mirror_git_url(内部镜像;IDC worker 够不到
#      github,macmini 泵 provider-mirror-sync.sh 每 10 分钟保鲜)
#   3) 其余(macmini/开发机) → .git_url(github 真源,保持直拉最新)
# role 优先取 env;缺省时读 bridge/jarvis.env 落盘值,手动 shell 也能解析对。
ROLE="${JARVIS_BRIDGE_ROLE:-}"
if [ -z "$ROLE" ] && [ -f "$ROOT/bridge/jarvis.env" ]; then
  ROLE=$(grep -E '^(export )?JARVIS_BRIDGE_ROLE=' "$ROOT/bridge/jarvis.env" | tail -1 | sed 's/.*=//' | tr -d '"' || true)
fi
if [ -n "${JARVIS_PROVIDER_GIT_URL:-}" ]; then
  REMOTE="$JARVIS_PROVIDER_GIT_URL"
elif [ "$ROLE" = "worker" ]; then
  REMOTE=$(jq -r '.workspaces.terraform_provider.mirror_git_url // empty' "$ROOT/config/workspaces.json" 2>/dev/null)
fi
REMOTE="${REMOTE:-$(jq -r '.workspaces.terraform_provider.git_url // empty' "$ROOT/config/workspaces.json" 2>/dev/null)}"
REMOTE="${REMOTE:-https://github.com/aliyun/terraform-provider-alicloud.git}"

if [ ! -d "$REPO_DIR/.git" ]; then
  echo "[sync-provider] cloning into $REPO_DIR (from $REMOTE) ..."
  mkdir -p "$(dirname "$REPO_DIR")"
  git clone --depth 1 "$REMOTE" "$REPO_DIR"
else
  echo "[sync-provider] updating $REPO_DIR ..."
  # 源切换(github→内部镜像)后旧 clone 的 origin URL 会过期——对齐到当前解析值
  cur_url=$(git -C "$REPO_DIR" remote get-url origin 2>/dev/null || echo "")
  if [ -n "$REMOTE" ] && [ "$cur_url" != "$REMOTE" ]; then
    echo "[sync-provider] origin url → $REMOTE"
    git -C "$REPO_DIR" remote set-url origin "$REMOTE"
  fi
  DEF=$(git -C "$REPO_DIR" remote show origin | sed -n 's/.*HEAD branch: //p')
  DEF="${DEF:-master}"
  git -C "$REPO_DIR" fetch --depth 1 origin "$DEF"
  # 强制对齐 upstream HEAD——避免 workspace 的本地 branch 停留在旧 HEAD 而 grep 到 stale 代码
  git -C "$REPO_DIR" reset --hard FETCH_HEAD
fi

# 镜像新鲜度守卫: 泵每次同步会力推 jarvis-sync-heartbeat 分支(合成空树提交,
# committer 时间=同步时刻)。滞后只告警不阻断——PD 查证据此降置信,避免拿两周前
# 的镜像得出「master 未修复」的错误结论。github 直连时该 ref 不存在,静默跳过。
if git -C "$REPO_DIR" fetch -q --depth 1 origin refs/heads/jarvis-sync-heartbeat 2>/dev/null; then
  hb_ct=$(git -C "$REPO_DIR" log -1 --format=%ct FETCH_HEAD 2>/dev/null || echo 0)
  if [ "${hb_ct:-0}" -gt 0 ]; then
    lag_min=$(( ( $(date +%s) - hb_ct ) / 60 ))
    if [ "$lag_min" -gt 120 ]; then
      echo "[sync-provider] WARNING: mirror last synced ${lag_min}min ago (>120min) — source conclusions are LOW-CONFIDENCE (mirror may lag github upstream; check macmini provider-mirror-sync cron)" >&2
    else
      echo "[sync-provider] mirror freshness: synced ${lag_min}min ago"
    fi
  fi
fi
echo "[sync-provider] ready: $REPO_DIR @ $(git -C "$REPO_DIR" rev-parse --short HEAD)"
