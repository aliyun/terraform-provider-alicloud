---
name: terraform-changelog
description: >-
  Generate CHANGELOG.md entries for terraform-provider-alicloud by aggregating merged GitHub PRs
  since the last release, cut the release (stamp the date, bump providerVersion, open the next
  Unreleased block), and submit a PR with proper conventions (single commit, rebase, title CHANGELOG).
  Categorizes PRs into New Resources / Data Sources, ENHANCEMENTS, and BUG FIXES.
  Trigger on: generate changelog, update changelog, changelog entry, release notes,
  cut release, 发布版本, 生成changelog, 更新changelog.
metadata:
  version: "1.1.0"
  domain: terraform-provider
  triggers: generate changelog, update changelog, changelog entry, release notes, cut release, 发布版本, 生成changelog, 更新changelog
---

# Changelog Generator

Aggregate merged PRs from `aliyun/terraform-provider-alicloud` and produce CHANGELOG.md entries that match the existing house style.

## Prerequisites

All file paths in this skill (`CHANGELOG.md`, `alicloud/connectivity/client.go`, etc.) are relative to the **terraform-provider-alicloud repository root**. Before running any step, resolve the local path:

```bash
PROVIDER_DIR="$(bootstrap/workspace.sh dir terraform_provider)"
```

Then work inside `$PROVIDER_DIR`. If `bootstrap/workspace.sh` reports the workspace is missing, stop and escalate (`missing_capability`) -- do not guess paths.

This skill only reads GitHub via `gh pr list` (no write operations), so `gh` can be used directly without `bootstrap/github-identity.sh`.

## When to Use

- A new release is about to be cut and the maintainer needs to fill in the `## 1.X.Y (Unreleased)` block
- The user provides a since-date (e.g., the date of the last released entry) or a GitHub Pulls search URL
- The user says "generate changelog" / "update changelog" / "release notes"

## Inputs

| Input | How to obtain | Example |
|---|---|---|
| `since_date` | Prefer the user's value. If absent, parse the most recent **dated** version header in `CHANGELOG.md` and use that date. | `2026-04-26` |
| `repo` | Fixed: `aliyun/terraform-provider-alicloud` | -- |
| `unreleased_header` | The first line of CHANGELOG.md if it matches `^## \d+\.\d+\.\d+ \(Unreleased\)$`. If missing, propose a new version number = bump minor of the most recent released entry. | `## 1.278.0 (Unreleased)` |

If the user supplies a URL like `https://github.com/aliyun/terraform-provider-alicloud/pulls?q=...merged%3A%3E%3D2026-04-26`, parse the `merged:>=YYYY-MM-DD` filter to derive `since_date`. Ignore irrelevant filters like `language:javascript`.

## Step 1: Fetch Merged PRs

```bash
gh pr list \
  --repo aliyun/terraform-provider-alicloud \
  --state merged \
  --search "merged:>=<since_date>" \
  --limit 100 \
  --json number,title,mergedAt,labels,url,author \
  > /tmp/changelog-prs.json
```

Sort the result ascending by PR number -- the existing CHANGELOG lists entries roughly in PR-number order within each section.

If 100 results are returned, there may be more PRs. Increase `--limit` or paginate to ensure completeness.

## Step 2: Filter Out Noise

Drop any PR matching **any** of these:

1. **Already in CHANGELOG.md** -- `grep -c "issues/<num>)" CHANGELOG.md` returns >= 1
2. **Release-cut PRs** -- title is exactly `CHANGELOG` / `changelog` / `Release v*`. Such PRs typically modify only `CHANGELOG.md` plus the version constant in `alicloud/connectivity/client.go`; do **not** require "only CHANGELOG.md changed" -- title is the reliable signal.
3. **Reverted/draft** -- closed without merge (already excluded by `--state merged`)
4. **Bot-only mechanical edits** if they don't reflect user-visible behavior (rare; flag for review rather than silently drop)

For ambiguous cases, **list them under "needs review"** at the end of the report rather than silently dropping.

## Step 3: Categorize Each PR

Use **labels first, title as fallback**:

| Category | Label | Title fallback (case-insensitive) |
|---|---|---|
| **New Resource** | `new-resource` | starts with `New Resource:` or `New Resource/` |
| **New Data Source** | `new-data-source` | starts with `New Data Source:` or `New Data-Source:` |
| **Bug Fix** | `bug` | first segment after `:` starts with `Fix`, `Fixed`, `Fixes`, or title starts with `fix:` / `fix(...)`: |
| **Enhancement** | `enhancement`, `documentation`, `technical-debt` (default) | everything else |

Caveats:
- A PR can carry both `new-resource` and `enhancement` labels (e.g., it adds a resource AND modifies sibling resources). Place it under **New Resource** as the primary category -- sibling-resource changes are typically described in the PR title and don't need a duplicate ENHANCEMENT line.
- A PR titled `New Resource:` but missing the label still counts as a new resource -- the title is authoritative for this case.
- Title starting with Conventional Commit prefixes (`feat(scope):`, `fix(scope):`, `chore(scope):`) -- strip the prefix when formatting (see Step 4) and choose the category from the verb (`feat`/`chore` -> enhancement, `fix` -> bug fix).
- A PR labeled `breaking-change` is **not** a separate category in this CHANGELOG -- but the skill MUST surface it to the user with a warning note ("breaking change -- confirm wording") before writing the file. Breaking changes still go in ENHANCEMENTS (or BUG FIXES if also a bugfix) but must be flagged in the Step 8 report so the maintainer can add extra notes manually.

## Step 4: Format Each Entry

### New Resource / Data Source

Extract the resource name from the title using regex `alicloud_[a-z0-9_]+`. Format:

```
- **New Resource:** `alicloud_xxx` ([#NNNN](https://github.com/aliyun/terraform-provider-alicloud/issues/NNNN))
- **New Data Source:** `alicloud_xxx` ([#NNNN](https://github.com/aliyun/terraform-provider-alicloud/issues/NNNN))
```

Note: the existing CHANGELOG uses `/issues/N`, not `/pull/N` -- keep that convention.

### Enhancement / Bug Fix

1. **Strip Conventional Commit prefix** if present:
   - `feat(oss): foo` -> keep "foo" but rewrite the prefix to a CHANGELOG-style scope: `provider:` (if scope is generic) or `resource/alicloud_<scope>:` if scope is a resource. When unsure, use `provider:` and flag for review.
2. **Normalize doc prefix**: `doc:` -> `docs:` (existing CHANGELOG always uses `docs:`).
3. **Preserve the rest of the title verbatim** -- multi-resource titles separated by `;` are kept on one line.
4. **Append a period** if the result doesn't already end in `.`, `!`, or `?`.
5. **Append the issue link**:

```
- <normalized title>. ([#NNNN](https://github.com/aliyun/terraform-provider-alicloud/issues/NNNN))
```

## Step 5: Order Within Each Section

Within each section, sort by PR number ascending. This matches the existing CHANGELOG's loose ordering.

The section ordering inside a release block is fixed:

```
## <version> (<status>)

- **New Resource:** ... (one bullet per new resource)
- **New Data Source:** ... (one bullet per new data source)

ENHANCEMENTS:

- ...

BUG FIXES:

- ...
```

Omit a section header (`ENHANCEMENTS:` / `BUG FIXES:`) entirely if it has no entries.

## Step 6: Write to CHANGELOG.md

1. **Locate the Unreleased block** at the top of the file. Pattern: `^## \d+\.\d+\.\d+ \(Unreleased\)$`.
2. If the Unreleased block is empty (next non-blank line is the previous release header), insert the generated content right after the Unreleased header, separated by a blank line, with another blank line before the previous release header.
3. If the Unreleased block already has entries, **merge by section**: add new bullets in sort order to existing sections, creating sections that don't yet exist.
4. **Never overwrite or remove existing entries** -- additive only.

Use the `Edit` tool with a unique anchor (the Unreleased header itself + the next existing release header) to make the edit safe.

## Step 7: Self-Check

After writing, verify:

| Check | How |
|---|---|
| All target PR numbers appear in the file | `for n in <list>; do grep -q "issues/$n)" CHANGELOG.md \|\| echo MISSING $n; done` |
| No duplicate entries | `grep -c "issues/<n>)" CHANGELOG.md` returns 1 for each |
| Section headers correct casing | `grep -n "^ENHANCEMENTS:\|^BUG FIXES:" CHANGELOG.md \| head` |
| File still parses as Markdown | open the file, scan for ragged headers / orphan bullets |

If any check fails, revert with `git checkout CHANGELOG.md` and report the error to the user.

## Step 8: Report

Output a summary:

```
## Changelog Update Summary

Since: <since_date>
Total PRs found: <N>
  - Already in CHANGELOG: <N1>
  - Release-cut PRs skipped: <N2>
  - Added to "<unreleased_header>": <N3>
    * New Resource: <count>
    * New Data Source: <count>
    * Enhancements: <count>
    * Bug Fixes: <count>

Items needing manual review:
  - #<num>: <reason>
  - ...
```

## Step 9: Cut the Release (default)

**Always perform this step** unless the user explicitly says they only want the Unreleased block populated without cutting (e.g., "don't cut the release" / "just add entries"). After populating the Unreleased block (Steps 1-8), ask the user to confirm the release date (default: today) and proceed. This step can also run on its own when the Unreleased block is already populated by hand.

### 9.1 Inputs

| Input | How to obtain | Example |
|---|---|---|
| `release_version` | Version on the current Unreleased line -- `## <release_version> (Unreleased)` | `1.278.0` |
| `release_date` | User-supplied; otherwise today's date in the local timezone | `May 09, 2026` |
| `next_version` | Bump the **minor** segment of `release_version`. Patch releases (e.g., `1.272.1`) bump back to a `.0` minor -- confirm with the user when in doubt. | `1.278.0` -> `1.279.0` |

### 9.2 Date format

Existing CHANGELOG dates are `Month DD, YYYY`. Day padding is inconsistent in the existing file (`April 03, 2026` vs `March 2, 2026`) -- use whatever the user provides verbatim, otherwise default to `Month D, YYYY` with no leading zero.

### 9.3 Edits

All paths below are relative to the provider repo root (`$PROVIDER_DIR`).

**A. CHANGELOG.md** -- replace the Unreleased line with the dated header and prepend a fresh Unreleased block above it. Both go on adjacent lines (no blank line between), matching the file's existing convention:

```diff
-## <release_version> (Unreleased)
+## <next_version> (Unreleased)
+## <release_version> (<release_date>)
```

**B. alicloud/connectivity/client.go** -- update the `providerVersion` constant to `release_version`. Anchor pattern:

```go
// The main version number that is being run at the moment.
var providerVersion = "<release_version>"
```

Apply with a single `Edit`. The previous value is whatever currently sits in that line; do **not** assume it is one minor below `release_version` (a patch release means it could already match).

### 9.4 Verify

| Check | How |
|---|---|
| Unreleased moved to `next_version` | `head -2 CHANGELOG.md` shows `## <next_version> (Unreleased)` then `## <release_version> (<release_date>)` |
| `providerVersion` matches the released version | `grep '^var providerVersion' alicloud/connectivity/client.go` shows `release_version` |
| Diff is small and surgical | `git diff --stat` shows only `CHANGELOG.md` (a few lines) and `alicloud/connectivity/client.go` (1 line) |

Pre-existing diagnostics in unrelated files (e.g., `go.mod` vendor warnings) are out of scope -- flag them in the report but do not attempt to fix them as part of the release cut.

### 9.5 Do NOT

- Do **not** tag -- tagging is the maintainer's job
- Do **not** modify any other version-bearing file unless the user names it explicitly. The only auto-edit targets are `CHANGELOG.md` and `alicloud/connectivity/client.go`

## Step 10: Submit PR

After all edits (CHANGELOG entries + release cut) are verified, submit a PR to the upstream repo. Follow these conventions strictly -- the repo CI enforces them:

### 10.1 Commit

Squash all changes into a **single commit**. The repo has a `Pull Request Max Commits` CI check that fails if the PR contains more than one commit.

```bash
git checkout -b changelog/<release_version>
git add CHANGELOG.md alicloud/connectivity/client.go
git commit -m "CHANGELOG"
```

Commit message must be exactly `CHANGELOG` -- no extra description needed. Do **not** add `Co-Authored-By` or any AI signature to the commit.

### 10.2 Rebase

Before pushing, **always rebase** onto the latest upstream master to avoid merge conflicts and ensure CI runs against current code:

```bash
git fetch origin master
git rebase origin/master
```

If rebase produces conflicts, resolve them (CHANGELOG conflicts are typically additive and straightforward), then continue.

### 10.3 Push and Create PR

Push via `bootstrap/github-identity.sh` (required for Jarvis identity). The PR head must be on `api-tool-agent`'s fork:

```bash
bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud changelog/<version> changelog/<version>
```

Create the PR with title **exactly `CHANGELOG`** -- this matches the repo's convention and avoids triggering the `Pull Request Title` CI check:

```bash
bootstrap/github-identity.sh gh pr create \
  --repo aliyun/terraform-provider-alicloud \
  --head api-tool-agent:changelog/<version> \
  --base master \
  --title "CHANGELOG" \
  --body "<brief summary: N new resources, N enhancements, N bug fixes>"
```

### 10.4 Post-push fixes

If CI fails (e.g., because of a force-push race or new commits landed on master):

1. Rebase again: `git fetch origin master && git rebase origin/master`
2. Force-push: `bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud +changelog/<version> changelog/<version>` (note the `+` prefix for force push)

Always verify the PR has exactly 1 commit after any rebase/squash operation.

## Important Notes

1. **Single commit, title `CHANGELOG`** -- the repo CI enforces max 1 commit per PR and checks the PR title. Always squash and use the exact title.
2. **Issue link format**: always `/issues/N`, not `/pull/N`. GitHub redirects, but stay consistent with the existing CHANGELOG
3. **Don't invent entries**: every bullet must trace back to a real merged PR in the JSON output
4. **Don't paraphrase titles**: the maintainer expects the PR title (lightly normalized). Only rewrite Conventional Commit prefixes -- leave the rest verbatim
5. **Breaking changes**: still go in ENHANCEMENTS unless they are also bugfixes -- but flag them in Step 8 so the maintainer can add an extra note manually
6. **No trailing blank line spam**: match the surrounding file's blank-line style (one blank line between sections, two between release blocks is the existing convention)

## Acceptance Criteria

1. All non-trivial merged PRs since `since_date` appear in CHANGELOG.md exactly once
2. Sections, ordering, and bullet formatting match existing entries
3. PR links use `/issues/N` and resolve correctly
4. Release-cut PRs and PRs already listed are excluded
5. Breaking changes and ambiguous items are surfaced in the report (not silently merged)
6. Release is cut by default (Unreleased stamped with date, providerVersion bumped, next Unreleased opened)
7. PR submitted with exactly 1 commit, title `CHANGELOG`, rebased on latest master, pushed to `api-tool-agent` fork

## Example Workflow

User: "generate changelog since last release"

1. Resolve provider path: `PROVIDER_DIR="$(bootstrap/workspace.sh dir terraform_provider)"`
2. Read `$PROVIDER_DIR/CHANGELOG.md` line 1 -> `## 1.278.0 (Unreleased)` (the target block)
3. Parse the most recent dated header to derive `since_date = 2026-04-26`
4. `gh pr list --repo aliyun/terraform-provider-alicloud --state merged --search "merged:>=2026-04-26" --json ...`
5. Drop already-in-CHANGELOG PRs and the release-cut PR (title `CHANGELOG`)
6. Categorize the rest by labels + title
7. Format entries; insert under the Unreleased header in `$PROVIDER_DIR/CHANGELOG.md`
8. Run self-check, report summary with any items flagged for manual review
