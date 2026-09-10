#!/usr/bin/env bash

set -euo pipefail

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."

# Get changed Go files from the latest commit
# Priority:
# 1. Uncommitted changes (if any)
# 2. Latest commit changes (HEAD~1 HEAD)
# 3. All changes from master branch (for CI/PR checks)

# HACK: If we seem to be running inside a GitHub Actions pull request check
# then we'll use the PR's target branch from this variable instead.
if [[ -n "${GITHUB_BASE_REF:-}" ]]; then
  # CI mode: check entire PR branch vs target branch
  base_branch="origin/$GITHUB_BASE_REF"
  target_files=$(git diff --name-only ${base_branch} --diff-filter=MA 2>/dev/null | grep "\.go$" | grep -v ".pb.go" | grep -v ".go-version" || true)
else
  # Local mode: check latest commit or uncommitted changes
  # First try uncommitted changes
  target_files=$(git diff --name-only HEAD --diff-filter=MA 2>/dev/null | grep "\.go$" | grep -v ".pb.go" | grep -v ".go-version" || true)

  if [[ -z "$target_files" ]]; then
    # No uncommitted changes, check latest commit
    target_files=$(git diff --name-only HEAD~1 HEAD --diff-filter=MA 2>/dev/null | grep "\.go$" | grep -v ".pb.go" | grep -v ".go-version" || true)
  fi

  # If still no changes, don't fall back to master branch
  # Only check what's actually changed in the latest commit
  if [[ -z "$target_files" ]]; then
    echo "No Go files have changed in the latest commit, so there's nothing to check!"
    exit 0
  fi
fi

# Check each file for formatting issues
incorrect_files=""

while IFS= read -r filename; do
  if [[ -z "$filename" ]]; then
    continue
  fi

  if [[ ! -f "$filename" ]]; then
    continue
  fi

  # Check if file needs formatting
  if gofmt -l "$filename" 2>/dev/null | grep -q .; then
    incorrect_files="$incorrect_files$filename"$'\n'
  fi
done <<< "$target_files"

if [[ -n "$incorrect_files" ]]; then
  echo 'gofmt needs running on the following files:'
  echo "$incorrect_files" | sed 's/^/ - /'
  echo "You can use the command: \`make fmt\` to reformat code."
  exit 1
fi

echo 'All of the changed files look good!'
exit 0
