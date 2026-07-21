#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
mkdir -p "$TMP/config/jarvis"
RUNTIME="$TMP/config/jarvis/runtime.env"

printf '%s\n' \
  'JARVIS_CONTROL_PLANE_BASE_URL=https://machine.example' \
  'JARVIS_CONTROL_PLANE_TOKEN=machine-secret-value' >"$RUNTIME"
chmod 600 "$RUNTIME"

out="$(env -u JARVIS_CONTROL_PLANE_BASE_URL -u JARVIS_CONTROL_PLANE_TOKEN \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash "$ROOT/bootstrap/runtime-config.sh" diagnose)"
grep -q "runtime_source=$RUNTIME" <<<"$out"
grep -q 'control_plane_base=https://machine.example' <<<"$out"
grep -q 'control_plane_token=configured' <<<"$out"
! grep -q 'machine-secret-value' <<<"$out"

out="$(JARVIS_CONTROL_PLANE_BASE_URL=https://explicit.example \
  JARVIS_CONTROL_PLANE_TOKEN=explicit-secret \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash "$ROOT/bootstrap/runtime-config.sh" diagnose)"
grep -q 'control_plane_base=https://explicit.example' <<<"$out"
! grep -q 'explicit-secret\|machine-secret-value' <<<"$out"

chmod 644 "$RUNTIME"
if env -u JARVIS_CONTROL_PLANE_TOKEN XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
    bash "$ROOT/bootstrap/runtime-config.sh" diagnose >"$TMP/out" 2>"$TMP/err"; then
  echo "expected insecure runtime config to be rejected" >&2
  exit 1
fi
grep -q 'refusing insecure config' "$TMP/err"
! grep -q 'machine-secret-value' "$TMP/out" "$TMP/err"

echo "runtime_config_test: PASS"
