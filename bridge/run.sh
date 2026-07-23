#!/usr/bin/env bash
# bridge/run.sh — Jarvis DingTalk bridge 单一入口(进程管理器)。
#
#   bridge/run.sh <start|stop|restart|status|logs|dry-run|daemon>
#
# start 是唯一点火入口(仓库主人 2026-07-06:「运行 bridge/run.sh start 即可,不需额外点火」):
#   · 幂等: pidfile 活着则提示已运行并退 0。
#   · 自动 source bootstrap/.env + bridge/jarvis.env(存在才 source, 变量自动 export)。
#   · 模式判定: 有 DINGTALK_APP_KEY/SECRET → 全功能; 缺 → 自动 JARVIS_NO_DINGTALK=1 降级点火
#     (自动派发+各调度器照常, 卡片/播报落 bot.log; 配好凭证后 run.sh restart 即回全功能)。
#   · nohup 守护 + 写 pidfile, 启动后自检(进程仍活 + 日志首行非 ERROR), 失败退非零并打印日志尾。
# stop 优雅停止: 发 SIGTERM → bot 自身 SIGTERM handler 收尾在跑 worker 并 release
#   其已认领工单；scheduler 角色会先 drain 已准入 Job，再停止执行器；drain 超时
#   则拒绝替换，避免打断已准入的周期任务。
#
# 依赖 F 线 JARVIS_NO_DINGTALK(bridge 降级模式, 分支 worktree-f3-nodingtalk / MR-4)。若该能力
# 尚未合并, 缺凭证时旧 bot 会忽略 flag、因缺凭证退 2 —— 本脚本的启动自检会如实报错并提示先合 MR-4。
#
# 可覆盖(测试/部署): JARVIS_BRIDGE_PYTHON(默认 python3) JARVIS_BRIDGE_BOT(默认本目录 bot)
#   JARVIS_BRIDGE_STATE_DIR(默认 <repo>/.my-day/bridge) JARVIS_BRIDGE_BOOTSTRAP_ENV
#   JARVIS_BRIDGE_ENV JARVIS_BRIDGE_START_WAIT(默认 2s)
#   JARVIS_BRIDGE_STOP_WAIT / JARVIS_STOP_GRACE (默认 30s)。
#   JARVIS_BRIDGE_SUPERVISOR=launchd 时 start/stop/restart/status 委托给 launchctl；可覆盖
#   JARVIS_BRIDGE_LAUNCHCTL、JARVIS_BRIDGE_LAUNCHD_LABEL/DOMAIN/PLIST（测试/定制安装）。
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

if [ -n "${JARVIS_BRIDGE_PYTHON:-}" ]; then
  PYTHON="$JARVIS_BRIDGE_PYTHON"
elif [ -x "$REPO_ROOT/.venv/bridge/bin/python" ]; then
  PYTHON="$REPO_ROOT/.venv/bridge/bin/python"
else
  PYTHON="python3"
fi
BOT="${JARVIS_BRIDGE_BOT:-$SCRIPT_DIR/jarvis_dingtalk_bot.py}"
TASK_WORKER="${JARVIS_TASK_WORKER:-$SCRIPT_DIR/task_worker.py}"
STATE_DIR="${JARVIS_BRIDGE_STATE_DIR:-$REPO_ROOT/.my-day/bridge}"
# PIDFILE/LOG default to scheduler role; _resolve_paths_by_role() re-derives after
# _source_env so JARVIS_BRIDGE_ROLE from env files applies. Worker role uses
# bot-worker.{pid,log} so scheduler and worker can coexist on the same host without
# stepping on each other's pidfile/log (though normal deployment runs one role per host).
PIDFILE="$STATE_DIR/bot.pid"
LOG="$STATE_DIR/bot.log"
BOOTSTRAP_ENV="${JARVIS_BRIDGE_BOOTSTRAP_ENV:-$REPO_ROOT/bootstrap/.env}"
BRIDGE_ENV="${JARVIS_BRIDGE_ENV:-$SCRIPT_DIR/jarvis.env}"
START_WAIT="${JARVIS_BRIDGE_START_WAIT:-2}"
SCHEDULER_PIDFILE="$STATE_DIR/scheduler.pid"
SCHEDULER_LOG="$STATE_DIR/scheduler.log"
SCHEDULER_READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
SCHEDULER_DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-600}"
TASK_WORKER_PIDFILE="$STATE_DIR/task-worker.pid"
TASK_WORKER_LOG="$STATE_DIR/task-worker.log"
PRESERVE_TASK_WORKER_ONCE="$STATE_DIR/preserve-task-worker-once"
say()  { printf '%s\n' "$*"; }
err()  { printf '%s\n' "$*" >&2; }

# -- pid helpers -----------------------------------------------------------
_read_pid() {
  [ -f "$PIDFILE" ] || return 1
  local p; p="$(cat "$PIDFILE" 2>/dev/null || true)"
  [ -n "$p" ] || return 1
  printf '%s' "$p"
}
_alive() { [ -n "${1:-}" ] && kill -0 "$1" 2>/dev/null; }
_running_pid() {  # echo pid iff pidfile points at a live process
  local p; p="$(_read_pid)" || return 1
  _alive "$p" && { printf '%s' "$p"; return 0; }
  return 1
}

_unmanaged_bot_pids() {
  # A manually launched bridge has no pidfile, so blindly starting another
  # instance would create duplicate scanners and Task producers.  The default
  # production bot name is stable; test/custom bots opt out by overriding
  # JARVIS_BRIDGE_BOT.
  [ -n "${JARVIS_BRIDGE_BOT:-}" ] && return 0
  command -v pgrep >/dev/null 2>&1 || return 0
  pgrep -f '[j]arvis_dingtalk_bot' 2>/dev/null || true
}

_tail_log() {  # $1 = n
  local n="${1:-20}"
  if [ -f "$LOG" ]; then tail -n "$n" "$LOG"; else say "(暂无日志: $LOG)"; fi
}

_bridge_stop_wait() {
  printf '%s' "${JARVIS_BRIDGE_STOP_WAIT:-${JARVIS_STOP_GRACE:-30}}"
}

_remove_pidfile_if_matches() { # $1 = pidfile, $2 = pid
  local pidfile="$1" expected_pid="$2" current_pid=""
  [ -f "$pidfile" ] || return 0
  current_pid="$(cat "$pidfile" 2>/dev/null || true)"
  [ "$current_pid" = "$expected_pid" ] && rm -f "$pidfile"
}

_rollback_started_process() { # $1 = pid, $2 = pidfile, $3 = label
  # A process that never reached READY owns no durable lifecycle yet.  Always
  # reap the exact child started by this shell; do not route this through the
  # graceful runtime stop functions, which may intentionally preserve leased
  # work and therefore leave a failed-start orphan behind.
  local pid="$1" pidfile="$2" label="$3"
  local i=0 rollback_wait="${JARVIS_BRIDGE_START_ROLLBACK_WAIT:-5}"
  local deadline=$(( rollback_wait * 10 ))
  if _alive "$pid"; then
    kill -TERM "$pid" 2>/dev/null || true
    while [ "$i" -lt "$deadline" ] && _alive "$pid"; do
      sleep 0.1
      i=$((i + 1))
    done
  fi
  if _alive "$pid"; then
    err "$label 启动回滚超时，强制终止 pid $pid"
    kill -KILL "$pid" 2>/dev/null || true
  fi
  # The process is a direct child of this run.sh invocation. wait both reaps a
  # killed child and observes an already-exited startup failure.
  wait "$pid" 2>/dev/null || true
  if _alive "$pid"; then
    err "$label 启动回滚失败: pid $pid 仍在运行"
    return 1
  fi
  _remove_pidfile_if_matches "$pidfile" "$pid"
  return 0
}

_bridge_ready_in_log() { # $1 = launchd-owned pid
  local pid="$1"
  [ -f "$LOG" ] || return 1
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
    grep -F "Task worker READY pid=$pid " "$LOG" >/dev/null 2>&1
  else
    grep -F "Bridge READY pid=$pid role=supervisor" "$LOG" >/dev/null 2>&1
  fi
}

# -- role-aware paths (call AFTER _source_env; env files may set JARVIS_BRIDGE_ROLE) --
_resolve_paths_by_role() {
  local role="${JARVIS_BRIDGE_ROLE:-scheduler}"
  case "$role" in
    scheduler)
      PIDFILE="$STATE_DIR/bot.pid";        LOG="$STATE_DIR/bot.log" ;;
    worker)
      PIDFILE="$TASK_WORKER_PIDFILE"; LOG="$TASK_WORKER_LOG" ;;
    *)
      err "unsupported JARVIS_BRIDGE_ROLE=$role (accept: scheduler|worker)"
      return 2 ;;
  esac
  return 0
}

_scheduler_enabled() { [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ]; }

_component_read_pid() {
  local pidfile="$1" p
  [ -f "$pidfile" ] || return 1
  p="$(cat "$pidfile" 2>/dev/null || true)"
  [ -n "$p" ] || return 1
  printf '%s' "$p"
}

_component_running_pid() {
  local p; p="$(_component_read_pid "$1")" || return 1
  _alive "$p" && { printf '%s' "$p"; return 0; }
  return 1
}

_scheduler_validate() {
  _scheduler_enabled || return 0
  PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" "$PYTHON" -c '
from bridge.scheduler.jobs import JOBS
if not JOBS:
    raise SystemExit("scheduler registry is empty")
' || return $?
  if [ -z "${JARVIS_CONTROL_PLANE_TOKEN:-}" ] && [ -z "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
    err "scheduler control-plane token is required"
    return 2
  fi
}

_task_worker_validate() {
  if [ -z "${JARVIS_CONTROL_PLANE_TOKEN:-}" ] && [ -z "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
    err "task worker control-plane token is required"
    return 2
  fi
}

_component_start() {
  local label="$1" pidfile="$2" logfile="$3" ready="$4" ready_wait="$5" validate="$6"
  shift 6
  local pid i=0 deadline=$(( ready_wait * 10 ))
  COMPONENT_STARTED_PID=""
  if pid="$(_component_running_pid "$pidfile")"; then
    say "$label 已在运行 (pid $pid)"
    return 0
  fi
  rm -f "$pidfile"
  "$validate" || return $?
  mkdir -p "$STATE_DIR"
  nohup "$@" >>"$logfile" 2>&1 &
  pid=$!
  COMPONENT_STARTED_PID="$pid"
  printf '%s\n' "$pid" >"$pidfile"
  while [ "$i" -lt "$deadline" ]; do
    if ! _alive "$pid"; then
      err "$label 启动失败"
      tail -n 20 "$logfile" >&2
      _rollback_started_process "$pid" "$pidfile" "$label" || true
      return 1
    fi
    if grep -F "$ready pid=$pid " "$logfile" >/dev/null 2>&1; then
      say "$label 已启动 (pid $pid, log=$logfile)"
      return 0
    fi
    sleep 0.1; i=$((i + 1))
  done
  err "$label 未在 ${ready_wait}s 内 READY"
  _rollback_started_process "$pid" "$pidfile" "$label" || true
  return 1
}

_scheduler_start() {
  _scheduler_enabled || return 0
  _component_start \
    "scheduler" "$SCHEDULER_PIDFILE" "$SCHEDULER_LOG" "Scheduler READY" \
    "$SCHEDULER_READY_WAIT" _scheduler_validate \
    env "PYTHONPATH=$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
    "$PYTHON" -m bridge.main
}

_scheduler_stop() {
  _scheduler_enabled || return 0
  local pid i=0 deadline=$(( SCHEDULER_DRAIN_WAIT * 10 ))
  pid="$(_component_read_pid "$SCHEDULER_PIDFILE")" || { say "scheduler 未运行"; return 0; }
  _alive "$pid" || { rm -f "$SCHEDULER_PIDFILE"; say "scheduler 已停止"; return 0; }
  kill -TERM "$pid"
  say "scheduler 正在 drain (pid $pid, 宽限 ${SCHEDULER_DRAIN_WAIT}s)"
  while [ "$i" -lt "$deadline" ] && _alive "$pid"; do sleep 0.1; i=$((i + 1)); done
  if _alive "$pid"; then
    err "scheduler 仍有已准入 Job 未完成；保持当前进程，不停止 bridge。"
    return 1
  fi
  rm -f "$SCHEDULER_PIDFILE"
  say "scheduler 已停止 (graceful)"
}

_task_worker_start() {
  TASK_WORKER_STARTED_PID=""
  if [ -n "${JARVIS_TASK_WORKER:-}" ]; then
    _component_start \
      "task worker" "$TASK_WORKER_PIDFILE" "$TASK_WORKER_LOG" "Task worker READY" \
      "$SCHEDULER_READY_WAIT" _task_worker_validate \
      "$PYTHON" "$TASK_WORKER"
  else
    _component_start \
      "task worker" "$TASK_WORKER_PIDFILE" "$TASK_WORKER_LOG" "Task worker READY" \
      "$SCHEDULER_READY_WAIT" _task_worker_validate \
      env "PYTHONPATH=$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
      "$PYTHON" -m bridge.task_worker
  fi
  local rc=$?
  TASK_WORKER_STARTED_PID="${COMPONENT_STARTED_PID:-}"
  return "$rc"
}

_task_worker_stop() {
  local pid i=0 stop_wait deadline
  pid="$(_component_read_pid "$TASK_WORKER_PIDFILE")" || { say "task worker 未运行"; return 0; }
  _alive "$pid" || { rm -f "$TASK_WORKER_PIDFILE"; say "task worker 已停止"; return 0; }
  stop_wait="$(_bridge_stop_wait)"
  deadline=$(( stop_wait * 10 ))
  say "停止 task worker (pid $pid): 发 SIGTERM，宽限 ${stop_wait}s…"
  kill -TERM "$pid"
  while [ "$i" -lt "$deadline" ] && _alive "$pid"; do sleep 0.1; i=$((i + 1)); done
  if _alive "$pid"; then
    err "task worker 在 ${stop_wait}s 内未停止；保持进程以避免中断 Session。"
    return 1
  fi
  rm -f "$TASK_WORKER_PIDFILE"
  say "task worker 已停止 (graceful)"
}

_component_status() {
  local label="$1" pidfile="$2" logfile="$3" pid
  if pid="$(_component_running_pid "$pidfile")"; then
    say "$label: RUNNING (pid $pid, log=$logfile)"
  else
    say "$label: STOPPED"
  fi
}

_task_worker_status() {
  _component_status "task worker" "$TASK_WORKER_PIDFILE" "$TASK_WORKER_LOG"
}

_scheduler_status() {
  _scheduler_enabled || return 0
  _component_status "scheduler" "$SCHEDULER_PIDFILE" "$SCHEDULER_LOG"
}

# -- env sourcing (variables auto-exported for the bot) --------------------
_source_env() {
  # All Jarvis entrypoints share one machine-level credential loader. Explicit
  # bridge paths remain compatibility inputs for installations that override
  # the legacy file locations.
  export JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$BOOTSTRAP_ENV"
  export JARVIS_INTERACTIVE_BRIDGE_ENV="$BRIDGE_ENV"
  # shellcheck disable=SC1091
  source "$REPO_ROOT/bootstrap/runtime-config.sh"
  jarvis_load_runtime_config || return $?
  START_WAIT="${JARVIS_BRIDGE_START_WAIT:-2}"
  SCHEDULER_READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
  SCHEDULER_DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-600}"
  # Non-interactive SSH/launch shells on macOS commonly omit both the user-local
  # installer directory (a1) and Homebrew (claude). Normalize the standard tool
  # locations before spawning the daemon so registration implies an executable
  # worker instead of a process that can heartbeat but cannot run a Task.
  local tool_dir tool_dirs
  PATH="${PATH:-/usr/bin:/bin:/usr/sbin:/sbin}"
  # Space-separated override is intentionally test/deployment-only; standard paths
  # contain no spaces on supported macOS hosts.
  tool_dirs="${JARVIS_BRIDGE_TOOL_DIRS:-/usr/local/bin /opt/homebrew/sbin /opt/homebrew/bin ${HOME:-}/.local/bin}"
  for tool_dir in $tool_dirs; do
    [ -d "$tool_dir" ] || continue
    case ":$PATH:" in
      *":$tool_dir:"*) ;;
      *) PATH="$tool_dir:$PATH" ;;
    esac
  done
  export PATH

  # bootstrap/lib.sh deliberately resolves git worktrees back to the common main
  # repository unless JARVIS_ROOT is explicit. A bridge launched from a review
  # worktree must use that worktree's guarded wrappers and configuration as one
  # coherent version, while still allowing an operator-provided override.
  if [ -z "${JARVIS_ROOT:-}" ]; then
    export JARVIS_ROOT="$REPO_ROOT"
  fi

  # Role-aware PIDFILE/LOG must be re-derived AFTER env files are sourced (they
  # may set JARVIS_BRIDGE_ROLE). Failure to resolve exits early via the || guard.
  _resolve_paths_by_role || return $?
}

# -- mode decision: 0 = full, 1 = degraded (also exports the flag) ---------
_decide_mode() {
  # Workers never run the DingTalk stream client (scheduler-only); inherited
  # DINGTALK_* creds from the credential bundle must not flip a worker to full.
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
    export JARVIS_NO_DINGTALK=1
    return 1
  fi
  if [ -n "${DINGTALK_APP_KEY:-}" ] && [ -n "${DINGTALK_APP_SECRET:-}" ] \
     && [ "${JARVIS_NO_DINGTALK:-}" != "1" ]; then
    return 0
  fi
  export JARVIS_NO_DINGTALK=1
  return 1
}

_mode_from_log() {  # infer running mode from the latest mode-defining banner in the log
  # Classify only on mode-UNIQUE banners: degraded bot prefixes its lines with
  # [NO-DINGTALK] and never opens a stream; full bot logs "starting DingTalk stream
  # listener". "scan scheduler started" is emitted by BOTH modes, so it must NOT be a
  # discriminator (else the last such line would mask the degraded banner).
  [ -f "$LOG" ] || { printf 'unknown'; return; }
  local last
  last="$(grep -E '\[NO-DINGTALK\]|starting DingTalk stream listener' "$LOG" 2>/dev/null | tail -n1)"
  case "$last" in
    *'[NO-DINGTALK]'*) printf '降级(no-dingtalk)' ;;
    '')                printf 'unknown' ;;
    *)                 printf '全功能(dingtalk)' ;;
  esac
}

# -- launchd supervisor ----------------------------------------------------
# 配置可以来自 shell，也可以来自 bootstrap/.env / bridge/jarvis.env。所有值都在调用时读取，
# 避免脚本启动早期固化默认值后漏掉 env 文件中的 JARVIS_BRIDGE_SUPERVISOR。
_launchd_config() {
  LAUNCHCTL_BIN="${JARVIS_BRIDGE_LAUNCHCTL:-launchctl}"
  LAUNCHD_LABEL="${JARVIS_BRIDGE_LAUNCHD_LABEL:-com.jarvis.dingtalk}"
  LAUNCHD_DOMAIN="${JARVIS_BRIDGE_LAUNCHD_DOMAIN:-gui/$(id -u)}"
  LAUNCHD_PLIST="${JARVIS_BRIDGE_LAUNCHD_PLIST:-${HOME:-}/Library/LaunchAgents/${LAUNCHD_LABEL}.plist}"
  LAUNCHD_SERVICE="${LAUNCHD_DOMAIN}/${LAUNCHD_LABEL}"
}

_launchd_require() {
  _launchd_config
  if ! command -v "$LAUNCHCTL_BIN" >/dev/null 2>&1; then
    err "launchd 模式不可用: 找不到 launchctl ($LAUNCHCTL_BIN)。"
    return 1
  fi
}

_launchd_loaded() {
  "$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" >/dev/null 2>&1
}

cmd_launchd_start() {
  local detail state pid
  _launchd_require || return 1
  if [ ! -f "$LAUNCHD_PLIST" ]; then
    err "launchd plist 不存在: $LAUNCHD_PLIST"
    err "请先运行: $SCRIPT_DIR/install-launchd.sh"
    return 1
  fi
  if _launchd_loaded; then
    detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null || true)"
    state="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*state = //p' | head -n1)"
    pid="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
    if [ "$state" = "running" ] && [ -n "$pid" ]; then
      if _bridge_ready_in_log "$pid"; then
        say "bridge 已由 launchd 运行 ($LAUNCHD_SERVICE, pid $pid)。"
        return 0
      fi
      err "launchd 中的 bridge 尚未 READY；请检查日志或执行受控 restart。"
      return 1
    fi
  else
    say "注册 launchd service: $LAUNCHD_SERVICE"
    "$LAUNCHCTL_BIN" bootstrap "$LAUNCHD_DOMAIN" "$LAUNCHD_PLIST" || {
      err "launchd bootstrap 失败: $LAUNCHD_PLIST"
      return 1
    }
  fi
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ]; then
    "$LAUNCHCTL_BIN" kickstart "$LAUNCHD_SERVICE" || {
      err "launchd kickstart 失败: $LAUNCHD_SERVICE"
      return 1
    }
  else
    "$LAUNCHCTL_BIN" kickstart -k "$LAUNCHD_SERVICE" || {
      err "launchd kickstart 失败: $LAUNCHD_SERVICE"
      return 1
    }
  fi
  say "bridge 已交由 launchd 启动 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_stop() {
  _launchd_require || return 1
  if ! _launchd_loaded; then
    say "bridge 未由 launchd 加载 ($LAUNCHD_SERVICE)。"
    return 0
  fi
  "$LAUNCHCTL_BIN" bootout "$LAUNCHD_SERVICE" || {
    err "launchd bootout 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge launchd service 已停止并卸载 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_restart() {
  _launchd_require || return 1
  if ! _launchd_loaded; then
    cmd_launchd_start
    return
  fi

  # The launchd daemon owns the Bot and standalone Scheduler. Replace those
  # together while the PID-bound marker keeps the independent Task worker
  # alive across the planned restart.
  local detail old_pid="" i=0 stop_wait preserve_worker=0
  detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null || true)"
  old_pid="$(printf '%s\n' "$detail" \
    | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  stop_wait="$(_bridge_stop_wait)"
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ] && [ -n "$old_pid" ]; then
    # The foreground daemon owns bot+scheduler+worker. A planned restart must
    # replace only bot+scheduler so leased Task Sessions retain their worker
    # process/fence. Bind the one-shot marker to the exact daemon PID so a
    # stale file can never weaken a later full stop.
    mkdir -p "$STATE_DIR"
    printf '%s\n' "$old_pid" >"$PRESERVE_TASK_WORKER_ONCE"
    preserve_worker=1
  fi
  "$LAUNCHCTL_BIN" disable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || {
    [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_TASK_WORKER_ONCE"
    err "launchd restart 失败: 无法禁用 KeepAlive ($LAUNCHD_SERVICE)"
    return 1
  }
  if [ -n "$old_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$LAUNCHD_SERVICE" || {
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_TASK_WORKER_ONCE"
      err "launchd restart 失败: 无法请求旧 bridge 优雅停止"
      return 1
    }
    local deadline=$(( stop_wait * 10 ))
    while [ "$i" -lt "$deadline" ] && _alive "$old_pid"; do
      sleep 0.1
      i=$((i + 1))
    done
    if _alive "$old_pid"; then
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_TASK_WORKER_ONCE"
      "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || true
      err "bridge 在 ${stop_wait}s 内未停止；本次 restart 已取消。"
      return 1
    fi
  fi
  # The old daemon normally consumes the marker in its TERM handler. If it
  # exited before doing so, remove the stale PID-bound marker now.
  [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_TASK_WORKER_ONCE"
  "$LAUNCHCTL_BIN" enable "$LAUNCHD_SERVICE" >/dev/null 2>&1 || {
    err "launchd restart 失败: 无法重新启用 $LAUNCHD_SERVICE"
    return 1
  }
  "$LAUNCHCTL_BIN" kickstart "$LAUNCHD_SERVICE" || {
    err "launchd restart 失败: $LAUNCHD_SERVICE"
    return 1
  }
  say "bridge 已重启 ($LAUNCHD_SERVICE)。"
}

cmd_launchd_status() {
  _launchd_require || return 1
  local detail state pid
  if ! detail="$("$LAUNCHCTL_BIN" print "$LAUNCHD_SERVICE" 2>/dev/null)"; then
    say "bridge: STOPPED  (launchd 未加载 $LAUNCHD_SERVICE)"
    return 0
  fi
  state="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*state = //p' | head -n1)"
  pid="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  case "$state" in
    running) say "bridge: RUNNING  (launchd $LAUNCHD_SERVICE, pid ${pid:-?})" ;;
    *)       say "bridge: LOADED  (launchd $LAUNCHD_SERVICE, state ${state:-unknown})" ;;
  esac
  say "  plist: $LAUNCHD_PLIST"
  say "  log:   $LOG"
}

_supervised_command() { # $1 = start|stop|restart|status
  local command="$1" supervisor
  _source_env
  supervisor="${JARVIS_BRIDGE_SUPERVISOR:-local}"
  case "$supervisor" in
    local|'') "cmd_${command}" ;;
    launchd)  "cmd_launchd_${command}" ;;
    *)
      err "不支持的 JARVIS_BRIDGE_SUPERVISOR=$supervisor (可选: local, launchd)"
      return 2 ;;
  esac
}

# -- commands --------------------------------------------------------------
cmd_start() {
  local pid unmanaged worker_started_pid=""
  # Human start = intent to run; clear the manual-stop sentinel so the cron
  # watchdog resumes its keep-alive duty.
  rm -f "$PIDFILE.manual-stop" 2>/dev/null || true
  # A worker host runs only the independently supervised Task worker.  It must
  # never create a DingTalk listener or a SchedulerEngine process.
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
    _task_worker_start
    return $?
  fi
  if pid="$(_running_pid)"; then
    say "bridge 已在运行 (pid $pid) — 无需重复 start。"
    # The listener and durable worker have independent lifecycles.  A watchdog
    # start must repair a missing worker without replacing the healthy bot.
    _task_worker_start || return $?
    worker_started_pid="${TASK_WORKER_STARTED_PID:-}"
    if ! _scheduler_start; then
      # Preserve only a worker that existed before this start invocation.  A
      # freshly-created companion must not survive a partially failed start.
      if [ -n "$worker_started_pid" ]; then
        _rollback_started_process "$worker_started_pid" "$TASK_WORKER_PIDFILE" "task worker" || true
      fi
      return 1
    fi
    return 0
  fi
  unmanaged="$(_unmanaged_bot_pids)"
  if [ -n "$unmanaged" ]; then
    err "检测到未被当前 pidfile 管理的 bridge 进程: $(printf '%s' "$unmanaged" | tr '\n' ' ')"
    err "为避免双 Scanner/双 Task 生产者，拒绝启动。请先确认并停止这些进程。"
    return 1
  fi
  [ -f "$PIDFILE" ] && rm -f "$PIDFILE"   # stale pidfile

  mkdir -p "$STATE_DIR"

  local mode
  if _decide_mode; then
    mode="full"
    say "钉钉凭证就绪 → 全功能模式启动。"
  else
    mode="degraded"
    say "=================================================================="
    if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
      say "  role=worker → 强制降级模式 (worker 不跑钉钉 stream, 发声统一归 scheduler)。"
      say "  被动接单执行照常; 本机播报落 bot.log ([BROADCAST]) 作审计轨迹。"
    else
      say "  无钉钉凭证 (DINGTALK_APP_KEY/SECRET 未设) → 以降级模式点火。"
      say "  自动派发 + 各调度器照常运行; 卡片/播报落 bot.log ([BROADCAST])。"
      say "  配好凭证后运行  bridge/run.sh restart  即回全功能模式。"
    fi
    say "=================================================================="
  fi

  local pre_sz=0
  [ -f "$LOG" ] && pre_sz="$(wc -c <"$LOG" 2>/dev/null | tr -d ' ')"
  [ -n "$pre_sz" ] || pre_sz=0

  say "启动 bridge: $PYTHON ${JARVIS_BRIDGE_BOT:-bridge.jarvis_dingtalk_bot}  (mode=$mode, role=${JARVIS_BRIDGE_ROLE:-scheduler}, log=$LOG)"
  if [ -n "${JARVIS_BRIDGE_BOT:-}" ]; then
    nohup "$PYTHON" "$BOT" >>"$LOG" 2>&1 &
  else
    PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
      nohup "$PYTHON" -m bridge.jarvis_dingtalk_bot >>"$LOG" 2>&1 &
  fi
  local newpid=$!
  printf '%s\n' "$newpid" >"$PIDFILE"
  sleep "$START_WAIT"

  if ! _alive "$newpid"; then
    err "bridge 启动失败: 进程 $newpid 已退出。日志尾 ↓"
    _tail_log 20 >&2
    if [ "$mode" = "degraded" ]; then
      err "提示: 降级点火依赖 bot 支持 JARVIS_NO_DINGTALK(F 线 MR-4, 分支 worktree-f3-nodingtalk)。"
      err "      若该能力尚未合并, 旧 bot 会因缺凭证退 2 —— 请先合 MR-4, 或配置钉钉凭证后重试。"
    fi
    rm -f "$PIDFILE"
    return 1
  fi

  local first_new
  first_new="$(tail -c "+$((pre_sz + 1))" "$LOG" 2>/dev/null | head -n1)"
  case "$first_new" in
    *ERROR*)
      err "bridge 启动异常: 日志首行为 ERROR。日志尾 ↓"
      _tail_log 20 >&2
      cmd_stop >/dev/null 2>&1 || true
      return 1 ;;
  esac

  say "bridge 已启动 (pid $newpid, mode=$mode)。日志: $LOG"
  if ! _task_worker_start; then
    err "task worker 启动失败；回滚刚启动的 bridge。"
    _stop_bridge_only >/dev/null 2>&1 || true
    return 1
  fi
  worker_started_pid="${TASK_WORKER_STARTED_PID:-}"
  if ! _scheduler_start; then
    err "scheduler 启动失败；回滚本轮新启动的组件。"
    if [ -n "$worker_started_pid" ]; then
      _rollback_started_process "$worker_started_pid" "$TASK_WORKER_PIDFILE" "task worker" || true
    fi
    _stop_bridge_only >/dev/null 2>&1 || true
    return 1
  fi
  return 0
}

cmd_watchdog() {
  # cron entry (@reboot + */10 keep-alive on worker hosts): start the bridge
  # ONLY when the operator has not deliberately stopped it. cmd_stop drops the
  # manual-stop sentinel; only a human `run.sh start` clears it. Silent no-op
  # while stopped-on-purpose or already running (cmd_start is pidfile-guarded).
  if [ -f "$PIDFILE.manual-stop" ]; then
    return 0
  fi
  cmd_start
}

_stop_bridge_only() {
  local pid stop_wait
  stop_wait="$(_bridge_stop_wait)"
  # Operator-intent sentinel: `run.sh stop` means STAY stopped. The cron
  # watchdog (run.sh watchdog) refuses to start while this file exists; only
  # a human `run.sh start` clears it. Dropped even when nothing is running —
  # stop expresses intent, not just process teardown. (Without this, the
  # */10 watchdog resurrected a deliberately stopped worker within 10min.)
  mkdir -p "$STATE_DIR" 2>/dev/null || true
  touch "$PIDFILE.manual-stop" 2>/dev/null || true
  pid="$(_read_pid)" || { say "bridge 未在运行 (无 pidfile)。已落 manual-stop 哨兵，watchdog 不会拉起。"; return 0; }
  if ! _alive "$pid"; then
    say "bridge 未在运行 (pid $pid 已退) — 清理 pidfile。"
    rm -f "$PIDFILE"
    return 0
  fi
  # 发 SIGTERM → bot 的 _graceful_stop handler；超时后保留 SIGKILL 兜底。
  say "停止 bridge (pid $pid): 发 SIGTERM (bot 自杀 worker + release claim), 宽限 ${stop_wait}s…"
  kill -TERM "$pid" 2>/dev/null || true
  local i=0 deadline=$(( stop_wait * 10 ))
  while [ "$i" -lt "$deadline" ] && _alive "$pid"; do sleep 0.1; i=$((i + 1)); done
  if _alive "$pid"; then
    say "TERM ${stop_wait}s 未退 → SIGKILL 兜底(清理可能不全, reconcile 收敛)。"
    kill -KILL "$pid" 2>/dev/null || true
    sleep 0.2
    rm -f "$PIDFILE"
    say "bridge 已停止 (forced)。"
    return 0
  fi
  rm -f "$PIDFILE"
  say "bridge 已停止 (graceful)。"
  return 0
}

cmd_stop() {
  _scheduler_stop || return $?
  _task_worker_stop || return $?
  _stop_bridge_only
}

cmd_restart() {
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "worker" ]; then
    _task_worker_stop || return $?
    _task_worker_start
    return $?
  fi
  # Scheduler restarts intentionally leave the independent Task worker alive;
  # this keeps its leased Sessions fenced to the same worker process.
  _scheduler_stop || return $?
  _stop_bridge_only || return $?
  cmd_start
}

cmd_status() {
  local pid
  if pid="$(_running_pid)"; then
    local uptime mode
    uptime="$(ps -o etime= -p "$pid" 2>/dev/null | tr -d ' ')"
    mode="$(_mode_from_log)"
    say "bridge: RUNNING  (pid $pid, uptime ${uptime:-?}, mode ${mode}, role ${JARVIS_BRIDGE_ROLE:-scheduler})"
    say "  pidfile: $PIDFILE"
    say "  log:     $LOG"
    say "  --- 最近 5 行日志 ---"
    _tail_log 5
    _scheduler_status
    _task_worker_status
    return 0
  fi
  if [ -f "$PIDFILE" ]; then
    say "bridge: STOPPED  (pidfile 存在但进程已退, 建议 run.sh start)"
  else
    say "bridge: STOPPED  (无 pidfile)"
  fi
  _scheduler_status
  _task_worker_status
  return 0
}

cmd_logs() {
  _source_env
  say "bridge 日志: $LOG"
  say "实时跟随:  tail -f \"$LOG\""
  say "--- 最近 20 行 ---"
  _tail_log 20
  if _scheduler_enabled; then
    say "--- scheduler 最近 20 行 ($SCHEDULER_LOG) ---"
    tail -n 20 "$SCHEDULER_LOG" 2>/dev/null || true
  fi
  say "--- task worker 最近 20 行 ($TASK_WORKER_LOG) ---"
  tail -n 20 "$TASK_WORKER_LOG" 2>/dev/null || true
  return 0
}

cmd_dryrun() {
  _source_env
  _decide_mode || true   # 与 start 同样判定/导出降级 flag, 仅不写 pidfile
  say "dry-run: $PYTHON ${JARVIS_BRIDGE_BOT:-bridge.jarvis_dingtalk_bot} --dry-run-once"
  if [ -n "${JARVIS_BRIDGE_BOT:-}" ]; then
    exec "$PYTHON" "$BOT" --dry-run-once
  fi
  PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
    exec "$PYTHON" -m bridge.jarvis_dingtalk_bot --dry-run-once
}

# launchd/systemd 等外部 supervisor 使用的前台入口。这里不 fork、不 nohup、不写 pidfile；
# exec 后 bot 与 supervisor 看到的是同一 PID，崩溃退出可被 KeepAlive 准确拉起。
cmd_daemon() {
  _source_env
  mkdir -p "$STATE_DIR"
  local mode bot_pid="" shutdown_started=0 preserve_task_worker=0 rc=0
  _daemon_shutdown() {
    local shutdown_rc=0 step_rc=0 preserve_pid=""
    # A supervisor may stop us while the standalone components are still
    # starting. Install this handler before spawning anything and make it
    # idempotent so both the signal and normal-exit paths can share it.
    if [ "$shutdown_started" -eq 1 ]; then
      return 0
    fi
    shutdown_started=1
    if [ -f "$PRESERVE_TASK_WORKER_ONCE" ]; then
      preserve_pid="$(cat "$PRESERVE_TASK_WORKER_ONCE" 2>/dev/null || true)"
      rm -f "$PRESERVE_TASK_WORKER_ONCE"
      if [ "$preserve_pid" = "$$" ]; then
        preserve_task_worker=1
      fi
    fi
    _scheduler_stop
    step_rc=$?
    [ "$step_rc" -eq 0 ] || shutdown_rc="$step_rc"
    if [ -n "$bot_pid" ]; then
      kill -TERM "$bot_pid" 2>/dev/null || true
      wait "$bot_pid" 2>/dev/null || true
    fi
    if [ "$preserve_task_worker" -eq 0 ]; then
      _task_worker_stop
      step_rc=$?
      if [ "$shutdown_rc" -eq 0 ] && [ "$step_rc" -ne 0 ]; then
        shutdown_rc="$step_rc"
      fi
    else
      say "受控 restart：保留 task worker 及其已租约 Session。"
    fi
    return "$shutdown_rc"
  }
  trap '_daemon_shutdown; shutdown_rc=$?; trap - TERM INT; exit "$shutdown_rc"' TERM INT
  if _decide_mode; then mode="full"; else mode="degraded"; fi
  say "bridge foreground daemon 启动 (mode=$mode, role=${JARVIS_BRIDGE_ROLE:-scheduler}, pid=$$): $PYTHON ${JARVIS_BRIDGE_BOT:-bridge.jarvis_dingtalk_bot}"
  if ! _scheduler_enabled; then
    trap - TERM INT
    if [ -n "${JARVIS_TASK_WORKER:-}" ]; then
      exec "$PYTHON" "$TASK_WORKER"
    fi
    PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
      exec "$PYTHON" -m bridge.task_worker
  fi
  local worker_started_pid=""
  _task_worker_start || return $?
  worker_started_pid="${TASK_WORKER_STARTED_PID:-}"
  if ! _scheduler_start; then
    if [ -n "$worker_started_pid" ]; then
      _rollback_started_process "$worker_started_pid" "$TASK_WORKER_PIDFILE" "task worker" || true
    fi
    return 1
  fi
  if [ -n "${JARVIS_BRIDGE_BOT:-}" ]; then
    "$PYTHON" "$BOT" &
  else
    PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
      "$PYTHON" -m bridge.jarvis_dingtalk_bot &
  fi
  bot_pid=$!
  sleep "$START_WAIT"
  if ! _alive "$bot_pid"; then
    wait "$bot_pid" 2>/dev/null
    rc=$?
    _daemon_shutdown
    return "$rc"
  fi
  say "Bridge READY pid=$$ role=supervisor"
  wait "$bot_pid"
  rc=$?
  trap - TERM INT
  _daemon_shutdown
  local shutdown_rc=$?
  [ "$shutdown_rc" -eq 0 ] || return "$shutdown_rc"
  return "$rc"
}

usage() {
  err "usage: bridge/run.sh <start|stop|restart|status|logs|dry-run|daemon>"
  return 2
}

case "${1:-}" in
  start)          _supervised_command start ;;
  stop)           _supervised_command stop ;;
  restart)        _supervised_command restart ;;
  watchdog)       _supervised_command watchdog ;;
  status)         _supervised_command status ;;
  logs)           cmd_logs ;;
  dry-run|dryrun) cmd_dryrun ;;
  daemon)         cmd_daemon ;;
  *)              usage ;;
esac
