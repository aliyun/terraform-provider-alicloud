#!/usr/bin/env bash

# Assert the root alicloud production dependency graph does not depend on
# github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource, directly or
# transitively.
#
# Unlike scripts/ci/test-framework-guard.sh, which scans only newly ADDED
# import lines in the change range (a diff-level check), this queries the
# COMPLETE production dependency graph of the root provider main package with
# the Go toolchain (`go list -deps`). Both already-migrated and not-yet-migrated
# files are therefore covered: if any production source path introduces the
# forbidden package, the link graph of `.` (root main) will contain it and this
# check fails, regardless of whether the offending file is in the PR diff.
#
# Background: the test line runs terraform-plugin-testing since the flip. Both
# SDK v2 helper/resource and plugin-testing helper/resource register a -sweep
# flag in init(), so a provider binary that links both dies at startup with
# "flag redefined: sweep". Production code must use
# terraform-plugin-sdk/v2/helper/retry (flag-free) for the retry helpers, and
# tests must use terraform-plugin-testing/helper/resource for the harness.
#
# This is a real production dependency-graph query, not a vet rename and not a
# grep over a diff. The linked-graph assertion pairs with `make provider-build`
# (`go build -o /dev/null .`), which is the entry that proves the root main and
# its production dependencies compile and link. The cloudauth init probe
# (probe_init) only compiles+inits one service package via a no-such-test-case
# run and does NOT cover the root main, so it must not be cited as root-main
# build evidence.
#
# Usage:
#   make check-production-deps
#   scripts/ci/check-production-deps.sh
#   PROD_DEPS_INPUT=<file> scripts/ci/check-production-deps.sh
#       # use a fixed graph (one ImportPath per line) instead of `go list`,
#       # so the assertion can be exercised without a Go toolchain.
#   FORBIDDEN_DEPS="<path1> <path2>" scripts/ci/check-production-deps.sh
#
# Exit codes:
#   0  no forbidden package is a production dependency
#   1  at least one forbidden package is a production dependency
#   2  the dependency graph could not be queried (toolchain/config error)

set -euo pipefail

# Space-separated list of import paths that must never appear in the root
# main production dependency graph. Default: the SDK v2 helper/resource
# package whose -sweep init() clashes with the plugin-testing harness.
FORBIDDEN_DEPS="${FORBIDDEN_DEPS:-github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource}"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
project_root="$(cd "$script_dir/../.." && pwd)"

die() {
        echo "check-production-deps: $*" >&2
        exit 2
}

# Read the production dependency graph of the root main package. When
# PROD_DEPS_INPUT is set, use that file (one import path per line) instead of
# invoking `go list`, so the assertion logic can be tested without a Go
# toolchain. Otherwise query the real graph from the module root.
graph_file=""
graph_source=""
cleanup_graph=false
if [[ -n "${PROD_DEPS_INPUT:-}" ]]; then
        [[ -f "$PROD_DEPS_INPUT" ]] || die "PROD_DEPS_INPUT not found: $PROD_DEPS_INPUT"
        graph_file="$PROD_DEPS_INPUT"
        graph_source="PROD_DEPS_INPUT=$PROD_DEPS_INPUT"
else
        command -v go >/dev/null 2>&1 || die "go toolchain not found in PATH"
        graph_file="$(mktemp)"
        cleanup_graph=true
        graph_source="go list -deps -f '{{.ImportPath}}' ."
        # `go list -deps` resolves the import graph without compiling; it still
        # loads package metadata, so this check runs in CI / remote entries,
        # never as a local provider compile.
        if ! (cd "$project_root" && go list -deps -f '{{.ImportPath}}' . >"$graph_file" 2>"$graph_file.err"); then
                err="$(cat "$graph_file.err" 2>/dev/null || true)"
                rm -f "$graph_file" "$graph_file.err"
                die "failed to query production dependency graph: ${err:-unknown go list error}"
        fi
        rm -f "$graph_file.err"
fi
trap '[[ "$cleanup_graph" = true ]] && rm -f "$graph_file"' EXIT

total_deps="$(grep -c . "$graph_file" || true)"

# Match a forbidden import path exactly or as a parent package (path + "/").
# go list prints one ImportPath per line with no trailing slash, so the
# alternation `^forbidden(/|$)` covers the package and any sub-packages.
violations=0
found=""
for forbidden in $FORBIDDEN_DEPS; do
        [[ -n "$forbidden" ]] || continue
        # Escape regex metacharacters in the import path before anchoring.
        esc="$(printf '%s' "$forbidden" | sed 's/[.[\*^$()+?{|\\]/\\&/g')"
        matches="$(grep -E "^${esc}(/|$)" "$graph_file" || true)"
        if [[ -n "$matches" ]]; then
                violations=$((violations + 1))
                found="${found}${matches}"$'\n'
        fi
done

echo "check-production-deps: source=${graph_source}"
echo "check-production-deps: production dependency count=${total_deps}"
echo "check-production-deps: forbidden packages=${FORBIDDEN_DEPS}"

if ((violations > 0)); then
        echo "check-production-deps: FAILED — forbidden production dependency found:" >&2
        printf '%s\n' "$found" | sed '/^$/d; s/^/  /' >&2
        echo "check-production-deps: production code must not import" >&2
        echo "  github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource" >&2
        echo "check-production-deps: use terraform-plugin-sdk/v2/helper/retry for retry helpers," >&2
        echo "  and terraform-plugin-testing/helper/resource for the test harness." >&2
        exit 1
fi

echo "check-production-deps: PASS — no forbidden SDK v2 helper/resource production dependency"
