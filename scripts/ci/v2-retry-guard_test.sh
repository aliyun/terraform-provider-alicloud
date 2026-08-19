#!/usr/bin/env bash

# Fixture tests for scripts/ci/v2-retry-guard.sh. Each case builds a throwaway
# git repository with a base commit and one or more head commits, then runs the
# guard over the resulting range and asserts the exit code and output.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
guard="$script_dir/v2-retry-guard.sh"

die() {
        echo "FAIL: $*" >&2
        exit 1
}

setup_repo() {
        local repository="$1"
        mkdir -p "$repository/alicloud"
        git init -q -b master "$repository"
        git -C "$repository" config core.hooksPath /dev/null
        git -C "$repository" config user.name "V2 Retry Guard Test"
        git -C "$repository" config user.email "v2-retry-guard-test@example.invalid"
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

run_guard() {
        # run_guard <repository> <base> <head> -> captures output, returns the
        # guard exit code.
        local repository="$1"
        local base="$2"
        local head="$3"
        guard_output="$(cd "$repository" &&
                DIFF_BASE="$base" DIFF_HEAD="$head" bash "$guard" 2>&1)" || return $?
        return 0
}

run_guard_github_actions() {
        # Same as run_guard but with GITHUB_ACTIONS=true so the guard emits
        # workflow-command annotations instead of plain lines.
        local repository="$1"
        local base="$2"
        local head="$3"
        guard_output="$(cd "$repository" &&
                DIFF_BASE="$base" DIFF_HEAD="$head" \
                        env GITHUB_ACTIONS=true bash "$guard" 2>&1)" || return $?
        return 0
}

assert_output_contains() {
        local needle="$1"
        [[ "$guard_output" == *"$needle"* ]] ||
                die "expected output to contain '$needle', got:
$guard_output"
}

test_new_retry_call_is_rejected() {
        local repository="$1/new-retry-call"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_example.go \
                $'package alicloud\n\nfunc example() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        commit_file "$repository" alicloud/service_example.go \
                $'package alicloud\n\nimport (\n\t"context"\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example(ctx context.Context) {\n\tresource.Retry(60, func() *resource.RetryError {\n\t\treturn resource.RetryableError("try again")\n\t})\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        if run_guard "$repository" "$base" "$head"; then
                die "new resource.Retry call should be rejected
$guard_output"
        fi
        assert_output_contains "service_example.go:5: deprecated import \"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource\" -> replace with \"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry\""
        assert_output_contains "service_example.go:9: deprecated symbol: resource.Retry -> retry.Retry"
        assert_output_contains "service_example.go:9: deprecated symbol: resource.RetryError -> retry.RetryError"
        assert_output_contains "service_example.go:10: deprecated symbol: resource.RetryableError -> retry.RetryableError"
}

test_new_import_is_rejected() {
        local repository="$1/new-import"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_import.go \
                $'package alicloud\n\nfunc example() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        # Import added but the (existing) retry calls are untouched: the
        # import itself is the regression signal.
        commit_file "$repository" alicloud/service_import.go \
                $'package alicloud\n\nimport (\n\t_ "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example() {}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        if run_guard "$repository" "$base" "$head"; then
                die "new helper/resource import should be rejected
$guard_output"
        fi
        assert_output_contains "service_import.go:4: deprecated import \"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource\" -> replace with"
        assert_output_contains "no go get needed"
}

test_helper_retry_is_accepted() {
        local repository="$1/helper-retry"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_retry.go \
                $'package alicloud\n\nfunc example() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        commit_file "$repository" alicloud/service_retry.go \
                $'package alicloud\n\nimport (\n\t"context"\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"\n)\n\nfunc example(ctx context.Context) {\n\tretry.Retry(ctx, 60, func() *retry.RetryError {\n\t\treturn retry.RetryableError("try again")\n\t})\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        run_guard "$repository" "$base" "$head" ||
                die "helper/retry usage should be accepted
$guard_output"
}

test_untouched_legacy_usage_is_accepted() {
        local repository="$1/legacy-untouched"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_legacy.go \
                $'package alicloud\n\nimport (\n\t"context"\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example(ctx context.Context) {\n\tresource.Retry(60, func() *resource.RetryError {\n\t\treturn resource.RetryableError("try again")\n\t})\n}\n\nfunc other() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        # Edit the file without touching the retry lines: unmigrated files
        # must stay editable until the Route 1 migration reaches them.
        commit_file "$repository" alicloud/service_legacy.go \
                $'package alicloud\n\nimport (\n\t"context"\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example(ctx context.Context) {\n\tresource.Retry(60, func() *resource.RetryError {\n\t\treturn resource.RetryableError("try again")\n\t})\n}\n\nfunc other() {\n\t// updated\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        run_guard "$repository" "$base" "$head" ||
                die "untouched legacy usage should be accepted
$guard_output"
}

test_aliased_import_is_rejected() {
        local repository="$1/aliased-import"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_alias.go \
                $'package alicloud\n\nfunc example() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        commit_file "$repository" alicloud/service_alias.go \
                $'package alicloud\n\nimport (\n\t"context"\n\tsdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example(ctx context.Context) {\n\tsdkresource.Retry(60, func() *sdkresource.RetryError {\n\t\treturn sdkresource.RetryableError("try again")\n\t})\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        if run_guard "$repository" "$base" "$head"; then
                die "aliased helper/resource usage should be rejected
$guard_output"
        fi
        assert_output_contains "service_alias.go:9: deprecated symbol: sdkresource.Retry -> retry.Retry"
        assert_output_contains "service_alias.go:10: deprecated symbol: sdkresource.RetryableError -> retry.RetryableError"
}

test_github_actions_annotations() {
        local repository="$1/annotations"
        setup_repo "$repository"
        commit_file "$repository" alicloud/service_annotations.go \
                $'package alicloud\n\nfunc example() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        commit_file "$repository" alicloud/service_annotations.go \
                $'package alicloud\n\nimport (\n\t"context"\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc example(ctx context.Context) {\n\tresource.Retry(60, func() *resource.RetryError {\n\t\treturn resource.RetryableError("try again")\n\t})\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        if run_guard_github_actions "$repository" "$base" "$head"; then
                die "annotations run should be rejected
$guard_output"
        fi
        assert_output_contains "::error file=alicloud/service_annotations.go,line=9,title=v2-retry-guard::deprecated symbol: resource.Retry -> retry.Retry"
        # Workflow commands are consumed by the runner; the plain-text copy
        # keeps the violation with its line number visible in the job log.
        assert_output_contains "v2-retry-guard: alicloud/service_annotations.go:9: deprecated symbol: resource.Retry -> retry.Retry"
        assert_output_contains "v2-retry-guard: 4 deprecated helper/resource usage(s) introduced."
}

test_test_files_are_exempt() {
        local repository="$1/test-files"
        setup_repo "$repository"
        commit_file "$repository" alicloud/resource_alicloud_example_test.go \
                $'package alicloud\n\nfunc TestExample() {}' "base"
        local base head
        base="$(git -C "$repository" rev-parse HEAD)"
        commit_file "$repository" alicloud/resource_alicloud_example_test.go \
                $'package alicloud\n\nimport (\n\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"\n)\n\nfunc TestExample() {\n\tresource.Retry(60, func() *resource.RetryError {\n\t\treturn resource.RetryableError("try again")\n\t})\n}' "head"
        head="$(git -C "$repository" rev-parse HEAD)"

        run_guard "$repository" "$base" "$head" ||
                die "test files should stay exempt
$guard_output"
}

main() {
        temp_dir="$(mktemp -d)"
        trap 'rm -rf "${temp_dir:-}"' EXIT

        test_new_retry_call_is_rejected "$temp_dir"
        test_new_import_is_rejected "$temp_dir"
        test_aliased_import_is_rejected "$temp_dir"
        test_github_actions_annotations "$temp_dir"
        test_helper_retry_is_accepted "$temp_dir"
        test_untouched_legacy_usage_is_accepted "$temp_dir"
        test_test_files_are_exempt "$temp_dir"

        echo "v2-retry-guard tests: all passed"
}

main "$@"
