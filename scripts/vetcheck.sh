#!/usr/bin/env bash

set -euo pipefail

# Check go vet
echo "==> Checking for suspicious constructs with go vet..."

# Only check files that were changed relative to the main branch
base_branch="origin/master"

# HACK: If we seem to be running inside a GitHub Actions pull request check
# then we'll use the PR's target branch from this variable instead.
if [[ -n "${GITHUB_BASE_REF:-}" ]]; then
  base_branch="origin/$GITHUB_BASE_REF"
fi

# Get changed Go files (compatible with Bash 3)
target_files=$(git diff --name-only ${base_branch} --diff-filter=MA 2>/dev/null | grep "\.go$" | grep -v "_test\.go" | grep -v ".pb.go" | grep -v ".go-version" || true)

if [[ -z "$target_files" ]]; then
  echo "No Go files have changed relative to branch ${base_branch}, so there's nothing to check!"
  exit 0
fi

# Build package list from changed files
packages_map=""

while IFS= read -r filename; do
  if [[ -z "$filename" ]]; then
    continue
  fi

  if [[ ! -f "$filename" ]]; then
    continue
  fi

  # Get package path from file (add ./ prefix for relative paths)
  pkg_dir=$(dirname "$filename")
  # Ensure package path starts with ./
  if [[ "$pkg_dir" != ./* ]]; then
    pkg_dir="./$pkg_dir"
  fi
  if [[ -z "$packages_map" ]] || [[ "$packages_map" != *"$pkg_dir"* ]]; then
    packages_map="$packages_map $pkg_dir"
  fi
done <<< "$target_files"

if [[ -z "$packages_map" ]]; then
  echo "No packages to check!"
  exit 0
fi

# Run go vet on changed packages
echo "Checking packages: $packages_map"
if go vet $packages_map 2>&1; then
  echo 'All of the changed packages look good!'
  exit 0
else
  echo ""
  echo "Vet found suspicious constructs. Please check the reported constructs"
  echo "and fix them if necessary before submitting the code for review."
  exit 1
fi
