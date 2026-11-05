#!/usr/bin/env bash

# Detect whether the PR change range touches build-input files (CI workflows,
# GNUmakefile, CI/check scripts). Used by the ProductionBuild CI job's route to
# decide whether to run the real production build + production dependency-graph
# check on PRs that the Basic Policy classifier marks NOT relevant (e.g. pure
# CI/Make/script changes). Product-Go PRs already trigger ProductionBuild via
# route.outputs.relevant; this adds the build-input path so a base CI/Make/script
# PR itself is covered instead of being skipped.
#
# Prints "true" on the first matching changed file, otherwise "false". Always
# exits 0 so the calling GitHub Actions step can capture the output.
#
# Inputs:
#   DIFF_BASE, DIFF_HEAD   change range (required when BUILD_INPUT_FILES unset)
#   BUILD_INPUT_FILES      file with one changed path per line (testing override;
#                          when set, git is not invoked)

set -euo pipefail

# Monitored build-input paths. Entries ending in "/" are directory prefixes
# (any changed file under them matches); entries without a trailing slash are
# matched exactly. Keep in sync with detect-build-input-changes_test.sh.
monitored=(
  "GNUmakefile"
  ".github/workflows/"
  "scripts/ci/"
  "scripts/local-ci-check.sh"
  "scripts/vetcheck.sh"
  "scripts/gofmtcheck.sh"
  "scripts/goimportscheck.sh"
  "scripts/errcheck.sh"
)

match_one() {
        local f="$1"
        for p in "${monitored[@]}"; do
                if [[ "$p" == */ ]]; then
                        [[ "$f" == "$p"* ]] && return 0
                else
                        [[ "$f" == "$p" ]] && return 0
                fi
        done
        return 1
}

if [[ -n "${BUILD_INPUT_FILES:-}" ]]; then
        if [[ ! -f "$BUILD_INPUT_FILES" ]]; then
                echo "false"
                echo "detect-build-input-changes: BUILD_INPUT_FILES not found: $BUILD_INPUT_FILES" >&2
                exit 0
        fi
        files_file="$BUILD_INPUT_FILES"
else
        : "${DIFF_BASE:?DIFF_BASE required when BUILD_INPUT_FILES is unset}"
        : "${DIFF_HEAD:?DIFF_HEAD required when BUILD_INPUT_FILES is unset}"
        files_file="$(mktemp)"
        trap 'rm -f "$files_file"' EXIT
        git diff --name-only --no-renames "$DIFF_BASE" "$DIFF_HEAD" >"$files_file" 2>/dev/null || true
fi

changed=false
while IFS= read -r f; do
        [[ -n "$f" ]] || continue
        if match_one "$f"; then
                changed=true
                break
        fi
done <"$files_file"

echo "$changed"
