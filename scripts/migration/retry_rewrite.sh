#!/usr/bin/env bash
# Migrate non-test package alicloud code off terraform-plugin-sdk helper/resource onto
# helper/retry + helper/id.
#
# Why: helper/resource registers a -sweep flag in init, and so does
# terraform-plugin-testing/helper/resource. A test binary linking both dies at startup
# with "flag redefined: sweep". Once no non-test file imports helper/resource, subpackage
# test binaries (alicloud/service/..., alicloud/function, ...) can link
# terraform-plugin-testing alongside package alicloud. Test files are never touched: a
# test file is linked only into its own package's binary.
#
# When: run once to migrate the tree, and re-run after every upstream merge — the merge
# reintroduces helper/resource in the files it touches, and this script converts exactly
# those. Idempotent: with nothing left to migrate it exits 0 and changes nothing.
#
# Covers both SDK import paths: v1 (terraform-plugin-sdk/helper/resource, what
# upstream/master still uses) and v2 (terraform-plugin-sdk/v2/helper/resource).
#
# Requires: grep, perl, goimports. Verify afterwards:
#   go build ./... . && go vet ./alicloud/...
set -euo pipefail
cd "$(dirname "$0")/../.."

IMPORT_V1='github.com/hashicorp/terraform-plugin-sdk/helper/resource'
IMPORT_V2='github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource'

FILES="$(grep -rlF -e "\"$IMPORT_V1\"" -e "\"$IMPORT_V2\"" --include='*.go' alicloud | grep -v '_test\.go$' || true)"

if [ -z "$FILES" ]; then
    echo "nothing to migrate: no non-test file imports helper/resource"
    exit 0
fi

echo "files to migrate: $(printf '%s\n' "$FILES" | wc -l | tr -d ' ')"

printf '%s\n' "$FILES" | while IFS= read -r f; do
    # Rewrite the eight helper/resource symbols actually in use, longest name first so
    # e.g. RetryableError is settled before the resource.Retry rule runs. Anchored on
    # "resource." on purpose: resource_alicloud_cs_serverless_kubernetes.go has a local
    # variable named resource whose fields must not be rewritten.
    perl -pi -e '
        s/\bresource\.PrefixedUniqueId\b/id.PrefixedUniqueId/g;
        s/\bresource\.UniqueId\b/id.UniqueId/g;
        s/\bresource\.RetryableError\b/retry.RetryableError/g;
        s/\bresource\.NonRetryableError\b/retry.NonRetryableError/g;
        s/\bresource\.RetryError\b/retry.RetryError/g;
        s/\bresource\.StateRefreshFunc\b/retry.StateRefreshFunc/g;
        s/\bresource\.StateChangeConf\b/retry.StateChangeConf/g;
        s/\bresource\.Retry\b/retry.Retry/g;
    ' "$f"
    # Drop the helper/resource import line (either SDK version).
    perl -ni -e 'print unless m{^\s*"github\.com/hashicorp/terraform-plugin-sdk(/v2)?/helper/resource"\s*$}' "$f"
done

# Re-add helper/retry / helper/id imports where the rewritten symbols need them. This
# also keeps files that already imported helper/retry from ending up with it twice.
# Batched: one huge goimports invocation can die silently mid-way; -n 200 keeps
# every batch observable and failures attributable.
printf '%s\n' "$FILES" | xargs -n 200 goimports -w

echo "done: verify with 'go build ./... .' and 'go vet ./alicloud/...'"
