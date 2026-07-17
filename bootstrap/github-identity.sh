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

# Verify the token belongs to the api-tool-agent account. Exit codes:
#   0  identity verified (login == api-tool-agent)
#   0  token present but GitHub transiently unreachable after retries — identity
#      could NOT be confirmed, but the pinned JARVIS_GITHUB_TOKEN is still used
#      for the real op (no ambient fallback), so proceed rather than mask a
#      transient GitHub 5xx/HTML/network blip as an identity failure (emits a
#      loud WARNING).
#   2  token missing (hard fail)
#   3  token rejected (HTTP 401 / bad credentials) or login mismatch (hard fail)
#
# The probe distinguishes a *definitive* auth rejection (401 / Bad credentials)
# from a *transient* failure (5xx, network error, or a non-JSON HTML error page
# that makes `--jq` emit "invalid character '<'"). Only the former is a real
# identity failure; the latter is retried, then downgraded to a non-blocking
# warning. Retries/backoff tunable via JARVIS_GHID_CHECK_RETRIES (default 3).
check_identity() {
    if [ -z "${JARVIS_GITHUB_TOKEN:-}" ]; then
        echo "github-identity: JARVIS_GITHUB_TOKEN is required" >&2
        return 2
    fi

    local expected="api-tool-agent"
    local attempts="${JARVIS_GHID_CHECK_RETRIES:-3}"
    local i out status
    for (( i=1; i<=attempts; i++ )); do
        # `|| status=$?` both captures the real probe exit status AND shields the
        # failing command substitution from `set -e` (a bare assignment would
        # abort the script). status stays 0 on the success path.
        status=0
        out="$(GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh api user --jq .login 2>&1)" || status=$?
        if [ "$status" -eq 0 ]; then
            if [ "$out" != "$expected" ]; then
                echo "github-identity: GitHub login mismatch: expected '$expected', got '$out'" >&2
                return 3
            fi
            printf '%s\n' "$out"
            return 0
        fi

        # Definitive auth rejection → real identity failure, do not retry.
        if printf '%s' "$out" | grep -qiE 'Bad credentials|HTTP 401'; then
            echo "github-identity: token rejected (HTTP 401 / bad credentials)" >&2
            printf '%s\n' "$out" >&2
            return 3
        fi

        # Transient: 5xx / network / non-JSON HTML page (jq 'invalid character').
        echo "github-identity: gh api user probe failed (attempt $i/$attempts, exit $status): $(printf '%s' "$out" | head -1)" >&2
        if [ "$i" -lt "$attempts" ]; then
            sleep "$(( i * 2 ))"
        fi
    done

    # Exhausted retries without a definitive rejection → treat as transient
    # unavailability, not an identity failure. The op still runs under the pinned
    # token; if GitHub stays down the op surfaces its own error for the caller's
    # existing retry (e.g. PrWatch "keep watching") to handle.
    echo "github-identity: WARNING: could not verify GitHub identity after $attempts probes (GitHub 5xx/HTML/network); proceeding with pinned token — identity unconfirmed." >&2
    return 0
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
