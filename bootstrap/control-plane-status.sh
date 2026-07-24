#!/usr/bin/env bash
# control-plane-status.sh — Jarvis 控制面可观测 CLI（人工排查入口）。
#   workers            列出全部注册 worker（key/client/activityStatus/assignment aone id）
#   ready [--limit N]  列出 READY 任务及无 eligible worker 的原因
#   task <aone_id>     按 Aone ID 全链路查询（task 状态/current session/fence/
#                      最近事件/operations 回执）
#   operation <id>     单 operation + Task/Session/fence/readbackSpec point-read
#   discard-resume <task_id> <session_id> --reason TEXT --yes
#                      人工确认后丢弃一个精确的遗留恢复上下文，使任务可由新 Worker 接管
#   legacy-cleanup [--yes]
#                      预览（默认）/删除 task_type 为已废弃 kind（如 field_repair）的
#                      残留僵尸 Task 及其 session/event/operation 行；带 --yes 才删，
#                      服务端按精确快照+active 守护，改动后 409。
#
# 环境加载与 run-interactive-worker-hook.sh 同源：主仓 gitignored bootstrap/.env +
# bridge/jarvis.env；token 缺省回退 JARVIS_HTML_REPORT_TOKEN；控制面 base url 可由
# JARVIS_CONTROL_PLANE_BASE_URL / JARVIS_HTML_REPORT_BASE_URL 覆盖，默认预发。
# 实现体在同目录 control-plane-status.py。

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
# shellcheck disable=SC1091
source "$script_dir/runtime-config.sh"
jarvis_load_runtime_config || exit $?

python_bin="${JARVIS_CONTROL_PLANE_PYTHON:-python3}"
exec "$python_bin" "$script_dir/control-plane-status.py" "$@"
