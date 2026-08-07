#!/usr/bin/env bash
# Shared machine-level Jarvis runtime configuration loader.
#
# Precedence (highest first): existing process environment, JARVIS_RUNTIME_ENV,
# ~/.config/jarvis/runtime.env, then the main checkout's legacy gitignored
# bootstrap/.env and bridge/jarvis.env.  Worktrees therefore consume one machine
# credential source without copying secrets into a checkout.

_jarvis_runtime_repo_root() {
  local anchor common
  anchor="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." 2>/dev/null && pwd)" || return 1
  common="$(git -C "$anchor" rev-parse --git-common-dir 2>/dev/null)" || {
    printf '%s\n' "$anchor"
    return 0
  }
  case "$common" in
    /*) ;;
    *) common="$anchor/$common" ;;
  esac
  common="$(cd "$common" 2>/dev/null && pwd)" || return 1
  case "$common" in
    */.git) printf '%s\n' "${common%/.git}" ;;
    *) printf '%s\n' "$anchor" ;;
  esac
}

_jarvis_runtime_secure_file() {
  local path="$1" mode
  [ -f "$path" ] || return 1
  mode="$(stat -f '%Lp' "$path" 2>/dev/null || stat -c '%a' "$path" 2>/dev/null || true)"
  case "$mode" in
    ''|*[!0-7]*)
      printf 'runtime-config: cannot verify permissions for %s\n' "$path" >&2
      return 2 ;;
  esac
  if [ $((8#$mode & 077)) -ne 0 ]; then
    printf 'runtime-config: refusing insecure config %s (mode %s; require 0600 or stricter)\n' \
      "$path" "$mode" >&2
    return 2
  fi
  return 0
}

jarvis_load_runtime_config() {
  [ "${JARVIS_RUNTIME_CONFIG_LOADED:-}" = "1" ] && return 0
  local main_root machine_env explicit_env path rc name saved_file
  main_root="$(_jarvis_runtime_repo_root)" || return 2
  machine_env="${XDG_CONFIG_HOME:-${HOME:?HOME is required}/.config}/jarvis/runtime.env"
  explicit_env="${JARVIS_RUNTIME_ENV:-}"
  saved_file="$(mktemp "${TMPDIR:-/tmp}/jarvis-runtime-env.XXXXXX")" || return 2
  # Preserve caller-supplied values. The file contains declarations, never values
  # printed to stdout/stderr, and is removed before returning.
  for name in $(env | sed -n 's/^\(JARVIS_[A-Za-z0-9_]*\)=.*/\1/p'); do
    [ -n "${!name}" ] || continue
    printf 'export %s=%q\n' "$name" "${!name}" >>"$saved_file"
  done
  # Legacy checkout-local files remain a compatibility fallback. Their historic
  # modes vary across installations, so strict permission enforcement applies
  # to the new machine-level/explicit secret sources below.
  for path in \
      "${JARVIS_INTERACTIVE_BOOTSTRAP_ENV:-$main_root/bootstrap/.env}" \
      "${JARVIS_INTERACTIVE_BRIDGE_ENV:-$main_root/bridge/jarvis.env}"; do
    [ -f "$path" ] || continue
    set -a; set +u
    # shellcheck disable=SC1090
    . "$path"
    set -u; set +a
    JARVIS_RUNTIME_CONFIG_SOURCE="$path"
  done
  if [ -f "$machine_env" ]; then
    _jarvis_runtime_secure_file "$machine_env"; rc=$?
    if [ "$rc" -eq 2 ]; then rm -f "$saved_file"; return 2; fi
    set -a; set +u
    # shellcheck disable=SC1090
    . "$machine_env"
    set -u; set +a
    JARVIS_RUNTIME_CONFIG_SOURCE="$machine_env"
  fi
  if [ -n "$explicit_env" ]; then
    _jarvis_runtime_secure_file "$explicit_env"; rc=$?
    if [ "$rc" -ne 0 ]; then rm -f "$saved_file"; return 2; fi
    set -a; set +u
    # shellcheck disable=SC1090
    . "$explicit_env"
    set -u; set +a
    JARVIS_RUNTIME_CONFIG_SOURCE="$explicit_env"
  fi
  # Process environment is the final override.
  # shellcheck disable=SC1090
  . "$saved_file"
  rm -f "$saved_file"
  export JARVIS_RUNTIME_CONFIG_SOURCE="${JARVIS_RUNTIME_CONFIG_SOURCE:-none}"
  export JARVIS_RUNTIME_CONFIG_LOADED=1
  if [ -z "${JARVIS_CONTROL_PLANE_TOKEN:-}" ] && [ -n "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
    export JARVIS_CONTROL_PLANE_TOKEN="$JARVIS_HTML_REPORT_TOKEN"
  fi
  if [ -z "${JARVIS_CONTROL_PLANE_BASE_URL:-}" ] && [ -n "${JARVIS_HTML_REPORT_BASE_URL:-}" ]; then
    export JARVIS_CONTROL_PLANE_BASE_URL="$JARVIS_HTML_REPORT_BASE_URL"
  fi
}

jarvis_runtime_config_diagnose() {
  jarvis_load_runtime_config || return $?
  printf 'runtime_source=%s\n' "${JARVIS_RUNTIME_CONFIG_SOURCE:-none}"
  printf 'control_plane_base=%s\n' "${JARVIS_CONTROL_PLANE_BASE_URL:-default}"
  if [ -n "${JARVIS_CONTROL_PLANE_TOKEN:-}" ]; then
    printf 'control_plane_token=configured\n'
  else
    printf 'control_plane_token=missing\n'
  fi
  if [ -n "${JARVIS_CONTROL_PLANE_ADMIN_TOKEN:-}" ]; then
    printf 'control_plane_admin_token=configured\n'
  else
    printf 'control_plane_admin_token=missing\n'
  fi
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then
  case "${1:-diagnose}" in
    diagnose) jarvis_runtime_config_diagnose ;;
    *) printf 'usage: %s diagnose\n' "$0" >&2; exit 64 ;;
  esac
fi
