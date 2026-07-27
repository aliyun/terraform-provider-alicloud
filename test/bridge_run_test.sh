#!/usr/bin/env bash
# test/bridge_run_test.sh — hermetic tests for bridge/run.sh (单一入口进程管理器).
#
# 全程 mock: 用一个假 bot(桩 python)顶替 `python3 bridge/jarvis_dingtalk_bot.py`,
# 状态目录/env 文件/python/bot 路径全部走 JARVIS_BRIDGE_* 覆盖到临时目录, 不碰真
# .my-day、不连钉钉、不起真 claude。
#
# 覆盖: bash -n 语法; start 幂等/pidfile/模式判定(有无凭证 env + 走 jarvis.env source);
#       stop 清理/重复 stop; restart 换 pid; status 运行/停止/模式; logs 非阻塞;
#       dry-run 透传; 启动失败检测(立即退出 / 首行 ERROR 但驻留).
#
# Run: bash test/bridge_run_test.sh   (exit 0 = all pass)

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
RUNSH="$repo_root/bridge/run.sh"
INSTALLER="$repo_root/bridge/install-launchd.sh"
PY_YAML_SITE="$(python3 -c \
  'from pathlib import Path; import yaml; print(Path(yaml.__file__).resolve().parents[1])')"

if [ ! -f "$RUNSH" ]; then
  echo "SKIP bridge_run_test: $RUNSH not found"
  exit 0
fi

TMP="$(mktemp -d 2>/dev/null || mktemp -d -t bridge_run)"
FAKE_HOME="$TMP/home"
FAKE_BREW="$TMP/opt/homebrew"
STATE="$TMP/state"
FAKEPY="$TMP/fakepy"
BOTFILE="$TMP/bot.py"       # dummy path handed to the fake python; content unused
BSENV="$TMP/bootstrap.env"  # stand-in for bootstrap/.env
JENV="$TMP/jarvis.env"      # stand-in for bridge/jarvis.env
FAKECTL="$TMP/launchctl"
FAKECTL_STATE="$TMP/launchctl-state"
LPLIST="$TMP/com.jarvis.test.plist"
: >"$BOTFILE"; : >"$BSENV"; : >"$JENV"
mkdir -p "$FAKE_HOME/.local/bin" "$FAKE_BREW/bin" "$FAKE_BREW/sbin"
printf '%s\n' '<?xml version="1.0"?><plist version="1.0"><dict><key>Label</key><string>com.jarvis.test</string></dict></plist>' >"$LPLIST"
cat >"$FAKE_BREW/bin/pgrep" <<'FAKE'
#!/usr/bin/env bash
# Lock-cleanup tests must not depend on unrelated host a1 processes.
exit 1
FAKE
cat >"$FAKE_BREW/bin/ps" <<'FAKE'
#!/usr/bin/env bash
# Watchdog process age is deterministic; liveness is still checked with kill -0.
case " $* " in
  *" etime="*) printf '00:01\n' ;;
  *) /bin/ps "$@" ;;
esac
FAKE
chmod +x "$FAKE_BREW/bin/pgrep" "$FAKE_BREW/bin/ps"

pass=0; fail=0
ok() { echo "PASS $1"; pass=$((pass+1)); }
no() { echo "FAIL $1"; fail=$((fail+1)); }
has() { case "$2" in *"$1"*) ok "$3";; *) no "$3 [missing '$1']";; esac; }
hasnot() { case "$2" in *"$1"*) no "$3 [unexpected '$1']";; *) ok "$3";; esac; }

# The fake bot: stands in for the real python bot. Behaviour via FAKE_BOT_MODE.
cat >"$FAKEPY" <<'FAKE'
#!/usr/bin/env bash
ts="$(date '+%Y-%m-%d %H:%M:%S')"
if [ "${1:-}" = "-c" ]; then exit 0; fi
if [ "${1:-}" = "-m" ] && [ "${2:-}" = "bridge.main" ]; then
  if [ "${3:-}" = "--validate" ]; then
    [ "${FAKE_VALIDATE_MODE:-ok}" = "ok" ] && exit 0
    echo "$ts ERROR [MainThread] fake bridge validation failed"
    exit 2
  fi
  [ -n "${FAKE_PROCESS_STATE:-}" ] && printf '%s\n' "$$" >"$FAKE_PROCESS_STATE/scheduler-started"
  case "${FAKE_SCHEDULER_MODE:-ready}" in
    exit)
      echo "$ts ERROR [MainThread] fake scheduler immediate exit"
      exit 2 ;;
    no-ready)
      echo "$ts INFO [MainThread] fake bridge supervisor deliberately not READY"
      exec sleep 300 ;;
    once)
      echo "$ts INFO [MainThread] Bridge READY pid=$$ role=supervisor components=fake JARVIS_NO_DINGTALK=${JARVIS_NO_DINGTALK:-0}"
      exit 0 ;;
    cleanup-child-pidfile)
      echo "$ts INFO [MainThread] Bridge READY pid=$$ role=supervisor components=fake"
      trap 'rm -f "$FAKE_PROCESS_STATE/persistent-worker.pid" "$FAKE_PROCESS_STATE/persistent-worker.pid.identity"; exit 0' TERM INT
      while true; do sleep 1; done ;;
    *)
      if [ "${JARVIS_NO_DINGTALK:-}" = "1" ]; then
        echo "$ts WARNING [MainThread] [NO-DINGTALK] fake supervised bot"
      else
        echo "$ts INFO [MainThread] starting DingTalk stream listener… (fake supervised bot)"
      fi
      echo "$ts INFO [MainThread] Bridge READY pid=$$ role=supervisor components=fake"
      exec sleep 300 ;;
  esac
fi
if [ "${1:-}" = "-m" ] && [ "${2:-}" = "bridge.persistent_worker" ]; then
  [ -n "${FAKE_PROCESS_STATE:-}" ] && printf '%s\n' "$$" >"$FAKE_PROCESS_STATE/persistent-worker-started"
  case "${FAKE_PERSISTENT_WORKER_MODE:-ready}" in
    exit)
      echo "$ts ERROR [MainThread] fake persistent worker immediate exit"
      exit 2 ;;
    no-ready)
      echo "$ts INFO [MainThread] fake persistent worker deliberately not READY"
      exec sleep 300 ;;
    once)
      echo "$ts INFO [MainThread] foreground once JARVIS_NO_DINGTALK=${JARVIS_NO_DINGTALK:-0}"
      exit 0 ;;
    *)
      echo "$ts INFO [MainThread] Persistent worker READY pid=$$ worker=fake"
      exec sleep 300 ;;
  esac
fi
for a in "$@"; do
  if [ "$a" = "--dry-run-once" ]; then
    echo "$ts INFO [MainThread] dry-run ok (fake bot) PATH=$PATH JARVIS_ROOT=${JARVIS_ROOT:-}"
    exit 0
  fi
done
case "${FAKE_BOT_MODE:-stay}" in
  exit)
    echo "$ts ERROR [MainThread] fake bot immediate exit (missing creds sim)"
    exit 2 ;;
  errfirst)
    echo "$ts ERROR [MainThread] fake bot fatal startup line (stays alive)"
    exec sleep 300 ;;
  deaf)
    # 忽略 SIGTERM(模拟清理超时/卡死的 bot) → run.sh 宽限超时后必须 SIGKILL 兜底。
    # 用 while 循环(而非 exec sleep)才能保住 trap; SIGKILL 不可捕获, 故最终会被杀。
    trap '' TERM
    echo "$ts INFO [MainThread] deaf fake bot (ignores SIGTERM)"
    while true; do sleep 1; done ;;
  once)
    echo "$ts INFO [MainThread] foreground once JARVIS_NO_DINGTALK=${JARVIS_NO_DINGTALK:-0}"
    exit 0 ;;
  *)
    [ -n "${FAKE_PROCESS_STATE:-}" ] && printf '%s\n' "$$" >"$FAKE_PROCESS_STATE/bot-started"
    # Mirror the real bot's dual-line startup: the degraded bot prefixes EVERY line with
    # [NO-DINGTALK]; both modes log "scan scheduler started" (an ambiguous non-discriminator).
    if [ "${JARVIS_NO_DINGTALK:-}" = "1" ]; then
      echo "$ts WARNING [MainThread] [NO-DINGTALK] 降级模式启动 (fake bot)"
      echo "$ts INFO [MainThread] [NO-DINGTALK] scan scheduler started (fake bot)"
    else
      echo "$ts INFO [MainThread] starting DingTalk stream listener… (fake bot)"
      echo "$ts INFO [MainThread] scan scheduler started (fake bot)"
    fi
    exec sleep 300 ;;
esac
FAKE
chmod +x "$FAKEPY"

# Minimal launchctl state machine. Every call is recorded and scoped to TMP; it never
# touches the host's launchd domain.
cat >"$FAKECTL" <<'FAKE'
#!/usr/bin/env bash
set -u
state="${FAKE_LAUNCHCTL_STATE:?}"
mkdir -p "$state"
printf '%s\n' "$*" >>"$state/calls"
case "${1:-}" in
  print)
    [ -f "$state/loaded" ] || exit 113
    pid="$(cat "$state/pid" 2>/dev/null || true)"
    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
      printf '%s\n' '{' '    state = running' "    pid = $pid" '}'
    else
      printf '%s\n' '{' '    state = waiting' '}'
    fi ;;
  bootstrap)
    touch "$state/loaded" ;;
  bootout)
    pid="$(cat "$state/pid" 2>/dev/null || true)"
    [ -n "$pid" ] && kill -TERM "$pid" 2>/dev/null || true
    rm -f "$state/loaded" "$state/pid" ;;
  enable|disable)
    : ;;
  kill)
    pid="$(cat "$state/pid" 2>/dev/null || true)"
    [ -n "$pid" ] && kill -TERM "$pid" 2>/dev/null || true
    rm -f "$state/pid" ;;
  kickstart)
    [ -f "$state/loaded" ] || exit 1
    pid="$(cat "$state/pid" 2>/dev/null || true)"
    if [ "${2:-}" = "-k" ] && [ -n "$pid" ]; then
      kill -TERM "$pid" 2>/dev/null || true
      rm -f "$state/pid"
    fi
    sleep 300 </dev/null >/dev/null 2>&1 &
    pid=$!
    printf '%s\n' "$pid" >"$state/pid"
    if [ -n "${FAKE_LAUNCHCTL_LOG:-}" ]; then
      mkdir -p "$(dirname "$FAKE_LAUNCHCTL_LOG")"
      printf '%s INFO [MainThread] Bridge READY pid=%s role=supervisor\n' \
        "$(date '+%Y-%m-%d %H:%M:%S')" "$pid" >>"$FAKE_LAUNCHCTL_LOG"
    fi ;;
  *)
    echo "unexpected fake launchctl command: $*" >&2
    exit 2 ;;
esac
FAKE
chmod +x "$FAKECTL"

# Invoke run.sh with all state/env/deps redirected into the hermetic tmp sandbox.
run() {
  env JARVIS_BRIDGE_PYTHON="$FAKEPY" \
      JARVIS_BRIDGE_BOT="$BOTFILE" \
      JARVIS_BRIDGE_STATE_DIR="$STATE" \
      JARVIS_BRIDGE_BOOTSTRAP_ENV="$BSENV" \
      JARVIS_BRIDGE_ENV="$JENV" \
      JARVIS_BRIDGE_START_WAIT="1.0" \
      JARVIS_BRIDGE_ROLE="${TEST_ROLE:-scheduler}" \
      JARVIS_CONTROL_PLANE_TOKEN="test-token" \
      JARVIS_BRIDGE_STOP_WAIT="1" \
      JARVIS_BRIDGE_START_ROLLBACK_WAIT="1" \
      JARVIS_SCHEDULER_READY_WAIT="${TEST_READY_WAIT:-1}" \
      JARVIS_WATCHDOG_STARTUP_GRACE_SEC="${TEST_WATCHDOG_STARTUP_GRACE:-300}" \
      JARVIS_WATCHDOG_BEACON_STALE_SEC="${TEST_WATCHDOG_STALE:-300}" \
      JARVIS_WATCHDOG_TREE_TERM_SEC="${TEST_WATCHDOG_TERM:-1}" \
      JARVIS_WATCHDOG_CHILD_TERM_SEC="${TEST_WATCHDOG_CHILD_TERM:-0}" \
      JARVIS_A1_LOCK_STALE_SEC="${TEST_LOCK_STALE:-0}" \
      JARVIS_EXECUTOR_BEACON_PREFIX="$STATE/heartbeat.persistent-worker" \
      JARVIS_BRIDGE_TOOL_DIRS="$FAKE_BREW/sbin $FAKE_BREW/bin $FAKE_HOME/.local/bin" \
      JARVIS_BRIDGE_SUPERVISOR="${TEST_SUPERVISOR-}" \
      JARVIS_BRIDGE_LAUNCHCTL="$FAKECTL" \
      JARVIS_BRIDGE_LAUNCHD_LABEL="com.jarvis.test" \
      JARVIS_BRIDGE_LAUNCHD_DOMAIN="gui/4242" \
      JARVIS_BRIDGE_LAUNCHD_PLIST="$LPLIST" \
      FAKE_LAUNCHCTL_STATE="$FAKECTL_STATE" \
      FAKE_LAUNCHCTL_LOG="$STATE/bot.log" \
      HOME="$FAKE_HOME" \
      PATH="$FAKE_BREW/bin:$FAKE_BREW/sbin:/usr/bin:/bin:/usr/sbin:/sbin" \
      DINGTALK_APP_KEY="${TEST_KEY-}" \
      DINGTALK_APP_SECRET="${TEST_SECRET-}" \
      JARVIS_NO_DINGTALK="" \
      FAKE_BOT_MODE="${TEST_BOT_MODE:-stay}" \
      FAKE_VALIDATE_MODE="${TEST_VALIDATE_MODE:-ok}" \
      FAKE_SCHEDULER_MODE="${TEST_SCHEDULER_MODE:-ready}" \
      FAKE_PERSISTENT_WORKER_MODE="${TEST_PERSISTENT_WORKER_MODE:-ready}" \
      FAKE_PROCESS_STATE="$STATE" \
      bash "$RUNSH" "$@"
}

pidval() { cat "$STATE/bot.pid" 2>/dev/null || echo ""; }
state_pid() { cat "$STATE/$1" 2>/dev/null || echo ""; }
kill_test_processes() {
  local name p
  # Only pidfiles prove current ownership.  Diagnostic "*-started" markers may
  # outlive a clean stop and their numeric pid can later be reused by this test.
  for name in bot.pid scheduler.pid persistent-worker.pid dingtalk-bot.pid; do
    p="$(state_pid "$name")"
    [ -n "$p" ] && kill -9 "$p" 2>/dev/null || true
  done
}
fresh() {
  kill_test_processes
  rm -rf "$STATE"; mkdir -p "$STATE"
}
cleanup() {
  local launchd_pid
  kill_test_processes
  launchd_pid="$(cat "$FAKECTL_STATE/pid" 2>/dev/null || true)"
  [ -n "$launchd_pid" ] && kill -9 "$launchd_pid" 2>/dev/null || true
  rm -rf "$TMP"
}
trap cleanup EXIT

# --- T0: syntax ------------------------------------------------------------
if bash -n "$RUNSH" 2>/dev/null; then ok "bash -n: run.sh syntax clean"; else no "bash -n: run.sh syntax clean"; fi
if bash -n "$INSTALLER" 2>/dev/null; then ok "bash -n: install-launchd.sh syntax clean"; else no "bash -n: install-launchd.sh syntax clean"; fi

# --- T1: degraded start (no creds) -----------------------------------------
fresh
out="$(TEST_KEY= TEST_SECRET= run start 2>&1)"; rc=$?
has "降级" "$out" "start(degraded): prints prominent degraded hint"
[ "$rc" = 0 ] && ok "start(degraded): exit 0" || no "start(degraded): exit 0 (got $rc)"
[ -f "$STATE/bot.pid" ] && ok "start(degraded): pidfile written" || no "start(degraded): pidfile written"
p1="$(pidval)"; kill -0 "$p1" 2>/dev/null && ok "start(degraded): bot alive after start-check" || no "start(degraded): bot alive"
st="$(run status 2>&1)"
has "RUNNING" "$st" "status: reports RUNNING"
has "降级" "$st" "status: mode read from log banner = degraded"
run stop >/dev/null 2>&1

# --- T2: full mode via creds in env ----------------------------------------
fresh
out="$(TEST_KEY=k TEST_SECRET=s run start 2>&1)"; rc=$?
has "全功能" "$out" "start(full/env): reports full mode"
hasnot "降级" "$out" "start(full/env): no degraded hint when creds present"
[ "$rc" = 0 ] && ok "start(full/env): exit 0" || no "start(full/env): exit 0 (got $rc)"
st="$(run status 2>&1)"
has "全功能" "$st" "status(full): mode = full"
run stop >/dev/null 2>&1

# --- T2b: scheduler role owns both processes through the same run.sh -------
fresh
out="$(TEST_ROLE=scheduler run start 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "start(scheduler): exit 0" || no "start(scheduler): exit 0 (got $rc)"
[ -f "$STATE/bot.pid" ] && ok "start(scheduler): supervisor pidfile written" || no "start(scheduler): supervisor pidfile"
[ ! -f "$STATE/scheduler.pid" ] && ok "run.sh leaves component pidfiles to main" || no "run.sh unexpectedly supervises scheduler child"
st="$(TEST_ROLE=scheduler run status 2>&1)"
has "bridge: RUNNING" "$st" "status(scheduler): reports one supervisor"
TEST_ROLE=scheduler run stop >/dev/null 2>&1
[ ! -f "$STATE/bot.pid" ] && ok "stop(scheduler): supervisor pidfile removed" || no "stop(scheduler): supervisor pidfile removed"

# --- T3: full mode via sourcing bridge/jarvis.env (creds NOT in env) --------
fresh
printf 'DINGTALK_APP_KEY=fromfile\nDINGTALK_APP_SECRET=fromfile\n' >"$JENV"
out="$(TEST_KEY= TEST_SECRET= run start 2>&1)"
has "全功能" "$out" "start: sources creds from jarvis.env → full mode"
run stop >/dev/null 2>&1
: >"$JENV"

# --- T4: start idempotency --------------------------------------------------
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
out="$(run start 2>&1)"; rc=$?
has "已在运行" "$out" "start(idempotent): reports already running"
[ "$rc" = 0 ] && ok "start(idempotent): exit 0" || no "start(idempotent): exit 0 (got $rc)"
p2="$(pidval)"
[ -n "$p1" ] && [ "$p1" = "$p2" ] && ok "start(idempotent): pid unchanged" || no "start(idempotent): pid unchanged ($p1 vs $p2)"
run stop >/dev/null 2>&1

# --- T5: stop cleanup + repeat stop ----------------------------------------
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
out="$(run stop 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "stop: exit 0" || no "stop: exit 0 (got $rc)"
[ ! -f "$STATE/bot.pid" ] && ok "stop: pidfile removed" || no "stop: pidfile removed"
kill -0 "$p1" 2>/dev/null && no "stop: process terminated" || ok "stop: process terminated"
out2="$(run stop 2>&1)"; rc2=$?
[ "$rc2" = 0 ] && ok "stop(no pid): clean exit 0" || no "stop(no pid): exit 0 (got $rc2)"
has "未在运行" "$out2" "stop(no pid): reports not running"

# --- T6: supervised dispatch mode starts and waits for authorization --------
fresh
out="$(JARVIS_AUTO_DISPATCH=0 run start 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "start(auto-dispatch=0): supervised mode starts" || no "start(auto-dispatch=0): exit 0 (got $rc)"
[ -f "$STATE/bot.pid" ] && ok "start(auto-dispatch=0): supervisor pidfile written" || no "start(auto-dispatch=0): supervisor pidfile"
st="$(JARVIS_AUTO_DISPATCH=0 run status 2>&1)"
has "bridge: RUNNING" "$st" "start(auto-dispatch=0): waits for authorization while running"
JARVIS_AUTO_DISPATCH=0 run stop >/dev/null 2>&1

# --- T8: dry-run passes --dry-run-once through -----------------------------
fresh
out="$(run dry-run 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "dry-run: exit 0" || no "dry-run: exit 0 (got $rc)"
has "dry-run ok" "$out" "dry-run: --dry-run-once forwarded to bot"
has "PATH=$FAKE_HOME/.local/bin:$FAKE_BREW/bin:$FAKE_BREW/sbin:/usr/bin:/bin:/usr/sbin:/sbin" "$out" \
  "dry-run: non-interactive PATH includes user-local and Homebrew tool directories"
has "JARVIS_ROOT=$repo_root" "$out" \
  "dry-run: review worktree stays on its own wrappers and configuration"

# --- T8b: supervisor READY failure is rolled back before service ownership --
fresh
out="$(TEST_SCHEDULER_MODE=no-ready TEST_READY_WAIT=1 run start 2>&1)"; rc=$?
scheduler_started="$(state_pid scheduler-started)"
[ "$rc" != 0 ] && ok "start(scheduler READY timeout): non-zero exit" || no "start(scheduler READY timeout): non-zero exit"
[ ! -f "$STATE/bot.pid" ] && ok "start(supervisor READY timeout): pidfile removed" || no "start(supervisor READY timeout): pidfile remains"
[ -n "$scheduler_started" ] && ! kill -0 "$scheduler_started" 2>/dev/null \
  && ok "start(supervisor READY timeout): child terminated" \
  || no "start(supervisor READY timeout): child terminated"

# --- T9: status when stopped -----------------------------------------------
fresh
out="$(run status 2>&1)"
has "STOPPED" "$out" "status(stopped): reports STOPPED"

# --- T10: logs is non-blocking (must return) -------------------------------
fresh
run start >/dev/null 2>&1
out="$(run logs 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "logs: returns without blocking (no tail -f)" || no "logs: non-blocking (got $rc)"
has "$STATE/bot.log" "$out" "logs: prints log path"
run stop >/dev/null 2>&1

# --- T11: restart swaps the pid --------------------------------------------
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
out="$(run restart 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "restart: exit 0" || no "restart: exit 0 (got $rc)"
p2="$(pidval)"
[ -n "$p2" ] && [ "$p1" != "$p2" ] && ok "restart: new pid" || no "restart: new pid ($p1 -> $p2)"
kill -0 "$p1" 2>/dev/null && no "restart: old process gone" || ok "restart: old process gone"
run stop >/dev/null 2>&1

# --- T11c: watchdog uses three independent loop beacons --------------------
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
printf '%s\n' "$p1" >"$STATE/persistent-worker.pid"
now="$(date -u +%s)"
for kind in worker lease session; do
  printf '%s\n' "$now" >"$STATE/heartbeat.persistent-worker.$kind.epoch"
done
out="$(TEST_WATCHDOG_STARTUP_GRACE=0 run watchdog 2>&1)"; rc=$?
p2="$(pidval)"
[ "$rc" = 0 ] && [ "$p1" = "$p2" ] \
  && ok "watchdog: fresh worker/lease/session beacons keep supervisor" \
  || no "watchdog: fresh beacons unexpectedly restarted supervisor ($p1 -> $p2, rc=$rc, out=$out)"

# Missing/invalid beacons are healthy only during startup grace.
rm -f "$STATE/heartbeat.persistent-worker.lease.epoch"
out="$(TEST_WATCHDOG_STARTUP_GRACE=3600 run watchdog 2>&1)"; rc=$?
[ "$rc" = 0 ] && [ "$(pidval)" = "$p1" ] \
  && ok "watchdog: missing beacon tolerated during startup grace" \
  || no "watchdog: missing startup beacon was not tolerated"
printf '%s\n' "invalid" >"$STATE/heartbeat.persistent-worker.lease.epoch"
out="$(TEST_WATCHDOG_STARTUP_GRACE=3600 run watchdog 2>&1)"; rc=$?
[ "$rc" = 0 ] && [ "$(pidval)" = "$p1" ] \
  && ok "watchdog: invalid beacon tolerated during startup grace" \
  || no "watchdog: invalid startup beacon was not tolerated"
run stop >/dev/null 2>&1

# Outside grace, one unhealthy loop forces the supervisor and every private
# component process group down before starting the replacement.
fresh
TEST_SCHEDULER_MODE=cleanup-child-pidfile run start >/dev/null 2>&1
p1="$(pidval)"
grandchild_pidfile="$STATE/stuck-grandchild.pid"
/usr/bin/python3 -c '
import os
import signal
import sys
import time

os.setsid()
child = os.fork()
if child:
    raise SystemExit(0)
signal.signal(signal.SIGTERM, signal.SIG_IGN)
with open(sys.argv[1], "w", encoding="utf-8") as stream:
    stream.write(str(os.getpid()))
time.sleep(300)
' "$grandchild_pidfile" &
stuck_group=$!
wait "$stuck_group"
for _attempt in 1 2 3 4 5 6 7 8 9 10; do
  [ -s "$grandchild_pidfile" ] && break
  sleep 0.1
done
stuck_child="$(cat "$grandchild_pidfile")"
printf '%s\n' "$stuck_group" >"$STATE/persistent-worker.pid"
printf '%s|%s|dead-original-leader\n' \
  "$stuck_group" "$stuck_group" >"$STATE/persistent-worker.pid.identity"
now="$(date -u +%s)"
printf '%s\n' "$now" >"$STATE/heartbeat.persistent-worker.worker.epoch"
printf '%s\n' "$now" >"$STATE/heartbeat.persistent-worker.lease.epoch"
printf '%s\n' "invalid" >"$STATE/heartbeat.persistent-worker.session.epoch"
lock_root="$FAKE_HOME/.config/a1/identities/jarvis"
mkdir -p "$lock_root"
touch "$lock_root/auth.yaml.lock"
out="$(TEST_WATCHDOG_STARTUP_GRACE=0 run watchdog 2>&1)"; rc=$?
p2="$(pidval)"
[ "$rc" = 0 ] && [ -n "$p2" ] && [ "$p1" != "$p2" ] \
  && ok "watchdog: overdue session loop restarts supervisor" \
  || no "watchdog: overdue session loop failed to restart ($p1 -> $p2, rc=$rc)"
kill -0 "$stuck_child" 2>/dev/null \
  && no "watchdog: grandchild survives after its group leader exited" \
  || ok "watchdog: pre-TERM child PGID snapshot survives supervisor pidfile cleanup"
[ ! -e "$lock_root/auth.yaml.lock" ] \
  && ok "watchdog: restart clears stale a1 lock" \
  || no "watchdog: restart left stale a1 lock"
has "session beacon missing/invalid" "$out" \
  "watchdog: unhealthy loop reason identifies session beacon"
run stop >/dev/null 2>&1

# A stale pidfile whose numeric PID has been reused must fail closed before the
# supervisor or the unrelated private group is signalled.
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
unrelated_pidfile="$STATE/unrelated.pid"
/usr/bin/python3 -c '
import os
import pathlib
import signal
import sys
import time

os.setsid()
signal.signal(signal.SIGTERM, signal.SIG_IGN)
pathlib.Path(sys.argv[1]).write_text(str(os.getpid()))
time.sleep(300)
' "$unrelated_pidfile" &
for _attempt in 1 2 3 4 5 6 7 8 9 10; do
  [ -s "$unrelated_pidfile" ] && break
  sleep 0.1
done
unrelated_pid="$(cat "$unrelated_pidfile")"
printf '%s\n' "$unrelated_pid" >"$STATE/scheduler.pid"
printf '%s|%s|definitely-not-current-birth\n' \
  "$unrelated_pid" "$unrelated_pid" >"$STATE/scheduler.pid.identity"
out="$(TEST_WATCHDOG_STARTUP_GRACE=0 run watchdog 2>&1)"; rc=$?
[ "$rc" != 0 ] && [ "$(pidval)" = "$p1" ] \
  && ok "watchdog: reused child PGID fails closed before supervisor restart" \
  || no "watchdog: reused child PGID did not fail closed (rc=$rc, out=$out)"
kill -0 "$unrelated_pid" 2>/dev/null \
  && ok "watchdog: reused child PGID is not signalled" \
  || no "watchdog: reused child PGID was killed"
has "refuses unverified component pidfile" "$out" \
  "watchdog: reused child PGID reports ownership failure"
rm -f "$STATE/scheduler.pid" "$STATE/scheduler.pid.identity"
kill -KILL -- "-$unrelated_pid" 2>/dev/null || true
run stop >/dev/null 2>&1

# --- T11b: replacement validation is a pre-downtime gate -------------------
fresh
run start >/dev/null 2>&1
p1="$(pidval)"
out="$(TEST_VALIDATE_MODE=fail run restart 2>&1)"; rc=$?
p2="$(pidval)"
[ "$rc" != 0 ] && ok "restart(validation failure): non-zero exit" || no "restart(validation failure): exit 0"
[ -n "$p1" ] && [ "$p1" = "$p2" ] \
  && ok "restart(validation failure): keeps current supervisor pid" \
  || no "restart(validation failure): supervisor changed ($p1 -> $p2)"
kill -0 "$p1" 2>/dev/null \
  && ok "restart(validation failure): current supervisor stays alive" \
  || no "restart(validation failure): current supervisor stopped"
run stop >/dev/null 2>&1

# --- T12: graceful stop delegates child shutdown to Python supervisor -------
fresh
run start >/dev/null 2>&1
gout="$(run stop 2>&1)"
has "graceful" "$gout" "stop(normal): bot 收 SIGTERM 即退 → graceful"

# --- T13: daemon is a true foreground entrypoint ---------------------------
fresh
out="$(TEST_ROLE=worker TEST_SCHEDULER_MODE=once TEST_KEY=k TEST_SECRET=s run daemon 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "daemon: foreground child exit is propagated" || no "daemon: foreground child exit is propagated (got $rc)"
has "foreground daemon" "$out" "daemon: reports foreground mode"
has "JARVIS_NO_DINGTALK=1" "$out" "daemon(worker): forces degraded mode"
[ ! -f "$STATE/bot.pid" ] && ok "daemon: does not write local pidfile" || no "daemon: does not write local pidfile"

# --- T14: launchd-supervised lifecycle uses only fake launchctl ------------
fresh
rm -rf "$FAKECTL_STATE"; mkdir -p "$FAKECTL_STATE"
out="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run start 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "launchd start: exit 0" || no "launchd start: exit 0 (got $rc)"
calls="$(cat "$FAKECTL_STATE/calls" 2>/dev/null)"
has "bootstrap gui/4242 $LPLIST" "$calls" "launchd start: bootstraps configured plist"
has "kickstart gui/4242/com.jarvis.test" "$calls" "launchd start: kickstarts service"
[ ! -f "$STATE/bot.pid" ] && ok "launchd start: no local bot process/pidfile" || no "launchd start: no local pidfile"
st="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run status 2>&1)"
has "RUNNING" "$st" "launchd status: parses fake launchctl state"
launchd_pid="$(cat "$FAKECTL_STATE/pid")"
has "pid $launchd_pid" "$st" "launchd status: reports launchd pid"
: >"$FAKECTL_STATE/calls"
out="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run restart 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "launchd restart: exit 0" || no "launchd restart: exit 0 (got $rc)"
calls="$(cat "$FAKECTL_STATE/calls")"
has "enable gui/4242/com.jarvis.test" "$calls" "launchd restart: keeps service enabled"
has "disable gui/4242/com.jarvis.test" "$calls" "launchd restart: disables KeepAlive during drain"
has "kill SIGTERM gui/4242/com.jarvis.test" "$calls" "launchd restart: requests graceful drain"
has "kickstart gui/4242/com.jarvis.test" "$calls" "launchd restart: starts replacement without overlap"
hasnot "kickstart -k gui/4242/com.jarvis.test" "$calls" "launchd restart: avoids overlapping replacement"
hasnot "bootout gui/4242/com.jarvis.test" "$calls" "launchd restart: keeps service registered"

# A broken replacement runtime is rejected before launchd is disabled or the
# current process is signalled.
launchd_pid="$(cat "$FAKECTL_STATE/pid")"
: >"$FAKECTL_STATE/calls"
out="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd TEST_VALIDATE_MODE=fail run restart 2>&1)"; rc=$?
[ "$rc" != 0 ] && ok "launchd restart(validation failure): non-zero exit" || no "launchd restart(validation failure): exit 0"
calls="$(cat "$FAKECTL_STATE/calls")"
hasnot "disable gui/4242/com.jarvis.test" "$calls" "launchd restart(validation failure): leaves KeepAlive enabled"
hasnot "kill SIGTERM gui/4242/com.jarvis.test" "$calls" "launchd restart(validation failure): leaves current process running"
[ "$(cat "$FAKECTL_STATE/pid")" = "$launchd_pid" ] && kill -0 "$launchd_pid" 2>/dev/null \
  && ok "launchd restart(validation failure): preserves current pid" \
  || no "launchd restart(validation failure): current pid changed"

out="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run stop 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "launchd stop: exit 0" || no "launchd stop: exit 0 (got $rc)"
st="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run status 2>&1)"
has "STOPPED" "$st" "launchd status: reports stopped after bootout"

# --- T15: installer renders absolute plist and converges idempotently -------
rm -rf "$FAKECTL_STATE" "$TMP/home"; mkdir -p "$FAKECTL_STATE" "$TMP/home"
INSTALLED_PLIST="$TMP/home/Library/LaunchAgents/com.jarvis.test.plist"
install_fake() {
  env HOME="$TMP/home" \
      JARVIS_LAUNCHD_ALLOW_NON_DARWIN=1 \
      JARVIS_BRIDGE_LAUNCHCTL="$FAKECTL" \
      JARVIS_BRIDGE_LAUNCHD_LABEL="com.jarvis.test" \
      JARVIS_BRIDGE_LAUNCHD_DOMAIN="gui/4242" \
      JARVIS_BRIDGE_LAUNCHD_PLIST="$INSTALLED_PLIST" \
      JARVIS_BRIDGE_STATE_DIR="$TMP/install-state" \
      JARVIS_BRIDGE_BOOTSTRAP_ENV="$BSENV" \
      JARVIS_BRIDGE_ENV="$JENV" \
      JARVIS_BRIDGE_PYTHON="/usr/bin/python3" \
      JARVIS_BRIDGE_ROLE="${INSTALL_ROLE:-worker}" \
      JARVIS_BOOT_ID="test-boot" \
      JARVIS_CONTROL_PLANE_TOKEN="test-token" \
      JARVIS_AUTO_DISPATCH="${INSTALL_AUTO_DISPATCH:-1}" \
      JARVIS_SCHEDULER_READY_WAIT="${INSTALL_READY_WAIT:-1}" \
      JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS="1" \
      PYTHONPATH="$repo_root:$PY_YAML_SITE" \
      FAKE_LAUNCHCTL_STATE="$FAKECTL_STATE" \
      FAKE_LAUNCHCTL_LOG="$TMP/install-state/bot.log" \
      bash "$INSTALLER"
}
iout="$(install_fake 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "install-launchd: first install exits 0" || no "install-launchd: first install exits 0 (got $rc)"
[ -f "$INSTALLED_PLIST" ] && ok "install-launchd: plist created" || no "install-launchd: plist created"
rendered="$(cat "$INSTALLED_PLIST" 2>/dev/null)"
has "<string>$RUNSH</string>" "$rendered" "install-launchd: renders absolute run.sh path"
has "<string>daemon</string>" "$rendered" "install-launchd: plist invokes daemon"
has "<key>Umask</key>" "$rendered" "install-launchd: launchd child uses private umask"
has "<integer>63</integer>" "$rendered" "install-launchd: private umask is 0077"
hasnot "__RUN_SH_PATH__" "$rendered" "install-launchd: no unresolved path placeholder"
state_mode="$(python3 - "$TMP/install-state" <<'PY'
import os
import stat
import sys
print(oct(stat.S_IMODE(os.stat(sys.argv[1]).st_mode)))
PY
)"
[ "$state_mode" = "0o700" ] && ok "install-launchd: state directory is 0700" || no "install-launchd: state directory is 0700 (got $state_mode)"
calls="$(cat "$FAKECTL_STATE/calls")"
has "bootstrap gui/4242 $INSTALLED_PLIST" "$calls" "install-launchd: bootstraps rendered plist"
has "kickstart gui/4242/com.jarvis.test" "$calls" "install-launchd: kickstarts rendered service without overlap"
iout2="$(install_fake 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "install-launchd: repeat install exits 0" || no "install-launchd: repeat install exits 0 (got $rc)"
has "unchanged" "$iout2" "install-launchd: repeat render is idempotent"
calls="$(cat "$FAKECTL_STATE/calls")"
has "disable gui/4242/com.jarvis.test" "$calls" "install-launchd: loaded upgrade uses safe drain"
has "kill SIGTERM gui/4242/com.jarvis.test" "$calls" "install-launchd: loaded upgrade signals current main"
hasnot "bootout gui/4242/com.jarvis.test" "$calls" "install-launchd: loaded upgrade never bootouts adopted worker"

: >"$FAKECTL_STATE/calls"
schedule_install_out="$(
  INSTALL_ROLE=scheduler INSTALL_AUTO_DISPATCH=0 install_fake 2>&1
)"
schedule_install_rc=$?
if [ "$schedule_install_rc" -eq 0 ]; then
  ok "install-launchd: supervised scheduler config is supported"
else
  no "install-launchd: supervised scheduler config is supported [$schedule_install_out]"
fi
calls="$(cat "$FAKECTL_STATE/calls" 2>/dev/null)"
has "disable gui/4242/com.jarvis.test" "$calls" "install-launchd: supervised upgrade uses safe drain"
has "kill SIGTERM gui/4242/com.jarvis.test" "$calls" "install-launchd: supervised upgrade signals current main"

echo ""
echo "bridge_run_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ] && { echo "ALLPASS"; exit 0; } || { echo "FAILED"; exit 1; }
