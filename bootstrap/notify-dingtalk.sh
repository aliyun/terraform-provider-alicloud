#!/usr/bin/env bash
# bootstrap/notify-dingtalk.sh — 转单/关单动作的钉钉私信通知
#
# 用法:
#   notify-dingtalk.sh <staffId> <title> <body>
#   notify-dingtalk.sh --result-json [--out-track-id <id>] <staffId> <title> <body>
#   notify-dingtalk.sh <staffId> <title> --body-file <path>
#   notify-dingtalk.sh <staffId> <title> --body-stdin
#   notify-dingtalk.sh --dry-run <staffId> <title> <body>   # 打印会发的消息,不实际发
#
# 静默降级(所有下列条件都退 0 且不阻断 caller):
#   1. JARVIS_NOTIFY_DINGTALK=0       全局关闭
#   2. 缺 DINGTALK_APP_KEY/SECRET/TEMPLATE_ID  凭据未配
#   3. staffId ∈ config/dingtalk-optout.txt   个人 opt-out(gitignored)
#   4. 依赖 skill 缺失(streaming.py 不在)
#
# 失败(网络/API 错):
#   stderr 打红,落 escalation/notify-fail-<ts>-<staffId>.md,退 0(不阻断 bookend)

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"

# 凭据自装载:DINGTALK_* 只配在 bridge/jarvis.env,装载点原本只有 bridge/run.sh
# (bridge/headless 链路继承);交互会话不经过 run.sh 会落在盲区 → 缺任一变量时
# 按 run.sh _source_env 同款姿势补 source(存在才 source,行为与 bridge 链路一致)。
if [ -z "${DINGTALK_APP_KEY:-}" ] || [ -z "${DINGTALK_APP_SECRET:-}" ] || [ -z "${DINGTALK_TEMPLATE_ID:-}" ]; then
  set -a; set +u
  # shellcheck disable=SC1090
  [ -f "$jarvis_root/bootstrap/.env" ] && . "$jarvis_root/bootstrap/.env"
  # shellcheck disable=SC1090
  [ -f "$jarvis_root/bridge/jarvis.env" ] && . "$jarvis_root/bridge/jarvis.env"
  set -u; set +a
fi

dry_run=0
result_json=0
out_track_id=""
while [ $# -gt 0 ]; do
  case "${1:-}" in
    --dry-run)
      dry_run=1
      shift
      ;;
    --result-json)
      result_json=1
      shift
      ;;
    --out-track-id)
      [ -n "${2:-}" ] || { echo "notify-dingtalk.sh: --out-track-id requires a value" >&2; exit 2; }
      out_track_id="$2"
      shift 2
      ;;
    *)
      break
      ;;
  esac
done

if [ $# -lt 3 ]; then
  echo "notify-dingtalk.sh: usage: notify-dingtalk.sh [--dry-run] [--result-json] [--out-track-id ID] <staffId> <title> <body|--body-file PATH|--body-stdin>" >&2
  exit 2
fi

staff_id="$1"
title="$2"
shift 2

# body 来源
body=""
case "${1:-}" in
  --body-file)
    [ -f "${2:-}" ] || { echo "notify-dingtalk.sh: --body-file $2 not found" >&2; exit 2; }
    body="$(cat "$2")"
    ;;
  --body-stdin)
    body="$(cat)"
    ;;
  "")
    echo "notify-dingtalk.sh: body missing" >&2
    exit 2
    ;;
  *)
    body="$1"
    ;;
esac

log_skip() { echo "[NOTIFY-SKIP: $1] staffId=$staff_id title=\"$title\"" >&2; }
emit_result() {
  [ "$result_json" = "1" ] || return 0
  python3 - "$1" "$2" "$3" <<'PY'
import json, sys
print(json.dumps({
    "status": sys.argv[1],
    "reason": sys.argv[2],
    "receipt": sys.argv[3] or None,
}, ensure_ascii=False))
PY
}

# 全局开关: 0 = 关闭
if [ "${JARVIS_NOTIFY_DINGTALK:-1}" = "0" ]; then
  log_skip "JARVIS_NOTIFY_DINGTALK=0"
  emit_result "skipped" "disabled" "$out_track_id"
  exit 0
fi

# opt-out 名单(gitignored)
optout_file="$jarvis_root/config/dingtalk-optout.txt"
if [ -f "$optout_file" ] && grep -qxF "$staff_id" "$optout_file"; then
  log_skip "staff in $optout_file"
  emit_result "skipped" "optout" "$out_track_id"
  exit 0
fi

# 组合消息: 标题一行加粗,空行,正文
message="**${title}**

${body}"

# dry-run 只受全局开关 + opt-out 影响(上面已判),不需要凭据/skill 路径
if [ "$dry_run" = "1" ]; then
  echo "[NOTIFY-DRY-RUN] staffId=$staff_id" >&2
  echo "----- message -----"
  printf '%s\n' "$message"
  echo "----- /message -----"
  emit_result "skipped" "dry_run" "$out_track_id"
  exit 0
fi

# 凭据必备三件套(仅实际发送时校验)
missing=""
[ -z "${DINGTALK_APP_KEY:-}" ]     && missing="$missing DINGTALK_APP_KEY"
[ -z "${DINGTALK_APP_SECRET:-}" ]  && missing="$missing DINGTALK_APP_SECRET"
[ -z "${DINGTALK_TEMPLATE_ID:-}" ] && missing="$missing DINGTALK_TEMPLATE_ID"
if [ -n "$missing" ]; then
  log_skip "missing env:$missing"
  emit_result "failed" "missing_credentials" "$out_track_id"
  exit 0
fi

# 定位 dingtalk-ai-card streaming.py:
#   优先 DINGTALK_SKILL(bridge/jarvis.env 已支持)
#   其次仓内 .claude/skills/dingtalk-ai-card/scripts/streaming.py
#   再次 ~/.claude/skills/dingtalk-ai-card/scripts/streaming.py
streaming=""
for cand in \
  "${DINGTALK_SKILL:-}" \
  "$jarvis_root/.claude/skills/dingtalk-ai-card/scripts/streaming.py" \
  "$HOME/.claude/skills/dingtalk-ai-card/scripts/streaming.py"; do
  [ -z "$cand" ] && continue
  if [ -f "$cand" ]; then
    streaming="$cand"
    break
  fi
done
if [ -z "$streaming" ]; then
  log_skip "streaming.py not found (set DINGTALK_SKILL to override)"
  emit_result "failed" "missing_sender" "$out_track_id"
  exit 0
fi

# 实际发送: --no-stream(通知类一次性投递,无打字机效果)
send_args=(
  "$streaming"
  --to "$staff_id"
  --template-id "$DINGTALK_TEMPLATE_ID"
  --no-stream
  -m "$message"
)
[ -n "$out_track_id" ] && send_args+=(--out-track-id "$out_track_id")
send_out=""
if ! send_out="$(python3 "${send_args[@]}" 2>&1)"; then
  ts="$(date +%Y%m%d-%H%M%S)"
  esc_dir="$(lib_escalation_dir)"
  mkdir -p "$esc_dir"
  fail_file="$esc_dir/notify-fail-${ts}-${staff_id}.md"
  {
    echo "# 钉钉私信失败: $staff_id"
    echo ""
    echo "- 时间: $(date -u +%FT%TZ)"
    echo "- 标题: $title"
    echo ""
    echo "## 消息内容"
    echo ""
    printf '%s\n' "$message"
  } > "$fail_file"
  echo "[NOTIFY-FAIL] staffId=$staff_id -> $fail_file" >&2
  receipt="$(printf '%s\n' "$send_out" | python3 -c 'import json,sys
value=""
for line in sys.stdin:
    try:
        data=json.loads(line)
        value=data.get("outTrackId") or value
    except Exception:
        pass
print(value)' 2>/dev/null || true)"
  emit_result "failed" "send_failed" "${receipt:-$out_track_id}"
  exit 0
fi

echo "[NOTIFY-OK] staffId=$staff_id title=\"$title\"" >&2
receipt="$(printf '%s\n' "$send_out" | python3 -c 'import json,sys
value=""
for line in sys.stdin:
    try:
        data=json.loads(line)
        value=data.get("outTrackId") or value
    except Exception:
        pass
print(value)' 2>/dev/null || true)"
emit_result "sent" "delivered" "${receipt:-$out_track_id}"
exit 0
