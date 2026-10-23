#!/usr/bin/env bash

# Guard the testing-framework boundary on master: no new imports of the
# terraform-plugin-sdk v2 helper/resource package, in test files or not.
#
# Background: the test line flipped to terraform-plugin-testing
# (helper/resource, helper/acctest, terraform). Every alicloud test binary now
# links the plugin-testing harness, and both SDK v2 helper/resource and
# plugin-testing helper/resource register a -sweep flag in init(), so a binary
# linking both dies at startup with "flag redefined: sweep". This guard stops
# new SDK v2 helper/resource imports from creeping back in:
#   - non-test files: the retry helpers live in
#     terraform-plugin-sdk/v2/helper/retry (flag-free)
#   - test files under alicloud/: the harness is
#     terraform-plugin-testing/helper/resource
#
# The check inspects only lines ADDED by the change range, so files that have
# not been migrated yet can still be edited. Test files outside alicloud/
# (e.g. scripts/ fixtures) are exempt: they build standalone binaries that
# never link the plugin-testing harness.
#
# Every violation is reported with its real file line number and the exact
# replacement. Under GitHub Actions the violations are additionally emitted
# as ::error workflow commands so they show up as inline annotations in the
# pull request diff; a plain-text copy is printed as well because workflow
# commands are hidden from the job log.
#
# Usage: DIFF_BASE=<sha> DIFF_HEAD=<sha> test-framework-guard.sh
#
# When DIFF_BASE/DIFF_HEAD are unset and TEST_FRAMEWORK_GUARD_AUTO_RANGE=true,
# the range is resolved against origin/${GITHUB_BASE_REF:-master} the same way
# scripts/basic-check.sh does.

set -euo pipefail

die() {
        echo "test-framework-guard: $*" >&2
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

if [[ -z "${DIFF_BASE:-}" || -z "${DIFF_HEAD:-}" ]]; then
        [[ "${TEST_FRAMEWORK_GUARD_AUTO_RANGE:-false}" == "true" ]] ||
                die "DIFF_BASE and DIFF_HEAD are required (or set TEST_FRAMEWORK_GUARD_AUTO_RANGE=true)"

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

retry_symbols="Retry|RetryableError|NonRetryableError|RetryError|StateRefreshFunc|StateChangeConf"
import_pattern='"github\.com/hashicorp/terraform-plugin-sdk(/v2)?/helper/resource"'
retry_import='"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"'
testing_import='"github.com/hashicorp/terraform-plugin-testing/helper/resource"'

violations=0

# Report one violation. Under GitHub Actions the same message becomes an
# inline annotation on the offending line of the pull request diff.
emit_violation() {
        local file="$1" lineno="$2" message="$3"

        if [[ "${GITHUB_ACTIONS:-}" == "true" ]]; then
                echo "::error file=${file},line=${lineno},title=test-framework-guard::${message}"
                # Workflow commands are hidden from the job log; repeat the
                # violation as plain text so it is visible there too.
                echo "test-framework-guard: ${file}:${lineno}: ${message}"
        else
                echo "test-framework-guard: ${file}:${lineno}: ${message}"
        fi
        violations=$((violations + 1))
}

check_file() {
        local file="$1"
        local added_lines identifiers symbol symbols lineno line_text replacement is_test

        [[ "$file" == *.go ]] || return 0
        [[ -f "$file" ]] || return 0

        if [[ "$file" == *_test.go ]]; then
                # Only alicloud test binaries link the plugin-testing harness;
                # test files elsewhere (scripts/ fixtures) are standalone.
                [[ "$file" == alicloud/* ]] || return 0
                is_test="true"
                replacement="$testing_import"
        else
                is_test="false"
                replacement="$retry_import"
        fi

        # Added lines with their real line numbers in the head version,
        # derived from the @@ hunk headers of a zero-context diff.
        added_lines="$(git diff --unified=0 "$diff_base" "$diff_head" -- "$file" |
                awk '/^\+\+\+/ { next }
                        /^@@/ { match($0, /\+[0-9]+/); n = substr($0, RSTART + 1, RLENGTH - 1) + 0; next }
                        /^\+/ { printf "%d %s\n", n, substr($0, 2); n++ }')" || true
        [[ -n "$added_lines" ]] || return 0

        # Identifiers bound to helper/resource by added import lines: the
        # default package name plus any explicit aliases. Only the non-test
        # retry-symbol check needs them.
        identifiers="resource"
        while IFS= read -r alias; do
                [[ -n "$alias" ]] && identifiers="${identifiers}|${alias}"
        done < <(printf '%s\n' "$added_lines" |
                { grep -oE "[A-Za-z_][A-Za-z0-9_]*[[:space:]]+${import_pattern}" || true; } |
                { grep -oE "^[A-Za-z_][A-Za-z0-9_]*" || true; })

        while IFS= read -r added; do
                [[ -n "$added" ]] || continue
                lineno="${added%% *}"
                line_text="${added#* }"

                if grep -Eq "^[[:space:]]*(_[[:space:]]+|[A-Za-z_][A-Za-z0-9_]*[[:space:]]+)?${import_pattern}" <<<"$line_text"; then
                        deprecated_import="$(grep -oE "${import_pattern}" <<<"$line_text" | head -1)"
                        emit_violation "$file" "$lineno" \
                                "deprecated import ${deprecated_import} -> replace with ${replacement} (already in go.mod; no go get needed, just change the import)"
                fi

                if [[ "$is_test" == "true" ]]; then
                        # plugin-testing helper/resource exports the retry
                        # family too (deprecated copies there); only the
                        # import decides which fork a binary links, so the
                        # symbol-level check stays non-test only.
                        continue
                fi

                symbols="$(grep -oE "(${identifiers})\\.(${retry_symbols})" <<<"$line_text")" || true
                [[ -n "$symbols" ]] || continue
                while IFS= read -r symbol; do
                        [[ -n "$symbol" ]] || continue
                        emit_violation "$file" "$lineno" \
                                "deprecated symbol: ${symbol} -> retry.${symbol#*.}"
                done <<<"$symbols"
        done <<<"$added_lines"
}

changed_files="$(git diff --name-only --diff-filter=MA "$diff_base" "$diff_head")"
if [[ -n "$changed_files" ]]; then
        while IFS= read -r changed_file; do
                [[ -n "$changed_file" ]] && check_file "$changed_file"
        done <<<"$changed_files"
fi

if ((violations > 0)); then
        echo "test-framework-guard: ${violations} deprecated SDK v2 helper/resource usage(s) introduced."
        echo "        non-test code: use terraform-plugin-sdk/v2/helper/retry for the retry helpers"
        echo "        (scripts/migration/retry_rewrite.sh applies the rewrite);"
        echo "        test files: use terraform-plugin-testing/helper/resource for the harness."
        exit 1
fi

echo "test-framework-guard: no new SDK v2 helper/resource usage found"
