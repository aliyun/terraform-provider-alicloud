#!/usr/bin/env bash
# Run policy from CHECKER_SHA against the candidate checkout and Git history.
set -euo pipefail
[[ "$CHECKER_SHA" =~ ^[0-9a-f]{40}$ ]]
test "$(git rev-parse HEAD)" = "$CODE_HEAD"

checker=$(mktemp -d)
trap 'rm -rf "$checker"' EXIT
git archive "$CHECKER_SHA" scripts/ci scripts/title scripts/internal/splitlayout go.mod | tar -x -C "$checker"
# Receipt reuse must also account for the Go resolver outside scripts/ci.
resolver_inputs=$(git ls-tree "$CHECKER_SHA" -- scripts/title scripts/internal/splitlayout go.mod scripts/ci/run-integration-checker.sh)
PROVIDER_CHECKER_INPUTS=$(printf '%s' "$resolver_inputs" | git hash-object --stdin)
export PROVIDER_CHECKER_INPUTS
candidate=$(pwd)
GIT_DIR=$(git rev-parse --absolute-git-dir)
export GIT_DIR GIT_WORK_TREE="$candidate"
ln -s "$candidate/alicloud" "$checker/alicloud"
cd "$checker"
export GOWORK=off GOFLAGS=-mod=readonly
echo "Provider checker=$CHECKER_SHA candidate=$CODE_HEAD base=${CODE_MR_BASE:-${CODE_PUSH_BEFORE:-auto}} check=$CHECK_KIND" >&2
python3 -Es scripts/ci/acube-integration.py --kind "$CHECK_KIND" --output "$candidate/.acube-ci" "$@"
