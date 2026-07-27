#!/usr/bin/env bash

# Runs tfproviderlint over ./alicloud/... and reports only findings in the files
# the current PR changed. See scripts/run-tflint-selftest.sh for the behaviour
# this is expected to have.

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Lint scope. Change detection and output filtering stay in sync with this path
# so we only ever consider files that are actually linted.
LINT_TARGET="./alicloud/..."
LINT_PATHSPEC="alicloud/*.go"

ANALYZERS=(
  -AT001 -AT001.ignored-filename-prefixes data_source_alicloud_
  -AT005 -AT006 -AT007
  -R001 -R002 -R003 -R004 -R006
  -S{001..033}
)

function ensure_tfproviderlint() {
  command -v tfproviderlint >/dev/null 2>&1 && return 0
  echo -e "${YELLOW} tfproviderlint not installed... ${NC}"
  exit 1
}

# Base-detection diagnostics go to stderr so they survive the command
# substitution that captures the resolved base.
function note() {
  echo -e "${YELLOW} tfproviderlint: base detection — $1${NC}" >&2
}

function fail_with() {
  echo -e "${RED} $1${NC}"
  echo "$2"
  echo -e "${RED} Fix the above before committing.${NC}"
}

# Print $1 as a commit sha, or fail if this clone does not have it.
function commit_sha() {
  git rev-parse --verify -q "${1}^{commit}" 2>/dev/null
}

function base_sha_from_payload() {
  local payload="${GITHUB_EVENT_PATH:-}"
  [ -n "$payload" ] && [ -r "$payload" ] || return 1
  if command -v jq >/dev/null 2>&1; then
    jq -r '.pull_request.base.sha // empty' "$payload" 2>/dev/null
  elif command -v python3 >/dev/null 2>&1; then
    python3 -c 'import json,sys; print((json.load(open(sys.argv[1])).get("pull_request") or {}).get("base",{}).get("sha") or "")' \
      "$payload" 2>/dev/null
  fi
}

# Print the commit this PR diverged from. The layers keep the whole-repo fallback
# in run_lint — which reports pre-existing issues the PR never touched —
# unreachable in practice; returns 1 only if all of them fail.
function resolve_base_commit() {
  local sha

  # 1. PR context: GITHUB_BASE_REF is the target branch. The fetch is a safety
  #    net (it also restores origin/<base> on a shallow checkout), not required.
  if [ -n "${GITHUB_BASE_REF:-}" ]; then
    git fetch --no-tags origin "${GITHUB_BASE_REF}" 2>/dev/null || true
    commit_sha "origin/${GITHUB_BASE_REF}" && return 0
    note "origin/${GITHUB_BASE_REF} unavailable, trying the event payload"
  fi

  # 2. The event payload holds the exact base commit GitHub computed, whatever
  #    refs the checkout happened to fetch.
  if sha=$(base_sha_from_payload) && [ -n "$sha" ]; then
    commit_sha "$sha" && return 0
    git fetch --no-tags --depth=1 origin "$sha" 2>/dev/null || true
    commit_sha "$sha" && return 0
    note "base sha ${sha} from the event payload is not in this clone"
  fi

  # 3. actions/checkout checks out refs/pull/N/merge for PRs — a merge commit
  #    whose first parent IS the base, needing no network or extra refs. Gated on
  #    the event: on a plain branch checkout HEAD^1 is just the previous commit.
  case "${GITHUB_EVENT_NAME:-}" in
    pull_request | pull_request_target)
      if commit_sha 'HEAD^2' >/dev/null && sha=$(commit_sha 'HEAD^1'); then
        note "using the merge commit's first parent ${sha:0:9} as base"
        echo "$sha"
        return 0
      fi
      ;;
  esac

  # 4. Local runs: whatever main-line ref this clone happens to have.
  local ref
  for ref in origin/main origin/master upstream/main upstream/master main master; do
    commit_sha "$ref" && return 0
  done
  return 1
}

# core.quotePath=false keeps non-ASCII paths unquoted ("alicloud/é.go", not
# "alicloud/\303\251.go") so they still match tfproviderlint's output.
function diff_names() {
  git -c core.quotePath=false diff --name-only --diff-filter=ACMR \
    "$1" HEAD -- "$LINT_PATHSPEC" 2>/dev/null
}

# Print the changed alicloud/*.go files (added, copied, modified or renamed —
# a deleted file cannot be linted). Returns 1 only when the base is unknown; a
# known base with no matching changes returns 0 and no output.
function get_changed_go_files() {
  local base_commit diff_base changed
  base_commit=$(resolve_base_commit) || return 1
  # The merge-base is correct PR semantics (only what this branch added); fall
  # back to the base itself when no common ancestor is reachable.
  diff_base=$(git merge-base "$base_commit" HEAD 2>/dev/null || echo "$base_commit")

  if changed=$(diff_names "$diff_base"); then
    printf '%s\n' "$changed"
    return 0
  fi

  # A shallow clone can know the base commit yet lack its objects, failing the
  # diff. Deepen once and retry, since giving up means the whole-repo fallback.
  [ "$(git rev-parse --is-shallow-repository 2>/dev/null)" = "true" ] || return 1
  note "shallow clone — deepening history to diff against the base"
  git fetch --no-tags --unshallow origin 2>/dev/null ||
    git fetch --no-tags --deepen=100 origin 2>/dev/null || true
  changed=$(diff_names "$diff_base") || return 1
  printf '%s\n' "$changed"
}

# Turn the changed-file list into an anchored, regex-escaped alternation:
# "alicloud/x.go" -> "(^|/)alicloud/x\.go:". Matching the whole path (never a
# basename or substring) keeps alicloud/foo.go and alicloud/sub/foo.go distinct,
# and the (^|/) anchor covers both the relative and absolute forms
# tfproviderlint prints.
function build_filter_pattern() {
  printf '%s\n' "$1" | sed -e '/^$/d' \
    -e 's/[.[\*^$()+?{}|]/\\&/g' -e 's#^#(^|/)#' -e 's#$#:#' |
    paste -sd'|' -
}

function run_lint {
  # Resolved up front so a PR touching no linted Go file skips the scan entirely.
  local base_known=1 changed_files=""
  changed_files=$(get_changed_go_files) || base_known=0

  # Nothing changed, nothing to report. This also makes a run on the base branch
  # a no-op (its diff against itself is empty): this is a PR gate, not an audit.
  if [ "$base_known" -eq 1 ] && [ -z "$changed_files" ]; then
    echo -e "${GREEN} tfproviderlint: No ${LINT_PATHSPEC} files changed — skipping lint.${NC}"
    return 0
  fi

  echo -e "==> ${GREEN} Checking source code against tfproviderlint...${NC}"

  local full_output="" lint_exit=0
  full_output=$(tfproviderlint "${ANALYZERS[@]}" "$LINT_TARGET" 2>&1) || lint_exit=$?
  if [ $lint_exit -eq 0 ]; then
    echo -e "${GREEN} tfproviderlint: No issues found.${NC}"
    return 0
  fi

  # A non-zero exit is a lint *result* only if the output carries analyzer
  # findings ("path.go:line:col: CODE: msg"). Build errors, panics and bad flags
  # also exit non-zero but carry no code — filtering those would be a false pass.
  local findings
  findings=$(printf '%s\n' "$full_output" | grep -E '\.go:[0-9]+:[0-9]+: [A-Z][A-Z0-9]+:' || true)
  if [ -z "$findings" ]; then
    fail_with "tfproviderlint did not complete normally (no analyzable findings — likely a build or load error):" \
      "$full_output"
    return $lint_exit
  fi

  # Last resort: report everything and say why. Reaching this in CI is a checkout
  # problem (needs fetch-depth:0 / a reachable base), not a code problem.
  if [ "$base_known" -eq 0 ]; then
    echo -e "${RED} tfproviderlint: could not determine this PR's base commit — see the base detection notes above.${NC}"
    fail_with "Found issues (could not determine changed files, showing all):" "$findings"
    return $lint_exit
  fi

  local filtered
  filtered=$(printf '%s\n' "$findings" | grep -E "$(build_filter_pattern "$changed_files")" || true)
  if [ -z "$filtered" ]; then
    echo -e "${GREEN} tfproviderlint: No issues in changed files.${NC}"
    echo -e "${YELLOW} Note: pre-existing issues in unchanged files were ignored.${NC}"
    return 0
  fi

  fail_with "Found issues in changed files:" "$filtered"
  return $lint_exit
}

function main() {
  ensure_tfproviderlint
  run_lint
}

main
