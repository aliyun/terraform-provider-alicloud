#!/usr/bin/env bash
# bootstrap/claim.sh – claim/release/finish primitive for concurrent multi-machine agents.
# Usage:
#   claim.sh claim   <workitem_id> <project_id>
#   claim.sh release <workitem_id> <project_id>
#   claim.sh finish  <workitem_id> <project_id>
#
# claim:   tags the workitem jarvis-claimed, freezes a `jarvis-claim <host> <utc>` prefix
#          into .my-day/claim-prefix-<id>.txt (consumed by wrap.sh sync/done into the
#          first business comment to avoid an independent claim-only comment), then reads
#          back the tag list to verify we won any concurrent race. Exits 0 on success,
#          1 if the workitem is not found in the readback (lost race; prefix cleaned).
# release: tags the workitem jarvis-idle —— 本轮 jarvis 处理完毕、释放锁、等待人或下一个 jarvis
#          接手；不动 Aone status。cleans up any leftover prefix. Exits 0.
# finish:  tags the workitem jarvis-done + 改 Aone status 为 .claim.done_status（默认"已发布
#          待需求排期"）—— jarvis 判断工单真完成。cleans up any leftover prefix. Exits 0.
#
# Reads tag names from config/pools.json (.claim.tag / .claim.idle_tag / .claim.done_tag /
# .claim.done_status). Respects JARVIS_ROOT env override for the repo root.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
pools_cfg="$jarvis_root/config/pools.json"

if [ ! -f "$pools_cfg" ]; then
    echo "claim.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

# Read tag names from config
CLAIM_TAG=$(jq -r '.claim.tag' "$pools_cfg")
IDLE_TAG=$(jq -r '.claim.idle_tag' "$pools_cfg")
DONE_TAG=$(jq -r '.claim.done_tag' "$pools_cfg")
DONE_STATUS=$(jq -r '.claim.done_status' "$pools_cfg")

cmd="${1:-}"
workitem_id="${2:-}"
project_id="${3:-}"

if [ -z "$cmd" ] || [ -z "$workitem_id" ] || [ -z "$project_id" ]; then
    echo "Usage: claim.sh <claim|release|finish> <workitem_id> <project_id>" >&2
    exit 1
fi

myday_dir="$jarvis_root/.my-day"
today_date="$(date -u +%F)"
ledger_file="$myday_dir/claims-${today_date}.json"

_claim_prefix_path() { echo "$myday_dir/claim-prefix-$1.txt"; }

# Atomically append/update entry in the claims ledger.
# Usage: _ledger_upsert <id> <done_value>   (done_value: true|false)
_ledger_upsert() {
    local id="$1"
    local done_val="$2"
    mkdir -p "$myday_dir"
    local tmp
    tmp="$(mktemp "$myday_dir/.claims-tmp.XXXXXX")"
    # Seed tmp with existing ledger or empty array; never create a partial ledger_file first
    if [ -f "$ledger_file" ]; then
        jq --arg id "$id" --argjson done "$done_val" \
            'if any(.[]; .id == $id) then
                 map(if .id == $id then .done = $done else . end)
             else
                 . + [{"id": $id, "done": $done}]
             end' \
            "$ledger_file" > "$tmp" && mv "$tmp" "$ledger_file"
    else
        printf '[{"id":"%s","done":%s}]' "$id" "$done_val" > "$tmp" && mv "$tmp" "$ledger_file"
    fi
}

case "$cmd" in
    claim)
        # 1. Tag the workitem as claimed
        a1 project workitem update "$workitem_id" --tag "$CLAIM_TAG"

        # 2. 冻结 hostname + utc 到 prefix 文件；wrap.sh sync/done 第一次发业务评论时自动 prefix
        # 到顶部并消费，避免 claim 时刻额外刷一条独立留痕评论。
        # 优先 ComputerName(用户自定义显示名)，回落 hostname -s 短名，最后 hostname；避免 mac 默认带 .local 后缀
        utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        host=$(scutil --get ComputerName 2>/dev/null || hostname -s 2>/dev/null || hostname)
        mkdir -p "$myday_dir"
        printf 'jarvis-claim %s %s' "$host" "$utcnow" > "$(_claim_prefix_path "$workitem_id")"

        # 3. Readback (with retries for tag-index lag): own id under claimed tag = won
        for _ in 1 2 3; do
            readback=$(a1 project workitem list --project "$project_id" --tag "$CLAIM_TAG" -f json 2>/dev/null || echo "[]")
            if echo "$readback" | jq -e --arg id "$workitem_id" 'any(.[]; ((.identifier // .id)|tostring)==$id)' > /dev/null 2>&1; then
                echo "claim.sh: claimed workitem $workitem_id in project $project_id"
                _ledger_upsert "$workitem_id" false
                exit 0
            fi
            sleep 2
        done
        # 抢锁失败：清理本机 prefix 痕迹（不留给别人）
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: lost race for workitem $workitem_id (not found in readback)" >&2
        exit 1
        ;;

    release)
        # Tag as idle: 本轮 jarvis 处理完，释放锁；后续等待人或下一个 jarvis 接手
        # (no untag exists; idle_tag supersedes claimed in tag list)
        a1 project workitem update "$workitem_id" --tag "$IDLE_TAG"
        # 兜底：如果整个 jarvis 周期没发过业务评论，prefix 文件仍在，此处清理
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: released workitem $workitem_id in project $project_id"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    finish)
        # Tag as done + 改 status：jarvis 判断工单真完成
        a1 project workitem update "$workitem_id" --tag "$DONE_TAG" --status "$DONE_STATUS"
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: finished workitem $workitem_id in project $project_id (status=$DONE_STATUS)"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim, release, or finish." >&2
        exit 1
        ;;
esac
