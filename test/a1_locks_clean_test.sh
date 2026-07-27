#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cleaner="$repo_root/bootstrap/a1-locks-clean.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
bin="$tmp/bin"
root="$tmp/a1-root"
mkdir -p "$bin" "$root/identities/jarvis"
printf 'test-auth\n' >"$root/identities/jarvis/auth.yaml"

cat >"$bin/pgrep" <<'EOF'
#!/usr/bin/env bash
case "${STUB_PGREP_MODE:-none}" in
  none) exit 1 ;;
  live|zombie) echo 4242; exit 0 ;;
  race)
    n=0
    [ -f "$STUB_PGREP_COUNT" ] && n="$(cat "$STUB_PGREP_COUNT")"
    n=$((n + 1))
    printf '%s\n' "$n" >"$STUB_PGREP_COUNT"
    [ "$n" -eq 1 ] && exit 1
    echo 4242
    ;;
esac
EOF
cat >"$bin/ps" <<'EOF'
#!/usr/bin/env bash
case "${STUB_PGREP_MODE:-none}" in
  zombie) printf 'Z+\n' ;;
  *) printf 'S+\n' ;;
esac
EOF
chmod +x "$bin/pgrep" "$bin/ps"

lock="$root/identities/jarvis/auth.yaml.lock"

touch "$lock"
PATH="$bin:$PATH" A1ID_ROOT="$root" JARVIS_A1_LOCK_STALE_SEC=0 \
  STUB_PGREP_MODE=live bash "$cleaner" >/dev/null 2>&1
[ -f "$lock" ] || { echo "live a1 lock was removed" >&2; exit 1; }

PATH="$bin:$PATH" A1ID_ROOT="$root" JARVIS_A1_LOCK_STALE_SEC=0 \
  STUB_PGREP_MODE=zombie bash "$cleaner" >/dev/null 2>&1
[ ! -e "$lock" ] || { echo "zombie-only lock was not removed" >&2; exit 1; }

touch "$lock"
count="$tmp/pgrep-count"
PATH="$bin:$PATH" A1ID_ROOT="$root" JARVIS_A1_LOCK_STALE_SEC=0 \
  STUB_PGREP_MODE=race STUB_PGREP_COUNT="$count" \
  bash "$cleaner" >/dev/null 2>&1
[ -f "$lock" ] || { echo "check-delete race removed a newly owned lock" >&2; exit 1; }

rm -f "$lock" "$count"
other="$tmp/wrong-root"
mkdir -p "$other"
touch "$root/identities/jarvis/telemetry-queue.jsonl.lock"
PATH="$bin:$PATH" A1ID_ROOT="$root" A1_HOME="$other" \
  JARVIS_A1_LOCK_STALE_SEC=0 STUB_PGREP_MODE=none \
  bash "$cleaner" >/dev/null 2>&1
[ ! -e "$root/identities/jarvis/telemetry-queue.jsonl.lock" ] \
  || { echo "A1ID_ROOT was not used" >&2; exit 1; }

# Deterministically enter the old final signature→unlink window.  The rm shim
# starts a real bin/a1id invocation whose a1 stub replaces auth.yaml.lock as
# soon as it passes the shared gate.  It must remain blocked until the cleaner
# has unlinked the stale file and released its exclusive descriptor; the new
# owner's replacement lock must then survive.
fault_bin="$tmp/fault-bin"
mkdir -p "$fault_bin"
entered="$tmp/new-a1-entered"
entered_early="$tmp/new-a1-entered-before-unlink"
replacement_done="$tmp/replacement-done"
spawned_pid="$tmp/new-a1.pid"
cleaner_done="$tmp/cleaner-done"

cat >"$fault_bin/a1" <<'EOF'
#!/usr/bin/env bash
if [ "${FAULT_HOLD_SHARED:-0}" = "1" ]; then
  printf 'holding\n' >"$FAULT_HOLD_ENTERED"
  while [ ! -e "$FAULT_HOLD_RELEASE" ]; do sleep 0.05; done
  exit 0
fi
printf 'entered\n' >"$FAULT_ENTERED"
printf 'new-owner\n' >"$FAULT_LOCK"
printf 'done\n' >"$FAULT_REPLACEMENT_DONE"
EOF
cat >"$fault_bin/rm" <<'EOF'
#!/usr/bin/env bash
if [ "${FAULT_INJECT_ON_RM:-0}" = "1" ]; then
  FAULT_INJECT_ON_RM=0 "$FAULT_A1ID" -- noop \
    >"$FAULT_A1ID_LOG" 2>&1 &
  child=$!
  printf '%s\n' "$child" >"$FAULT_SPAWNED_PID"
  sleep 0.25
  [ ! -e "$FAULT_ENTERED" ] \
    || printf 'entered-too-early\n' >"$FAULT_ENTERED_EARLY"
fi
exec /bin/rm "$@"
EOF
chmod +x "$fault_bin/a1" "$fault_bin/rm"

# An a1id that entered first retains a shared descriptor through the real a1
# stub.  Exclusive cleanup must skip the entire pass and preserve its lock.
hold_entered="$tmp/hold-entered"
hold_release="$tmp/hold-release"
A1ID_ROOT="$root" A1_BIN="$fault_bin/a1" FAULT_HOLD_SHARED=1 \
  FAULT_HOLD_ENTERED="$hold_entered" FAULT_HOLD_RELEASE="$hold_release" \
  "$repo_root/bin/a1id" -- hold >"$tmp/hold.log" 2>&1 &
holder=$!
deadline=$((SECONDS + 5))
while [ ! -e "$hold_entered" ] && [ "$SECONDS" -lt "$deadline" ]; do
  sleep 0.05
done
[ -e "$hold_entered" ] \
  || { echo "shared a1id gate holder did not enter" >&2; exit 1; }
touch "$lock"
active_out="$(
  PATH="$fault_bin:$bin:$PATH" A1ID_ROOT="$root" \
    JARVIS_A1_LOCK_STALE_SEC=0 STUB_PGREP_MODE=none \
    bash "$cleaner" 2>&1
)"
[ -e "$lock" ] \
  || { echo "cleaner removed lock owned by active gated a1id" >&2; exit 1; }
case "$active_out" in
  *"active a1id gate owner detected"*) ;;
  *) echo "cleaner did not report active a1id gate owner" >&2; exit 1 ;;
esac
touch "$hold_release"
wait "$holder"
/bin/rm -f "$lock"

touch "$lock"
PATH="$fault_bin:$bin:$PATH" A1ID_ROOT="$root" \
  A1_BIN="$fault_bin/a1" JARVIS_A1_LOCK_STALE_SEC=0 \
  STUB_PGREP_MODE=none FAULT_INJECT_ON_RM=1 \
  FAULT_A1ID="$repo_root/bin/a1id" FAULT_A1ID_LOG="$tmp/new-a1.log" \
  FAULT_ENTERED="$entered" FAULT_ENTERED_EARLY="$entered_early" \
  FAULT_LOCK="$lock" FAULT_REPLACEMENT_DONE="$replacement_done" \
  FAULT_SPAWNED_PID="$spawned_pid" \
  bash "$cleaner" >/dev/null 2>&1
touch "$cleaner_done"

[ ! -e "$entered_early" ] \
  || { echo "new a1id entered during cleaner's final unlink window" >&2; exit 1; }
[ -s "$spawned_pid" ] \
  || { echo "final-window a1id fault was not triggered" >&2; exit 1; }
deadline=$((SECONDS + 5))
while [ ! -e "$replacement_done" ] && [ "$SECONDS" -lt "$deadline" ]; do
  sleep 0.05
done
[ -e "$replacement_done" ] \
  || { echo "new a1id did not acquire the gate after cleanup" >&2; exit 1; }
[ "$(cat "$lock")" = "new-owner" ] \
  || { echo "cleaner removed the new owner's replacement lock" >&2; exit 1; }

echo "a1 lock cleanup tests passed"
