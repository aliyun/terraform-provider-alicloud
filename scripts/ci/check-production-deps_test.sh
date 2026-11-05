#!/usr/bin/env bash

# Fixture tests for scripts/ci/check-production-deps.sh. Each case feeds a
# fixed production dependency graph via PROD_DEPS_INPUT (so the assertion runs
# without a Go toolchain) and asserts the exit code and output.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
check="$script_dir/check-production-deps.sh"

die() {
        echo "FAIL: $*" >&2
        exit 1
}

# graph_test <name> <graph-content> <expected-exit> <expected-substring>
graph_test() {
        local name="$1" graph="$2" expect_exit="$3" expect_match="$4"
        local tmp
        tmp="$(mktemp)"
        printf '%s\n' "$graph" >"$tmp"
        set +e
        output="$(PROD_DEPS_INPUT="$tmp" "$check" 2>&1)"
        rc=$?
        set -e
        rm -f "$tmp"
        [[ "$rc" -eq "$expect_exit" ]] ||
                die "$name: expected exit $expect_exit, got $rc; output: $output"
        [[ -z "$expect_match" || "$output" == *"$expect_match"* ]] ||
                die "$name: expected output to contain '$expect_match'; got: $output"
        echo "ok - $name"
}

# A clean graph: SDK v2 helper/retry is a production dependency, but the
# forbidden helper/resource test package is absent.
graph_test "clean graph passes" \
        $'github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server\ngithub.com/hashicorp/terraform-plugin-sdk/v2/helper/retry\ngithub.com/aliyun/terraform-provider-alicloud/alicloud\n' \
        0 "PASS"

# The forbidden package is a direct production dependency.
graph_test "direct forbidden dep fails" \
        $'github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource\n' \
        1 "FAILED"

# A transitive forbidden dependency (anywhere in the graph) fails.
graph_test "transitive forbidden dep fails" \
        $'github.com/aliyun/terraform-provider-alicloud/alicloud\ngithub.com/hashicorp/terraform-plugin-sdk/v2/helper/resource\n' \
        1 "forbidden production dependency found"

# A sub-package of the forbidden path is also rejected (anchoring).
graph_test "forbidden sub-package fails" \
        $'github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource/sub\n' \
        1 "FAILED"

# An unrelated package whose path merely contains the forbidden substring as
# a prefix is NOT rejected (anchoring prevents false positives).
graph_test "unrelated look-alike passes" \
        $'evilsuite/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource\n' \
        0 "PASS"

# FORBIDDEN_DEPS can be overridden; a custom forbidden path is then enforced,
# while the default scope does not forbid it.
custom_tmp="$(mktemp)"
printf '%s\n' 'example.com/forbidden/pkg' >"$custom_tmp"
set +e
out="$(PROD_DEPS_INPUT="$custom_tmp" "$check" 2>&1)"; rc=$?
set -e
rm -f "$custom_tmp"
[[ "$rc" -eq 0 ]] || die "custom forbidden: default scope should pass, got rc=$rc; output: $out"
echo "ok - custom forbidden not triggered under default scope"

custom_tmp="$(mktemp)"
printf '%s\n' 'example.com/forbidden/pkg' >"$custom_tmp"
set +e
out="$(FORBIDDEN_DEPS="example.com/forbidden/pkg" PROD_DEPS_INPUT="$custom_tmp" "$check" 2>&1)"; rc=$?
set -e
rm -f "$custom_tmp"
[[ "$rc" -eq 1 ]] || die "custom forbidden: override should fail, got rc=$rc; output: $out"
echo "ok - custom forbidden dep override enforced"

echo "all check-production-deps tests passed"
