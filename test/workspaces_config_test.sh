#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

jq -e '.workspaces.terraform_generator_v4.git_url=="git@gitlab.alibaba-inc.com:opensource-tools/terraform-generator-v4.git"' \
  "$repo_root/config/workspaces.json" >/dev/null
jq -e '.workspaces.terraform_generator_v4.default_branch=="main"' \
  "$repo_root/config/workspaces.json" >/dev/null

mkdir -p "$tmpdir/terraform-generator-v4"
resolved="$(JARVIS_WORKSPACE_ROOT="$tmpdir" bash "$repo_root/bootstrap/workspace.sh" dir terraform_generator_v4)"

if [ "$resolved" != "$tmpdir/terraform-generator-v4" ]; then
  echo "expected $tmpdir/terraform-generator-v4, got $resolved" >&2
  exit 1
fi

echo "workspaces_config_test: PASS"
