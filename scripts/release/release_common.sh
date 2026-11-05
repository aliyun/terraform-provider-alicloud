#!/usr/bin/env bash
# Copyright (c) Alibaba, Inc.
# SPDX-License-Identifier: MPL-2.0
# Shared helpers for the release pipeline scripts.
#
# Sourced by generate_changelog.py callers and the release-pipeline workflow
# steps. Provides version validation, last-release-tag resolution, comparison
# base computation and change-range detection that mirrors the contract of
# scripts/ci/resolve-change-range.sh but operates against the full release
# range (last release tag -> HEAD) rather than a single PR/push event.
#
# Usage:
#   source scripts/release/release_common.sh
#   rc_validate_version "2.2.0"
#   rc_last_release_tag           # prints e.g. v2.1.0
#   rc_resolve_tag_sha "v2.1.0"   # prints full 40-char commit SHA

set -euo pipefail

# ---------------------------------------------------------------------------
# Error handling
# ---------------------------------------------------------------------------

rc_die() {
	echo "release-common: $*" >&2
	exit 1
}

rc_require_commit_sha() {
	local name="$1"
	local value="$2"
	[[ "$value" =~ ^[0-9a-f]{40}$ ]] ||
		rc_die "$name must be a full 40-character commit SHA, got: ${value:-<empty>}"
}

# ---------------------------------------------------------------------------
# Version validation
# ---------------------------------------------------------------------------

# Match the release.yml tag filter: v2.X.Y with optional -betaN / -rcN.
# Accepts both bare (2.2.0) and tagged (v2.2.0) forms.
RC_VERSION_RE='^v?[0-9]+\.[0-9]+\.[0-9]+(-(beta|rc)[1-9][0-9]*)?$'

rc_validate_version() {
	local version="$1"
	[[ "$version" =~ $RC_VERSION_RE ]] ||
		rc_die "invalid version '$version' (expected v2.X.Y or 2.X.Y, optional -betaN/-rcN)"
}

# Strip a leading 'v' so callers can pass either form and still get the
# bare numeric version (e.g. for goreleaser .Version and CHANGELOG headers).
rc_strip_v() {
	local version="$1"
	printf '%s' "${version#v}"
}

# Tag form always carries the leading 'v' (release.yml triggers on v2.*).
rc_tag_form() {
	local version="$1"
	rc_validate_version "$version"
	local stripped="${version#v}"
	printf 'v%s' "$stripped"
}

# ---------------------------------------------------------------------------
# Last release tag resolution
# ---------------------------------------------------------------------------

# Print the most recent release tag (v*) by creation date. Returns non-zero
# when no release tag exists, which is expected for first-time releases.
rc_last_release_tag() {
	local tags
	tags="$(git tag --list 'v*' --sort=-creatordate 2>/dev/null || true)"
	[[ -n "$tags" ]] || return 1
	printf '%s\n' "$tags" | head -n1
}

# Resolve any tag/ref to a full 40-char commit SHA. Verifies the object is a
# commit (annotated tags are peeled via ^{commit}, lightweight tags directly).
rc_resolve_tag_sha() {
	local ref="$1"
	[[ -n "$ref" ]] || rc_die "rc_resolve_tag_sha: ref is required"
	local sha
	sha="$(git rev-parse --verify --quiet "${ref}^{commit}" || true)"
	[[ -n "$sha" ]] || rc_die "rc_resolve_tag_sha: ref '$ref' does not resolve to a commit"
	rc_require_commit_sha "$ref" "$sha"
	printf '%s\n' "$sha"
}

# ---------------------------------------------------------------------------
# Comparison base for the release range
# ---------------------------------------------------------------------------

# Compute the merge-base of the last release tag and HEAD. This is the
# CI_CHECK_BASE used by `make ci-check-quick` and the --from anchor for
# generate_changelog.py. When no prior tag exists, falls back to the first
# commit so a maiden release still covers the whole history.
rc_release_diff_base() {
	local from_ref="${1:-}"
	if [[ -z "$from_ref" ]]; then
		from_ref="$(rc_last_release_tag 2>/dev/null || true)"
	fi
	if [[ -z "$from_ref" ]]; then
		# No prior release tag: use the root commit as the base so the whole
		# history is treated as the change range.
		local root
		root="$(git rev-list --max-parents=0 HEAD 2>/dev/null | head -n1 || true)"
		[[ -n "$root" ]] || rc_die "rc_release_diff_base: cannot determine a base commit"
		rc_require_commit_sha "root" "$root"
		printf '%s\n' "$root"
		return
	fi
	local base
	base="$(git merge-base "${from_ref}" HEAD 2>/dev/null || true)"
	[[ -n "$base" ]] || rc_die "rc_release_diff_base: no common ancestor between '$from_ref' and HEAD"
	rc_require_commit_sha "diff-base" "$base"
	printf '%s\n' "$base"
}

# ---------------------------------------------------------------------------
# Change detection (release-range variant)
# ---------------------------------------------------------------------------

# Print the files changed since the release base, one path per line. Mirrors
# the diff semantics of resolve-change-range.sh (three-point merge-base diff)
# so additions, modifications, deletions and renames are all covered.
rc_changed_files() {
	local base="${1:-}"
	[[ -n "$base" ]] || base="$(rc_release_diff_base)"
	rc_require_commit_sha "base" "$base"
	git diff --no-renames --name-only --diff-filter=AMRC "${base}"..HEAD 2>/dev/null || true
}

# Detect whether any module example/website test is affected by the release
# range. Returns 0 (affected) / 1 (not affected) and prints a summary line.
# The full module test execution is orchestrated by the workflow; this helper
# only answers "is there module work to do?" so the job can be skipped when
# nothing changed under the watched paths.
rc_module_tests_affected() {
	local base="${1:-}"
	[[ -n "$base" ]] || base="$(rc_release_diff_base)"
	local paths
	paths="$(rc_changed_files "$base" | grep -E '^(examples/|modules/|website/)' || true)"
	if [[ -z "$paths" ]]; then
		echo "no module/example/website changes detected"
		return 1
	fi
	echo "module/example/website changes detected:"
	printf '%s\n' "$paths"
	return 0
}
