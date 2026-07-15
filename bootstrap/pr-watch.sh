#!/usr/bin/env bash
# bootstrap/pr-watch.sh — PR 观察登记表 (方案A)。
#
# skill/persona 提交 PR 后按自治边界 release 成 jarvis-idle，工单永久停在 idle 永不推到
# 「已完成」——RevisitScheduler 只捞标题/描述含特定词的 idle 单，terraform 发布单不含。
# 本登记表 + 后台 PrWatchScheduler(bridge/jarvis_dingtalk_bot.py) 补这个缺口：PR 合并后
# 自动 claim.sh finish 收尾；PR 未合并即关闭则评论 + escalate 交人工。
#
# Usage:
#   pr-watch.sh add    <ticket> <pr_url> <project>   登记一条 PR 观察
#   pr-watch.sh remove <ticket>                      摘除一条（终态/收尾/放弃观察后）
#   pr-watch.sh list                                 列出 <ticket>\t<pr_url>\t<project>（TSV）
#
# 登记表落 .my-day/bridge/pr-watch.json，**对象键**结构（非数组）：
#   {"<ticket>": {"pr_url": "...", "project": "...", "submitted_at": "..."}}
# add/remove 是 read-modify-write：套 mkdir-lock（macOS 无 flock，克隆 claim.sh 的锁范式）；
# 写用同目录 mktemp 兄弟文件 + mv 原子替换；文件不存在时先 jq -n '{}' 播种（对象，不是 []）。
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"; R="$(jarvis_root)"
F="$R/.my-day/bridge/pr-watch.json"; mkdir -p "$(dirname "$F")"; umask 077

LOCK="$(dirname "$F")/.pr-watch.lock"

# mkdir-lock（克隆 claim.sh _ledger_upsert 的并发范式，macOS 无 flock）：
# 竞得锁 echo 1；一把 >10s 的陈旧锁(持有者可能死在写一半)被偷；~5s 争用后放行不锁，
# 宁可无锁写(mv 仍原子)也不挂死一个 add/remove。stat 读不到 mtime 一律不偷(锁刚消失/瞬态)。
_acquire_lock() {
    local acquired="" waited=0 lock_mtime
    while [ "$waited" -lt 50 ]; do
        if mkdir "$LOCK" 2>/dev/null; then acquired=1; break; fi
        lock_mtime="$(stat -f %m "$LOCK" 2>/dev/null || stat -c %Y "$LOCK" 2>/dev/null || echo "")"
        if [ -n "$lock_mtime" ] && [ "$(( $(date +%s) - lock_mtime ))" -gt 10 ]; then
            rm -rf "$LOCK" 2>/dev/null
        fi
        sleep 0.1; waited=$((waited + 1))
    done
    [ -n "$acquired" ] && echo 1 || echo ""
}

_seed() {  # 文件不存在 → 播种为空对象 {}（绝不是 []）
    [ -f "$F" ] || jq -n '{}' > "$F"
}

cmd="${1:-}"
case "$cmd" in
    add)
        ticket="${2:-}"; pr_url="${3:-}"; project="${4:-}"
        [ -n "$ticket" ] && [ -n "$pr_url" ] && [ -n "$project" ] \
            || { echo "Usage: pr-watch.sh add <ticket> <pr_url> <project>" >&2; exit 1; }
        # 校验 pr_url 形如完整 GitHub PR URL——防止 bare number 被 gh 解析到错仓
        # (jarvis worktree 的默认仓)，那会导致 PrWatchScheduler 观察到无关 PR。
        if ! printf '%s' "$pr_url" | grep -Eq '^https://github\.com/.+/pull/[0-9]+$'; then
            echo "pr-watch.sh: invalid pr_url '$pr_url' — expect a full GitHub PR URL like https://github.com/<owner>/<repo>/pull/<n>" >&2
            exit 1
        fi
        ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
        acquired="$(_acquire_lock)"
        _seed
        tmp="$(mktemp "$(dirname "$F")/.pr-watch-tmp.XXXXXX")"
        jq --arg id "$ticket" --arg url "$pr_url" --arg pj "$project" --arg ts "$ts" \
            '. + {($id): {pr_url:$url, project:$pj, submitted_at:$ts}}' "$F" > "$tmp" \
            && mv "$tmp" "$F" || rm -f "$tmp"
        [ -n "$acquired" ] && rmdir "$LOCK" 2>/dev/null
        echo "pr-watch.sh: watching #$ticket → $pr_url (project $project)"
        ;;
    remove)
        ticket="${2:-}"
        [ -n "$ticket" ] || { echo "Usage: pr-watch.sh remove <ticket>" >&2; exit 1; }
        acquired="$(_acquire_lock)"
        _seed
        tmp="$(mktemp "$(dirname "$F")/.pr-watch-tmp.XXXXXX")"
        jq --arg id "$ticket" 'del(.[$id])' "$F" > "$tmp" && mv "$tmp" "$F" || rm -f "$tmp"
        [ -n "$acquired" ] && rmdir "$LOCK" 2>/dev/null
        echo "pr-watch.sh: unwatched #$ticket"
        ;;
    list)
        # 只读免锁：单次 jq 读，缺文件静默出空。
        [ -f "$F" ] || exit 0
        jq -r 'to_entries[] | [.key, .value.pr_url, .value.project] | @tsv' "$F"
        ;;
    *)
        echo "Usage: pr-watch.sh {add <ticket> <pr_url> <project> | remove <ticket> | list}" >&2
        exit 2
        ;;
esac
