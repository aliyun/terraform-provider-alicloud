#!/usr/bin/env bash

# Route-regression fixture tests for detect-build-input-changes.sh. Feeds
# fixed changed-file lists via BUILD_INPUT_FILES (no git needed) and asserts
# the detector triggers on build-input changes (CI workflows, GNUmakefile,
# scripts/ci/*, the local-ci-check and *check.sh scripts) and does NOT trigger
# for doc-only or product-Go changes (the product-Go path is covered separately
# by route.outputs.relevant).

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
check="$script_dir/detect-build-input-changes.sh"

die() {
        echo "FAIL: $*" >&2
        exit 1
}

case_test() {
        local name="$1" files="$2" expect="$3"
        local tmp
        tmp="$(mktemp)"
        printf '%s\n' "$files" >"$tmp"
        set +e
        output="$(BUILD_INPUT_FILES="$tmp" "$check" 2>&1)"
        rc=$?
        set -e
        rm -f "$tmp"
        [[ "$rc" -eq 0 ]] || die "$name: expected exit 0, got $rc; output: $output"
        [[ "$output" == "$expect" ]] || die "$name: expected '$expect', got '$output'"
        echo "ok - $name"
}

# A CI workflow change triggers the build + dep-graph check.
case_test "workflow file triggers" \
        $'.github/workflows/acctest-terraform-basic.yml\n' "true"

# GNUmakefile change triggers.
case_test "makefile triggers" $'GNUmakefile\n' "true"

# Any file under scripts/ci/ triggers (CI/check scripts).
case_test "ci script triggers" \
        $'scripts/ci/check-production-deps.sh\nscripts/ci/detect-build-input-changes.sh\n' "true"

# The local-ci-check entry script triggers.
case_test "local-ci-check triggers" $'scripts/local-ci-check.sh\n' "true"

# A vet/gofmt/goimports/errcheck helper script triggers.
case_test "check helper script triggers" $'scripts/vetcheck.sh\n' "true"

# A doc-only change does NOT trigger (doc PRs are not build-input changes).
case_test "docs only does not trigger" \
        $'website/docs/r/alicloud_vpc.html.markdown\n' "false"

# A product Go file change does NOT trigger the detector: product-Go PRs reach
# ProductionBuild via route.outputs.relevant, not via this detector. The
# detector only covers the CI/Make/script path that relevance skips.
case_test "product go file does not trigger detector" \
        $'alicloud/resource_alicloud_vpc.go\n' "false"

# Mixed doc + CI change still triggers.
case_test "mixed doc and ci triggers" \
        $'website/docs/r/x.md\nscripts/ci/vetcheck.sh\n' "true"

# Empty / blank input does not trigger.
case_test "empty changes do not trigger" $'\n' "false"

# A look-alike path that merely contains a monitored prefix as a substring but
# is not under it does not trigger (prefix anchoring).
case_test "look-alike suffix does not trigger" \
        $'docs/GNUmakefile-notes.md\n' "false"

echo "all detect-build-input-changes tests passed"
