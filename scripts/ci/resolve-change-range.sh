#!/usr/bin/env bash

set -euo pipefail

die() {
	echo "change-range: $*" >&2
	exit 1
}

require_commit_sha() {
	local name="$1"
	local value="$2"
	[[ "$value" =~ ^[0-9a-f]{40}$ ]] ||
		die "$name must be a full 40-character commit SHA"
}

fetch_base_branch() {
	local repository_url="$1"
	local base_ref="$2"
	local destination_ref="$3"

	[[ -n "$repository_url" && "$repository_url" != -* ]] ||
		die "CI_BASE_REPOSITORY_URL is invalid"
	git check-ref-format "refs/heads/$base_ref" >/dev/null 2>&1 ||
		die "base ref is invalid: $base_ref"

	git fetch --no-tags --no-recurse-submodules --force \
		"$repository_url" \
		"+refs/heads/$base_ref:$destination_ref"
}

head_sha="$(git rev-parse --verify 'HEAD^{commit}')"
require_commit_sha HEAD "$head_sha"

case "${GITHUB_EVENT_NAME:-}" in
pull_request)
	base_sha="${PR_BASE_SHA:-}"
	base_ref="${PR_BASE_REF:-}"
	base_repository_url="${CI_BASE_REPOSITORY_URL:-}"
	require_commit_sha PR_BASE_SHA "$base_sha"
	[[ -n "$base_ref" ]] || die "PR_BASE_REF is required"

	if [[ "$(git rev-parse --is-shallow-repository)" == "true" ]]; then
		die "pull request checkout is shallow; use fetch-depth: 0"
	fi

	fetched_base_ref="refs/remotes/ci-base/base"
	fetch_base_branch \
		"$base_repository_url" \
		"$base_ref" \
		"$fetched_base_ref"

	git cat-file -e "${base_sha}^{commit}" 2>/dev/null ||
		die "PR base commit is unavailable after fetching the base branch"
	git merge-base --is-ancestor "$base_sha" "$fetched_base_ref" ||
		die "PR base commit is not part of the fetched base branch"

	if ! diff_base="$(git merge-base "$base_sha" "$head_sha")"; then
		die "PR base and head do not have a common ancestor"
	fi
	[[ -n "$diff_base" ]] ||
		die "PR base and head do not have a common ancestor"
	require_commit_sha diff_base "$diff_base"
	;;
push)
	diff_base="${PUSH_BASE_SHA:-}"
	if [[ -z "$diff_base" ||
		"$diff_base" == "0000000000000000000000000000000000000000" ]]; then
		diff_base="$(git rev-parse --verify 'HEAD^{commit}^')"
	fi
	require_commit_sha PUSH_BASE_SHA "$diff_base"

	if ! git cat-file -e "${diff_base}^{commit}" 2>/dev/null; then
		base_repository_url="${CI_BASE_REPOSITORY_URL:-}"
		base_ref="${CI_BASE_REF:-}"
		[[ -n "$base_ref" ]] ||
			die "CI_BASE_REF is required when the push base commit is unavailable"
		fetch_base_branch \
			"$base_repository_url" \
			"$base_ref" \
			"refs/remotes/ci-base/base"
	fi
	git cat-file -e "${diff_base}^{commit}" 2>/dev/null ||
		die "push base commit is unavailable"
	;;
*)
	die "unsupported event: ${GITHUB_EVENT_NAME:-<empty>}"
	;;
esac

printf 'diff_base=%s\n' "$diff_base"
printf 'diff_head=%s\n' "$head_sha"
