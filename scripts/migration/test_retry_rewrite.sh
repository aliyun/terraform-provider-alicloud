#!/usr/bin/env bash
# Migrate the retry symbols in package alicloud's *test* files off
# terraform-plugin-sdk/v2/helper/resource onto helper/retry.
#
# This is the test-file half of what retry_rewrite.sh (Route 1, #10255) did to the ~1900
# non-test files. It touches nothing but the retry names: the acceptance-test harness
# (resource.Test, TestCase, TestStep, ComposeTestCheckFunc, ParallelTest, AddTestSweepers,
# the TestMain family) keeps coming from helper/resource, and the import stays.
#
# Why it is safe: in terraform-plugin-sdk/v2 the retry names in helper/resource are not a
# second implementation, they are aliases and forwarders to helper/retry
# (helper/resource/aliases.go) —
#
#     type RetryError = retry.RetryError            // and RetryFunc, StateRefreshFunc,
#     type StateChangeConf = retry.StateChangeConf  // NotFoundError, TimeoutError, ...
#     func Retry(t time.Duration, f RetryFunc) error { return retry.Retry(t, f) }
#
# so resource.Retry and retry.Retry are the same function over the same types, and this
# rewrite cannot change behaviour. Each alias also carries "Deprecated: Use helper/retry
# package instead", so the spelling this produces is the upstream-supported one.
#
# It does NOT resolve the -sweep flag collision: test files still import helper/resource
# for the harness, so a package alicloud test binary still cannot also link
# terraform-plugin-testing/helper/resource. That swap is a separate change; this one
# shrinks it to an import line per file.
#
# Selection is by import rather than by filename, so a file that only mentions
# resource.Retry in a comment is not a candidate and needs no exclusion list. Only the v2
# import path is handled — nothing in the tree imports v1, whose retry types are its own
# rather than aliases of helper/retry.
#
# When: run once, and re-run after every upstream merge — the merge reintroduces
# resource.Retry in the test files it touches, and this converts exactly those. Idempotent:
# with nothing left to migrate it exits 0 having changed nothing.
#
# Requires: grep, perl, gofmt — the one import added has a single possible path, so there is
# nothing for goimports to resolve. Verify afterwards:
#   go vet ./alicloud/...
set -euo pipefail
cd "$(dirname "$0")/../.."

RESOURCE_IMPORT='github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource'
RETRY_IMPORT='github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry'

# Every retry name helper/resource aliases or forwards. Only five occur in the tree today
# (Retry, RetryableError, NonRetryableError, RetryError, StateRefreshFunc); the rest are
# listed so a later upstream merge that introduces one is converted too.
RETRY_IDS='RetryContext|RetryableError|NonRetryableError|RetryError|RetryFunc|StateRefreshFunc|StateChangeConf|NotFoundError|UnexpectedStateError|TimeoutError|Retry'

# Candidates must both import helper/resource and spell a retry symbol against it.
FILES="$(grep -rlF "\"$RESOURCE_IMPORT\"" --include='*_test.go' alicloud |
    xargs grep -lE "\bresource\.($RETRY_IDS)\b" || true)"

if [ -z "$FILES" ]; then
    echo "nothing to migrate: no test file spells retry against helper/resource"
    exit 0
fi

echo "files to migrate: $(printf '%s\n' "$FILES" | wc -l | tr -d ' ')"

dropped=0
while IFS= read -r f; do
    [ -n "$f" ] || continue

    # Longest name first, matching retry_rewrite.sh. The trailing \b already stops
    # resource.Retry from eating resource.RetryableError, so the order is parity with the
    # sibling script, not correctness. Anchored on "resource." so only the qualified form
    # moves: package alicloud has its own unqualified NotFoundError() helper, which must
    # not be touched.
    perl -pi -e '
        s/\bresource\.UnexpectedStateError\b/retry.UnexpectedStateError/g;
        s/\bresource\.NonRetryableError\b/retry.NonRetryableError/g;
        s/\bresource\.StateRefreshFunc\b/retry.StateRefreshFunc/g;
        s/\bresource\.StateChangeConf\b/retry.StateChangeConf/g;
        s/\bresource\.RetryableError\b/retry.RetryableError/g;
        s/\bresource\.NotFoundError\b/retry.NotFoundError/g;
        s/\bresource\.TimeoutError\b/retry.TimeoutError/g;
        s/\bresource\.RetryContext\b/retry.RetryContext/g;
        s/\bresource\.RetryError\b/retry.RetryError/g;
        s/\bresource\.RetryFunc\b/retry.RetryFunc/g;
        s/\bresource\.Retry\b/retry.Retry/g;
    ' "$f"

    # Add helper/retry next to helper/resource unless it is already there, so a re-run
    # cannot duplicate it. "resource" < "retry" < "schema", so inserting on the line after
    # helper/resource is already in sorted order; the gofmt pass below only has to fix
    # alignment.
    if ! grep -qF "\"$RETRY_IMPORT\"" "$f"; then
        perl -pi -e "s{^(\\s*)\"\\Q$RESOURCE_IMPORT\\E\"(\\s*)\$}{\$1\"$RESOURCE_IMPORT\"\n\$1\"$RETRY_IMPORT\"\$2}" "$f"
    fi

    # A file whose only use of helper/resource was retry now has an unused import. None of
    # the current candidates are in that position — they all use the harness — but a future
    # merge could add one, and an unused import does not compile. Line comments are stripped
    # first so a commented-out selector cannot hold the import alive.
    #
    # Decided inside one perl process on purpose: the obvious `perl … | grep -q` form is
    # wrong under `set -o pipefail`, because grep -q exits on the first match and perl then
    # dies of SIGPIPE, so a pipeline that *found* a selector reports failure and the import
    # gets dropped from exactly the files that still need it. Here a match means exit 1.
    if perl -ne 's{//.*}{}; exit 1 if /\bresource\./' "$f"; then
        perl -ni -e "print unless m{^\\s*\"\\Q$RESOURCE_IMPORT\\E\"\\s*\$}" "$f"
        dropped=$((dropped + 1))
    fi
done <<EOF
$FILES
EOF

# Re-align the import groups the insert disturbed. Batched: one huge invocation can die
# silently mid-way; -n 200 keeps every batch observable and failures attributable.
printf '%s\n' "$FILES" | grep -v '^$' | xargs -n 200 gofmt -w

if [ "$dropped" -gt 0 ]; then
    echo "dropped the now-unused helper/resource import from ${dropped} file(s)"
fi
echo "done: verify with 'go vet ./alicloud/...'"
