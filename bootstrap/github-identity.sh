#!/usr/bin/env bash
# Verify and use the jarvis GitHub identity.

set -euo pipefail

prog="$(basename "$0")"

usage() {
    cat >&2 <<EOF
Usage:
  $prog check
  $prog gh <args...>
  $prog push <owner/repo> <local-ref> <remote-ref>
EOF
}

check_identity() {
    if [ -z "${JARVIS_GITHUB_TOKEN:-}" ]; then
        echo "github-identity: JARVIS_GITHUB_TOKEN is required" >&2
        return 2
    fi

    local expected actual
    expected="api-tool-agent"

    if actual="$(GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh api user --jq .login 2>&1)"; then
        :
    else
        local status=$?
        echo "github-identity: gh api user failed (exit $status)" >&2
        if [ -n "$actual" ]; then
            printf '%s\n' "$actual" >&2
        fi
        return "$status"
    fi

    if [ "$actual" != "$expected" ]; then
        echo "github-identity: GitHub login mismatch: expected '$expected', got '$actual'" >&2
        return 3
    fi

    printf '%s\n' "$actual"
}

run_gh() {
    if check_identity >/dev/null; then
        :
    else
        return $?
    fi

    if GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh "$@"; then
        return 0
    else
        local status=$?
        echo "github-identity: gh command failed (exit $status): gh $*" >&2
        return "$status"
    fi
}

run_push() {
    local repo="${1:-}"
    local local_ref="${2:-}"
    local remote_ref="${3:-}"

    if [ "$#" -ne 3 ] || [[ ! "$repo" =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]]; then
        echo "github-identity: push requires <owner/repo> <local-ref> <remote-ref>" >&2
        return 2
    fi

    if check_identity >/dev/null; then
        :
    else
        return $?
    fi

    local askpass status
    askpass="$(mktemp "${TMPDIR:-/tmp}/jarvis-git-askpass.XXXXXX")"
    cat > "$askpass" <<'EOF'
#!/usr/bin/env bash
case "${1:-}" in
    *Username*) printf '%s\n' "x-access-token" ;;
    *Password*) printf '%s\n' "${JARVIS_GITHUB_TOKEN:?}" ;;
    *) printf '\n' ;;
esac
EOF
    chmod 700 "$askpass"

    if JARVIS_GITHUB_TOKEN="$JARVIS_GITHUB_TOKEN" \
        GIT_ASKPASS="$askpass" \
        GIT_TERMINAL_PROMPT=0 \
            git push "https://github.com/${repo}.git" "${local_ref}:${remote_ref}"; then
        status=0
    else
        status=$?
    fi
    rm -f "$askpass"

    if [ "$status" -ne 0 ]; then
        echo "github-identity: git push failed (exit $status): $repo $local_ref:$remote_ref" >&2
    fi
    return "$status"
}

mode="${1:-}"
case "$mode" in
    check)
        shift
        if [ "$#" -ne 0 ]; then
            usage
            exit 2
        fi
        check_identity
        ;;
    gh)
        shift
        run_gh "$@"
        ;;
    push)
        shift
        run_push "$@"
        ;;
    *)
        usage
        exit 2
        ;;
esac
