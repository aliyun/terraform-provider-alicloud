---
name: provider-resource-acceptance
description: Automated acceptance workflow for Terraform Provider resource releases. Supports two modes - Aone-driven (from workitem link) and local-branch (already on feature branch). Covers code review, acceptance testing, fixing failures, and CI checks.
metadata:
  version: "3.1.0"
  domain: terraform-provider
  triggers: resource acceptance, 资源验收, 发布验收, 自动审核, resource release review, acceptance test workflow, 验收当前分支
---

# Provider Resource Acceptance

End-to-end acceptance workflow: review → test → fix → commit → CI.

## Modes

| Mode | Trigger | Aone interaction |
|---|---|---|
| **A: Aone-driven** | User provides workitem ID or link | ✅ Fetch info, post results as comments |
| **B: Local branch** | User says "验收当前分支" or is on a `feature/*` branch | ❌ None |

## Step 1: Obtain Branch

**Mode A:** Use `link-info-extractor` skill to fetch workitem details. Extract branch name from workitem **comments** (NOT description) — pattern: `feature/<Product>-<Resource>-<action>-<timestamp>`.

```bash
git checkout <branch_name>
git pull --rebase alicloud master
```

**Mode B:** Confirm on a `feature/*` branch, proceed.

## Step 2: Identify Resource

From `git diff --name-only master..HEAD`, extract resource name:
- `alicloud/resource_alicloud_<X>.go` → `RESOURCE=alicloud_<X>`
- Multiple resource files → separate test runs per resource

## Step 3: Code Review

Execute `provider-resource-review` skill (`load_skills=["provider-resource-review"]`).

- `provider-resource-review` returns ❌ → fix all issues, record fixes. Mode A: post fixes to Aone via `coop_add_comment`.
- `provider-resource-review` returns ✅ (possibly with ⚠️) → fix warnings if straightforward, proceed.

## Step 4: Run Single Test

Pick a **new** `TestAcc*_basic` test case from this branch. If none, pick simplest existing `_basic`.

```bash
make test-resource-debug RESOURCE=alicloud_<X> TESTCASE=<TestName> LOGLEVEL=TRACE LOGFILE=<X>-test.log
```

Fix failures, re-run. Mode A: post errors/fixes to Aone.

## Step 5: Run All Tests

```bash
make test-resource-debug RESOURCE=alicloud_<X>
```

For each failure:
- **Our change caused it** → fix, re-run the fixed case, then re-run all
- **Pre-existing** (unmodified test, mock credentials, cross-account) → document, proceed

## Step 6: Commit

```bash
make commit    # NOT git commit — this auto-generates the correct message format
```

Squash to exactly 1 commit on the branch:
```bash
git log --oneline master..HEAD    # if >1 commit:
git reset --soft HEAD~N && git commit -m "<original message>"
```

## Step 7: CI Check

```bash
make ci-check
```

Fix failures, re-squash, re-run. If integration tests timeout but all quality checks passed AND Step 5 tests passed → acceptable.

**Mode A success:** post "所有测试用例通过，CI检查通过，验收成功"
**Mode B success:** report to user.

## Rules

1. Only run `TestAcc*` tests, skip unit tests
2. All test failures related to our changes must be resolved
3. **Mode A:** post every significant action to Aone as a comment
4. **3-strike rule:** if 3 consecutive fix attempts fail on the same issue → `git stash`, report full context to user, STOP

## Troubleshooting

| Problem | Error signature | Fix |
|---|---|---|
| Quota exceeded | `QuotaExceeded`, `limit` | `make sweep REGION=cn-hangzhou RESOURCE=alicloud_<resource>` (match the exhausted resource, not necessarily the one under test) |
| Region unsupported | `not available in this region` | Add `testAccPreCheckWithRegions` in test PreCheck |
| API throttling | `Throttling`, `request throttling` | Wait 30s, re-run |
| Test timeout | `context deadline exceeded` | Increase `StateChangeConf.Timeout` or add `TIMEOUT=60m` to make command |
| Dependency creation failed | VPC/VSwitch/SG errors | Check quota (sweep), region (PreCheck), or dependency function correctness |
| Pre-existing failure | Unmodified test fails | Check for mock credentials (`cross_account_user_id = 1`), document, proceed |

## Acceptance Criteria

1. Code review ✅ (no ❌ items)
2. All acceptance tests pass (pre-existing failures documented)
3. `make ci-check` passes
4. Single squashed commit

**Exception:** dependent resource type doesn't exist → "unable to continue" (not pass/fail).
