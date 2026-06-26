#!/bin/sh
# Jarvis DingTalk bridge process manager. Portable: avoids bashisms, runs on /bin/sh.
set -eu

BRIDGE_DIR=$(cd "$(dirname "$0")" && pwd)
PY=${JARVIS_PY:-python3}
BOT="$BRIDGE_DIR/jarvis_dingtalk_bot.py"
STATE=${JARVIS_STATE:-$HOME/.jarvis}
PIDFILE="$STATE/dingtalk.pid"
LOGFILE="$STATE/dingtalk.log"
ENVFILE=${JARVIS_ENV_FILE:-$BRIDGE_DIR/jarvis.env}

mkdir -p "$STATE"

is_running() {
  [ -f "$PIDFILE" ] || return 1
  pid=$(cat "$PIDFILE" 2>/dev/null || echo)
  [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null
}

start() {
  if is_running; then
    echo "already running (pid $(cat "$PIDFILE"))"; exit 0
  fi
  [ -f "$ENVFILE" ] && { set -a; . "$ENVFILE"; set +a; }
  echo "starting bot, log -> $LOGFILE"
  nohup "$PY" "$BOT" >>"$LOGFILE" 2>&1 &
  echo $! > "$PIDFILE"
  ID=$(bash "$BRIDGE_DIR/../bootstrap/coord.sh" register dispatch 2>/dev/null||true); [ -n "$ID" ] && nohup bash "$BRIDGE_DIR/../bootstrap/heartbeat.sh" "$ID" "$(cat "$PIDFILE")" >/dev/null 2>&1 &
  sleep 1
  if is_running; then echo "started (pid $(cat "$PIDFILE"))"; else
    echo "FAILED to start; tail log:"; tail -n 20 "$LOGFILE"; exit 1; fi
}

stop() {
  if is_running; then
    pid=$(cat "$PIDFILE"); kill "$pid" 2>/dev/null || true
    sleep 1; kill -9 "$pid" 2>/dev/null || true
    echo "stopped (pid $pid)"
  else echo "not running"; fi
  rm -f "$PIDFILE"
}

status() {
  if is_running; then echo "running (pid $(cat "$PIDFILE"))"; else echo "stopped"; fi
}

case "${1:-}" in
  start)   start ;;
  stop)    stop ;;
  restart) stop; start ;;
  status)  status ;;
  logs)    tail -n "${2:-50}" -f "$LOGFILE" ;;
  *) echo "usage: $0 {start|stop|status|restart|logs}"; exit 2 ;;
esac
