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

# fenced（回执）路径用：只读不删。回执收敛前删掉 prefix 会让失败重跑生成不同正文
# → 不同 operationKey → 与遗留槽 conflict；成功后由 claim_prefix_clear 清理。
claim_prefix_peek() {
    local tid="$1"; local dir="${JARVIS_ROOT:-$jarvis_root}/.my-day"
    local f="$dir/claim-prefix-$tid.txt"
    [ -f "$f" ] || return 0
    cat "$f"
}

claim_prefix_clear() {
    local tid="$1"; local dir="${JARVIS_ROOT:-$jarvis_root}/.my-day"
    rm -f "$dir/claim-prefix-$tid.txt"
}

# --- Aone 外部写 operation receipt（docs/aone-operation-receipts.md） ---
# fenced = database-fenced 交互会话且持有该工单（has-current）；其余上下文（含
# bookend 顺带回填的非当前工单）完全走原有裸写路径。

_receipt_fenced() {
    local tid="$1"
    jarvis_interactive_context || return 1
    jarvis_interactive_worker_cli has-current "$tid" >/dev/null 2>&1
}

# 规范化（去 \r、去每行尾空白、去末尾空白）后 sha256。begin 的 material 与
# readback 匹配必须用同一实现，故收敛在这一个 python 前缀里。
_norm_digest_py='
import hashlib, sys
def norm_digest(text):
    text = text.replace("\r", "")
    text = "\n".join(line.rstrip() for line in text.split("\n")).rstrip()
    return hashlib.sha256(text.encode("utf-8")).hexdigest()
'

_normalized_digest() {
    printf '%s' "$1" | python3 -c "${_norm_digest_py}
print(norm_digest(sys.stdin.read()))"
}

# 在评论流里按规范化 digest 找我们那条评论。readback 三态：
#   rc=0 且 stdout 非空 → found（输出 comment id）
#   rc=0 且 stdout 为空 → definitely-not-found（评论流完整拉取且解析成功，确认无匹配）
#   rc≠0               → unavailable（a1 非零退出 / 非法 JSON / 命中但无 id）
# 只有 definitely-not-found 才允许 reconcile --not-found；unavailable ≠ 副作用不存在。
_find_comment_by_digest() {
    local tid="$1" digest="$2" out
    if ! out="$($A1 project workitem comment list "$tid" -f json 2>/dev/null)"; then
        return 2
    fi
    printf '%s' "$out" | python3 -c "${_norm_digest_py}
import json
want = sys.argv[1]
try:
    comments = json.loads(sys.stdin.read())
except Exception:
    sys.exit(3)   # 解析失败 ≠ 合法空列表：readback 不可用
if not isinstance(comments, list):
    sys.exit(3)
for c in comments:
    if isinstance(c, dict) and norm_digest(str(c.get(\"content\") or \"\")) == want:
        cid = c.get(\"id\")
        if not cid:
            sys.exit(3)   # 命中却无 id：无法构造 externalRef，不可当 not-found
        print(cid)
        break
" "$digest" 2>/dev/null
}

# 带回执发评论（fenced 专用）。协议：begin →（ACKED 跳过 / needsReadback 按 digest
# 匹配评论流走 reconcile / proceed 发送）→ ack。返回 0=已确认恰好一次落地；非 0=
# fail-closed，调用方必须阻断退出。
_receipted_comment() {
    local tid="$1" text="$2"
    local digest payload begin_json needs_readback proceed cid out rc
    local readback_rc op_status
    digest="$(_normalized_digest "$text")"
    payload="$(printf '%s' "$text" | python3 -c "
import json, sys
print(json.dumps({\"digest\": sys.argv[1], \"preview\": sys.stdin.read()[:120]},
                 ensure_ascii=False))" "$digest")"
    if ! begin_json="$(jarvis_interactive_worker_cli operation-begin "$tid" comment "$digest" --payload-json "$payload")"; then
        echo "wrap.sh: comment 回执 begin 失败（id=${tid}），fail-closed 不写 Aone" >&2
        return 2
    fi
    needs_readback="$(printf '%s' "$begin_json" | jq -r '.needsReadback // false' 2>/dev/null)"
    proceed="$(printf '%s' "$begin_json" | jq -r '.proceed // false' 2>/dev/null)"
    if [ "$needs_readback" = "true" ]; then
        # SENDING/UNKNOWN 存量回执：comment 不幂等，先 readback 再 reconcile 收敛。
        cid="$(_find_comment_by_digest "$tid" "$digest")"
        readback_rc=$?
        if [ "$readback_rc" -ne 0 ]; then
            # readback 不可用（a1 失败/坏 JSON）≠ 副作用不存在：绝不 --not-found 重发。
            # 槽仍是 SENDING 时冻结为 UNKNOWN；槽已是 UNKNOWN 则保持原状不再 abort。
            op_status="$(printf '%s' "$begin_json" | jq -r '.operationStatus // ""' 2>/dev/null)"
            if [ "$op_status" != "UNKNOWN" ]; then
                jarvis_interactive_worker_cli operation-abort "$tid" "comment readback unavailable" --unknown >/dev/null 2>&1 || true
            fi
            echo "wrap.sh: 评论 readback 不可用（id=${tid}），无法判定副作用是否存在；保持 UNKNOWN，稍后重跑 wrap 收敛" >&2
            return 1
        fi
        if [ -n "$cid" ]; then
            jarvis_interactive_worker_cli operation-reconcile "$tid" --found "aone:$tid:comment:$cid" >/dev/null || return 2
            return 0
        fi
        jarvis_interactive_worker_cli operation-reconcile "$tid" --not-found >/dev/null || return 2
        if ! begin_json="$(jarvis_interactive_worker_cli operation-begin "$tid" comment "$digest" --payload-json "$payload")"; then
            echo "wrap.sh: comment 回执重试 begin 失败（id=${tid}），fail-closed 不写 Aone" >&2
            return 2
        fi
        proceed="$(printf '%s' "$begin_json" | jq -r '.proceed // false' 2>/dev/null)"
    fi
    if [ "$proceed" != "true" ]; then
        # ACKED：上轮已恰好一次落地，跳过发送。
        return 0
    fi
    out="$($A1 project workitem comment create "$tid" -m "$text" -f json 2>/dev/null)"
    rc=$?
    if [ "$rc" -ne 0 ]; then
        jarvis_interactive_worker_cli operation-abort "$tid" "a1 comment create failed (rc=$rc)" >/dev/null 2>&1 || true
        echo "wrap.sh: a1 comment 失败（id=${tid}, rc=$rc），回执已终结（可重试）" >&2
        return 1
    fi
    cid="$(printf '%s' "$out" | jq -r '.id // empty' 2>/dev/null)"
    [ -n "$cid" ] || cid="$(printf '%s' "$out" | sed -n '1s/^ID:[[:space:]]*\([0-9][0-9]*\).*/\1/p')"
    if [ -z "$cid" ]; then
        # 退出码 0 但拿不到 comment id → 结果不明：冻结 UNKNOWN，等重跑 readback 收敛。
        jarvis_interactive_worker_cli operation-abort "$tid" "comment result indeterminate (no id)" --unknown >/dev/null 2>&1 || true
        echo "wrap.sh: 评论结果不明（id=${tid}），回执已冻结 UNKNOWN，重跑 wrap 收敛" >&2
        return 1
    fi
    if ! jarvis_interactive_worker_cli operation-ack "$tid" "aone:$tid:comment:$cid" >/dev/null; then
        echo "wrap.sh: 评论已落 Aone 但回执 ACK 失败（id=${tid}），重跑 wrap 收敛" >&2
        return 1
    fi
    return 0
}

# 带回执改状态（fenced 专用）。status 是幂等 set → --replay-safe，SENDING+本地意图
# 直接重放，无需 readback。返回值语义同 _receipted_comment。
_receipted_status() {
    local tid="$1" status="$2" begin_json proceed
    if ! begin_json="$(jarvis_interactive_worker_cli operation-begin "$tid" status "$status" --replay-safe)"; then
        echo "wrap.sh: status 回执 begin 失败（id=${tid}），fail-closed 不写 Aone" >&2
        return 2
    fi
    proceed="$(printf '%s' "$begin_json" | jq -r '.proceed // false' 2>/dev/null)"
    if [ "$proceed" != "true" ]; then
        return 0   # ACKED：上轮已落地
    fi
    if ! $A1 project workitem update "$tid" --status "$status"; then
        jarvis_interactive_worker_cli operation-abort "$tid" "a1 status update failed" >/dev/null 2>&1 || true
        echo "wrap.sh: a1 status 更新失败（id=${tid}），回执已终结（可重试）" >&2
        return 1
    fi
    if ! jarvis_interactive_worker_cli operation-ack "$tid" "aone:$tid:status:$status" >/dev/null; then
        echo "wrap.sh: 状态已落 Aone 但回执 ACK 失败（id=${tid}），重跑 wrap 收敛" >&2
        return 1
    fi
    return 0
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
    # fenced 会话（database-fenced 且持有本单）走 operation receipt 协议，其余现状裸写。
    local fenced=0
    if _receipt_fenced "$tid"; then fenced=1; fi
    local prefix local_summary
    if [ "$fenced" = "1" ]; then
        prefix="$(claim_prefix_peek "$tid")"   # 回执收敛前不消费，保正文可重放
    else
        prefix="$(claim_prefix_pop "$tid")"
    fi
    local_summary="$summary"
    [ -n "$prefix" ] && local_summary="${prefix}"$'\n\n'"${local_summary}"
    local_summary="${local_summary}$(code_footer)"
    local_summary="$(format_comment "$local_summary")"
    if [ "$fenced" = "1" ]; then
        if ! _receipted_comment "$tid" "$local_summary"; then
            exit 1
        fi
        claim_prefix_clear "$tid"
    else
        $A1 project workitem comment create "$tid" -m "$local_summary"
    fi
    # 3) 默认改状态；无状态收尾用于当前状态不可转/无需转时避免半失败
    if [ "$update_status" = "1" ]; then
        if [ "$fenced" = "1" ]; then
            # comment 与 status 是两个串行回执（单槽约束）；comment 先行。
            if ! _receipted_status "$tid" "$status"; then
                exit 1
            fi
        else
            $A1 project workitem update "$tid" --status "$status"
        fi
    fi
    bash "$script_dir/cache.sh" bust "wi-$tid"  # 收尾改动后详情已变，丢缓存
}

# 位置正文 / 状态槽收到 `--flag` 形态的 token = 调用方把命令行参数放错了位置
# （如 `done <id> --status X --summary-stdin`：--status 落入 else 分支被当 summary 静默贴出，
# 真 summary 从 stdin 丢弃，避免把命令行参数误当正文。
# 静默吞下会污染工单 + 丢真内容 + 静默改错状态；这里响亮报错而非当正文。
# allow（可选）= 该槽合法的 `--` 值白名单（仅状态槽的 --no-status）。
reject_flag_as_text() {  # args: value slot usage [allow]
    local val="$1" slot="$2" usage="$3" allow="${4:-}"
    [ -n "$allow" ] && [ "$val" = "$allow" ] && return 0
    case "$val" in
        --*) printf 'wrap.sh: %s 位置收到疑似命令行参数 "%s"（wrap 无此 flag）——请检查参数顺序。\n%s\n' \
                    "$slot" "$val" "$usage" >&2; exit 2 ;;
    esac
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
            reject_flag_as_text "$text" "进展正文" "Usage: wrap.sh sync <id> \"<progress>\" | --summary-file <path> | --summary-stdin"
        fi
        [ -n "$id" ] && [ -n "$text" ] || { echo "Usage: wrap.sh sync <id> \"<progress>\" | --summary-file <path> | --summary-stdin" >&2; exit 1; }
        reject_literal_newline "$text"
        touch_ledger "$id"
        fenced=0
        if _receipt_fenced "$id"; then fenced=1; fi
        if [ "$fenced" = "1" ]; then
            prefix="$(claim_prefix_peek "$id")"   # 回执收敛前不消费，保正文可重放
        else
            prefix="$(claim_prefix_pop "$id")"
        fi
        [ -n "$prefix" ] && text="${prefix}"$'\n\n'"${text}"
        text="${text}$(code_footer)"
        text="$(format_comment "$text")"
        if [ "$fenced" = "1" ]; then
            # fenced 会话的写必须有回执：begin/发送失败一律阻断（sync 也不降级为告警）。
            if ! _receipted_comment "$id" "$text"; then
                echo "wrap.sh: fenced sync 未取得回执（id=${id}），进展未落 Aone，fail-closed" >&2
                exit 1
            fi
            claim_prefix_clear "$id"
        else
            $A1 project workitem comment create "$id" -m "$text" \
                || echo "wrap.sh: a1 comment 失败（id=${id}），进展未落 Aone，请人工补" >&2
        fi
        bash "$script_dir/cache.sh" bust "wi-$id"  # 评论后详情已变，丢缓存
        ;;
    done)
        done_usage="Usage: wrap.sh done <id> \"<summary>\"|--summary-file <path>|--summary-stdin [<status>|--status <status>|--no-status]"
        # status 可用 flag 或位置参数给：--status <值>/--status=<值>/--no-status 从任意位置抽出
        # （对齐 a1 update --status 与直觉，见 cap-wrap-status-alias），其余 token 按原位置语义
        # 决定 summary 来源。flag 与位置状态互斥；误放的其它 --flag 仍由 reject_flag_as_text 拦。
        status=""; status_set=0; rest=(); args=("${@:3}"); i=0
        while [ "$i" -lt "${#args[@]}" ]; do
            case "${args[$i]}" in
                --status)
                    i=$((i+1))
                    [ "$i" -lt "${#args[@]}" ] || { echo "wrap.sh: --status 需要一个状态值" >&2; echo "$done_usage" >&2; exit 2; }
                    status="${args[$i]}"; status_set=1 ;;
                --status=*) status="${args[$i]#--status=}"; status_set=1 ;;
                --no-status) status="--no-status"; status_set=1 ;;
                *) rest+=("${args[$i]}") ;;
            esac
            i=$((i+1))
        done
        # 从 rest[] 定 summary 来源；pos_status = 正文之后的位置状态参数（若有）
        if [ "${rest[0]:-}" = "--summary-file" ]; then
            summary="$(summary_from_file "${rest[1]:-}")"; pos_status="${rest[2]:-}"; consumed=2
        elif [ "${rest[0]:-}" = "--summary-stdin" ]; then
            summary="$(cat)"; pos_status="${rest[1]:-}"; consumed=1
        else
            summary="${rest[0]:-}"
            reject_flag_as_text "$summary" "summary" "$done_usage"
            pos_status="${rest[1]:-}"; consumed=1
        fi
        if [ "$status_set" -eq 1 ]; then
            # flag 已给状态：正文之后不应再有多余位置参数（防重复指定/参数错位）
            [ "${#rest[@]}" -le "$consumed" ] || { echo "wrap.sh: 状态既用 --status/--no-status 指定又出现多余位置参数：'$pos_status'" >&2; echo "$done_usage" >&2; exit 2; }
        else
            status="$pos_status"
        fi
        # 状态槽兜底：仅 --no-status 是合法 `--` 值；其余 --flag = 参数错位，响亮报错。
        reject_flag_as_text "$status" "status" "$done_usage" "--no-status"
        [ -n "$id" ] && [ -n "$summary" ] && [ -n "$status" ] || { echo "$done_usage" >&2; exit 1; }
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
            reject_flag_as_text "$summary" "summary" "Usage: wrap.sh done-no-status <id> \"<summary>\" | --summary-file <path> | --summary-stdin"
        fi
        [ -n "$id" ] && [ -n "$summary" ] || { echo "Usage: wrap.sh done-no-status <id> \"<summary>\" | --summary-file <path> | --summary-stdin" >&2; exit 1; }
        write_done "$id" "$summary" "" 0
        ;;
    *)
        echo "Usage: wrap.sh {sync|done|done-no-status} <id> \"<text>\"|--summary-file <path>|--summary-stdin [status|--no-status]" >&2
        exit 1
        ;;
esac
