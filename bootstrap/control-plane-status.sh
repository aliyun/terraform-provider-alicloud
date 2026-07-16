#!/usr/bin/env bash
# control-plane-status.sh — Jarvis 控制面可观测 CLI（人工排查入口）。
#   workers            列出全部注册 worker（key/client/activityStatus/assignment aone id）
#   task <aone_id>     按 Aone ID 全链路查询（task 状态/current session/fence/
#                      最近事件/operations 回执）
#
# 环境加载与 run-interactive-worker-hook.sh 同源：主仓 gitignored bootstrap/.env +
# bridge/jarvis.env；token 缺省回退 JARVIS_HTML_REPORT_TOKEN（base url 回退
# JARVIS_HTML_REPORT_BASE_URL）。实现体在同目录 control-plane-status.py。

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

# Worktrees share the main repository's gitignored bootstrap/.env and
# bridge/jarvis.env.  Explicit path overrides keep tests and unusual installs
# deterministic.
main_root="$repo_root"
git_common="$(git -C "$repo_root" rev-parse --path-format=absolute --git-common-dir 2>/dev/null || true)"
case "$git_common" in
  */.git) main_root="${git_common%/.git}" ;;
esac
bootstrap_env="${JARVIS_INTERACTIVE_BOOTSTRAP_ENV:-$main_root/bootstrap/.env}"
bridge_env="${JARVIS_INTERACTIVE_BRIDGE_ENV:-$main_root/bridge/jarvis.env}"

set -a
set +u
# shellcheck disable=SC1090
[ -f "$bootstrap_env" ] && . "$bootstrap_env"
# shellcheck disable=SC1090
[ -f "$bridge_env" ] && . "$bridge_env"
set -u
set +a

# The control plane intentionally reuses the established report credential.
# Keep the fallback in-process only; nothing is persisted.
if [ -z "${JARVIS_CONTROL_PLANE_TOKEN:-}" ] && [ -n "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
  export JARVIS_CONTROL_PLANE_TOKEN="$JARVIS_HTML_REPORT_TOKEN"
fi
if [ -z "${JARVIS_CONTROL_PLANE_BASE_URL:-}" ] && [ -n "${JARVIS_HTML_REPORT_BASE_URL:-}" ]; then
  export JARVIS_CONTROL_PLANE_BASE_URL="$JARVIS_HTML_REPORT_BASE_URL"
fi

python_bin="${JARVIS_CONTROL_PLANE_PYTHON:-python3}"
exec "$python_bin" "$script_dir/control-plane-status.py" "$@"
