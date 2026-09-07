#!/usr/bin/env bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

LINT_TARGET='./alicloud/...'
LINT_PATHSPEC=':(glob)alicloud/**/*.go'

ANALYZERS=(
	-AT001
	-AT001.ignored-filename-prefixes data_source_alicloud_
	-AT005 -AT006 -AT007
	-R001 -R002 -R003 -R004 -R006
	-S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010
	-S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020
	-S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030
	-S031 -S032 -S033
)

die() {
	echo -e "${RED}tfproviderlint: $*${NC}" >&2
	exit 1
}

ensure_tfproviderlint() {
	if ! command -v tfproviderlint >/dev/null 2>&1; then
		echo -e "${YELLOW}tfproviderlint not installed.${NC}" >&2
		exit 1
	fi
}

require_commit_sha() {
	local name="$1"
	local value="$2"
	[[ "$value" =~ ^[0-9a-f]{40}$ ]] ||
		die "$name must be a full 40-character commit SHA"
	git cat-file -e "${value}^{commit}" 2>/dev/null ||
		die "$name is not a valid commit in this checkout"
}

validate_change_range() {
	if [[ -z "${DIFF_BASE:-}" || -z "${DIFF_HEAD:-}" ]]; then
		die "DIFF_BASE and DIFF_HEAD are required"
	fi
	if [[ "$(git rev-parse --is-shallow-repository)" == 'true' ]]; then
		die "checkout is shallow; use fetch-depth: 0"
	fi

	require_commit_sha DIFF_BASE "$DIFF_BASE"
	require_commit_sha DIFF_HEAD "$DIFF_HEAD"

	local checkout_head
	checkout_head="$(git rev-parse --verify 'HEAD^{commit}')"
	[[ "$DIFF_HEAD" == "$checkout_head" ]] ||
		die "DIFF_HEAD does not match HEAD"
	git merge-base --is-ancestor "$DIFF_BASE" "$DIFF_HEAD" ||
		die "DIFF_BASE is not an ancestor of DIFF_HEAD"
}

load_changed_go_files() {
	local changed_list
	changed_list="$(mktemp "${TMPDIR:-/tmp}/tfproviderlint-changed.XXXXXX")"
	if ! git diff --name-only -z --diff-filter=ACMR \
		"$DIFF_BASE" "$DIFF_HEAD" -- "$LINT_PATHSPEC" >"$changed_list"; then
		rm -f "$changed_list"
		die "could not determine changed Go files"
	fi

	# Bash 3.2 treats an empty array expansion as unset under `set -u`, so keep a
	# sentinel that cannot be a repository path.
	CHANGED_GO_FILES=('__NO_CHANGED_GO_FILES__')
	while IFS= read -r -d '' changed_file; do
		CHANGED_GO_FILES+=("$changed_file")
	done <"$changed_list"
	rm -f "$changed_list"
}

is_changed_go_file() {
	local candidate="$1"
	local changed_file
	for changed_file in "${CHANGED_GO_FILES[@]}"; do
		[[ "$candidate" == "$changed_file" ]] && return 0
	done
	return 1
}

fail_with_output() {
	local message="$1"
	local output="$2"
	echo -e "${RED}${message}${NC}"
	printf '%s\n' "$output"
	echo -e "${RED}Fix the above issues before committing.${NC}"
}

run_lint() {
	echo -e "==> ${GREEN}Checking source code against tfproviderlint...${NC}"

	local full_output=''
	local lint_exit=0
	full_output="$(tfproviderlint "${ANALYZERS[@]}" "$LINT_TARGET" 2>&1)" ||
		lint_exit=$?

	if [[ "$lint_exit" -eq 0 ]]; then
		echo -e "${GREEN}tfproviderlint: No issues found.${NC}"
		return 0
	fi

	# golang.org/x/tools/go/analysis uses exit 3 for completed analysis with
	# diagnostics and exit 1 for package loading or analyzer failures. Only a
	# completed diagnostics result is safe to scope to the files in this PR.
	if [[ "$lint_exit" -ne 3 ]]; then
		fail_with_output \
			"tfproviderlint did not complete normally (exit ${lint_exit}):" \
			"$full_output"
		return "$lint_exit"
	fi

	local repo_root
	repo_root="$(git rev-parse --show-toplevel)"
	local line diagnostic_path
	local finding_count=0
	local filtered_output=''
	while IFS= read -r line; do
		if [[ "$line" =~ ^(.+\.go):[0-9]+:[0-9]+:[[:space:]]+[A-Z][A-Z0-9]+: ]]; then
			finding_count=$((finding_count + 1))
			diagnostic_path="${BASH_REMATCH[1]}"
			case "$diagnostic_path" in
			"$repo_root"/*)
				diagnostic_path="${diagnostic_path#"$repo_root"/}"
				;;
			./*)
				diagnostic_path="${diagnostic_path#./}"
				;;
			esac
			if is_changed_go_file "$diagnostic_path"; then
				if [[ -n "$filtered_output" ]]; then
					filtered_output+=$'\n'
				fi
				filtered_output+="$line"
			fi
		fi
	done <<<"$full_output"

	if [[ "$finding_count" -eq 0 ]]; then
		fail_with_output \
			"tfproviderlint returned diagnostics status but no analyzable findings:" \
			"$full_output"
		return "$lint_exit"
	fi

	if [[ -z "$filtered_output" ]]; then
		echo -e "${GREEN}tfproviderlint: No issues in changed files.${NC}"
		echo -e "${YELLOW}Note: pre-existing issues in unchanged files were ignored.${NC}"
		return 0
	fi

	fail_with_output "Found issues in changed files:" "$filtered_output"
	return "$lint_exit"
}

main() {
	ensure_tfproviderlint
	validate_change_range
	load_changed_go_files
	run_lint
}

main "$@"
