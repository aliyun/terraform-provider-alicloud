# Release Cutting + PR Submission

> 本 reference 是 terraform-changelog skill Step 9/10「切版本 + 提 PR」的详细流程,主文件引用至此。默认在 Step 8 report 之后进入。

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
# 注意:并发 triage 的 sync-provider.sh 会对主目录 reset --hard——本分支切出后尽快 commit+push,
# 不留未推工作;要彻底隔离可改在 worktree 做(git worktree add -b changelog/<ver> ../changelog-<ver> origin/master)
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

