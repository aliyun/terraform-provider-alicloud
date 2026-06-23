#!/bin/bash
# scan.sh – pull assigned Aone work items, emit [{id,title,type,status}] JSON to stdout.
# Exits non-zero if a1 fails.

set -uo pipefail

account=$(a1 auth whoami | awk '/Account:/{print $2}')
if [ -z "$account" ]; then
  echo "scan.sh: could not determine account from 'a1 auth whoami'" >&2
  exit 1
fi
a1 project workitem list --assignee "$account" -f json \
  | jq '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status}]'
