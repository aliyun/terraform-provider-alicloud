#!/usr/bin/env bash
# Verify and use the jarvis GitHub identity.

set -euo pipefail

prog="$(basename "$0")"

# Commit identity for jarvis-authored contributions. CLA-assistant keys on the
# COMMIT AUTHOR email (not the push token / PR opener), so commits must be
# authored as the CLA-signed api-tool-agent identity or the license/cla check
# fails. Overridable via env for other identities.
: "${JARVIS_GIT_AUTHOR_NAME:=api-tool-agent}"
: "${JARVIS_GIT_AUTHOR_EMAIL:=cloudspec_bot@alibaba-inc.com}"

usage() {
    cat >&2 <<EOF
Usage:
  $prog check
  $prog gh <args...>
  $prog commit <git-commit-args...>
  $prog push <owner/repo> <local-ref> <remote-ref>
EOF
}

# Commit with the CLA-signed author+committer identity. Pass through any git
# commit args (-m / --amend / --no-edit / -F ...). Local op, no token needed.
run_commit() {
    GIT_AUTHOR_NAME="$JARVIS_GIT_AUTHOR_NAME" \
    GIT_AUTHOR_EMAIL="$JARVIS_GIT_AUTHOR_EMAIL" \
    GIT_COMMITTER_NAME="$JARVIS_GIT_AUTHOR_NAME" \
    GIT_COMMITTER_EMAIL="$JARVIS_GIT_AUTHOR_EMAIL" \
        git commit --author="$JARVIS_GIT_AUTHOR_NAME <$JARVIS_GIT_AUTHOR_EMAIL>" "$@"
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

    # Backstop: CLA-assistant keys on the commit author email. Warn (do not
    # block) if the tip commit of the ref being pushed is not authored by the
    # CLA-signed identity, so license/cla failures are caught before review.
    local tip_ref tip_email
    tip_ref="${local_ref#+}"        # strip force-push '+' prefix (e.g. +HEAD)
    tip_email="$(git log -1 --format='%ae' "$tip_ref" 2>/dev/null || true)"
    if [ -n "$tip_email" ] && [ "$tip_email" != "$JARVIS_GIT_AUTHOR_EMAIL" ]; then
        echo "github-identity: WARNING: tip commit author <$tip_email> != <$JARVIS_GIT_AUTHOR_EMAIL>." >&2
        echo "github-identity: CLA-assistant keys on the commit author; license/cla will likely fail." >&2
        echo "github-identity: re-author before pushing, e.g.:" >&2
        echo "  $prog commit --amend --no-edit    # or: git commit --amend --author=\"$JARVIS_GIT_AUTHOR_NAME <$JARVIS_GIT_AUTHOR_EMAIL>\" --no-edit" >&2
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
    commit)
        shift
        run_commit "$@"
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
