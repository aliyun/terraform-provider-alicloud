#!/usr/bin/env bash
# Render and install the Jarvis bridge LaunchAgent for the current macOS user.
#
# Safe/idempotent properties:
#   - only writes the configured user plist (never sudo/system domain);
#   - refuses symlink targets and a target owned by a different Label;
#   - renders to a temporary file, validates XML, then atomically replaces;
#   - bootout/bootstrap/kickstart converges both first install and reinstall.
set -euo pipefail

SCRIPT_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -P "$SCRIPT_DIR/.." && pwd)"
LABEL="${JARVIS_BRIDGE_LAUNCHD_LABEL:-com.jarvis.dingtalk}"
DOMAIN="${JARVIS_BRIDGE_LAUNCHD_DOMAIN:-gui/$(id -u)}"
SERVICE="$DOMAIN/$LABEL"
LAUNCHCTL_BIN="${JARVIS_BRIDGE_LAUNCHCTL:-launchctl}"
PYTHON_BIN="${JARVIS_BRIDGE_PYTHON:-python3}"
TEMPLATE="${JARVIS_BRIDGE_LAUNCHD_TEMPLATE:-$SCRIPT_DIR/com.jarvis.dingtalk.plist.example}"
TARGET="${JARVIS_BRIDGE_LAUNCHD_PLIST:-${HOME:-}/Library/LaunchAgents/$LABEL.plist}"
ENV_FILE="${JARVIS_BRIDGE_ENV:-$SCRIPT_DIR/jarvis.env}"
STATE_DIR="${JARVIS_BRIDGE_STATE_DIR:-$REPO_ROOT/.my-day/bridge}"
LOG_PATH="${JARVIS_BRIDGE_LOG:-$STATE_DIR/bot.log}"
LAUNCHD_PATH="${JARVIS_BRIDGE_LAUNCHD_PATH:-/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:${HOME:-}/.local/bin}"
RUN_SH="$SCRIPT_DIR/run.sh"
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

if [ -f "$TARGET" ] && cmp -s "$TMP_PLIST" "$TARGET"; then
  rm -f "$TMP_PLIST"
  TMP_PLIST=""
  printf 'launchd plist unchanged: %s\n' "$TARGET"
else
  mv -f "$TMP_PLIST" "$TARGET"
  TMP_PLIST=""
  printf 'launchd plist installed: %s\n' "$TARGET"
fi

if "$LAUNCHCTL_BIN" print "$SERVICE" >/dev/null 2>&1; then
  "$LAUNCHCTL_BIN" bootout "$SERVICE"
fi
"$LAUNCHCTL_BIN" bootstrap "$DOMAIN" "$TARGET"
"$LAUNCHCTL_BIN" enable "$SERVICE" >/dev/null 2>&1 || true
"$LAUNCHCTL_BIN" kickstart -k "$SERVICE"

printf 'Jarvis bridge is managed by launchd: %s\n' "$SERVICE"
printf 'status: JARVIS_BRIDGE_SUPERVISOR=launchd %s status\n' "$RUN_SH"
printf 'logs:   %s\n' "$LOG_PATH"
