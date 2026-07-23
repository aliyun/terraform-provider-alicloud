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
for a in "$@"; do
  if [ "${a##*/}" = "main.py" ]; then
    [ -n "${FAKE_PROCESS_STATE:-}" ] && printf '%s\n' "$$" >"$FAKE_PROCESS_STATE/scheduler-started"
    case "${FAKE_SCHEDULER_MODE:-ready}" in
      exit)
        echo "$ts ERROR [MainThread] fake scheduler immediate exit"
        exit 2 ;;
      no-ready)
        echo "$ts INFO [MainThread] fake scheduler deliberately not READY"
        exec sleep 300 ;;
      *)
        echo "$ts INFO [MainThread] Scheduler READY pid=$$ worker=bridge-scheduler jobs=fake"
        exec sleep 300 ;;
    esac
  fi
  if [ "${a##*/}" = "task_worker.py" ]; then
    [ -n "${FAKE_PROCESS_STATE:-}" ] && printf '%s\n' "$$" >"$FAKE_PROCESS_STATE/task-worker-started"
    case "${FAKE_TASK_WORKER_MODE:-ready}" in
      exit)
        echo "$ts ERROR [MainThread] fake task worker immediate exit"
        exit 2 ;;
      no-ready)
        echo "$ts INFO [MainThread] fake task worker deliberately not READY"
        exec sleep 300 ;;
      once)
        echo "$ts INFO [MainThread] foreground once JARVIS_NO_DINGTALK=${JARVIS_NO_DINGTALK:-0}"
        exit 0 ;;
      *)
        echo "$ts INFO [MainThread] Task worker READY pid=$$ worker=fake"
        exec sleep 300 ;;
    esac
  fi
done
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
      printf '%s INFO [MainThread] Bridge READY pid=%s mode=fake\n' \
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
      FAKE_SCHEDULER_MODE="${TEST_SCHEDULER_MODE:-ready}" \
      FAKE_TASK_WORKER_MODE="${TEST_TASK_WORKER_MODE:-ready}" \
      FAKE_PROCESS_STATE="$STATE" \
      bash "$RUNSH" "$@"
}

pidval() { cat "$STATE/bot.pid" 2>/dev/null || echo ""; }
state_pid() { cat "$STATE/$1" 2>/dev/null || echo ""; }
kill_test_processes() {
  local name p
  # Only pidfiles prove current ownership.  Diagnostic "*-started" markers may
  # outlive a clean stop and their numeric pid can later be reused by this test.
  for name in bot.pid scheduler.pid task-worker.pid; do
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
[ -f "$STATE/bot.pid" ] && ok "start(scheduler): executor pidfile written" || no "start(scheduler): executor pidfile"
[ -f "$STATE/scheduler.pid" ] && ok "start(scheduler): scheduler pidfile written" || no "start(scheduler): scheduler pidfile"
st="$(TEST_ROLE=scheduler run status 2>&1)"
has "scheduler: RUNNING" "$st" "status(scheduler): reports Scheduler through run.sh"
TEST_ROLE=scheduler run stop >/dev/null 2>&1
[ ! -f "$STATE/scheduler.pid" ] && ok "stop(scheduler): scheduler pidfile removed" || no "stop(scheduler): scheduler pidfile removed"
[ ! -f "$STATE/task-worker.pid" ] && ok "stop(scheduler): task worker pidfile removed" || no "stop(scheduler): task worker pidfile removed"

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

# --- T6: startup failure — bot exits immediately ---------------------------
fresh
out="$(TEST_BOT_MODE=exit run start 2>&1)"; rc=$?
[ "$rc" != 0 ] && ok "start(fail/exit): non-zero exit" || no "start(fail/exit): non-zero exit (got $rc)"
has "失败" "$out" "start(fail/exit): reports failure + log tail"
[ ! -f "$STATE/bot.pid" ] && ok "start(fail/exit): stale pidfile cleaned" || no "start(fail/exit): stale pidfile cleaned"

# --- T7: startup failure — first log line is ERROR though process stays ----
fresh
out="$(TEST_BOT_MODE=errfirst run start 2>&1)"; rc=$?
[ "$rc" != 0 ] && ok "start(fail/errfirst): non-zero exit" || no "start(fail/errfirst): non-zero exit (got $rc)"
[ ! -f "$STATE/bot.pid" ] && ok "start(fail/errfirst): pidfile cleaned (proc stopped)" || no "start(fail/errfirst): pidfile cleaned"

# --- T8: dry-run passes --dry-run-once through -----------------------------
fresh
out="$(run dry-run 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "dry-run: exit 0" || no "dry-run: exit 0 (got $rc)"
has "dry-run ok" "$out" "dry-run: --dry-run-once forwarded to bot"
has "PATH=$FAKE_HOME/.local/bin:$FAKE_BREW/bin:$FAKE_BREW/sbin:/usr/bin:/bin:/usr/sbin:/sbin" "$out" \
  "dry-run: non-interactive PATH includes user-local and Homebrew tool directories"
has "JARVIS_ROOT=$repo_root" "$out" \
  "dry-run: review worktree stays on its own wrappers and configuration"

# --- T8b: standalone runtime READY failures are fully rolled back -----------
fresh
out="$(TEST_TASK_WORKER_MODE=no-ready TEST_READY_WAIT=1 run start 2>&1)"; rc=$?
worker_started="$(state_pid task-worker-started)"
[ "$rc" != 0 ] && ok "start(task worker READY timeout): non-zero exit" || no "start(task worker READY timeout): non-zero exit"
[ ! -f "$STATE/task-worker.pid" ] && ok "start(task worker READY timeout): pidfile removed" || no "start(task worker READY timeout): pidfile removed"
[ -n "$worker_started" ] && ! kill -0 "$worker_started" 2>/dev/null \
  && ok "start(task worker READY timeout): child terminated" \
  || no "start(task worker READY timeout): child terminated"
[ ! -f "$STATE/bot.pid" ] && ok "start(task worker READY timeout): bot rolled back" || no "start(task worker READY timeout): bot rolled back"

fresh
out="$(TEST_SCHEDULER_MODE=no-ready TEST_READY_WAIT=1 run start 2>&1)"; rc=$?
scheduler_started="$(state_pid scheduler-started)"
worker_started="$(state_pid task-worker-started)"
[ "$rc" != 0 ] && ok "start(scheduler READY timeout): non-zero exit" || no "start(scheduler READY timeout): non-zero exit"
[ ! -f "$STATE/scheduler.pid" ] && ok "start(scheduler READY timeout): scheduler pidfile removed" || no "start(scheduler READY timeout): scheduler pidfile removed"
[ -n "$scheduler_started" ] && ! kill -0 "$scheduler_started" 2>/dev/null \
  && ok "start(scheduler READY timeout): scheduler child terminated" \
  || no "start(scheduler READY timeout): scheduler child terminated"
[ ! -f "$STATE/task-worker.pid" ] && ok "start(scheduler READY timeout): new worker pidfile removed" || no "start(scheduler READY timeout): new worker pidfile removed"
[ -n "$worker_started" ] && ! kill -0 "$worker_started" 2>/dev/null \
  && ok "start(scheduler READY timeout): new worker terminated" \
  || no "start(scheduler READY timeout): new worker terminated"
[ ! -f "$STATE/bot.pid" ] && ok "start(scheduler READY timeout): bot rolled back" || no "start(scheduler READY timeout): bot rolled back"

# A worker that predates this start attempt may own leased Sessions and must be
# retained when the newly-created scheduler fails.
fresh
TEST_ROLE=worker run start >/dev/null 2>&1
worker_before="$(state_pid task-worker.pid)"
out="$(TEST_ROLE=scheduler TEST_SCHEDULER_MODE=no-ready TEST_READY_WAIT=1 run start 2>&1)"; rc=$?
worker_after="$(state_pid task-worker.pid)"
[ "$rc" != 0 ] && ok "start(scheduler timeout/pre-existing worker): non-zero exit" || no "start(scheduler timeout/pre-existing worker): non-zero exit"
[ -n "$worker_before" ] && [ "$worker_before" = "$worker_after" ] && kill -0 "$worker_after" 2>/dev/null \
  && ok "start(scheduler timeout/pre-existing worker): existing worker preserved" \
  || no "start(scheduler timeout/pre-existing worker): existing worker preserved"
TEST_ROLE=worker run stop >/dev/null 2>&1

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
worker_before="$(state_pid task-worker.pid)"
out="$(run restart 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "restart: exit 0" || no "restart: exit 0 (got $rc)"
p2="$(pidval)"
[ -n "$p2" ] && [ "$p1" != "$p2" ] && ok "restart: new pid" || no "restart: new pid ($p1 -> $p2)"
kill -0 "$p1" 2>/dev/null && no "restart: old process gone" || ok "restart: old process gone"
worker_after="$(state_pid task-worker.pid)"
[ -n "$worker_before" ] && [ "$worker_before" = "$worker_after" ] && kill -0 "$worker_after" 2>/dev/null \
  && ok "restart(scheduler): existing task worker remains alive" \
  || no "restart(scheduler): existing task worker remains alive"
run stop >/dev/null 2>&1

# --- T12: graceful stop semantics (吸收 master f7f1f72 优雅停止) --------------
# 正常 bot 收 SIGTERM 即退 → graceful; deaf bot 忽略 TERM → 宽限超时后 SIGKILL forced 兜底。
fresh
run start >/dev/null 2>&1
gout="$(run stop 2>&1)"
has "graceful" "$gout" "stop(normal): bot 收 SIGTERM 即退 → graceful"
fresh
TEST_BOT_MODE=deaf run start >/dev/null 2>&1
p1="$(pidval)"
fout="$(run stop 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "stop(deaf): exit 0" || no "stop(deaf): exit 0 (got $rc)"
has "forced" "$fout" "stop(deaf): 忽略 TERM → 宽限后 SIGKILL forced 兜底"
kill -0 "$p1" 2>/dev/null && no "stop(deaf): process gone after SIGKILL" || ok "stop(deaf): process gone after SIGKILL"
[ ! -f "$STATE/bot.pid" ] && ok "stop(deaf): pidfile removed" || no "stop(deaf): pidfile removed"

# --- T13: daemon is a true foreground entrypoint ---------------------------
fresh
out="$(TEST_ROLE=worker TEST_TASK_WORKER_MODE=once TEST_KEY=k TEST_SECRET=s run daemon 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "daemon: foreground child exit is propagated" || no "daemon: foreground child exit is propagated (got $rc)"
has "foreground daemon" "$out" "daemon: reports foreground mode"
has "JARVIS_NO_DINGTALK=1" "$out" "daemon(worker): forces degraded mode"
[ ! -f "$STATE/bot.pid" ] && ok "daemon: does not write local pidfile" || no "daemon: does not write local pidfile"

# Scheduler-role daemon owns all three components.  TERM must stop the
# separately-nohup'd worker as well as scheduler and bot.
fresh
(
  export JARVIS_BRIDGE_PYTHON="$FAKEPY"
  export JARVIS_BRIDGE_BOT="$BOTFILE"
  export JARVIS_BRIDGE_STATE_DIR="$STATE"
  export JARVIS_BRIDGE_BOOTSTRAP_ENV="$BSENV"
  export JARVIS_BRIDGE_ENV="$JENV"
  export JARVIS_BRIDGE_ROLE="scheduler"
  export JARVIS_CONTROL_PLANE_TOKEN="test-token"
  export JARVIS_BRIDGE_STOP_WAIT="1"
  export JARVIS_BRIDGE_START_ROLLBACK_WAIT="1"
  export JARVIS_SCHEDULER_READY_WAIT="1"
  export JARVIS_BRIDGE_TOOL_DIRS="$FAKE_BREW/sbin $FAKE_BREW/bin $FAKE_HOME/.local/bin"
  export HOME="$FAKE_HOME"
  export PATH="$FAKE_BREW/bin:$FAKE_BREW/sbin:/usr/bin:/bin:/usr/sbin:/sbin"
  export FAKE_BOT_MODE="stay"
  export FAKE_SCHEDULER_MODE="ready"
  export FAKE_TASK_WORKER_MODE="ready"
  export FAKE_PROCESS_STATE="$STATE"
  exec bash "$RUNSH" daemon
) >"$STATE/daemon.log" 2>&1 &
daemon_pid=$!
i=0
while [ "$i" -lt 30 ] && { [ ! -f "$STATE/task-worker.pid" ] || [ ! -f "$STATE/scheduler.pid" ]; }; do
  sleep 0.1
  i=$((i + 1))
done
daemon_worker="$(state_pid task-worker.pid)"
daemon_scheduler="$(state_pid scheduler.pid)"
kill -TERM "$daemon_pid" 2>/dev/null || true
wait "$daemon_pid" 2>/dev/null; daemon_rc=$?
[ "$daemon_rc" = 0 ] && ok "daemon(scheduler): TERM exits cleanly" || no "daemon(scheduler): TERM exits cleanly (got $daemon_rc)"
[ -n "$daemon_worker" ] && ! kill -0 "$daemon_worker" 2>/dev/null \
  && ok "daemon(scheduler): TERM stops task worker" \
  || no "daemon(scheduler): TERM stops task worker"
[ -n "$daemon_scheduler" ] && ! kill -0 "$daemon_scheduler" 2>/dev/null \
  && ok "daemon(scheduler): TERM stops scheduler" \
  || no "daemon(scheduler): TERM stops scheduler"
[ ! -f "$STATE/task-worker.pid" ] && ok "daemon(scheduler): worker pidfile removed" || no "daemon(scheduler): worker pidfile removed"

# A launchd-controlled restart uses a PID-bound one-shot marker.  The same
# daemon TERM must then replace bot+scheduler without fencing leased Task work.
fresh
(
  export JARVIS_BRIDGE_PYTHON="$FAKEPY"
  export JARVIS_BRIDGE_BOT="$BOTFILE"
  export JARVIS_BRIDGE_STATE_DIR="$STATE"
  export JARVIS_BRIDGE_BOOTSTRAP_ENV="$BSENV"
  export JARVIS_BRIDGE_ENV="$JENV"
  export JARVIS_BRIDGE_ROLE="scheduler"
  export JARVIS_CONTROL_PLANE_TOKEN="test-token"
  export JARVIS_BRIDGE_STOP_WAIT="1"
  export JARVIS_SCHEDULER_READY_WAIT="1"
  export JARVIS_BRIDGE_TOOL_DIRS="$FAKE_BREW/sbin $FAKE_BREW/bin $FAKE_HOME/.local/bin"
  export HOME="$FAKE_HOME"
  export PATH="$FAKE_BREW/bin:$FAKE_BREW/sbin:/usr/bin:/bin:/usr/sbin:/sbin"
  export FAKE_BOT_MODE="stay"
  export FAKE_SCHEDULER_MODE="ready"
  export FAKE_TASK_WORKER_MODE="ready"
  export FAKE_PROCESS_STATE="$STATE"
  exec bash "$RUNSH" daemon
) >"$STATE/daemon-restart.log" 2>&1 &
daemon_pid=$!
i=0
while [ "$i" -lt 30 ] && { [ ! -f "$STATE/task-worker.pid" ] || [ ! -f "$STATE/scheduler.pid" ]; }; do
  sleep 0.1
  i=$((i + 1))
done
daemon_worker="$(state_pid task-worker.pid)"
daemon_scheduler="$(state_pid scheduler.pid)"
printf '%s\n' "$daemon_pid" >"$STATE/preserve-task-worker-once"
kill -TERM "$daemon_pid" 2>/dev/null || true
wait "$daemon_pid" 2>/dev/null; daemon_rc=$?
[ "$daemon_rc" = 0 ] && ok "daemon(restart): TERM exits cleanly" || no "daemon(restart): TERM exits cleanly (got $daemon_rc)"
[ -n "$daemon_worker" ] && kill -0 "$daemon_worker" 2>/dev/null \
  && ok "daemon(restart): preserves task worker" \
  || no "daemon(restart): preserves task worker"
[ -n "$daemon_scheduler" ] && ! kill -0 "$daemon_scheduler" 2>/dev/null \
  && ok "daemon(restart): stops scheduler" \
  || no "daemon(restart): stops scheduler"
[ ! -f "$STATE/preserve-task-worker-once" ] \
  && ok "daemon(restart): consumes one-shot marker" \
  || no "daemon(restart): consumes one-shot marker"
TEST_ROLE=worker run stop >/dev/null 2>&1

# --- T14: launchd-supervised lifecycle uses only fake launchctl ------------
fresh
rm -rf "$FAKECTL_STATE"; mkdir -p "$FAKECTL_STATE"
out="$(TEST_ROLE=worker TEST_SUPERVISOR=launchd run start 2>&1)"; rc=$?
[ "$rc" = 0 ] && ok "launchd start: exit 0" || no "launchd start: exit 0 (got $rc)"
calls="$(cat "$FAKECTL_STATE/calls" 2>/dev/null)"
has "bootstrap gui/4242 $LPLIST" "$calls" "launchd start: bootstraps configured plist"
has "kickstart -k gui/4242/com.jarvis.test" "$calls" "launchd start: kickstarts service"
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
      JARVIS_BRIDGE_ENV="$JENV" \
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

echo ""
echo "bridge_run_test: $pass passed, $fail failed"
[ "$fail" -eq 0 ] && { echo "ALLPASS"; exit 0; } || { echo "FAILED"; exit 1; }
