#!/usr/bin/env bash
# Hermetic ordering/compatibility tests for claim.sh's interactive DB gate.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/config" "$tmp/bin" "$tmp/.my-day"

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
if [ "$1 $2 $3" = "project workitem get" ]; then
  if [ "${A1_GET_FAIL_ONCE:-0}" = "1" ] && [ ! -f "$TEST_DIR/a1-get-failed-once" ]; then
    : > "$TEST_DIR/a1-get-failed-once"
    exit 1
  fi
  [ "${A1_GET_RC:-0}" = "0" ] || exit "$A1_GET_RC"
  tags="$(cat "$A1_STATE" 2>/dev/null || true)"
  tags_ml="${tags//,/, }"
  printf '{"subject":"%s","fields":[{"identifier":"tag","displayValue":"%s","value":"%s"},{"identifier":"workitemType","displayValue":"需求"},{"identifier":"status","displayValue":"%s"}]}\n' "${A1_TITLE-Interactive test title}" "$tags_ml" "$tags_ml" "${A1_STATUS:-处理中}"
  exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
  args=("$@"); tag=""; status=""
  for ((i=0; i<${#args[@]}; i++)); do
    [ "${args[$i]}" = "--tag" ] && tag="${args[$((i+1))]}"
    [ "${args[$i]}" = "--status" ] && status="${args[$((i+1))]}"
  done
  if [ -n "$tag" ]; then
    [ "${A1_TAG_RC:-0}" = "0" ] || exit "$A1_TAG_RC"
    printf '%s' "$tag" > "$A1_STATE"
  fi
  if [ -n "$status" ] && [ "${A1_REJECT_STATUS:-}" = "$status" ]; then
    exit 1
  fi
  exit 0
fi
exit 0
STUB
chmod +x "$tmp/bin/a1"

cat > "$tmp/interactive-runner.sh" <<'STUB'
#!/usr/bin/env bash
echo "worker:$*" >> "$TEST_LOG"
[ "$1" = "cli" ] && shift
case "${1:-}" in
  prepare-claim)
    rc="${IW_PREPARE_RC:-0}"
    [ "$rc" = "0" ] || exit "$rc"
    printf '%s\n' "${IW_PREPARE_JSON:-{\"accepted\":true,\"proceed\":true,\"operationStatus\":\"SENDING\"}}"
    ;;
  operation-ack) exit "${IW_ACK_RC:-0}" ;;
  operation-fail) exit 0 ;;
  operation-begin)
    rc="${IW_BEGIN_RC:-0}"
    [ "$rc" = "0" ] || exit "$rc"
    default_begin='{"accepted":true,"proceed":true,"needsReadback":false,"operationStatus":"SENDING"}'
    printf '%s\n' "${IW_BEGIN_JSON:-$default_begin}"
    ;;
  operation-abort) exit "${IW_ABORT_RC:-0}" ;;
  operation-reconcile) printf '{"proceed":false}\n'; exit "${IW_RECONCILE_RC:-0}" ;;
  has-current) exit "${IW_CURRENT_RC:-0}" ;;
  suspend|complete) exit "${IW_TRANSITION_RC:-0}" ;;
  *) exit 64 ;;
esac
STUB
chmod +x "$tmp/interactive-runner.sh"

export JARVIS_ROOT="$tmp"
export JARVIS_A1="$tmp/bin/a1"
export JARVIS_INTERACTIVE_WORKER_RUNNER="$tmp/interactive-runner.sh"
# This suite deliberately simulates an interactive session and drives claim.sh's
# control-plane gate through the stub runner above. Clear any signal inherited
# from the host session first (canonical list in test/lib/hermetic.sh), then open
# exactly one, and pin the endpoint at the discard port so a broken stub fails
# fast instead of reaching production (Aone 85192197).
# shellcheck source=lib/hermetic.sh
source "$test_dir/lib/hermetic.sh"
jarvis_test_hermetic_isolate_signals
jarvis_test_hermetic_pin_control_plane
export CODEX_THREAD_ID="codex-test-thread"
export TEST_LOG="$tmp/events.log"
export TEST_DIR="$tmp"
export A1_STATE="$tmp/aone-tags"
export JARVIS_CLAIM_READBACK_SLEEP=0
export JARVIS_CLAIM_PROGRESS=0
export JARVIS_SKIP_MR_GATE=1

pass=0
fail=0
ok() { echo "PASS: $1"; pass=$((pass + 1)); }
bad() { echo "FAIL: $1"; fail=$((fail + 1)); }
reset_case() {
  : > "$TEST_LOG"
  rm -f "$TEST_DIR/a1-get-failed-once"
  printf '%s' "${1:-}" > "$A1_STATE"
  unset IW_PREPARE_RC IW_PREPARE_JSON IW_ACK_RC IW_CURRENT_RC IW_TRANSITION_RC
  unset IW_BEGIN_RC IW_BEGIN_JSON IW_ABORT_RC IW_RECONCILE_RC
  unset A1_TAG_RC A1_REJECT_STATUS A1_GET_RC A1_GET_FAIL_ONCE
  export A1_STATUS="${2:-处理中}"
  export A1_TITLE="Interactive test title"
}
run_claim() {
  bash "$repo_root/bootstrap/claim.sh" "$@" >/dev/null 2>&1
}

reset_case
export IW_PREPARE_RC=10
if run_claim claim 84345050 2100304; then
  bad "409 conflict returns nonzero"
elif grep -q 'a1:project workitem update' "$TEST_LOG"; then
  bad "409 conflict performs no Aone writes"
else
  ok "409 conflict fails after title read and before Aone writes"
fi

reset_case
export IW_PREPARE_RC=11
if run_claim claim 84345050 2100304; then
  bad "control-plane outage returns nonzero"
elif grep -q 'a1:project workitem update' "$TEST_LOG"; then
  bad "control-plane outage performs no Aone writes"
else
  ok "control-plane outage fails after title read and before Aone writes"
fi

reset_case
unset CODEX_THREAD_ID CLAUDE_CODE_SESSION_ID
export JARVIS_INTERACTIVE_CLIENT=claude
export JARVIS_INTERACTIVE_SESSION_ID=claude-env-file-session
export IW_PREPARE_RC=10
if run_claim claim 84345050 2100304 || grep -q 'a1:project workitem update' "$TEST_LOG"; then
  bad "persisted Claude hook context enters the DB gate"
else
  ok "persisted Claude hook context enters the DB gate"
fi
unset JARVIS_INTERACTIVE_CLIENT JARVIS_INTERACTIVE_SESSION_ID
export CODEX_THREAD_ID="codex-test-thread"

reset_case
if ! run_claim claim 84345050 2100304; then
  bad "interactive claim succeeds"
elif ! grep -q '^worker:cli prepare-claim 84345050 2100304 Interactive test title 处理中$' "$TEST_LOG"; then
  bad "interactive claim forwards the itemId title and pre-claim status"
else
  prepare_line="$(grep -n 'worker:cli prepare-claim' "$TEST_LOG" | head -1 | cut -d: -f1)"
  tag_line="$(grep -n 'a1:project workitem update .*--tag' "$TEST_LOG" | head -1 | cut -d: -f1)"
  ack_line="$(grep -n 'worker:cli operation-ack' "$TEST_LOG" | head -1 | cut -d: -f1)"
  if [ -n "$prepare_line" ] && [ -n "$tag_line" ] && [ -n "$ack_line" ] \
      && [ "$prepare_line" -lt "$tag_line" ] && [ "$tag_line" -lt "$ack_line" ]; then
    ok "claim order is DB prepare then Aone tag/readback then receipt ACK"
  else
    bad "claim ordering"
  fi
fi

reset_case "" "待认领"
if ! run_claim claim 84345050 2100304; then
  bad "interactive claim with observed status succeeds"
elif ! grep -q '^worker:cli prepare-claim 84345050 2100304 Interactive test title 待认领$' "$TEST_LOG"; then
  bad "interactive claim forwards the observed pre-claim status"
else
  ok "interactive claim forwards the observed pre-claim status"
fi

reset_case
export A1_TITLE="   "
if ! run_claim claim 84345050 2100304; then
  bad "blank title remains backward compatible"
elif ! grep -q '^worker:cli prepare-claim 84345050 2100304' "$TEST_LOG"; then
  bad "blank title still reaches control-plane claim"
else
  ok "blank title is omitted without blocking claim"
fi

reset_case
export A1_GET_FAIL_ONCE=1
if ! run_claim claim 84345050 2100304; then
  bad "title point-read failure remains backward compatible"
elif ! grep -q '^worker:cli prepare-claim 84345050 2100304' "$TEST_LOG"; then
  bad "title point-read failure still reaches control-plane claim"
else
  ok "title point-read failure does not block claim"
fi

reset_case
export IW_ACK_RC=13
if run_claim claim 84345050 2100304; then
  bad "ACK failure returns nonzero"
elif ! grep -q 'a1:project workitem update .*--tag' "$TEST_LOG"; then
  bad "ACK failure happens after Aone tag"
else
  ok "ACK failure stops after the persisted Aone effect"
fi

reset_case jarvis-claimed
if run_claim release 84345050 2100304 \
    && grep -q 'worker:cli suspend 84345050' "$TEST_LOG"; then
  tag_line="$(grep -n 'a1:project workitem update .*--tag' "$TEST_LOG" | head -1 | cut -d: -f1)"
  suspend_line="$(grep -n 'worker:cli suspend' "$TEST_LOG" | head -1 | cut -d: -f1)"
  if [ "$tag_line" -lt "$suspend_line" ]; then ok "release idles Aone then suspends session"; else bad "release order"; fi
else
  bad "release session transition"
fi

reset_case jarvis-claimed 已完成
if run_claim finish 84345050 2100304 \
    && grep -q 'worker:cli complete 84345050' "$TEST_LOG"; then
  ok "finish completes database session at terminal Aone state"
else
  bad "finish complete transition"
fi

reset_case jarvis-claimed 处理中
export A1_REJECT_STATUS=已完成
if run_claim finish 84345050 2100304 \
    && grep -q 'worker:cli suspend 84345050' "$TEST_LOG"; then
  ok "downgraded finish suspends database session"
else
  bad "downgraded finish suspend transition"
fi

# --- release/finish external-write receipts (docs/aone-operation-receipts.md) ---

reset_case jarvis-claimed
if run_claim release 84345050 2100304; then
  begin_line="$(grep -n 'worker:cli operation-begin 84345050 release-tag idle --replay-safe' "$TEST_LOG" | head -1 | cut -d: -f1)"
  tag_line="$(grep -n 'a1:project workitem update .*--tag' "$TEST_LOG" | head -1 | cut -d: -f1)"
  ack_line="$(grep -n 'worker:cli operation-ack 84345050 aone:2100304:84345050:tag:jarvis-idle' "$TEST_LOG" | head -1 | cut -d: -f1)"
  suspend_line="$(grep -n 'worker:cli suspend' "$TEST_LOG" | head -1 | cut -d: -f1)"
  if [ -n "$begin_line" ] && [ -n "$tag_line" ] && [ -n "$ack_line" ] && [ -n "$suspend_line" ] \
      && [ "$begin_line" -lt "$tag_line" ] && [ "$tag_line" -lt "$ack_line" ] \
      && [ "$ack_line" -lt "$suspend_line" ]; then
    ok "release receipt order is begin then tag then ACK then suspend"
  else
    bad "release receipt ordering"
  fi
else
  bad "release with receipt succeeds"
fi

reset_case jarvis-claimed
export IW_BEGIN_RC=10
if run_claim release 84345050 2100304; then
  bad "release begin conflict returns nonzero"
elif grep -q 'a1:project workitem update' "$TEST_LOG"; then
  bad "release begin failure performs no Aone writes"
elif grep -q 'worker:cli suspend' "$TEST_LOG"; then
  bad "release begin failure keeps session active"
else
  ok "release begin failure fails closed before Aone"
fi

reset_case jarvis-idle
export IW_BEGIN_JSON='{"accepted":true,"proceed":false,"needsReadback":false,"operationStatus":"ACKED"}'
if run_claim release 84345050 2100304 \
    && ! grep -q 'a1:project workitem update' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-ack' "$TEST_LOG" \
    && grep -q 'worker:cli suspend 84345050' "$TEST_LOG"; then
  ok "ACKED release skips duplicate tag write and still suspends"
else
  bad "ACKED release dedup"
fi

reset_case jarvis-claimed
export A1_TAG_RC=7
if run_claim release 84345050 2100304; then
  bad "release tag failure returns nonzero"
elif grep -q 'worker:cli operation-abort 84345050' "$TEST_LOG" \
    && ! grep -q 'worker:cli operation-abort 84345050 .*--unknown' "$TEST_LOG" \
    && ! grep -q 'worker:cli suspend' "$TEST_LOG"; then
  ok "release tag failure aborts receipt and keeps session"
else
  bad "release tag failure abort"
fi

reset_case jarvis-claimed
export A1_GET_RC=1
if run_claim release 84345050 2100304; then
  bad "release readback outage returns nonzero"
elif grep -q 'worker:cli operation-abort 84345050 .* --unknown' "$TEST_LOG" \
    && ! grep -q 'worker:cli suspend' "$TEST_LOG"; then
  ok "release inconclusive readback freezes receipt UNKNOWN"
else
  bad "release inconclusive readback"
fi

reset_case jarvis-claimed 已完成
if run_claim finish 84345050 2100304; then
  begin_line="$(grep -n 'worker:cli operation-begin 84345050 finish-tag done --replay-safe' "$TEST_LOG" | head -1 | cut -d: -f1)"
  tag_line="$(grep -n 'a1:project workitem update .*--tag' "$TEST_LOG" | head -1 | cut -d: -f1)"
  ack_line="$(grep -n 'worker:cli operation-ack 84345050 aone:2100304:84345050:tag:jarvis-done' "$TEST_LOG" | head -1 | cut -d: -f1)"
  complete_line="$(grep -n 'worker:cli complete' "$TEST_LOG" | head -1 | cut -d: -f1)"
  if [ -n "$begin_line" ] && [ -n "$tag_line" ] && [ -n "$ack_line" ] && [ -n "$complete_line" ] \
      && [ "$begin_line" -lt "$tag_line" ] && [ "$tag_line" -lt "$ack_line" ] \
      && [ "$ack_line" -lt "$complete_line" ]; then
    ok "finish receipt order is begin then tag then ACK then complete"
  else
    bad "finish receipt ordering"
  fi
else
  bad "finish with receipt succeeds"
fi

reset_case jarvis-claimed 已完成
export IW_BEGIN_RC=11
if run_claim finish 84345050 2100304; then
  bad "finish begin outage returns nonzero"
elif grep -q 'a1:project workitem update' "$TEST_LOG"; then
  bad "finish begin failure performs no Aone writes"
elif grep -q 'worker:cli complete' "$TEST_LOG"; then
  bad "finish begin failure keeps session active"
else
  ok "finish begin failure fails closed before Aone"
fi

echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
