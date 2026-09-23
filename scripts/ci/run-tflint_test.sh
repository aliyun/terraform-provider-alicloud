#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
runner="$repo_root/scripts/run-tflint.sh"

die() {
	echo "FAIL: $*" >&2
	exit 1
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

run_lint() {
	local repository="$1"
	shift
	set +e
	RUN_OUTPUT="$(
		cd "$repository"
		env PATH="$MOCK_BIN:$PATH" "$@" "$runner" 2>&1
	)"
	RUN_STATUS=$?
	set -e
}

assert_contains() {
	local haystack="$1"
	local needle="$2"
	local message="$3"
	[[ "$haystack" == *"$needle"* ]] ||
		die "$message (missing: $needle)"
}

assert_not_contains() {
	local haystack="$1"
	local needle="$2"
	local message="$3"
	[[ "$haystack" != *"$needle"* ]] ||
		die "$message (unexpected: $needle)"
}

test_range_validation_and_literal_filtering() {
	local temp_dir="$1"
	local repository="$temp_dir/repository"

	git init -q -b master "$repository"
	git -C "$repository" config core.hooksPath /dev/null
	git -C "$repository" config user.name "CI Lint Test"
	git -C "$repository" config user.email "ci-lint-test@example.invalid"

	commit_file "$repository" alicloud/changed/same.go \
		'package changed' base
	commit_file "$repository" alicloud/unchanged/same.go \
		'package unchanged' base-same-name
	commit_file "$repository" 'alicloud/resource_name[1].go' \
		'package alicloud' base-regex-name
	commit_file "$repository" alicloud/resource_name1.go \
		'package alicloud' base-regex-lookalike
	local base_sha
	base_sha="$(git -C "$repository" rev-parse HEAD)"

	git -C "$repository" switch -q -c feature
	commit_file "$repository" alicloud/changed/same.go \
		$'package changed\n\nvar changed = true' change-same-name
	commit_file "$repository" 'alicloud/resource_name[1].go' \
		$'package alicloud\n\nvar changed = true' change-regex-name
	local head_sha
	head_sha="$(git -C "$repository" rev-parse HEAD)"

	run_lint "$repository" \
		env -u DIFF_BASE -u DIFF_HEAD \
		MOCK_LINT_EXIT=0 MOCK_LINT_OUTPUT=
	[[ "$RUN_STATUS" -ne 0 ]] ||
		die "missing range must fail closed"
	assert_contains "$RUN_OUTPUT" "DIFF_BASE and DIFF_HEAD are required" \
		"missing range failure must explain the required inputs"

	run_lint "$repository" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$base_sha" \
		MOCK_LINT_EXIT=0 MOCK_LINT_OUTPUT=
	[[ "$RUN_STATUS" -ne 0 ]] ||
		die "DIFF_HEAD different from checkout HEAD must fail closed"
	assert_contains "$RUN_OUTPUT" "DIFF_HEAD does not match HEAD" \
		"head mismatch failure must be explicit"

	local unrelated_sha
	unrelated_sha="$(printf 'unrelated\n' | git -C "$repository" commit-tree "$(git -C "$repository" mktree </dev/null)")"
	run_lint "$repository" \
		DIFF_BASE="$unrelated_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=0 MOCK_LINT_OUTPUT=
	[[ "$RUN_STATUS" -ne 0 ]] ||
		die "unrelated range must fail closed"
	assert_contains "$RUN_OUTPUT" "DIFF_BASE is not an ancestor of DIFF_HEAD" \
		"unrelated range failure must be explicit"

	local unchanged_output
	unchanged_output=$'alicloud/unchanged/same.go:10:2: R001: stored issue\nalicloud/resource_name1.go:20:3: S006: stored lookalike'
	run_lint "$repository" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=3 MOCK_LINT_OUTPUT="$unchanged_output"
	[[ "$RUN_STATUS" -eq 0 ]] ||
		die "diagnostics in unchanged files must not fail the PR"
	assert_contains "$RUN_OUTPUT" "No issues in changed files" \
		"unchanged findings must be reported as ignored"
	assert_not_contains "$RUN_OUTPUT" "stored issue" \
		"same-basename finding must not match a changed nested path"
	assert_not_contains "$RUN_OUTPUT" "stored lookalike" \
		"regex-like changed path must be matched literally"

	local mixed_output
	mixed_output=$'alicloud/unchanged/same.go:10:2: R001: stored issue\nalicloud/changed/same.go:11:4: R001: changed issue\nalicloud/resource_name1.go:20:3: S006: stored lookalike\nalicloud/resource_name[1].go:21:5: S006: changed literal path'
	run_lint "$repository" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=3 MOCK_LINT_OUTPUT="$mixed_output"
	[[ "$RUN_STATUS" -eq 3 ]] ||
		die "diagnostics in changed files must preserve tfproviderlint exit 3"
	assert_contains "$RUN_OUTPUT" "alicloud/changed/same.go" \
		"changed same-basename finding must be retained"
	assert_contains "$RUN_OUTPUT" 'alicloud/resource_name[1].go' \
		"changed regex-like path must be retained literally"
	assert_not_contains "$RUN_OUTPUT" "alicloud/unchanged/same.go" \
		"unchanged same-basename finding must be filtered"
	assert_not_contains "$RUN_OUTPUT" "alicloud/resource_name1.go" \
		"regex lookalike finding must be filtered"

	assert_contains "$(cat "$MOCK_ARGS_FILE")" "./alicloud/..." \
		"tfproviderlint must continue to analyze the full provider package tree"

	local package_error
	package_error=$'alicloud/unchanged/same.go:10:2: R001: partial diagnostic\n-: package load failed'
	run_lint "$repository" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=1 MOCK_LINT_OUTPUT="$package_error"
	[[ "$RUN_STATUS" -eq 1 ]] ||
		die "package or tool failure must fail closed with its original status"
	assert_contains "$RUN_OUTPUT" "package load failed" \
		"package or tool failure output must not be filtered"
	assert_contains "$RUN_OUTPUT" "did not complete normally" \
		"package or tool failure must be distinguished from diagnostics"

	run_lint "$repository" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=3 MOCK_LINT_OUTPUT='unexpected diagnostic format'
	[[ "$RUN_STATUS" -eq 3 ]] ||
		die "unparseable diagnostic output must fail closed"
	assert_contains "$RUN_OUTPUT" "no analyzable findings" \
		"unparseable diagnostic failure must be explicit"

	commit_file "$repository" scripts/ci/marker.txt 'no Go changes' no-go-change
	local no_go_head_sha
	no_go_head_sha="$(git -C "$repository" rev-parse HEAD)"
	run_lint "$repository" \
		DIFF_BASE="$head_sha" DIFF_HEAD="$no_go_head_sha" \
		MOCK_LINT_EXIT=3 MOCK_LINT_OUTPUT="$unchanged_output"
	[[ "$RUN_STATUS" -eq 0 ]] ||
		die "a valid range without Go changes must ignore stored diagnostics"
	assert_contains "$RUN_OUTPUT" "No issues in changed files" \
		"empty changed-file scope must not fail on stored diagnostics"
}

test_shallow_checkout_fails_closed() {
	local temp_dir="$1"
	local source_repo="$temp_dir/shallow-source"
	local checkout_repo="$temp_dir/shallow-checkout"

	git init -q -b master "$source_repo"
	git -C "$source_repo" config core.hooksPath /dev/null
	git -C "$source_repo" config user.name "CI Lint Test"
	git -C "$source_repo" config user.email "ci-lint-test@example.invalid"
	commit_file "$source_repo" alicloud/example.go 'package alicloud' base
	local base_sha
	base_sha="$(git -C "$source_repo" rev-parse HEAD)"
	git -C "$source_repo" switch -q -c feature
	commit_file "$source_repo" alicloud/example.go \
		$'package alicloud\n\nvar changed = true' feature
	local head_sha
	head_sha="$(git -C "$source_repo" rev-parse HEAD)"

	git -c core.hooksPath=/dev/null clone -q --depth=1 --branch feature \
		"file://$source_repo" "$checkout_repo"
	run_lint "$checkout_repo" \
		DIFF_BASE="$base_sha" DIFF_HEAD="$head_sha" \
		MOCK_LINT_EXIT=0 MOCK_LINT_OUTPUT=
	[[ "$RUN_STATUS" -ne 0 ]] ||
		die "shallow checkout must fail closed"
	assert_contains "$RUN_OUTPUT" "checkout is shallow" \
		"shallow checkout failure must explain fetch-depth requirement"
}

main() {
	[[ -x "$runner" ]] || die "missing executable runner: $runner"

	local temp_dir
	temp_dir="$(mktemp -d)"
	TEST_TEMP_DIR="$temp_dir"
	trap 'rm -rf "$TEST_TEMP_DIR"' EXIT

	MOCK_BIN="$temp_dir/mock-bin"
	MOCK_ARGS_FILE="$temp_dir/tfproviderlint.args"
	mkdir -p "$MOCK_BIN"
	# The mock script must expand these variables when it runs, not while this
	# test creates it.
	# shellcheck disable=SC2016
	printf '%s\n' \
		'#!/usr/bin/env bash' \
		'printf "%s\\n" "$*" >"$MOCK_ARGS_FILE"' \
		'printf "%s\\n" "${MOCK_LINT_OUTPUT:-}"' \
		'exit "${MOCK_LINT_EXIT:-0}"' \
		>"$MOCK_BIN/tfproviderlint"
	chmod +x "$MOCK_BIN/tfproviderlint"
	export MOCK_BIN MOCK_ARGS_FILE

	test_range_validation_and_literal_filtering "$temp_dir"
	test_shallow_checkout_fails_closed "$temp_dir"
	echo "PASS: run-tflint regression tests"
}

main "$@"
