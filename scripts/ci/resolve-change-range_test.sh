#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
helper="$script_dir/resolve-change-range.sh"
basic_check="$repo_root/scripts/basic-check.sh"

die() {
	echo "FAIL: $*" >&2
	exit 1
}

assert_eq() {
	local expected="$1"
	local actual="$2"
	local message="$3"
	if [[ "$actual" != "$expected" ]]; then
		printf 'FAIL: %s\nexpected:\n%s\nactual:\n%s\n' \
			"$message" "$expected" "$actual" >&2
		exit 1
	fi
}

output_value() {
	local key="$1"
	local output="$2"
	local value
	value="$(sed -n "s/^${key}=//p" <<<"$output")"
	[[ -n "$value" ]] || die "missing $key in helper output"
	[[ "$(grep -c "^${key}=" <<<"$output")" -eq 1 ]] ||
		die "duplicate $key in helper output"
	printf '%s\n' "$value"
}

commit_file() {
	local repository="$1"
	local path="$2"
	local content="$3"
	local subject="$4"
	mkdir -p "$(dirname "$repository/$path")"
	printf '%s\n' "$content" >"$repository/$path"
	git -C "$repository" add "$path"
	git -C "$repository" commit -q -m "$subject"
}

test_diverged_multi_commit_pr() {
	local temp_dir="$1"
	local source_repo="$temp_dir/source"
	local checkout_repo="$temp_dir/checkout"

	git init -q -b master "$source_repo"
	git -C "$source_repo" config core.hooksPath /dev/null
	git -C "$source_repo" config user.name "CI Range Test"
	git -C "$source_repo" config user.email "ci-range-test@example.invalid"

	commit_file "$source_repo" alicloud/master_only.go \
		$'package alicloud\n\nfunc masterOnly() { fmt.Println("base-only") }' \
		common
	local common_sha
	common_sha="$(git -C "$source_repo" rev-parse HEAD)"

	git -C "$source_repo" switch -q -c feature
	commit_file "$source_repo" alicloud/pr_one.go \
		$'package alicloud\n\nfunc prOne() { fmt.Println("pr-one") }' \
		pr-one
	local pr_one_sha
	pr_one_sha="$(git -C "$source_repo" rev-parse HEAD)"
	commit_file "$source_repo" website/docs/r/pr_two.html.markdown \
		'See [documentation](https://help.aliyun.com/).' \
		pr-two
	local pr_head_sha
	pr_head_sha="$(git -C "$source_repo" rev-parse HEAD)"

	git -C "$source_repo" switch -q master
	commit_file "$source_repo" alicloud/master_only.go \
		$'package alicloud\n\nfunc masterOnly() {}' \
		master-only
	local base_tip_sha
	base_tip_sha="$(git -C "$source_repo" rev-parse HEAD)"

	git -c core.hooksPath=/dev/null clone -q \
		--branch feature "$source_repo" "$checkout_repo"
	local output
	output="$(
		cd "$checkout_repo"
		GITHUB_EVENT_NAME=pull_request \
			PR_BASE_SHA="$base_tip_sha" \
			PR_BASE_REF=master \
			CI_BASE_REPOSITORY_URL="$source_repo" \
			"$helper"
	)"

	local diff_base diff_head
	diff_base="$(output_value diff_base "$output")"
	diff_head="$(output_value diff_head "$output")"
	assert_eq "$common_sha" "$diff_base" \
		"PR range must start at the merge base"
	assert_eq "$pr_head_sha" "$diff_head" \
		"PR range must end at the checked-out head"

	local subjects
	subjects="$(git -C "$checkout_repo" log --reverse --format=%s \
		"$diff_base..$diff_head")"
	assert_eq $'pr-one\npr-two' "$subjects" \
		"PR range must retain every PR commit and exclude master-only commits"

	local changed_files
	changed_files="$(git -C "$checkout_repo" diff --name-only \
		"$diff_base" "$diff_head")"
	assert_eq \
		$'alicloud/pr_one.go\nwebsite/docs/r/pr_two.html.markdown' \
		"$changed_files" \
		"PR range must contain only files changed by the PR"

	local basic_output basic_status
	set +e
	basic_output="$(
		cd "$checkout_repo"
		GITHUB_ACTIONS=true \
			DIFF_BASE="$diff_base" \
			DIFF_HEAD="$diff_head" \
			"$basic_check" 2>&1
	)"
	basic_status=$?
	set -e
	[[ "$basic_status" -eq 1 ]] ||
		die "basic-check must report the two intentional PR violations"
	[[ "$basic_output" == *"alicloud/pr_one.go"* ]] ||
		die "basic-check must include the first PR commit"
	[[ "$basic_output" == *"website/docs/r/pr_two.html.markdown"* ]] ||
		die "basic-check must include the second PR commit"
	[[ "$basic_output" != *"alicloud/master_only.go"* ]] ||
		die "basic-check must exclude the master-only change"

	local legacy_output legacy_status
	set +e
	legacy_output="$(
		cd "$checkout_repo"
		env -u GITHUB_ACTIONS -u DIFF_BASE -u DIFF_HEAD \
			BASIC_CHECK_AUTO_RANGE=true "$basic_check" 2>&1
	)"
	legacy_status=$?
	set -e
	[[ "$legacy_status" -eq 1 &&
		"$legacy_output" == *"alicloud/pr_one.go"* &&
		"$legacy_output" == *"website/docs/r/pr_two.html.markdown"* &&
		"$legacy_output" != *"alicloud/master_only.go"* ]] ||
		die "non-GitHub compatibility mode must use the same merge-base scope"

	local stdin_output stdin_status
	set +e
	stdin_output="$(
		cd "$checkout_repo"
		printf '%s\n' \
			alicloud/pr_one.go \
			website/docs/r/pr_two.html.markdown |
			BASIC_CHECK_FILES_STDIN=true "$basic_check" 2>&1
	)"
	stdin_status=$?
	set -e
	[[ "$stdin_status" -eq 1 &&
		"$stdin_output" == *"alicloud/pr_one.go"* &&
		"$stdin_output" == *"website/docs/r/pr_two.html.markdown"* ]] ||
		die "local file-list compatibility mode must check every supplied file"

	local missing_output invalid_output missing_status invalid_status
	set +e
	missing_output="$(
		cd "$checkout_repo"
		GITHUB_ACTIONS=true \
			env -u DIFF_BASE -u DIFF_HEAD "$basic_check" 2>&1
	)"
	missing_status=$?
	invalid_output="$(
		cd "$checkout_repo"
		GITHUB_ACTIONS=true \
			DIFF_BASE=not-a-commit \
			DIFF_HEAD="$diff_head" \
			"$basic_check" 2>&1
	)"
	invalid_status=$?
	set -e
	[[ "$missing_status" -ne 0 &&
		"$missing_output" == *"DIFF_BASE and DIFF_HEAD are required"* ]] ||
		die "basic-check must fail closed when the range is missing"
	[[ "$invalid_status" -ne 0 &&
		"$invalid_output" == *"DIFF_BASE is not a valid commit"* ]] ||
		die "basic-check must fail closed when the range is invalid"

	local push_output
	push_output="$(
		cd "$checkout_repo"
		GITHUB_EVENT_NAME=push \
			PUSH_BASE_SHA="$pr_one_sha" \
			"$helper"
	)"
	assert_eq "$pr_one_sha" "$(output_value diff_base "$push_output")" \
		"push range must preserve before-to-head semantics"
	assert_eq "$pr_head_sha" "$(output_value diff_head "$push_output")" \
		"push range must end at the checked-out head"
}

test_no_common_ancestor_fails_closed() {
	local temp_dir="$1"
	local head_repo="$temp_dir/unrelated-head"
	local base_repo="$temp_dir/unrelated-base"

	git init -q -b feature "$head_repo"
	git -C "$head_repo" config core.hooksPath /dev/null
	git -C "$head_repo" config user.name "CI Range Test"
	git -C "$head_repo" config user.email "ci-range-test@example.invalid"
	commit_file "$head_repo" head.txt head head

	git init -q -b master "$base_repo"
	git -C "$base_repo" config core.hooksPath /dev/null
	git -C "$base_repo" config user.name "CI Range Test"
	git -C "$base_repo" config user.email "ci-range-test@example.invalid"
	commit_file "$base_repo" base.txt base base
	local base_sha
	base_sha="$(git -C "$base_repo" rev-parse HEAD)"

	if (
		cd "$head_repo"
		GITHUB_EVENT_NAME=pull_request \
			PR_BASE_SHA="$base_sha" \
			PR_BASE_REF=master \
			CI_BASE_REPOSITORY_URL="$base_repo" \
			"$helper"
	) >/dev/null 2>&1; then
		die "PR histories without a common ancestor must fail closed"
	fi
}

main() {
	[[ -x "$helper" ]] || die "missing executable helper: $helper"

	local temp_dir
	temp_dir="$(mktemp -d)"
	TEST_TEMP_DIR="$temp_dir"
	trap 'rm -rf "$TEST_TEMP_DIR"' EXIT

	test_diverged_multi_commit_pr "$temp_dir"
	test_no_common_ancestor_fails_closed "$temp_dir"
	echo "PASS: change-range regression tests"
}

main "$@"
