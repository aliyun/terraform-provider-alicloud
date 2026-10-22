#!/usr/bin/env bash
# Flip the root monolith test files (package alicloud, alicloud/*_test.go) from the
# SDKv2 test framework onto terraform-plugin-testing.
#
# Why: SDKv2 helper/resource and terraform-plugin-testing/helper/resource each
# register a -sweep flag in init; the root test binary can link only one of them.
# The split subpackages already run plugin-testing (framework sunk into
# alicloud/acctest), so the root binary is the last SDKv2 helper/resource
# consumer. One directory = one test binary: the flip must be atomic, never
# half-applied.
#
# Scope, deliberately narrow:
#   - alicloud/*_test.go in the repository root of the package ONLY (no
#     alicloud/service/... subpackage files: the 7 ephemeral tests alias
#     sdkacctest deliberately and disk_unit_test.go pairs SDKv2 terraform with
#     helper/schema — both stay).
#   - Exactly the three runner-framework import paths: helper/resource,
#     helper/acctest, terraform. NOT helper/retry and NOT helper/schema: both
#     stay on SDKv2 (no plugin-testing counterpart / production schema).
#   - Non-test files are never touched: common.go keeps SDKv2 terraform because
#     newInstanceDiff feeds SDKv2 helper/schema (schema.InternalMap(...).Data).
#   - helper/resource and helper/acctest flip in every matching file. terraform
#     flips only in files that carry neither a plugin-testing runner import nor
#     a schema-side terraform symbol (InstanceState / InstanceDiff /
#     ResourceAttrDiff / NewInstanceDiff / NewResourceConfigRaw): both mark a
#     file that must keep the SDKv2 terraform import — schema-side types for
#     newInstanceDiff and schemaMap.Data — and survive a re-run untouched.
#
# After the flip, a correction pass restores the SDKv2 terraform types in
# files that feed them into SDKv2 helper/schema APIs (schema.InternalMap(...).Data,
# Resource.Data/Apply/Diff/RefreshWithoutUpgrade): schema-side-only files keep
# the SDKv2 terraform import; mixed files (runner-side terraform.State /
# ImportStateCheck closures too) keep plugin-testing terraform under the plain
# name and take the SDKv2 import as sdktf, with the schema-side references
# rewritten to sdktf.*. 'go vet ./alicloud/' stays the arbiter: data crossing
# forks through variables rather than line-local calls surfaces there. Known
# hand-fix sites vet drove (re-check if a merge reintroduces their old shape):
# resource_alicloud_kvstore_instance_test.go (import state crosses to .Diff
# through a shallow copy) and
# resource_alicloud_vpn_gateway_vpn_attachment_test.go
# (vpnAttachmentStateAttrByIndex takes map[string]string so either fork feeds it).
#
# When: run once to flip the tree, and re-run after every upstream merge — a
# merge that copies SDKv2-era test files reintroduces the old paths and this
# script converts exactly those. Idempotent: with nothing left to flip it exits
# 0 and changes nothing.
#
# Requires: grep, perl, gofmt. Verify afterwards:
#   go vet ./alicloud/          # typechecks every _test.go without linking
#   go build ./alicloud/... .
set -euo pipefail
cd "$(dirname "$0")/../.."

SDK_RESOURCE='github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource'
SDK_ACCTEST='github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest'
SDK_TERRAFORM='github.com/hashicorp/terraform-plugin-sdk/v2/terraform'
TF_RESOURCE='github.com/hashicorp/terraform-plugin-testing/helper/resource'
TF_ACCTEST='github.com/hashicorp/terraform-plugin-testing/helper/acctest'
TF_TERRAFORM='github.com/hashicorp/terraform-plugin-testing/terraform'

# Root-of-package test files only: one path segment below alicloud/. Every
# pattern must carry its own -e: next to -e patterns a bare operand is a FILE
# to grep, not a pattern, and the path is silently dropped (the nonexistent
# "file" error is exactly what a 2>/dev/null would swallow).
RA_FILES="$(grep -rlF -e "\"$SDK_RESOURCE\"" -e "\"$SDK_ACCTEST\"" \
    --include='*_test.go' alicloud | grep -E '^alicloud/[^/]+$' || true)"

# helper/resource and helper/acctest have no legitimate SDKv2 future: flip
# every match. terraform is different — see the header — so flip it only in
# files that carry neither a plugin-testing runner import nor a schema-side
# terraform symbol: without the symbol check, files that test the SDKv2
# production stack directly (common_test.go, sdkv2_wrap_real_resource_test.go)
# look like fresh pre-flip files, since terraform may be their only framework
# import.
TF_FILES="$(grep -rlF "\"$SDK_TERRAFORM\"" --include='*_test.go' alicloud \
    | grep -E '^alicloud/[^/]+$' \
    | while IFS= read -r f; do
        [ -n "$f" ] || continue
        grep -qE 'terraform\.(InstanceState|InstanceDiff|ResourceAttrDiff|NewInstanceDiff|NewResourceConfigRaw)' "$f" && continue
        grep -qF -e "\"$TF_RESOURCE\"" -e "\"$TF_ACCTEST\"" -e "\"$TF_TERRAFORM\"" "$f" \
            || printf '%s\n' "$f"
    done || true)"

FILES="$(printf '%s\n%s\n' "$RA_FILES" "$TF_FILES" | grep -v '^$' | sort -u || true)"

if [ -n "$FILES" ]; then
    echo "files to flip: $(printf '%s\n' "$FILES" | wc -l | tr -d ' ')"
    printf '%s\n' "$FILES" | while IFS= read -r f; do
        perl -pi -e '
            s{"github\.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"}{"github.com/hashicorp/terraform-plugin-testing/helper/resource"}g;
            s{"github\.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"}{"github.com/hashicorp/terraform-plugin-testing/helper/acctest"}g;
            s{"github\.com/hashicorp/terraform-plugin-sdk/v2/terraform"}{"github.com/hashicorp/terraform-plugin-testing/terraform"}g;
        ' "$f"
    done
else
    echo "nothing to flip: no root test file imports the SDKv2 test framework unmigrated"
fi

# --- Schema-side terraform consumers ---------------------------------------
# The flip above moved every listed file to plugin-testing terraform, but files
# that feed terraform types into SDKv2 helper/schema APIs need the SDKv2 types
# at those call sites (see header). Two shapes:
#   - schema-side only: the SDKv2 terraform import comes back;
#   - mixed (runner-side terraform.State / []*terraform.InstanceState from
#     ImportStateCheck): plugin-testing terraform keeps the plain name, the
#     SDKv2 import joins as sdktf, and schema-side references rewrite to
#     sdktf.* — the []*terraform.InstanceState ImportStateCheck parameter line
#     is runner-side and is left alone.
# The sdktf import is added only when a schema-side reference outside the
# ImportStateCheck lines exists, so an ImportStateCheck-only file never gains
# an unused import.
CORRECTED="$(grep -rlF "\"$TF_TERRAFORM\"" --include='*_test.go' alicloud \
    | grep -E '^alicloud/[^/]+$' \
    | while IFS= read -r c; do
        [ -n "$c" ] || continue
        grep -qE 'terraform\.(InstanceState|InstanceDiff|ResourceAttrDiff|NewInstanceDiff|NewResourceConfigRaw)' "$c" || continue
        printf '%s\n' "$c"
    done || true)"

if [ -n "$CORRECTED" ]; then
    echo "schema-side consumers to correct: $(printf '%s\n' "$CORRECTED" | wc -l | tr -d ' ')"
    printf '%s\n' "$CORRECTED" | while IFS= read -r f; do
        if grep -qE 'terraform\.State[^A-Za-z]|terraform\.State$|\[\]\*terraform\.InstanceState' "$f"; then
            if grep -vE '\[\]\*terraform\.InstanceState' "$f" \
                | grep -qE 'terraform\.(NewInstanceDiff|ResourceAttrDiff|NewResourceConfigRaw|InstanceDiff|InstanceState)'; then
                grep -qF 'sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"' "$f" || \
                    perl -0777 -pi -e '
                        s{^(\t*)"github\.com/hashicorp/terraform-plugin-testing/terraform"\n}{$1sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"\n$1"github.com/hashicorp/terraform-plugin-testing/terraform"\n}m;
                    ' "$f"
                perl -pi -e '
                    next if /\[\]\*terraform\.InstanceState/;
                    s{terraform\.(NewInstanceDiff|ResourceAttrDiff|NewResourceConfigRaw|InstanceDiff|InstanceState)}{sdktf.$1}g;
                ' "$f"
            fi
        else
            perl -pi -e 's{"github\.com/hashicorp/terraform-plugin-testing/terraform"}{"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"}' "$f"
        fi
    done
fi

# Re-sort import blocks: the flipped paths move after the SDKv2 paths that stay
# (helper/retry, helper/schema), and the correction pass re-homes terraform.
# No import is added or removed beyond the sdktf line, so gofmt's block sort is
# sufficient — goimports' module resolution is not needed.
# Batched so failures stay attributable; gofmt accepts many files per run.
FMT_TARGETS="$(printf '%s\n%s\n' "$FILES" "$CORRECTED" | grep -v '^$' | sort -u || true)"
if [ -n "$FMT_TARGETS" ]; then
    printf '%s\n' "$FMT_TARGETS" | xargs -n 200 gofmt -w
fi

echo "done: verify with 'go vet ./alicloud/' and 'go build ./alicloud/... .'"
