#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
mkdir -p "$TMP/config/jarvis"
RUNTIME="$TMP/config/jarvis/runtime.env"

printf '%s\n' \
  'JARVIS_CONTROL_PLANE_BASE_URL=https://machine.example' \
  'JARVIS_CONTROL_PLANE_TOKEN=machine-secret-value' \
  'JARVIS_CONTROL_PLANE_ADMIN_TOKEN=machine-admin-secret-value' >"$RUNTIME"
chmod 600 "$RUNTIME"

out="$(env -u JARVIS_CONTROL_PLANE_BASE_URL -u JARVIS_CONTROL_PLANE_TOKEN \
  -u JARVIS_CONTROL_PLANE_ADMIN_TOKEN \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash "$ROOT/bootstrap/runtime-config.sh" diagnose)"
grep -q "runtime_source=$RUNTIME" <<<"$out"
grep -q 'control_plane_base=https://machine.example' <<<"$out"
grep -q 'control_plane_token=configured' <<<"$out"
grep -q 'control_plane_admin_token=configured' <<<"$out"
! grep -q 'machine-secret-value\|machine-admin-secret-value' <<<"$out"

printf '%s\n' \
  'JARVIS_CONTROL_PLANE_BASE_URL=https://machine.example' \
  'JARVIS_CONTROL_PLANE_TOKEN=machine-secret-value' >"$RUNTIME"
out="$(env -u JARVIS_CONTROL_PLANE_BASE_URL -u JARVIS_CONTROL_PLANE_TOKEN \
  -u JARVIS_CONTROL_PLANE_ADMIN_TOKEN \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash "$ROOT/bootstrap/runtime-config.sh" diagnose)"
grep -q 'control_plane_admin_token=missing' <<<"$out"

out="$(JARVIS_CONTROL_PLANE_BASE_URL=https://explicit.example \
  JARVIS_CONTROL_PLANE_TOKEN=explicit-secret \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash "$ROOT/bootstrap/runtime-config.sh" diagnose)"
grep -q 'control_plane_base=https://explicit.example' <<<"$out"
! grep -q 'explicit-secret\|machine-secret-value\|machine-admin-secret-value' <<<"$out"

# Caller-supplied AutoWonder credentials have the same highest precedence as
# JARVIS_* values. A blank/stale env file must not overwrite the launch token.
printf '%s\n' \
  'JARVIS_CONTROL_PLANE_BASE_URL=https://machine.example' \
  'AUTOWONDER_MCP_TOKEN=machine-mcp-secret' >"$RUNTIME"
AUTOWONDER_MCP_TOKEN=explicit-mcp-secret \
  JARVIS_INTERACTIVE_BOOTSTRAP_ENV="$TMP/missing-bootstrap.env" \
  JARVIS_INTERACTIVE_BRIDGE_ENV="$TMP/missing-bridge.env" \
  XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
  bash -c 'source "$1"; jarvis_load_runtime_config; \
    test "$AUTOWONDER_MCP_TOKEN" = explicit-mcp-secret' \
  _ "$ROOT/bootstrap/runtime-config.sh"

chmod 644 "$RUNTIME"
if env -u JARVIS_CONTROL_PLANE_TOKEN -u JARVIS_CONTROL_PLANE_ADMIN_TOKEN \
    XDG_CONFIG_HOME="$TMP/config" HOME="$TMP" \
    bash "$ROOT/bootstrap/runtime-config.sh" diagnose >"$TMP/out" 2>"$TMP/err"; then
  echo "expected insecure runtime config to be rejected" >&2
  exit 1
fi
grep -q 'refusing insecure config' "$TMP/err"
! grep -q 'machine-secret-value\|machine-admin-secret-value' "$TMP/out" "$TMP/err"

echo "runtime_config_test: PASS"
