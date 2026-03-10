---
name: provider-resource-acceptance
description: Automated acceptance workflow for Terraform Provider resource releases. Supports two modes - Aone-driven (from workitem link) and local-branch (already on feature branch). Covers code review, acceptance testing, fixing failures, and CI checks.
metadata:
  version: "2.0.0"
  domain: terraform-provider
  triggers: resource acceptance, 资源验收, 发布验收, 自动审核, resource release review, acceptance test workflow, 验收当前分支
---

# Provider Resource Acceptance

End-to-end automated acceptance workflow for Terraform Provider resource releases.

## Workflow Modes

This skill supports two trigger modes. Detect which mode applies and skip irrelevant steps.

| Mode | Trigger                                                           | Aone Steps | Git Steps |
|---|-------------------------------------------------------------------|---|---|
| **Mode A: Aone-driven** | User provides an Aone workitem ID or link                         | ✅ Fetch workitem, extract branch, report results | ✅ Checkout branch, rebase |
| **Mode B: Local branch** | User is already on a feature branch, or says "验收当前分支" or "验收当前资源" | ❌ Skip all Aone interaction | ❌ Skip checkout (already there) |

**Detection logic:**
- If user provides a workitem ID/link → Mode A
- If user says "验收当前分支" / "验收当前资源"" / "validate current resource" / "validate current branch" is already on a `feature/*` branch → Mode B
- If unclear → ask the user

## Step 1: Obtain Branch Context (Mode-dependent)

### Mode A: Aone-driven

1. Call `coop_query_workitem_detail(workitemId)` to get the title, description, and status
2. Call `coop_get_workitem_comments(workitemId)` to get all comments
3. **Extract the branch name** from comments — it follows the pattern: `feature/<Product>-<Resource>-<action>-<timestamp>` (e.g., `feature/SLS-Alert-update-1767947524953`)

> **Note:** The branch name is in the comment section, NOT in the workitem description.

4. Switch to the target branch and sync:

```bash
git checkout <branch_name>
git pull --rebase alicloud master
```

If there are rebase conflicts, resolve them file by file, then `git add <file>` and `git rebase --continue`.

**Checkpoint:** `git branch --show-current` shows the target branch, `git status` shows a clean tree.

### Mode B: Local branch

1. Confirm we're on a feature branch:

```bash
git branch --show-current
```

2. If already on a `feature/*` branch, proceed directly to Step 2.

## Step 2: Code Review

Execute the **provider-resource-review** skill against the changed files on this branch.

> **Important:** Load the review skill: `load_skills=["provider-resource-review"]`. Even if the Aone workitem has historical review records, always perform a fresh review.

The review skill covers:
- Schema validation (attribute tags, types, constraints)
- CRUD function correctness (ID construction, parameter encoding, error handling)
- Test quality (naming conventions, coverage, anti-patterns)
- Documentation accuracy

### If Review Finds Issues

1. Fix all issues found during review
2. Record what was fixed (for the commit message and Aone reporting)
3. **Mode A only:** Post the fix details as a comment on the Aone workitem via `coop_add_comment(workitemId, content)`

### If Review Passes with No Issues

Proceed to Step 3.

## Step 3: Locate and Verify Test Cases

Find the test file: `alicloud/resource_alicloud_<product>_<resource>_test.go`

Confirm that:
1. Test cases exist for this resource
2. New/modified attributes have test coverage
3. `TestAcc*` functions are present (ignore unit tests)

### If No Test Cases Found

- **Mode A:** Post comment: "找不到资源测试用例，验收失败" → **Stop**
- **Mode B:** Report to user: "No test cases found for this resource" → **Stop**

## Step 4: Run a Single Test Case First

Pick one test case — preferably a new test case added in this branch — and run it individually:

```bash
make test-resource-debug RESOURCE=alicloud_<product>_<resource> TESTCASE=<TestCaseName> LOGLEVEL=TRACE LOGFILE=<resource>-test.log
```

### If the Test Fails

1. Analyze the failure root cause from the log output
2. Fix the issue in the resource code or test code
3. **Record** the error details and fix description
4. **Mode A only:** Post the error info and fix process as a comment on the Aone workitem
5. Re-run the test to confirm it passes
6. Repeat until this test case passes

### If the Test Passes

Proceed to Step 5.

## Step 5: Run All Test Cases

Run the full test suite for this resource:

```bash
make test-resource-debug RESOURCE=alicloud_<product>_<resource>
```

### If Any Tests Fail

For each failing test:

1. **Determine if it's pre-existing or caused by our changes:**
   - If the failing test was NOT modified in our branch AND uses mock credentials or cross-account features → likely pre-existing
   - If the failing test was modified or tests new functionality → must fix

2. **For failures caused by our changes:**
   - Analyze root cause
   - Fix the issue
   - **Mode A only:** Post error details and fix to Aone
   - Re-run all tests after fixing

3. **For pre-existing failures:**
   - Document the test name and failure reason
   - Confirm it's unrelated to our changes
   - Proceed (do NOT fix pre-existing issues unless asked)

### If All Tests Pass

Proceed to Step 6.

## Step 6: Commit Changes

All changes (code fixes, test adjustments, etc.) must be squashed into **a single commit**.

```bash
make commit
```

> **Important:** Do NOT use `git commit` directly — always use `make commit`, which auto-generates the correct commit message format.

If there was already a commit on the branch and `make commit` creates a second one, squash them:

```bash
# Count commits on branch relative to master
git log --oneline master..HEAD

# If more than 1 commit, squash:
git reset --soft HEAD~N  # N = number of commits to squash
git commit -m "<original commit message from the first commit>"
```

**Checkpoint:** `git log --oneline master..HEAD` shows exactly 1 commit.

## Step 7: Run CI Check

```bash
make ci-check
```

This runs:
- `gofmt` and `goimports` checks
- `go vet`
- Basic checks (no `fmt.Println`, doc links)
- Testing coverage rate check (must be 100%)
- Documentation checks (terrafmt, spelling, content, consistency)
- Example tests (apply + destroy against real infrastructure)
- Resource integration tests (full acceptance test suite)

### If CI Check Fails

1. Read the error output carefully
2. Fix the reported issues
3. If code was changed, re-run the relevant tests (repeat Step 4/5 as needed)
4. Re-squash into one commit: squash with the original commit message
5. Re-run `make ci-check`
6. Repeat until CI check passes

### CI Check Timeout

The integration test phase re-runs all acceptance tests, which can take 30+ minutes. If the CI check times out:

1. Check that all quality checks before the integration tests passed (gofmt, goimports, go vet, coverage, docs)
2. If quality checks all passed AND we already ran all tests successfully in Step 5 → the timeout is acceptable
3. Report which tests completed before timeout

### If CI Check Passes

- **Mode A:** Post comment: "所有测试用例通过，CI检查通过，验收成功"
- **Mode B:** Report success to user

## Important Notes

1. **Skip unit tests** — Only run resource acceptance tests (`TestAcc*`), ignore unit tests
2. **All failures must be resolved** — API errors, inventory issues, and quota problems must all be fixed; never skip a failing test that's related to our changes
3. **Remove auto-generated comment** — Once you modify any resource code file, delete the first-line comment: `// Package alicloud. This file is generated automatically...`
4. **Mode A: Always report to Aone** — Every significant action (review failure, test failure, test success) must be posted as a comment to the Aone workitem
5. **Mode B: Report to user directly** — No Aone interaction needed

## Troubleshooting

### Quota Exceeded

If a test fails due to resource quota limits, clean up resources:

```bash
make sweep REGION=cn-hangzhou RESOURCE=alicloud_<resource>
```

> The `RESOURCE` in the sweep command should match the resource that exceeded quota, which may NOT be the resource under test.

### Region Not Supported

If a test fails because the region doesn't support the resource, add a region-specific `PreCheck`:

```go
PreCheck: func() {
    testAccPreCheck(t)
    testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
},
```

### Pre-existing Test Failures

If a test that was NOT modified in our branch fails:

1. Check if it uses mock credentials (`cross_account_user_id = 1`)
2. Check if it depends on external resources not available in the test environment
3. If confirmed pre-existing, document it and proceed
4. Do NOT fix pre-existing issues unless explicitly asked

## Acceptance Criteria

A release validation is **successful** when:
1. Code review passes (via `provider-resource-review` skill)
2. **All** acceptance test cases pass (pre-existing failures documented but not blocking)
3. `make ci-check` quality checks pass
4. All changes squashed into a single commit

### Exception

If the failure is due to a **dependent resource type not existing** (dependency resource unavailable), the validation task ends as "unable to continue" — this is NOT considered a pass or fail, just an early termination.

## Example Workflows

### Mode A: Aone-driven validation for SLS Alert

1. **Fetch Aone workitem** → Get details, extract branch `feature/SLS-Alert-update-1767947524953`
2. **Checkout branch** → `git checkout feature/SLS-Alert-update-1767947524953`
3. **Code review** → Load `provider-resource-review` skill, review changed files
4. **Review passed** → Proceed to testing
5. **Run single test** → `make test-resource-debug RESOURCE=alicloud_sls_alert TESTCASE=TestAccAliCloudSlsAlert_basic0001`
6. **Single test passed** → Run all tests
7. **Run all tests** → `make test-resource-debug RESOURCE=alicloud_sls_alert`
8. **2 tests failed** → Fix issues, post error details to Aone
9. **Re-run all tests** → All pass
10. **Commit** → `make commit`, squash to single commit
11. **CI check** → `make ci-check` passes
12. **Post success** → Comment on Aone: "所有测试用例通过，CI检查通过，验收成功"

### Mode B: Local branch validation for HBR PolicyBinding

1. **Confirm branch** → Already on `feature/HBR-PolicyBinding-update-1770810266269`
2. **Code review** → Load `provider-resource-review`, review 4 changed files
3. **Found 4 issues** → Fixed: auto-gen comment, Create ID construction, Delete encoding, hardcoded test names
4. **Run single test** → `make test-resource-debug RESOURCE=alicloud_hbr_policy_binding TESTCASE=TestAccAliCloudHbrPolicyBinding_basic6221` → PASS
5. **Run all tests** → 14/15 passed, 1 pre-existing failure (basic7232, cross-account mock)
6. **Commit** → `make commit`, squashed to single commit
7. **CI check** → All quality checks passed, example test passed
8. **Report to user** → "验收完成：14/15 测试通过，1个预存在失败，CI质量检查全部通过"
