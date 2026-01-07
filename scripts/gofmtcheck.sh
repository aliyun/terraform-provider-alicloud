#!/usr/bin/env bash

set -euo pipefail

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."

# Only check files that were changed relative to the main branch
base_branch="origin/master"

# HACK: If we seem to be running inside a GitHub Actions pull request check
# then we'll use the PR's target branch from this variable instead.
if [[ -n "${GITHUB_BASE_REF:-}" ]]; then
  base_branch="origin/$GITHUB_BASE_REF"
fi

# Get changed Go files (compatible with Bash 3)
target_files=$(git diff --name-only ${base_branch} --diff-filter=MA 2>/dev/null | grep "\.go$" | grep -v ".pb.go" | grep -v ".go-version" || true)

if [[ -z "$target_files" ]]; then
  echo "No Go files have changed relative to branch ${base_branch}, so there's nothing to check!"
  exit 0
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
