#!/usr/bin/env bash

set -euo pipefail

# Panic static check for alicloud/ Go files
# Modelled after scripts/vetcheck.sh:
#   - Diffs current branch against origin/master to get changed files
#   - Filters to .go files under alicloud/
#   - Runs the Go panic checker on changed files
#   - Exits non-zero if panics found in changed files

echo "==> Checking for panic() calls in changed Go files..."

base_branch="origin/master"

# If running inside a GitHub Actions pull request check, use the PR target branch
if [[ -n "${GITHUB_BASE_REF:-}" ]]; then
  base_branch="origin/$GITHUB_BASE_REF"
fi

# Get changed Go files under alicloud/ (exclude test files by default in the checker)
target_files=$(git diff --name-only "${base_branch}" --diff-filter=MA 2>/dev/null | grep "^alicloud/.*\.go$" || true)

if [[ -z "$target_files" ]]; then
  echo "No Go files under alicloud/ have changed relative to ${base_branch}, nothing to check!"
  exit 0
fi

# Write file list to a temp file for the checker
TEMP_FILELIST=$(mktemp)
echo "$target_files" > "$TEMP_FILELIST"

echo "Changed files to check:"
echo "$target_files" | sed 's/^/  /'
echo ""

# Run the panic checker on the changed files
exitCode=0
if go run scripts/panic_check/panic_check.go -fileNames="$TEMP_FILELIST"; then
  echo "==> PASS: No panic() calls found in changed files."
else
  echo ""
  echo "==> FAIL: panic() calls found in changed files."
  echo "Please replace panic() with proper error handling (log + return zero value)."
  exitCode=1
fi

rm -f "$TEMP_FILELIST"
exit $exitCode
