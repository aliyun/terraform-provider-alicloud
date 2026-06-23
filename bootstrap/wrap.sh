#!/usr/bin/env bash
# bootstrap/wrap.sh — 把研发进展同步回 Aone（唯一真源）+ 本地审计。
# Aone = 唯一真源：任何 jarvis 工作都得在 Aone 留痕；中途 sync，收尾 done。
# Usage:
#   wrap.sh sync <id> "<progress>"            — 中途报进展（评论），不改状态、不写 run_done
#   wrap.sh done <id> "<summary>" [status]    — 收尾：评论 + run_done + 可选改状态
#
# done 是每条 dev/adhoc/plugin-dev 路径在宣告 Done 前的必经收尾。
# a1 调用失败仅告警，本地审计照常落盘（Aone 真源不达也不丢账）。
# 评论自动追加代码落点（repo @ 分支 commit）：在开发库目录调用，或 CODE_DIR 指定。
# 标签/状态约定读 config/pools.json；JARVIS_ROOT 可覆盖仓库根。

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
pools_cfg="$jarvis_root/config/pools.json"

[ -f "$pools_cfg" ] || { echo "wrap.sh: config/pools.json not found at $pools_cfg" >&2; exit 1; }

# 代码落点页脚：从开发库当前 git 目录采 repo/分支/commit 自动追加，回填进展即定位代码。
# 取 CODE_DIR（默认 PWD）；非 git 目录则静默跳过，不阻断回填。
code_footer() {
    local dir="${CODE_DIR:-$PWD}"
    git -C "$dir" rev-parse --is-inside-work-tree >/dev/null 2>&1 || return 0
    local repo branch sha dirty
    repo=$(basename -s .git "$(git -C "$dir" remote get-url origin 2>/dev/null)" 2>/dev/null)
    [ -n "$repo" ] || repo=$(basename "$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)")
    branch=$(git -C "$dir" rev-parse --abbrev-ref HEAD 2>/dev/null)
    sha=$(git -C "$dir" rev-parse --short HEAD 2>/dev/null)
    [ -n "$(git -C "$dir" status --porcelain 2>/dev/null)" ] && dirty="+dirty"
    printf '\n\n代码：%s @ %s (%s%s)' "$repo" "$branch" "$sha" "${dirty:-}"
}

cmd="${1:-}"
id="${2:-}"

case "$cmd" in
    sync)
        text="${3:-}"
        [ -n "$id" ] && [ -n "$text" ] || { echo "Usage: wrap.sh sync <id> \"<progress>\"" >&2; exit 1; }
        text="${text}$(code_footer)"
        a1 project workitem comment create "$id" -m "$text" \
            || echo "wrap.sh: a1 comment 失败（id=${id}），进展未落 Aone，请人工补" >&2
        bash "$script_dir/cache.sh" bust "wi-$id"  # 评论后详情已变，丢缓存
        ;;
    done)
        summary="${3:-}"
        status="${4:-}"
        [ -n "$id" ] && [ -n "$summary" ] || { echo "Usage: wrap.sh done <id> \"<summary>\" [status]" >&2; exit 1; }
        # 1) 回填 Aone 进展评论（带代码落点页脚）
        summary="${summary}$(code_footer)"
        a1 project workitem comment create "$id" -m "$summary" \
            || echo "wrap.sh: a1 comment 失败（id=${id}），收尾未落 Aone，请人工补" >&2
        # 2) 本地审计（jq 校验 pools.json 可读，与 claim 同源）
        jq -e '.claim' "$pools_cfg" >/dev/null 2>&1 || echo "wrap.sh: pools.json .claim 缺失" >&2
        bash "$script_dir/log.sh" run_done "$id" "$summary"
        # 3) 可选改状态
        if [ -n "$status" ]; then
            a1 project workitem update "$id" --status "$status" \
                || echo "wrap.sh: a1 状态更新失败（id=${id} → ${status}），请人工核" >&2
        fi
        bash "$script_dir/cache.sh" bust "wi-$id"  # 收尾改动后详情已变，丢缓存
        ;;
    *)
        echo "Usage: wrap.sh {sync|done} <id> \"<text>\" [status]" >&2
        exit 1
        ;;
esac
