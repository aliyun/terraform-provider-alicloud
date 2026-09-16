#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

awk '
  /^      - name: Restore coverage checker for older pull requests$/ { step = 1; next }
  step && /^        run: \|$/ { body = 1; next }
  body && /^          / { sub(/^          /, ""); print; next }
  body { exit }
' "$repo_root/.github/workflows/collection-order-coverage.yml" > "$tmp/restore.sh"
test -s "$tmp/restore.sh"

git init -q "$tmp/repo"
cd "$tmp/repo"
git config core.hooksPath /dev/null
git config user.name "CI Test"
git config user.email "ci@example.invalid"
printf 'old schema\n' > schema
git add schema
git commit -qm initial
initial="$(git rev-parse HEAD)"
git commit --allow-empty -qm old-pr
old_pr="$(git rev-parse HEAD)"

git checkout -q --detach "$initial"
mkdir -p scripts/collection-order
printf 'base checker\n' > scripts/collection-order/marker
git add scripts
git commit -qm add-checker
guard_base="$(git rev-parse HEAD)"

# An older PR receives only the checker; HEAD and tracked files stay intact.
git checkout -q --detach "$old_pr"
export DIFF_BASE="$initial" DIFF_HEAD="$old_pr" PR_BASE_SHA="$guard_base"
bash "$tmp/restore.sh"
test "$(cat scripts/collection-order/marker)" = "base checker"
test "$(git rev-parse HEAD)" = "$old_pr"
test -z "$(git status --porcelain --untracked-files=no)"
rm -rf scripts/collection-order

# A PR containing its own checker keeps that version.
git checkout -q --detach "$guard_base"
printf 'PR checker\n' > scripts/collection-order/marker
git add scripts
git commit -qm change-checker
export DIFF_BASE="$guard_base" DIFF_HEAD="$(git rev-parse HEAD)"
bash "$tmp/restore.sh"
test "$(cat scripts/collection-order/marker)" = "PR checker"

# Removing the checker cannot trigger restoration.
git rm -qr scripts/collection-order
git commit -qm remove-checker
export DIFF_HEAD="$(git rev-parse HEAD)"
if bash "$tmp/restore.sh"; then
  echo "Deleted checker was silently restored" >&2
  exit 1
fi
test ! -e scripts/collection-order

# Missing fallback content must fail, including a failed archive pipeline.
git checkout -q --detach "$old_pr"
export DIFF_BASE="$initial" DIFF_HEAD="$old_pr" PR_BASE_SHA="$initial"
if bash "$tmp/restore.sh"; then
  echo "Missing fallback unexpectedly succeeded" >&2
  exit 1
fi

echo "TypeList Order Coverage compatibility checks passed"
