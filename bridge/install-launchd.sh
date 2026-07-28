#!/usr/bin/env bash
# Render and install the Jarvis bridge LaunchAgent for the current macOS user.
#
# Safe/idempotent properties:
#   - only writes the configured user plist (never sudo/system domain);
#   - refuses symlink targets and a target owned by a different Label;
#   - renders to a temporary file, validates XML, then atomically replaces;
#   - reinstall validates first, then safely restarts an already-loaded service
#     without bootout, waiting for the Worker to relinquish active Sessions;
#   - never uses kickstart -k to overlap an old Scheduler with its replacement.
set -euo pipefail

SCRIPT_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -P "$SCRIPT_DIR/.." && pwd)"
BOOTSTRAP_ENV="${JARVIS_BRIDGE_BOOTSTRAP_ENV:-$REPO_ROOT/bootstrap/.env}"
ENV_FILE="${JARVIS_BRIDGE_ENV:-$SCRIPT_DIR/jarvis.env}"
export JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$BOOTSTRAP_ENV"
export JARVIS_INTERACTIVE_BRIDGE_ENV="$ENV_FILE"
# shellcheck disable=SC1091
source "$REPO_ROOT/bootstrap/runtime-config.sh"
jarvis_load_runtime_config

LABEL="${JARVIS_BRIDGE_LAUNCHD_LABEL:-com.jarvis.dingtalk}"
DOMAIN="${JARVIS_BRIDGE_LAUNCHD_DOMAIN:-gui/$(id -u)}"
SERVICE="$DOMAIN/$LABEL"
LAUNCHCTL_BIN="${JARVIS_BRIDGE_LAUNCHCTL:-launchctl}"
if [ -n "${JARVIS_BRIDGE_PYTHON:-}" ]; then
  PYTHON_BIN="$JARVIS_BRIDGE_PYTHON"
else
  # Reconcile the pinned runtime on every install/upgrade. This is still
  # pre-downtime: a dependency download/install failure leaves the loaded
  # service untouched.
  JARVIS_BRIDGE_VENV="$REPO_ROOT/.venv/bridge" \
    bash "$REPO_ROOT/bootstrap/bridge-python.sh"
  PYTHON_BIN="$REPO_ROOT/.venv/bridge/bin/python"
fi
TEMPLATE="${JARVIS_BRIDGE_LAUNCHD_TEMPLATE:-$SCRIPT_DIR/com.jarvis.dingtalk.plist.example}"
TARGET="${JARVIS_BRIDGE_LAUNCHD_PLIST:-${HOME:-}/Library/LaunchAgents/$LABEL.plist}"
STATE_DIR="${JARVIS_BRIDGE_STATE_DIR:-$REPO_ROOT/.my-day/bridge}"
LOG_PATH="${JARVIS_BRIDGE_LOG:-$STATE_DIR/bot.log}"
LAUNCHD_PATH="${JARVIS_BRIDGE_LAUNCHD_PATH:-/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${HOME:-}/.local/bin}"
RUN_SH="$SCRIPT_DIR/run.sh"
READY_WAIT="${JARVIS_SCHEDULER_READY_WAIT:-30}"
DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-30}"
TMP_PLIST=""
REENABLE_SERVICE=0

die() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
cleanup() {
  if [ "$REENABLE_SERVICE" -eq 1 ]; then
    "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 || true
  fi
  if [ -n "$TMP_PLIST" ]; then rm -f "$TMP_PLIST"; fi
  return 0
}
trap cleanup EXIT

managed_child_snapshot() { # $1 = supervisor pid
  local parent_pid="$1" pidfile pid ppid
  for pidfile in \
      "$STATE_DIR/scheduler.pid" \
      "$STATE_DIR/dingtalk-bot.pid" \
      "$STATE_DIR/persistent-worker.pid"; do
    pid="$(cat "$pidfile" 2>/dev/null || true)"
    [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null || continue
    ppid="$(ps -o ppid= -p "$pid" 2>/dev/null | tr -d ' ')"
    [ "$ppid" = "$parent_pid" ] || continue
    printf '%s:%s:%s\n' "$pidfile" "$pid" "$parent_pid"
  done
}

snapshot_entry_owned() { # $1 = pidfile:pid:original-parent
  local entry="$1" rest current ppid
  CHILD_PARENT="${entry##*:}"
  rest="${entry%:*}"
  CHILD_PID="${rest##*:}"
  CHILD_PIDFILE="${rest%:*}"
  current="$(cat "$CHILD_PIDFILE" 2>/dev/null || true)"
  [ "$current" = "$CHILD_PID" ] && kill -0 "$CHILD_PID" 2>/dev/null || return 1
  ppid="$(ps -o ppid= -p "$CHILD_PID" 2>/dev/null | tr -d ' ')"
  [ "$ppid" = "$CHILD_PARENT" ]
}

force_managed_snapshot() { # $1 = newline-separated pidfile:pid:parent
  local snapshot="$1" entry rest current i=0 any_alive unresolved
  [ -n "$snapshot" ] || return 0
  while IFS= read -r entry; do
    [ -n "$entry" ] || continue
    snapshot_entry_owned "$entry" || continue
    kill -KILL "$CHILD_PID" 2>/dev/null || true
  done <<EOF
$snapshot
EOF
  while [ "$i" -lt 20 ]; do
    any_alive=0
    unresolved=0
    while IFS= read -r entry; do
      [ -n "$entry" ] || continue
      CHILD_PARENT="${entry##*:}"
      rest="${entry%:*}"
      CHILD_PID="${rest##*:}"
      CHILD_PIDFILE="${rest%:*}"
      current="$(cat "$CHILD_PIDFILE" 2>/dev/null || true)"
      if [ "$current" != "$CHILD_PID" ]; then
        continue
      elif ! kill -0 "$CHILD_PID" 2>/dev/null; then
        rm -f "$CHILD_PIDFILE"
      elif snapshot_entry_owned \
          "$CHILD_PIDFILE:$CHILD_PID:$CHILD_PARENT"; then
        any_alive=1
      else
        unresolved=1
      fi
    done <<EOF
$snapshot
EOF
    [ "$unresolved" -eq 1 ] && return 1
    [ "$any_alive" -eq 0 ] && return 0
    sleep 0.1
    i=$((i + 1))
  done
  return 1
}

case "$LABEL" in
  *[!A-Za-z0-9._-]*|'') die "invalid launchd label: $LABEL" ;;
esac
[ "$(uname -s)" = "Darwin" ] || [ "${JARVIS_LAUNCHD_ALLOW_NON_DARWIN:-0}" = "1" ] \
  || die "launchd installation is only supported on macOS"
[ -n "${HOME:-}" ] || die 'HOME is empty'
command -v "$LAUNCHCTL_BIN" >/dev/null 2>&1 || die "launchctl not found: $LAUNCHCTL_BIN"
command -v "$PYTHON_BIN" >/dev/null 2>&1 || die "python not found: $PYTHON_BIN"

absolute_path() {
  "$PYTHON_BIN" - "$1" <<'PY'
import os
import sys

print(os.path.abspath(os.path.expanduser(sys.argv[1])))
PY
}
TEMPLATE="$(absolute_path "$TEMPLATE")"
TARGET="$(absolute_path "$TARGET")"
ENV_FILE="$(absolute_path "$ENV_FILE")"
STATE_DIR="$(absolute_path "$STATE_DIR")"
LOG_PATH="$(absolute_path "$LOG_PATH")"

[ -f "$TEMPLATE" ] || die "plist template not found: $TEMPLATE"
[ -x "$RUN_SH" ] || die "bridge entrypoint is not executable: $RUN_SH"

# Provision and validate the exact runtime before touching a loaded service.
# This catches a missing PyYAML, malformed jobs.yaml, runner drift, and role/token
# errors before downtime. JARVIS_AUTO_DISPATCH=0 remains a supported supervised mode.
PYTHONPATH="$REPO_ROOT${PYTHONPATH:+:$PYTHONPATH}" \
  "$PYTHON_BIN" -m bridge.main --validate \
  || die "new bridge runtime validation failed; existing service was not stopped"

[ ! -L "$TARGET" ] || die "refusing to replace symlink plist: $TARGET"
if [ -f "$TARGET" ]; then
  existing_label="$("$PYTHON_BIN" - "$TARGET" <<'PY'
import plistlib
import sys

try:
    with open(sys.argv[1], "rb") as stream:
        print(plistlib.load(stream).get("Label", ""))
except Exception:
    raise SystemExit(1)
PY
  )" || die "existing plist is not a readable plist; refusing to replace: $TARGET"
  [ "$existing_label" = "$LABEL" ] \
    || die "existing plist has Label=$existing_label; refusing to replace: $TARGET"
fi

mkdir -p "$(dirname "$TARGET")" "$STATE_DIR"
chmod 700 "$STATE_DIR"
if [ -e "$LOG_PATH" ]; then
  [ ! -L "$LOG_PATH" ] || die "refusing symlink log path: $LOG_PATH"
  chmod 600 "$LOG_PATH"
fi
TMP_PLIST="$(mktemp "${TARGET}.tmp.XXXXXX")"

"$PYTHON_BIN" - "$TEMPLATE" "$TMP_PLIST" \
  "$LABEL" "$RUN_SH" "$REPO_ROOT" "$ENV_FILE" "$STATE_DIR" "$TARGET" "$LAUNCHD_PATH" "$LOG_PATH" <<'PY'
import sys
import xml.etree.ElementTree as ET
import re
from pathlib import Path
from xml.sax.saxutils import escape

template, output, *values = sys.argv[1:]
keys = (
    "__LABEL__",
    "__RUN_SH_PATH__",
    "__REPO_ROOT__",
    "__ENV_FILE__",
    "__STATE_DIR__",
    "__PLIST_PATH__",
    "__LAUNCHD_PATH__",
    "__LOG_PATH__",
)
if len(keys) != len(values):
    raise SystemExit("invalid renderer arguments")
text = Path(template).read_text(encoding="utf-8")
for key, value in zip(keys, values):
    text = text.replace(key, escape(value))
if re.search(r"__[A-Z0-9_]+__", text):
    raise SystemExit("unresolved plist placeholder remains")
Path(output).write_text(text, encoding="utf-8")
ET.parse(output)
PY
chmod 600 "$TMP_PLIST"

plist_changed=0
if [ -f "$TARGET" ] && cmp -s "$TMP_PLIST" "$TARGET"; then
  rm -f "$TMP_PLIST"
  TMP_PLIST=""
  printf 'launchd plist unchanged: %s\n' "$TARGET"
else
  mv -f "$TMP_PLIST" "$TARGET"
  TMP_PLIST=""
  plist_changed=1
  printf 'launchd plist installed: %s\n' "$TARGET"
fi

service_was_loaded=0
pre_log_size=0
[ -f "$LOG_PATH" ] && pre_log_size="$(wc -c <"$LOG_PATH" 2>/dev/null | tr -d ' ')"
[ -n "$pre_log_size" ] || pre_log_size=0

if detail="$("$LAUNCHCTL_BIN" print "$SERVICE" 2>/dev/null)"; then
  service_was_loaded=1
  old_pid="$(printf '%s\n' "$detail" \
    | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  managed_children=""
  child_cleanup_ok=1
  [ -n "$old_pid" ] && managed_children="$(managed_child_snapshot "$old_pid")"
  "$LAUNCHCTL_BIN" disable "$SERVICE" >/dev/null 2>&1 \
    || die "cannot disable loaded launchd service before drain: $SERVICE"
  REENABLE_SERVICE=1
  if [ -n "$old_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$SERVICE" \
      || {
        printf 'WARNING: graceful drain request failed; force-stopping old bridge pid %s\n' \
          "$old_pid" >&2
      }
    i=0
    deadline=$(( DRAIN_WAIT * 10 ))
    while [ "$i" -lt "$deadline" ] && kill -0 "$old_pid" 2>/dev/null; do
      sleep 0.1
      i=$((i + 1))
    done
    if kill -0 "$old_pid" 2>/dev/null; then
      printf 'loaded bridge did not drain in %ss; force-stopping before replacement\n' \
        "$DRAIN_WAIT"
      # Best effort before stopping the parent. A killed child may remain a
      # zombie until the supervisor is reaped, so only the post-parent check
      # below is authoritative.
      force_managed_snapshot "$managed_children" || true
      "$LAUNCHCTL_BIN" kill SIGKILL "$SERVICE" >/dev/null 2>&1 \
        || kill -KILL "$old_pid" 2>/dev/null || true
      i=0
      while [ "$i" -lt 20 ] && kill -0 "$old_pid" 2>/dev/null; do
        sleep 0.1
        i=$((i + 1))
      done
      kill -0 "$old_pid" 2>/dev/null \
        && die "cannot stop old bridge pid $old_pid before replacement"
    fi
    child_cleanup_ok=1
    force_managed_snapshot "$managed_children" || child_cleanup_ok=0
    [ "$child_cleanup_ok" -eq 1 ] \
      || die "managed bridge child ownership changed during loaded upgrade"
  fi
  "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 \
    || die "cannot re-enable launchd service after drain: $SERVICE"
  REENABLE_SERVICE=0
  "$LAUNCHCTL_BIN" kickstart "$SERVICE" \
    || die "cannot start replacement bridge: $SERVICE"
  if [ "$plist_changed" -eq 1 ]; then
    printf '%s\n' \
      "launchd plist changed while loaded; the safe restart preserved the active job configuration." \
      "The installed plist will take effect after the next full stop/start or login."
  fi
else
  "$LAUNCHCTL_BIN" bootstrap "$DOMAIN" "$TARGET"
  "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 || true
  "$LAUNCHCTL_BIN" kickstart "$SERVICE"
fi

i=0
deadline=$(( READY_WAIT * 10 ))
ready=0
while [ "$i" -lt "$deadline" ]; do
  detail="$("$LAUNCHCTL_BIN" print "$SERVICE" 2>/dev/null || true)"
  pid="$(printf '%s\n' "$detail" | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null && [ -f "$LOG_PATH" ]; then
    log_size="$(wc -c <"$LOG_PATH" 2>/dev/null | tr -d ' ')"
    [ -n "$log_size" ] || log_size=0
    [ "$log_size" -ge "$pre_log_size" ] || pre_log_size=0
    if tail -c "+$((pre_log_size + 1))" "$LOG_PATH" 2>/dev/null \
        | grep -F "Bridge READY pid=$pid " >/dev/null 2>&1; then
      ready=1
      break
    fi
  fi
  sleep 0.1
  i=$((i + 1))
done
if [ "$ready" -ne 1 ]; then
  detail="$("$LAUNCHCTL_BIN" print "$SERVICE" 2>/dev/null || true)"
  failed_pid="$(printf '%s\n' "$detail" \
    | sed -n 's/^[[:space:]]*pid = //p' | head -n1)"
  failed_children=""
  [ -n "$failed_pid" ] \
    && failed_children="$(managed_child_snapshot "$failed_pid")"
  worker_pid="$(cat "$STATE_DIR/persistent-worker.pid" 2>/dev/null || true)"
  "$LAUNCHCTL_BIN" disable "$SERVICE" >/dev/null 2>&1 || true
  [ "$service_was_loaded" -eq 1 ] && REENABLE_SERVICE=1
  if [ -n "$failed_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$SERVICE" >/dev/null 2>&1 || true
    i=0
    deadline=$(( DRAIN_WAIT * 10 ))
    while [ "$i" -lt "$deadline" ] && kill -0 "$failed_pid" 2>/dev/null; do
      sleep 0.1
      i=$((i + 1))
    done
    if kill -0 "$failed_pid" 2>/dev/null; then
      force_managed_snapshot "$failed_children" || true
      "$LAUNCHCTL_BIN" kill SIGKILL "$SERVICE" >/dev/null 2>&1 \
        || kill -KILL "$failed_pid" 2>/dev/null || true
    fi
    force_managed_snapshot "$failed_children" || true
  fi
  # Never boot out a loaded job or a failed first start whose Worker did not
  # finish its bounded shutdown. bootout may bypass Session relinquish.
  if [ "$service_was_loaded" -eq 0 ] \
      && { [ -z "$worker_pid" ] || ! kill -0 "$worker_pid" 2>/dev/null; }; then
    "$LAUNCHCTL_BIN" bootout "$SERVICE" >/dev/null 2>&1 || true
  fi
  die "launchd service did not report Bridge READY within ${READY_WAIT}s: $SERVICE"
fi

printf 'Jarvis bridge is managed by launchd: %s\n' "$SERVICE"
printf 'status: JARVIS_BRIDGE_SUPERVISOR=launchd %s status\n' "$RUN_SH"
printf 'logs:   %s\n' "$LOG_PATH"
