#!/usr/bin/env bash
# control-plane-status.sh — Jarvis 控制面可观测 CLI（人工排查入口）。
#   workers            列出全部注册 worker（key/client/activityStatus/assignment aone id）
#   ready [--limit N]  列出 READY 任务及无 eligible worker 的原因
#   task <aone_id>     按 Aone ID 全链路查询（task 状态/current session/fence/
#                      最近事件/operations 回执）
#   operation <id>     单 operation + Task/Session/fence/readbackSpec point-read
#   discard-resume <task_id> <session_id> --reason TEXT --yes
#                      人工确认后丢弃一个精确的遗留恢复上下文，使任务可由新 Worker 接管
#   force-release <task_id> [session_id] --reason TEXT --yes
#                      fresh-read Task timeline 后携带完整 CAS，人工强制解除精确 ownership；
#                      绝不由 scheduler/owner-health 自动调用
#   force-redispatch <task_id> [session_id]
#                      (--auto-target | --target-worker KEY | --target-host HOST)
#                      --target-runtime INTERACTIVE|PERSISTENT --reason TEXT --yes
#                      原子隔离旧 Session，并定向到另一台在线兼容 Worker；
#                      READY 只表示已定向排队，不表示目标 Worker 已开始执行
#   unresolvable-source-cleanup TASK_ID... [--reason TEXT] [--yes]
#                      TASK_ID 是 control-plane Task ID，不是 Aone ID。默认只预览；
#                      --yes 还必须带非空 reason，且仅删除服务端回显的规范化 ID 和
#                      完整 CAS 快照。active Task/Session 或 blocking required operation
#                      会阻断删除。
#   legacy-cleanup     已退役墓碑：只打印迁移用法，绝不发 HTTP。
#
# 环境加载与 run-interactive-worker-hook.sh 同源：主仓 gitignored bootstrap/.env +
# bridge/jarvis.env；普通命令 token 缺省回退 JARVIS_HTML_REPORT_TOKEN；
# unresolvable-source-cleanup 只接受 JARVIS_CONTROL_PLANE_ADMIN_TOKEN，绝不回退普通 token。
# 控制面 base url 可由 JARVIS_CONTROL_PLANE_BASE_URL /
# JARVIS_HTML_REPORT_BASE_URL 覆盖，默认生产控制面。
# 实现体在同目录 control-plane-status.py。

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
# shellcheck disable=SC1091
source "$script_dir/runtime-config.sh"
jarvis_load_runtime_config || exit $?

python_bin="${JARVIS_CONTROL_PLANE_PYTHON:-python3}"
exec "$python_bin" "$script_dir/control-plane-status.py" "$@"
