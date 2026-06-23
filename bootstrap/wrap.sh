#!/usr/bin/env bash
# bootstrap/wrap.sh — 把研发进展同步回 Aone（唯一真源）+ 本地审计。
# Aone = 唯一真源：任何 jarvis 工作都得在 Aone 留痕；中途 sync，收尾 done。
# Usage:
#   wrap.sh sync <id> "<progress>"            — 中途报进展（评论），不改状态、不写 run_done
#   wrap.sh done <id> "<summary>" [status]    — 收尾：评论 + run_done + 可选改状态
#
# done 是每条 dev/adhoc/plugin-dev 路径在宣告 Done 前的必经收尾。
# a1 调用失败仅告警，本地审计照常落盘（Aone 真源不达也不丢账）。
# 标签/状态约定读 config/pools.json；JARVIS_ROOT 可覆盖仓库根。

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
pools_cfg="$jarvis_root/config/pools.json"

[ -f "$pools_cfg" ] || { echo "wrap.sh: config/pools.json not found at $pools_cfg" >&2; exit 1; }

cmd="${1:-}"
id="${2:-}"

case "$cmd" in
    sync)
        text="${3:-}"
        [ -n "$id" ] && [ -n "$text" ] || { echo "Usage: wrap.sh sync <id> \"<progress>\"" >&2; exit 1; }
        a1 project workitem comment create "$id" -m "$text" \
            || echo "wrap.sh: a1 comment 失败（id=$id），进展未落 Aone，请人工补" >&2
        ;;
    done)
        summary="${3:-}"
        status="${4:-}"
        [ -n "$id" ] && [ -n "$summary" ] || { echo "Usage: wrap.sh done <id> \"<summary>\" [status]" >&2; exit 1; }
        # 1) 回填 Aone 进展评论
        a1 project workitem comment create "$id" -m "$summary" \
            || echo "wrap.sh: a1 comment 失败（id=$id），收尾未落 Aone，请人工补" >&2
        # 2) 本地审计（jq 校验 pools.json 可读，与 claim 同源）
        jq -e '.claim' "$pools_cfg" >/dev/null 2>&1 || echo "wrap.sh: pools.json .claim 缺失" >&2
        bash "$script_dir/log.sh" run_done "$id" "$summary"
        # 3) 可选改状态
        if [ -n "$status" ]; then
            a1 project workitem update "$id" --status "$status" \
                || echo "wrap.sh: a1 状态更新失败（id=$id → $status），请人工核" >&2
        fi
        ;;
    *)
        echo "Usage: wrap.sh {sync|done} <id> \"<text>\" [status]" >&2
        exit 1
        ;;
esac
