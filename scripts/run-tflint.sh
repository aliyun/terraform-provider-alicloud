#!/usr/bin/env bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

function ensure_tfproviderlint() {
  if ! command -v tfproviderlint >/dev/null 2>&1; then
    echo -e "${YELLOW} tfproviderlint not installed... ${NC}"
    exit 1
  fi
}

# Determine the merge-base ref and return a newline-separated list of changed
# .go files.  Prints nothing and returns 1 if the base cannot be determined.
function get_changed_go_files() {
  local base_ref=""

  # GitHub Actions PR context: GITHUB_BASE_REF is the target branch name (e.g. "main")
  if [ -n "${GITHUB_BASE_REF:-}" ]; then
    # Ensure the base branch is fetched (needed when fetch-depth > 1 but branch
    # was not previously fetched, or as a safety net).
    git fetch --no-tags --depth=1 origin "${GITHUB_BASE_REF}" 2>/dev/null || true
    base_ref="origin/${GITHUB_BASE_REF}"
  elif git rev-parse --verify origin/main >/dev/null 2>&1; then
    base_ref="origin/main"
  elif git rev-parse --verify origin/master >/dev/null 2>&1; then
    base_ref="origin/master"
  else
    return 1
  fi

  # Use symmetric diff (merge-base) so we only get commits on this branch
  git diff --name-only "${base_ref}...HEAD" -- '*.go' 2>/dev/null
}

function run_lint {
  echo -e "==> ${GREEN} Checking source code against tfproviderlint...${NC}"

  local full_output=""
  local lint_exit=0
  full_output=$(tfproviderlint \
    -AT001 \
    -AT001.ignored-filename-prefixes data_source_alicloud_ \
    -AT005 -AT006 -AT007 \
    -R001 -R002 -R003 -R004 -R006 \
    -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 \
    -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020 \
    -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 \
    -S031 -S032 -S033 \
    "./alicloud/..." 2>&1) || lint_exit=$?

  if [ $lint_exit -eq 0 ]; then
    echo -e "${GREEN} tfproviderlint: No issues found.${NC}"
    return 0
  fi

  # Try to get the list of changed files so we can report only new issues
  local changed_files=""
  changed_files=$(get_changed_go_files 2>/dev/null) || true

  if [ -z "$changed_files" ]; then
    # Cannot determine changed files — fall back to reporting everything
    echo -e "${RED} Found issues (could not determine changed files, showing all):${NC}"
    echo "$full_output"
    echo -e "${RED} Fix the above issues before committing.${NC}"
    return $lint_exit
  fi

  # Build a regex pattern that matches any of the changed file basenames
  local pattern=""
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    local base
    base=$(basename "$f")
    if [ -z "$pattern" ]; then
      pattern="$base"
    else
      pattern="${pattern}|${base}"
    fi
  done <<< "$changed_files"

  if [ -z "$pattern" ]; then
    echo -e "${GREEN} tfproviderlint: No .go files changed.${NC}"
    return 0
  fi

  local filtered_output=""
  filtered_output=$(echo "$full_output" | grep -E "($pattern)" || true)

  if [ -z "$filtered_output" ]; then
    echo -e "${GREEN} tfproviderlint: No issues in changed files.${NC}"
    echo -e "${YELLOW} Note: pre-existing issues in unchanged files were ignored.${NC}"
    return 0
  fi

  echo -e "${RED} Found issues in changed files:${NC}"
  echo "$filtered_output"
  echo -e "${RED} Fix the above issues before committing.${NC}"
  return $lint_exit
}

main() {
  ensure_tfproviderlint
  run_lint
}

main
