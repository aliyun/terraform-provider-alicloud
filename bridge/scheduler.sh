#!/usr/bin/env bash
# Standalone Scheduler entrypoint. It never starts jarvis_dingtalk_bot.py.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PYTHON="${JARVIS_SCHEDULER_PYTHON:-$REPO_ROOT/.venv/bridge/bin/python}"
STATE_DIR="${JARVIS_SCHEDULER_STATE_DIR:-$REPO_ROOT/.my-day/bridge}"
PIDFILE="$STATE_DIR/scheduler.pid"
LOG="$STATE_DIR/scheduler.log"
READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-600}"

say() { printf '%s\n' "$*"; }
err() { printf '%s\n' "$*" >&2; }

source_runtime() {
  source "$REPO_ROOT/bootstrap/runtime-config.sh"
  jarvis_load_runtime_config
  export JARVIS_ROOT="${JARVIS_ROOT:-$REPO_ROOT}"
}

read_pid() { [ -f "$PIDFILE" ] && cat "$PIDFILE"; }
alive() { [ -n "${1:-}" ] && kill -0 "$1" 2>/dev/null; }

validate() {
  [ -x "$PYTHON" ] || { err "scheduler Python 不存在: $PYTHON"; return 2; }
  PYTHONPATH="$SCRIPT_DIR:$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" "$PYTHON" -c '
import os
from scheduler.jobs import REGISTRY
if not REGISTRY.scheduler_job_keys():
    raise SystemExit("scheduler registry is empty")
if not (os.environ.get("JARVIS_CONTROL_PLANE_TOKEN") or os.environ.get("JARVIS_HTML_REPORT_TOKEN")):
    raise SystemExit("scheduler control-plane token is required")
' || return $?
}

start() {
  local pid i=0 deadline=$(( READY_WAIT * 10 ))
  source_runtime || return $?
  if pid="$(read_pid)" && alive "$pid"; then
    say "scheduler 已在运行 (pid $pid)"
    return 0
  fi
  rm -f "$PIDFILE"
  validate || return $?
  mkdir -p "$STATE_DIR"
  nohup "$PYTHON" "$SCRIPT_DIR/main.py" >>"$LOG" 2>&1 &
  pid=$!
  printf '%s\n' "$pid" >"$PIDFILE"
  while [ "$i" -lt "$deadline" ]; do
    alive "$pid" || { err "scheduler 启动失败"; tail -n 20 "$LOG" >&2; rm -f "$PIDFILE"; return 1; }
    if grep -F "Scheduler READY pid=$pid " "$LOG" >/dev/null 2>&1; then
      say "scheduler 已启动 (pid $pid, log=$LOG)"
      return 0
    fi
    sleep 0.1; i=$((i + 1))
  done
  err "scheduler 未在 ${READY_WAIT}s 内 READY"
  return 1
}

stop() {
  local pid i=0 deadline=$(( DRAIN_WAIT * 10 ))
  pid="$(read_pid)" || { say "scheduler 未运行"; return 0; }
  alive "$pid" || { rm -f "$PIDFILE"; say "scheduler 已停止"; return 0; }
  kill -TERM "$pid"
  say "scheduler 正在 drain (pid $pid, 宽限 ${DRAIN_WAIT}s)"
  while [ "$i" -lt "$deadline" ] && alive "$pid"; do
    sleep 0.1
    i=$((i + 1))
  done
  if alive "$pid"; then
    err "scheduler 仍有已准入 Job 未完成；保持当前进程，不启动替代实例。"
    return 1
  fi
  rm -f "$PIDFILE"
  say "scheduler 已停止 (graceful)"
}

status() {
  local pid
  pid="$(read_pid)" || { say "scheduler: STOPPED"; return 0; }
  if alive "$pid"; then say "scheduler: RUNNING (pid $pid, log=$LOG)"; else say "scheduler: STOPPED"; fi
}

logs() { say "scheduler 日志: $LOG"; tail -n 50 "$LOG" 2>/dev/null || true; }

case "${1:-}" in
  start) start ;;
  stop) stop ;;
  restart) stop && start ;;
  status) status ;;
  logs) logs ;;
  *) err "usage: bridge/scheduler.sh <start|stop|restart|status|logs>"; exit 2 ;;
esac
