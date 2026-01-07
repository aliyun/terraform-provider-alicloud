#!/usr/bin/env bash

set -euo pipefail

# Only check files that were changed relative to the main branch
base_branch="origin/master"

# HACK: If we seem to be running inside a GitHub Actions pull request check
# then we'll use the PR's target branch from this variable instead.
if [[ -n "${GITHUB_BASE_REF:-}" ]]; then
  base_branch="origin/$GITHUB_BASE_REF"
fi

diffFiles=$(git diff --name-only ${base_branch} --diff-filter=MA 2>/dev/null || git diff --name-only HEAD^ HEAD || true)

if [[ -z "$diffFiles" ]]; then
  echo "==> No files have changed, skipping basic checks"
  exit 0
fi

echo "==> Running basic checks on changed files..."
error=false

while IFS= read -r doc; do
  # Skip empty lines
  [[ -z "$doc" ]] && continue

  # Skip if file doesn't exist
  [[ ! -f "$doc" ]] && continue

  dirname=$(dirname "$doc") || true
  category=$(basename "$dirname") || true

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
done <<< "$diffFiles"

if $error; then
  exit 1
fi

echo "==> All basic checks passed!"
exit 0