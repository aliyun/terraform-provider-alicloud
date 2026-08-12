#!/usr/bin/env bash
# Git transport for Jarvis-owned Code repositories.
#
# The Code private token stays in the isolated jarvis a1 auth.yaml. Git receives
# it only through GIT_ASKPASS; it never appears in argv, a remote URL, or logs.
set -euo pipefail

prog="$(basename "$0")"
auth_root="${A1ID_ROOT:-${HOME:?HOME is required}/.config/a1}"
auth_file="${JARVIS_CODE_AUTH_FILE:-$auth_root/identities/jarvis/auth.yaml}"
python_bin="${JARVIS_CODE_PYTHON:-$(command -v python3 2>/dev/null || true)}"
git_bin="${JARVIS_CODE_GIT_BIN:-$(command -v git 2>/dev/null || true)}"
script_path="$(cd -P "$(dirname "$0")" && pwd)/$(basename "$0")"

die() {
  printf '%s: %s\n' "$prog" "$*" >&2
  exit 2
}

usage() {
  cat <<'EOF'
usage:
  jarvis-code-git.sh check
  jarvis-code-git.sh ls-remote <group/repo> [ref-pattern]
  jarvis-code-git.sh clone <group/repo> <feature-branch> <absolute-destination>
  jarvis-code-git.sh push <repo-dir> <group/repo> <local-ref> <feature-remote-ref>
EOF
}

read_code_auth() {
  local mode="$1"
  [ -n "$python_bin" ] || die "python3 is required"
  [ -f "$auth_file" ] || die "jarvis Code auth is missing: $auth_file"
  "$python_bin" - "$auth_file" "$mode" <<'PY'
import os
import stat
import sys

try:
    import yaml
except Exception:
    raise SystemExit("PyYAML is required to read the isolated a1 auth file")

path, mode = sys.argv[1:]
file_mode = stat.S_IMODE(os.stat(path).st_mode)
if file_mode & 0o077:
    raise SystemExit("jarvis Code auth file must not be group/world accessible")
with open(path, "r", encoding="utf-8") as stream:
    document = yaml.safe_load(stream) or {}
code = ((document.get("platforms") or {}).get("code") or {})
host = str(code.get("host") or "").strip()
auth_type = str(code.get("auth_type") or "").strip()
user = str(code.get("user") or "").strip()
token = str(code.get("token") or "").strip()
if host not in {"code.alibaba-inc.com", "gitlab.alibaba-inc.com"}:
    raise SystemExit("jarvis Code auth host is missing or unsupported")
if auth_type != "private_token":
    raise SystemExit("jarvis Code auth_type must be private_token")
if not user or not token:
    raise SystemExit("jarvis Code private token or user is missing")
if mode == "info":
    print(f"{host}\t{user}")
elif mode == "user":
    print(user)
elif mode == "token":
    print(token)
elif mode == "validate":
    pass
else:
    raise SystemExit("unsupported auth read mode")
PY
}

# Git invokes this same executable as GIT_ASKPASS. Do not add diagnostics to
# this path: stdout is the credential response channel.
if [ "${JARVIS_CODE_ASKPASS_MODE:-0}" = "1" ]; then
  case "${1:-}" in
    *Username*) read_code_auth user ;;
    *Password*) read_code_auth token ;;
    *) printf '\n' ;;
  esac
  exit 0
fi

validate_repo_path() {
  local repo="$1"
  [[ "$repo" =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] \
    || die "repository must be a group/repo path"
}

validate_feature_ref() {
  local ref="${1#refs/heads/}"
  case "$ref" in
    feature/*) ;;
    *) die "only feature/* branches are allowed" ;;
  esac
  case "$ref" in
    *'..'*|*'@{'*|*'//'*) die "invalid feature branch: $ref" ;;
  esac
  case "$ref" in
    *' '*|*'~'*|*'^'*|*':'*|*'?'*|*'*'*|*'['*|*'\\'*)
      die "invalid feature branch: $ref"
      ;;
  esac
}

auth_info="$(read_code_auth info)" || exit $?
IFS=$'\t' read -r code_host code_user <<<"$auth_info"
[ -n "$git_bin" ] || die "git is required"

run_git() {
  JARVIS_CODE_ASKPASS_MODE=1 \
  JARVIS_CODE_AUTH_FILE="$auth_file" \
  JARVIS_CODE_PYTHON="$python_bin" \
  GIT_ASKPASS="$script_path" \
  GIT_TERMINAL_PROMPT=0 \
    "$git_bin" "$@"
}

mode="${1:-}"
shift || true
case "$mode" in
  check)
    [ "$#" -eq 0 ] || die "check takes no arguments"
    printf 'jarvis-code-git: ready host=%s user=%s auth=private_token\n' \
      "$code_host" "$code_user"
    ;;
  ls-remote)
    [ "$#" -ge 1 ] && [ "$#" -le 2 ] || die "ls-remote requires <group/repo> [ref-pattern]"
    validate_repo_path "$1"
    repo_url="https://$code_host/$1.git"
    if [ "$#" -eq 2 ]; then
      run_git ls-remote "$repo_url" "$2"
    else
      run_git ls-remote "$repo_url"
    fi
    ;;
  clone)
    [ "$#" -eq 3 ] || die "clone requires <group/repo> <feature-branch> <absolute-destination>"
    validate_repo_path "$1"
    validate_feature_ref "$2"
    [[ "$3" = /* ]] || die "clone destination must be absolute"
    [ ! -e "$3" ] || die "clone destination already exists: $3"
    repo_url="https://$code_host/$1.git"
    run_git clone --branch "${2#refs/heads/}" --single-branch "$repo_url" "$3"
    ;;
  push)
    [ "$#" -eq 4 ] || die "push requires <repo-dir> <group/repo> <local-ref> <feature-remote-ref>"
    [ -d "$1/.git" ] || die "not a Git repository: $1"
    validate_repo_path "$2"
    validate_feature_ref "$4"
    repo_url="https://$code_host/$2.git"
    run_git -C "$1" push "$repo_url" "$3:refs/heads/${4#refs/heads/}"
    ;;
  ''|-h|--help|help)
    usage
    ;;
  *)
    usage >&2
    die "unknown mode: $mode"
    ;;
esac
