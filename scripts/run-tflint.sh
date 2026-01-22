#!/usr/bin/env bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# check tfproviderlint
function ensure_tfproviderlint() {
  if ! command -v tfproviderlint >/dev/null 2>&1; then
    echo -e "${YELLOW} tfproviderlint not installed... ${NC}"
    exit 1
  fi
}

function run_lint {
  echo -e "==> ${GREEN} Checking source code against tfproviderlint..."
  if output=$(tfproviderlint \
    -AT001\
    -AT001.ignored-filename-prefixes data_source_alicloud_\
    -AT005 -AT006 -AT007\
    -R001 -R002 -R003 -R004 -R006\
    -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020\
    -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033\
    "./alicloud/..." 2>&1); then
    if [ -z "$output" ]; then
      echo -e "${GREEN} tfproviderlint: No issues found.${NC}"
    else
      echo "$output"
      echo -e "${GREEN} tfproviderlint completed.${NC}"
    fi
    return 0
  else
    local exit_code=$?
    echo -e "${RED} Found issues:${NC}"
    echo "$output"
    echo -e "${RED} Fix the above issues before committing.${NC}"
    return $exit_code
  fi
}

main() {
  ensure_tfproviderlint
  run_lint
}

main

