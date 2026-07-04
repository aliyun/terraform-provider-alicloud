#!/usr/bin/env bash
# bootstrap/wrap.sh — 把研发进展同步回 Aone（唯一真源）+ 本地审计。
# Aone = 唯一真源：任何 jarvis 工作都得在 Aone 留痕；中途 sync，收尾 done。
# Usage:
#   wrap.sh sync <id> "<progress>"            — 中途报进展（评论），不改状态、不写 run_done
#   wrap.sh sync <id> --summary-file <path>   — 从文件读取多行进展
#   wrap.sh sync <id> --summary-stdin         — 从 stdin 读取多行进展
#   wrap.sh done <id> "<summary>" <status>     — 收尾：评论 + run_done + 改状态（status 必填）
#   wrap.sh done <id> --summary-file <path> <status|--no-status>
#   wrap.sh done <id> --summary-stdin <status|--no-status>
#   wrap.sh done-no-status <id> "<summary>"    — 收尾：评论 + run_done，不改状态
#   wrap.sh done-no-status <id> --summary-file <path>
#   wrap.sh done-no-status <id> --summary-stdin
#   wrap.sh done <id> "<summary>" --no-status  — 同 done-no-status
#
# done 是每条 dev/adhoc/plugin-dev 路径在宣告 Done 前的必经收尾。
# done: 本地审计先落盘（审计不丢），再 a1 调用——失败则 exit 1（Aone 真源强制）。
# sync: a1 调用失败仅告警，本地不写 run_done。
# 评论自动追加代码落点（repo @ 分支 commit）：在开发库目录调用，或经以下环境变量指定：
#   CODE_DIR_KEY=<workspace-key>  — 走 bootstrap/workspace.sh dir <key> 解析（对齐 CLAUDE.md 纪律 4，本地路径不手拼）
#   CODE_DIR=<abs-path>           — 直接绝对路径（向后兼容 / 临时路径）
# 未设 → 用 PWD。
# 标签/状态约定读 config/pools.json；JARVIS_ROOT 可覆盖仓库根。

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
pools_cfg="$jarvis_root/config/pools.json"

# a1 via bin/a1id → act as the jarvis identity regardless of ambient login (CLAUDE.md #6).
# Overridable via JARVIS_A1 (tests point it at a stubbed `a1` on PATH).
A1="${JARVIS_A1:-$jarvis_root/bin/a1id --}"

[ -f "$pools_cfg" ] || { echo "wrap.sh: config/pools.json not found at $pools_cfg" >&2; exit 1; }

# 代码落点页脚：从开发库当前 git 目录采 repo/分支/commit 自动追加，回填进展即定位代码。
# 取源优先级 CODE_DIR_KEY > CODE_DIR > PWD；非 git 目录静默跳过不阻断回填。
code_footer() {
    local dir=""
    if [ -n "${CODE_DIR_KEY:-}" ]; then
        # 走 workspace.sh 解析(对齐 CLAUDE.md 纪律 4:本地路径经 workspaces 登记,不手拼)
        dir="$(bash "$script_dir/workspace.sh" dir "$CODE_DIR_KEY" 2>/dev/null)" || return 0
    fi
    [ -n "$dir" ] || dir="${CODE_DIR:-$PWD}"
    git -C "$dir" rev-parse --is-inside-work-tree >/dev/null 2>&1 || return 0
    # jarvis 仓内调用属于查证/编排，不是开发库提交，footer 无意义且会污染评论；跳过
    local toplevel
    toplevel=$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)
    [ "$toplevel" = "$jarvis_root" ] && return 0
    local repo branch sha dirty
    repo=$(basename -s .git "$(git -C "$dir" remote get-url origin 2>/dev/null)" 2>/dev/null)
    [ -n "$repo" ] || repo=$(basename "$toplevel")
    branch=$(git -C "$dir" rev-parse --abbrev-ref HEAD 2>/dev/null)
    sha=$(git -C "$dir" rev-parse --short HEAD 2>/dev/null)
    [ -n "$(git -C "$dir" status --porcelain 2>/dev/null)" ] && dirty="+dirty"
    printf '\n\n代码：%s @ %s (%s%s)' "$repo" "$branch" "$sha" "${dirty:-}"
}

format_comment() {
    bash "$script_dir/aone-comment-format.sh" "$1"
}

summary_from_file() {
    local path="${1:-}"
    [ -n "$path" ] || { echo "wrap.sh: --summary-file requires a path" >&2; exit 1; }
    [ -f "$path" ] || { echo "wrap.sh: summary file not found: $path" >&2; exit 1; }
    cat -- "$path"
}

reject_literal_newline() {
    local text="$1"
    local literal_newline='\n'
    if [ "${JARVIS_ALLOW_LITERAL_NEWLINE:-0}" != "1" ] && [[ "$text" == *"$literal_newline"* ]]; then
        printf 'wrap.sh: summary contains literal \\n; use heredoc/file/stdin or --summary-file/--summary-stdin for real multiline text. Set JARVIS_ALLOW_LITERAL_NEWLINE=1 only when literal \\n text is intentional.\n' >&2
        exit 1
    fi
}

# 触碰台账：任何 sync/done 都记 id，wrap-check 据此抓"碰过但没收尾(无 run_done)"的盲区
# ——盲区指根本没 claim、claim 台账无记录的工单(see 83498126)。
touch_ledger() {
    local tid="$1"; local dir="${JARVIS_ROOT:-$jarvis_root}/.my-day"
    mkdir -p "$dir"; local f="$dir/touched-$(date -u +%F).json"
    local tmp; tmp="$(mktemp "$dir/.touched-tmp.XXXXXX")"
    if [ -f "$f" ]; then jq --arg id "$tid" 'if any(.[]; .==$id) then . else .+[$id] end' "$f" >"$tmp" && mv "$tmp" "$f"
    else echo "[\"$tid\"]" >"$f"; rm -f "$tmp"; fi
}

# 消费 claim.sh 冻结的 jarvis-claim 痕迹（.my-day/claim-prefix-<id>.txt）。
# 找到就输出内容并删除文件，让第一条业务评论吸收这条痕迹；找不到输出空。
claim_prefix_pop() {
    local tid="$1"; local dir="${JARVIS_ROOT:-$jarvis_root}/.my-day"
    local f="$dir/claim-prefix-$tid.txt"
    [ -f "$f" ] || return 0
    cat "$f"
    rm -f "$f"
}

write_done() {
    local tid="$1"
    local summary="$2"
    local status="${3:-}"
    local update_status="${4:-1}"

    reject_literal_newline "$summary"
    touch_ledger "$tid"
    # 1) 本地审计先写，审计绝不丢（即使 Aone 调用失败）
    jq -e '.claim' "$pools_cfg" >/dev/null 2>&1 || echo "wrap.sh: pools.json .claim 缺失" >&2
    bash "$script_dir/log.sh" run_done "$tid" "$summary"
    # 2) 回填 Aone 进展评论（带 claim 痕迹前缀 + 代码落点页脚 → 格式化）；失败则 exit 1
    local prefix local_summary
    prefix="$(claim_prefix_pop "$tid")"
    local_summary="$summary"
    [ -n "$prefix" ] && local_summary="${prefix}"$'\n\n'"${local_summary}"
    local_summary="${local_summary}$(code_footer)"
    local_summary="$(format_comment "$local_summary")"
    $A1 project workitem comment create "$tid" -m "$local_summary"
    # 3) 默认改状态；无状态收尾用于当前状态不可转/无需转时避免半失败
    if [ "$update_status" = "1" ]; then
        $A1 project workitem update "$tid" --status "$status"
    fi
    bash "$script_dir/cache.sh" bust "wi-$tid"  # 收尾改动后详情已变，丢缓存
}

cmd="${1:-}"
id="${2:-}"

case "$cmd" in
    sync)
        if [ "${3:-}" = "--summary-file" ]; then
            text="$(summary_from_file "${4:-}")"
        elif [ "${3:-}" = "--summary-stdin" ]; then
            text="$(cat)"
        else
            text="${3:-}"
        fi
        [ -n "$id" ] && [ -n "$text" ] || { echo "Usage: wrap.sh sync <id> \"<progress>\" | --summary-file <path> | --summary-stdin" >&2; exit 1; }
        reject_literal_newline "$text"
        touch_ledger "$id"
        prefix="$(claim_prefix_pop "$id")"
        [ -n "$prefix" ] && text="${prefix}"$'\n\n'"${text}"
        text="${text}$(code_footer)"
        text="$(format_comment "$text")"
        $A1 project workitem comment create "$id" -m "$text" \
            || echo "wrap.sh: a1 comment 失败（id=${id}），进展未落 Aone，请人工补" >&2
        bash "$script_dir/cache.sh" bust "wi-$id"  # 评论后详情已变，丢缓存
        ;;
    done)
        if [ "${3:-}" = "--summary-file" ]; then
            summary="$(summary_from_file "${4:-}")"
            status="${5:-}"
        elif [ "${3:-}" = "--summary-stdin" ]; then
            summary="$(cat)"
            status="${4:-}"
        else
            summary="${3:-}"
            status="${4:-}"
        fi
        [ -n "$id" ] && [ -n "$summary" ] && [ -n "$status" ] || { echo "Usage: wrap.sh done <id> \"<summary>\"|--summary-file <path>|--summary-stdin <status|--no-status>" >&2; exit 1; }
        if [ "$status" = "--no-status" ]; then
            write_done "$id" "$summary" "" 0
        else
            write_done "$id" "$summary" "$status" 1
        fi
        ;;
    done-no-status)
        if [ "${3:-}" = "--summary-file" ]; then
            summary="$(summary_from_file "${4:-}")"
        elif [ "${3:-}" = "--summary-stdin" ]; then
            summary="$(cat)"
        else
            summary="${3:-}"
        fi
        [ -n "$id" ] && [ -n "$summary" ] || { echo "Usage: wrap.sh done-no-status <id> \"<summary>\" | --summary-file <path> | --summary-stdin" >&2; exit 1; }
        write_done "$id" "$summary" "" 0
        ;;
    *)
        echo "Usage: wrap.sh {sync|done|done-no-status} <id> \"<text>\"|--summary-file <path>|--summary-stdin [status|--no-status]" >&2
        exit 1
        ;;
esac
