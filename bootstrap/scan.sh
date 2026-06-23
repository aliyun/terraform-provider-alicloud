#!/bin/bash
# scan.sh – pull assigned Aone work items, emit [{id,title,type,status}] JSON to stdout.
# Exits non-zero if a1 fails.

set -uo pipefail

account=$(a1 auth whoami | awk '/Account:/{print $2}')
a1 project workitem list --assignee "$account" -f json \
  | jq '[.[] | {id: .identifier, title, type: .categoryIdentifier, status}]'
