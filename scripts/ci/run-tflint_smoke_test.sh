#!/usr/bin/env bash
# Exercise real Go package loading, including cross-file references and errors.
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
runner="$(git -C "$script_dir" rev-parse --show-toplevel)/scripts/run-tflint.sh"
task_dir="$(mktemp -d)"
trap 'rm -rf "$task_dir"' EXIT
cd "$task_dir"
git init -q -b master
git config core.hooksPath /dev/null
git config user.name 'CI Lint Test'
git config user.email 'ci-lint-test@example.invalid'
printf 'module example.com/lintscope\n\ngo 1.23.0\n' > go.mod
mkdir -p alicloud/changed alicloud/unrelated
printf 'package changed\n\nconst value = 1\n' > alicloud/changed/value.go
printf 'package changed\n\nvar result = value\n' > alicloud/changed/use.go
# This package must never be loaded for a change to its sibling.
printf 'package unrelated\n\nvar broken = undefined\n' > alicloud/unrelated/broken.go
git add .
git commit -q -m base
base_sha="$(git rev-parse HEAD)"
printf 'package changed\n\nvar result = value + 1\n' > alicloud/changed/use.go
git add .
git commit -q -m change
DIFF_BASE="$base_sha" DIFF_HEAD="$(git rev-parse HEAD)" "$runner"
# Deleting a referenced declaration must still select the surviving package.
base_sha="$(git rev-parse HEAD)"
git rm -q alicloud/changed/value.go
git commit -q -m delete-declaration
if DIFF_BASE="$base_sha" DIFF_HEAD="$(git rev-parse HEAD)" "$runner"; then
	echo 'FAIL: deleting a referenced declaration must fail analysis' >&2
	exit 1
fi
printf 'package changed\n\nconst value = 1\n' > alicloud/changed/value.go
git add .
git commit -q -m restore-declaration
# A type error in the selected package must still fail the check.
base_sha="$(git rev-parse HEAD)"
printf 'package changed\n\nvar result = undefined\n' > alicloud/changed/use.go
git add .
git commit -q -m broken-change
if DIFF_BASE="$base_sha" DIFF_HEAD="$(git rev-parse HEAD)" "$runner"; then
	echo 'FAIL: package loading errors must fail the scoped analysis' >&2
	exit 1
fi
echo 'PASS: real tfproviderlint package scoping smoke test'
