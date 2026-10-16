#!/usr/bin/env bash

set -euo pipefail

die() {
	echo "basic-check: $*" >&2
	exit 2
}

resolve_commit() {
	local name="$1"
	local value="$2"
	local resolved

	[[ "$value" =~ ^[0-9a-f]{40}$ ]] ||
		die "$name is not a valid commit"
	if ! resolved="$(git rev-parse --verify "${value}^{commit}" 2>/dev/null)"; then
		die "$name is not a valid commit"
	fi
	printf '%s\n' "$resolved"
}

changed_files=()

if [[ "${BASIC_CHECK_FILES_STDIN:-false}" == "true" ]]; then
	while IFS= read -r changed_file; do
		[[ -n "$changed_file" ]] && changed_files+=("$changed_file")
	done
else
	if [[ -z "${DIFF_BASE:-}" || -z "${DIFF_HEAD:-}" ]]; then
		if [[ "${BASIC_CHECK_AUTO_RANGE:-false}" != "true" ]]; then
			die "DIFF_BASE and DIFF_HEAD are required"
		fi

		# Compatibility for the retained Aone workflow. Resolve a real merge
		# base and fail instead of falling back to a single commit when the
		# checkout history is incomplete.
		if [[ "$(git rev-parse --is-shallow-repository)" == "true" ]]; then
			git fetch --no-tags --unshallow origin ||
				die "failed to fetch complete history"
		fi
		base_ref="origin/${GITHUB_BASE_REF:-master}"
		base_tip="$(git rev-parse --verify "${base_ref}^{commit}" 2>/dev/null)" ||
			die "base branch is unavailable: $base_ref"
		head_tip="$(git rev-parse --verify 'HEAD^{commit}' 2>/dev/null)" ||
			die "HEAD is not a valid commit"
		DIFF_BASE="$(git merge-base "$base_tip" "$head_tip")" ||
			die "base branch and HEAD do not have a common ancestor"
		DIFF_HEAD="$head_tip"
	fi

	diff_base="$(resolve_commit DIFF_BASE "$DIFF_BASE")"
	diff_head="$(resolve_commit DIFF_HEAD "$DIFF_HEAD")"
	git merge-base --is-ancestor "$diff_base" "$diff_head" ||
		die "DIFF_BASE must be an ancestor of DIFF_HEAD"

	diff_file="$(mktemp)"
	trap 'rm -f "$diff_file"' EXIT
	if ! git diff --name-only -z --diff-filter=MA \
		"$diff_base" "$diff_head" >"$diff_file"; then
		die "failed to list changed files"
	fi
	while IFS= read -r -d '' changed_file; do
		changed_files+=("$changed_file")
	done <"$diff_file"
fi

if ((${#changed_files[@]} == 0)); then
	echo "==> No files have changed, skipping basic checks"
	exit 0
fi

echo "==> Running basic checks on changed files..."
error=false

for doc in "${changed_files[@]}"; do
	# Skip if file doesn't exist
	[[ ! -f "$doc" ]] && continue

	dirname="$(dirname "$doc")"
	category="$(basename "$dirname")"

	case "$category" in
	"d" | "r")
		# Check for incomplete aliyun.com links
		if grep "https://help.aliyun.com/)\.$" "$doc" >/dev/null 2>&1; then
			echo -e "\033[31mDoc :${doc}: Please input the exact link. Currently it is https://help.aliyun.com/. \033[0m"
			error=true
		fi
		;;
	"alicloud")
		# Check for fmt.Println usage
		if grep "fmt.Println" "$doc" >/dev/null 2>&1; then
			echo -e "\033[31mFile :${doc}: Please Remove the fmt.Println Method! \033[0m"
			error=true
		fi
		;;
	esac
done

if [[ "$error" == "true" ]]; then
	exit 1
fi

echo "==> All basic checks passed!"
exit 0
