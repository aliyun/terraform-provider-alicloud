#!/usr/bin/env bash
# Render and install the Jarvis bridge LaunchAgent for the current macOS user.
#
# Safe/idempotent properties:
#   - only writes the configured user plist (never sudo/system domain);
#   - refuses symlink targets and a target owned by a different Label;
#   - renders to a temporary file, validates XML, then atomically replaces;
#   - reinstall validates first, then safely restarts an already-loaded service
#     without bootout so an adopted Persistent Worker keeps its active leases;
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
DRAIN_WAIT="${JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS:-600}"
PRESERVE_WORKER_ONCE="$STATE_DIR/preserve-persistent-worker-once"
TMP_PLIST=""

die() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
cleanup() {
  if [ -n "$TMP_PLIST" ]; then rm -f "$TMP_PLIST"; fi
  return 0
}
trap cleanup EXIT

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
  preserve_worker=0
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ] && [ -n "$old_pid" ]; then
    printf '%s\n' "$old_pid" >"$PRESERVE_WORKER_ONCE"
    preserve_worker=1
  fi
  "$LAUNCHCTL_BIN" disable "$SERVICE" >/dev/null 2>&1 \
    || {
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_WORKER_ONCE"
      die "cannot disable loaded launchd service before drain: $SERVICE"
    }
  if [ -n "$old_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$SERVICE" \
      || {
        [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_WORKER_ONCE"
        "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 || true
        die "cannot request graceful drain from loaded launchd service: $SERVICE"
      }
    i=0
    deadline=$(( DRAIN_WAIT * 10 ))
    while [ "$i" -lt "$deadline" ] && kill -0 "$old_pid" 2>/dev/null; do
      sleep 0.1
      i=$((i + 1))
    done
    if kill -0 "$old_pid" 2>/dev/null; then
      [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_WORKER_ONCE"
      "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 || true
      die "loaded bridge did not drain in ${DRAIN_WAIT}s; upgrade cancelled without starting a replacement"
    fi
  fi
  [ "$preserve_worker" -eq 1 ] && rm -f "$PRESERVE_WORKER_ONCE"
  "$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 \
    || die "cannot re-enable launchd service after drain: $SERVICE"
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
  worker_pid="$(cat "$STATE_DIR/persistent-worker.pid" 2>/dev/null || true)"
  preserve_failed_worker=0
  if [ "${JARVIS_BRIDGE_ROLE:-scheduler}" = "scheduler" ] \
      && [ -n "$failed_pid" ] && [ -n "$worker_pid" ] \
      && kill -0 "$worker_pid" 2>/dev/null; then
    printf '%s\n' "$failed_pid" >"$PRESERVE_WORKER_ONCE"
    preserve_failed_worker=1
  fi
  "$LAUNCHCTL_BIN" disable "$SERVICE" >/dev/null 2>&1 || true
  if [ -n "$failed_pid" ]; then
    "$LAUNCHCTL_BIN" kill SIGTERM "$SERVICE" >/dev/null 2>&1 || true
    i=0
    deadline=$(( DRAIN_WAIT * 10 ))
    while [ "$i" -lt "$deadline" ] && kill -0 "$failed_pid" 2>/dev/null; do
      sleep 0.1
      i=$((i + 1))
    done
  fi
  [ "$preserve_failed_worker" -eq 1 ] && rm -f "$PRESERVE_WORKER_ONCE"
  # Never boot out a loaded job or a failed first start that already opened a
  # durable Worker. bootout may recursively terminate the adopted lease owner.
  if [ "$service_was_loaded" -eq 0 ] \
      && { [ -z "$worker_pid" ] || ! kill -0 "$worker_pid" 2>/dev/null; }; then
    "$LAUNCHCTL_BIN" bootout "$SERVICE" >/dev/null 2>&1 || true
  fi
  die "launchd service did not report Bridge READY within ${READY_WAIT}s: $SERVICE"
fi

printf 'Jarvis bridge is managed by launchd: %s\n' "$SERVICE"
printf 'status: JARVIS_BRIDGE_SUPERVISOR=launchd %s status\n' "$RUN_SH"
printf 'logs:   %s\n' "$LOG_PATH"
