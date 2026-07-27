#!/usr/bin/env bash
# Remove old a1 lock files only when no runnable a1 process exists.
set -uo pipefail

a1_root="${A1ID_ROOT:-${A1_HOME:-$HOME/.config/a1}}"
stale_seconds="${JARVIS_A1_LOCK_STALE_SEC:-30}"
script_dir="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
gate_helper="$script_dir/a1-file-gate.py"

# The exclusive descriptor covers the complete inventory/check/unlink pass.
# New a1id calls block on their shared descriptor; an already active a1id makes
# this opportunistic cleanup skip.  Kernel descriptor lifetime also makes
# TERM/KILL, zombies, and PID reuse safe without stale owner markers.
if [ "${JARVIS_A1_GATE_MODE:-}" != "exclusive" ] \
    || [ "${JARVIS_A1_GATE_ROOT:-}" != "$a1_root" ]; then
  gate_python="/usr/bin/python3"
  [ -x "$gate_python" ] || gate_python="$(command -v python3 2>/dev/null || true)"
  if [ -z "$gate_python" ] || [ ! -f "$gate_helper" ]; then
    echo "a1-locks-clean: skip (portable file gate unavailable)" >&2
    exit 1
  fi
  A1ID_ROOT="$a1_root" exec "$gate_python" -I "$gate_helper" \
    --root "$a1_root" exclusive -- "$0" "$@"
fi

case "$stale_seconds" in
  ''|*[!0-9]*)
    echo "a1-locks-clean: invalid JARVIS_A1_LOCK_STALE_SEC=$stale_seconds" >&2
    exit 1
    ;;
esac

if [ ! -d "$a1_root" ]; then
  echo "a1-locks-clean: skip (root not found: $a1_root)" >&2
  exit 0
fi

_runnable_a1_exists() {
  # Return 0=live, 1=none (zombies count as exited), 2=inspection failed.
  command -v pgrep >/dev/null 2>&1 || return 2
  command -v ps >/dev/null 2>&1 || return 2
  local output rc pid state
  output="$(pgrep -x a1 2>/dev/null)"
  rc=$?
  [ "$rc" -eq 1 ] && return 1
  [ "$rc" -eq 0 ] || return 2
  for pid in $output; do
    case "$pid" in
      ''|*[!0-9]*) continue ;;
    esac
    state="$(ps -o stat= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
    # A vanished PID is a normal inspection race.  A zombie has already
    # stopped executing and cannot retain an advisory/filesystem lock.
    [ -z "$state" ] && continue
    case "$state" in
      Z*) continue ;;
      *) return 0 ;;
    esac
  done
  return 1
}

_lock_signature() {
  # device:inode:size:mtime; supports BSD/macOS and GNU stat.
  stat -f '%d:%i:%z:%m' "$1" 2>/dev/null \
    || stat -c '%d:%i:%s:%Y' "$1" 2>/dev/null
}

if _runnable_a1_exists; then
  echo "a1-locks-clean: skip (live a1 process detected)" >&2
  exit 0
else
  inspect_rc=$?
  if [ "$inspect_rc" -eq 2 ]; then
    echo "a1-locks-clean: skip (cannot inspect a1 processes safely)" >&2
    exit 1
  fi
fi

locks=()
while IFS= read -r -d '' lock; do
  locks+=("$lock")
done < <(
  find "$a1_root" -type f \
    \( -name 'auth.yaml.lock' -o -name 'telemetry-queue.jsonl.lock' \) \
    -print0 2>/dev/null
)

now="$(date +%s)"
removed=0
for lock in "${locks[@]}"; do
  before="$(_lock_signature "$lock")" || continue
  modified="${before##*:}"
  case "$modified" in
    ''|*[!0-9]*) continue ;;
  esac
  age=$((now - modified))
  [ "$age" -ge "$stale_seconds" ] || continue

  # Close the useful check→delete race: a command may have started after the
  # inventory read, or may have replaced the lock.  Recheck both immediately
  # before unlinking and skip the rest of this pass as soon as a live a1 appears.
  if _runnable_a1_exists; then
    echo "a1-locks-clean: stop (a1 started during cleanup)" >&2
    break
  else
    inspect_rc=$?
    [ "$inspect_rc" -eq 1 ] || {
      echo "a1-locks-clean: stop (process recheck failed)" >&2
      break
    }
  fi
  after="$(_lock_signature "$lock")" || continue
  [ "$after" = "$before" ] || continue
  if rm -f -- "$lock" 2>/dev/null; then
    removed=$((removed + 1))
  fi
done

echo "a1-locks-clean: removed $removed stale lock(s) under $a1_root" >&2
exit 0
