#!/usr/bin/env bash
# Hermetic tests for wrap.sh's fenced external-write operation receipts
# (docs/aone-operation-receipts.md): begin→send→ack ordering, ACKED dedup,
# readback/reconcile convergence, fail-closed semantics, and the untouched
# non-fenced bare-write path.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/config" "$tmp/bin" "$tmp/.my-day" "$tmp/runs"

cat > "$tmp/config/pools.json" <<'JSON'
{
  "claim": {
    "tag": "jarvis-claimed",
    "idle_tag": "jarvis-idle",
    "done_tag": "jarvis-done",
    "done_status": "已完成",
    "done_statuses": ["已完成"],
    "progress_status": "处理中",
    "start_statuses": ["待处理"]
  },
  "pools": []
}
JSON

cat > "$tmp/bin/a1" <<'STUB'
#!/usr/bin/env bash
echo "a1:$*" >> "$TEST_LOG"
if [ "$1 $2 $3 ${4:-}" = "project workitem comment create" ]; then
  [ "${A1_COMMENT_RC:-0}" = "0" ] || exit "$A1_COMMENT_RC"
  if [ "${A1_COMMENT_NO_ID:-0}" = "1" ]; then echo '{}'; else echo '{"id": 987}'; fi
  exit 0
fi
if [ "$1 $2 $3 ${4:-}" = "project workitem comment list" ]; then
  cat "${A1_COMMENTS_FILE:-/dev/null}" 2>/dev/null || echo '[]'
  exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
  [ "${A1_STATUS_RC:-0}" = "0" ] || exit "$A1_STATUS_RC"
  exit 0
fi
exit 0
STUB
chmod +x "$tmp/bin/a1"

cat > "$tmp/interactive-runner.sh" <<'STUB'
#!/usr/bin/env bash
echo "worker:$*" >> "$TEST_LOG"
[ "$1" = "cli" ] && shift
default_begin='{"accepted":true,"proceed":true,"needsReadback":false,"operationStatus":"SENDING"}'
case "${1:-}" in
  has-current) exit "${IW_CURRENT_RC:-0}" ;;
  operation-begin)
    kind="${3:-}"
    case "$kind" in
      comment)
        n=$(( $(cat "$IW_BEGIN_COUNT" 2>/dev/null || echo 0) + 1 ))
        printf '%s' "$n" > "$IW_BEGIN_COUNT"
        rc="${IW_BEGIN_RC_COMMENT:-0}"
        [ "$rc" = "0" ] || exit "$rc"
        if [ "$n" -ge 2 ] && [ -n "${IW_BEGIN_JSON_COMMENT_2:-}" ]; then
          printf '%s\n' "$IW_BEGIN_JSON_COMMENT_2"
        else
          printf '%s\n' "${IW_BEGIN_JSON_COMMENT:-$default_begin}"
        fi
        ;;
      status)
        rc="${IW_BEGIN_RC_STATUS:-0}"
        [ "$rc" = "0" ] || exit "$rc"
        printf '%s\n' "${IW_BEGIN_JSON_STATUS:-$default_begin}"
        ;;
      *) exit 64 ;;
    esac
    ;;
  operation-ack) exit "${IW_ACK_RC:-0}" ;;
  operation-abort) exit "${IW_ABORT_RC:-0}" ;;
  operation-reconcile) printf '{"proceed":false}\n'; exit "${IW_RECONCILE_RC:-0}" ;;
  *) exit 64 ;;
esac
STUB
chmod +x "$tmp/interactive-runner.sh"

export JARVIS_ROOT="$tmp"
export JARVIS_A1="$tmp/bin/a1"
export JARVIS_RUNS_DIR="$tmp/runs"
export JARVIS_INTERACTIVE_WORKER_RUNNER="$tmp/interactive-runner.sh"
export TEST_LOG="$tmp/events.log"
export IW_BEGIN_COUNT="$tmp/begin-count"
export A1_COMMENTS_FILE="$tmp/comments.json"
export CODE_DIR="$tmp"   # 非 git 目录 → code_footer 为空，正文确定可复算
# fenced 标记：模拟 Codex 交互会话；先清可能来自外层 Claude 会话的标记
unset CLAUDE_CODE_SESSION_ID JARVIS_INTERACTIVE_CLIENT JARVIS_INTERACTIVE_SESSION_ID
export CODEX_THREAD_ID="codex-wrap-test-thread"

WRAP="$repo_root/bootstrap/wrap.sh"
FORMAT="$repo_root/bootstrap/aone-comment-format.sh"

pass=0
fail=0
ok() { echo "PASS: $1"; pass=$((pass + 1)); }
bad() { echo "FAIL: $1"; fail=$((fail + 1)); }
reset_case() {
  : > "$TEST_LOG"
  rm -f "$IW_BEGIN_COUNT"
  echo '[]' > "$A1_COMMENTS_FILE"
  unset IW_CURRENT_RC IW_ACK_RC IW_ABORT_RC IW_RECONCILE_RC
  unset IW_BEGIN_RC_COMMENT IW_BEGIN_JSON_COMMENT IW_BEGIN_JSON_COMMENT_2
  unset IW_BEGIN_RC_STATUS IW_BEGIN_JSON_STATUS
  unset A1_COMMENT_RC A1_COMMENT_NO_ID A1_STATUS_RC
}
line_of() { grep -n "$1" "$TEST_LOG" | head -1 | cut -d: -f1; }

# --- 1. fenced sync：begin < comment create < ack ---
reset_case
if bash "$WRAP" sync WI-1 "fenced progress note" >/dev/null 2>&1; then
  begin_line="$(line_of 'worker:cli operation-begin WI-1 comment')"
  create_line="$(line_of 'a1:project workitem comment create WI-1')"
  ack_line="$(line_of 'worker:cli operation-ack WI-1 aone:WI-1:comment:987')"
  if [ -n "$begin_line" ] && [ -n "$create_line" ] && [ -n "$ack_line" ] \
      && [ "$begin_line" -lt "$create_line" ] && [ "$create_line" -lt "$ack_line" ]; then
    ok "fenced sync order is begin then comment create then ACK"
  else
    bad "fenced sync receipt ordering"
  fi
else
  bad "fenced sync succeeds"
fi

# --- 2. ACKED 去重：不重发评论 ---
reset_case
export IW_BEGIN_JSON_COMMENT='{"accepted":true,"proceed":false,"needsReadback":false,"operationStatus":"ACKED"}'
if bash "$WRAP" sync WI-2 "already acked note" >/dev/null 2>&1 \
    && ! grep -q 'a1:project workitem comment create' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-ack' "$TEST_LOG"; then
  ok "ACKED comment receipt skips duplicate send"
else
  bad "ACKED comment dedup"
fi

# --- 3. needsReadback 命中：reconcile --found 且零重发 ---
reset_case
export IW_BEGIN_JSON_COMMENT='{"accepted":true,"proceed":false,"needsReadback":true,"operationStatus":"SENDING","operationId":"op-1"}'
expected="$(bash "$FORMAT" "sync progress for readback")"
python3 - "$expected" > "$A1_COMMENTS_FILE" <<'PY'
import json, sys
print(json.dumps([
    {"id": 111, "content": "unrelated comment"},
    {"id": 555, "content": sys.argv[1]},
], ensure_ascii=False))
PY
if bash "$WRAP" sync WI-3 "sync progress for readback" >/dev/null 2>&1 \
    && grep -q 'worker:cli operation-reconcile WI-3 --found aone:WI-3:comment:555' "$TEST_LOG" \
    && ! grep -q 'a1:project workitem comment create' "$TEST_LOG"; then
  ok "readback hit reconciles found with zero resend"
else
  bad "readback hit convergence"
fi

# --- 4. needsReadback 未命中：not-found → 重 begin → 发送 ---
reset_case
export IW_BEGIN_JSON_COMMENT='{"accepted":true,"proceed":false,"needsReadback":true,"operationStatus":"SENDING","operationId":"op-1"}'
export IW_BEGIN_JSON_COMMENT_2='{"accepted":true,"proceed":true,"needsReadback":false,"operationStatus":"SENDING"}'
if bash "$WRAP" sync WI-4 "sync progress not found" >/dev/null 2>&1; then
  notfound_line="$(line_of 'worker:cli operation-reconcile WI-4 --not-found')"
  create_line="$(line_of 'a1:project workitem comment create WI-4')"
  ack_line="$(line_of 'worker:cli operation-ack WI-4')"
  begin_count="$(cat "$IW_BEGIN_COUNT" 2>/dev/null || echo 0)"
  if [ -n "$notfound_line" ] && [ -n "$create_line" ] && [ -n "$ack_line" ] \
      && [ "$notfound_line" -lt "$create_line" ] && [ "$create_line" -lt "$ack_line" ] \
      && [ "$begin_count" = "2" ]; then
    ok "readback miss reconciles not-found then re-begins and sends once"
  else
    bad "readback miss convergence"
  fi
else
  bad "readback miss sync succeeds"
fi

# --- 5. begin 失败 fail-closed：零 Aone 写、退出非零 ---
reset_case
export IW_BEGIN_RC_COMMENT=10
if bash "$WRAP" sync WI-5 "must not land" >/dev/null 2>&1; then
  bad "begin failure returns nonzero"
elif grep -q '^a1:' "$TEST_LOG"; then
  bad "begin failure performs no Aone calls"
else
  ok "begin failure fails closed with zero Aone writes"
fi

# --- 6. done：comment 回执先行，status 回执随后 ---
reset_case
if bash "$WRAP" done WI-6 "finished work" "已完成" >/dev/null 2>&1; then
  comment_ack_line="$(line_of 'worker:cli operation-ack WI-6 aone:WI-6:comment:987')"
  status_begin_line="$(line_of 'worker:cli operation-begin WI-6 status 已完成 --replay-safe')"
  status_update_line="$(line_of 'a1:project workitem update WI-6 --status 已完成')"
  status_ack_line="$(line_of 'worker:cli operation-ack WI-6 aone:WI-6:status:已完成')"
  if [ -n "$comment_ack_line" ] && [ -n "$status_begin_line" ] \
      && [ -n "$status_update_line" ] && [ -n "$status_ack_line" ] \
      && [ "$comment_ack_line" -lt "$status_begin_line" ] \
      && [ "$status_begin_line" -lt "$status_update_line" ] \
      && [ "$status_update_line" -lt "$status_ack_line" ]; then
    ok "done runs serial comment receipt then status receipt"
  else
    bad "done receipt ordering"
  fi
  if ls "$JARVIS_RUNS_DIR"/*-WI-6.md >/dev/null 2>&1; then
    ok "done still writes local run_done audit"
  else
    bad "done run_done audit"
  fi
else
  bad "fenced done succeeds"
fi

# --- 7. status begin 失败：fail-closed，不写状态 ---
reset_case
export IW_BEGIN_RC_STATUS=11
if bash "$WRAP" done WI-7 "status blocked" "已完成" >/dev/null 2>&1; then
  bad "status begin failure returns nonzero"
elif grep -q 'a1:project workitem update WI-7 --status' "$TEST_LOG"; then
  bad "status begin failure performs no status write"
else
  ok "status begin failure fails closed before status write"
fi

# --- 8. a1 comment 失败：abort（可重试）+ 退出非零 ---
reset_case
export A1_COMMENT_RC=1
if bash "$WRAP" sync WI-8 "comment will fail" >/dev/null 2>&1; then
  bad "comment failure returns nonzero"
elif grep -q 'worker:cli operation-abort WI-8' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-abort WI-8 .*--unknown' "$TEST_LOG"; then
  ok "comment failure aborts receipt as retryable"
else
  bad "comment failure abort"
fi

# --- 9. 评论退出码 0 但拿不到 id：abort --unknown ---
reset_case
export A1_COMMENT_NO_ID=1
if bash "$WRAP" sync WI-9 "comment id missing" >/dev/null 2>&1; then
  bad "indeterminate comment returns nonzero"
elif grep -q 'worker:cli operation-abort WI-9 .* --unknown' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-ack' "$TEST_LOG"; then
  ok "indeterminate comment result freezes receipt UNKNOWN"
else
  bad "indeterminate comment freeze"
fi

# --- 10. 非 fenced 上下文：行为与现状完全一致（裸写、无 worker 调用） ---
reset_case
if env -u CODEX_THREAD_ID bash "$WRAP" sync WI-10 "plain sync" >/dev/null 2>&1 \
    && grep -q 'a1:project workitem comment create WI-10' "$TEST_LOG" \
    && ! grep -q 'a1:project workitem comment create WI-10 .* -f json' "$TEST_LOG" \
    && ! grep -q '^worker:' "$TEST_LOG"; then
  ok "non-fenced sync keeps the bare write path"
else
  bad "non-fenced sync bare path"
fi

reset_case
if env -u CODEX_THREAD_ID A1_COMMENT_RC=1 bash "$WRAP" sync WI-11 "warn only" >/dev/null 2>&1; then
  ok "non-fenced sync keeps warn-only failure semantics"
else
  bad "non-fenced sync warn-only semantics"
fi

# --- 11. fenced 但不持有该单（has-current 假）：同现状裸写 ---
reset_case
export IW_CURRENT_RC=1
if bash "$WRAP" sync WI-12 "other ticket backfill" >/dev/null 2>&1 \
    && grep -q 'a1:project workitem comment create WI-12' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-begin' "$TEST_LOG"; then
  ok "non-current ticket keeps the bare write path"
else
  bad "non-current ticket bare path"
fi

echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
