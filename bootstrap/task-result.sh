#!/usr/bin/env bash
# bootstrap/task-result.sh — 把本轮 Task 的结构化收尾结果写入文件通道。
#
# Usage:
#   task-result.sh <工单号> --file <path>    — 从文件读取结果 JSON
#   task-result.sh <工单号> --stdin          — 从 stdin 读取（heredoc 友好）
#   task-result.sh <工单号> < result.json    — 同 --stdin
#
# 为什么存在：stdout 上的 [[AONE_RESULT:{...}]] 和 run 的其它输出挤在同一条通道里，
# 长 reply_body + evidence/mr_cr_links 可能把 JSON 顶到输出上限、`]]` 永远等不到，
# 于是干完活的 run 被判成「什么都没产出」。写文件不受这个限制。
#
# 关键收益是**当场校验**：本脚本用 executor 同一套契约校验，不合规立刻非零退出并打印
# 具体原因，agent 在本轮内就能改；而不是跑完几小时后被拒、只能整轮重试。
#
# executor 优先读该文件，stdout 哨兵保留为兜底——两条通道都不是单点。
# 建议两条都给：先跑本脚本，再在最后一行输出同一份 [[AONE_RESULT:{...}]]。

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"

die() { printf '%s\n' "$*" >&2; exit 1; }

item_id="${1:-}"
[ -n "$item_id" ] || die "usage: task-result.sh <工单号> [--file <path>|--stdin]"
shift

source_file=""
case "${1:-}" in
  --file)  source_file="${2:-}"; [ -n "$source_file" ] || die "--file 需要路径" ;;
  --stdin|"") source_file="" ;;
  *) die "未知参数: $1（只接受 --file <path> 或 --stdin）" ;;
esac

if [ -n "$source_file" ]; then
  [ -f "$source_file" ] || die "结果文件不存在: $source_file"
  exec_input() { cat -- "$source_file"; }
else
  exec_input() { cat; }
fi

# python3 解析器与 executor 共用 bridge.persistent_tasks 的校验逻辑，避免两套契约漂移。
exec_input | (cd "$jarvis_root" && python3 -m bridge.task_result_file write "$item_id")
