---
name: provider-resource-acceptance
description: Automated acceptance workflow for Terraform Provider resource releases. Covers fetching Aone workitem, branch checkout, code review, running acceptance tests, fixing failures, and reporting results back to Aone.
metadata:
  version: "1.0.0"
  domain: terraform-provider
  triggers: resource acceptance, 资源验收, 发布验收, 自动审核, resource release review, acceptance test workflow
---

# Provider Resource Acceptance

End-to-end automated acceptance workflow for Terraform Provider resource releases triggered by Aone workitems.

## Applicable Scenarios

This skill applies when:
1. The Aone workitem title follows the pattern: `[Terraform 资源发布自动审核流程] 产品 [xxx] 资源 [xxx]`
2. The user is working in the Alibaba Cloud Terraform Provider repository

## Step 1: Fetch Aone Workitem Details

Use MCP tools to get the workitem details and comments:

1. Call `coop_query_workitem_detail(workitemId)` to get the title, description, and status
2. Call `coop_get_workitem_comments(workitemId)` to get all comments
3. **Extract the branch name** from comments — it follows the pattern: `feature/<Product>-<Resource>-<action>-<timestamp>` (e.g., `feature/SLS-Alert-update-1767947524953`)

> **Note:** The branch name is in the comment section, NOT in the workitem description.

## Step 2: Switch to the Target Branch and Sync with Master

Switch to the target branch:

```bash
git checkout <branch_name>
```

Then rebase the branch onto the latest master:

```bash
git pull --rebase alicloud master
```

### If There Are Conflicts

1. Resolve the merge conflicts in the affected files
2. Stage the resolved files: `git add <file_path>`
3. Continue the rebase: `git rebase --continue`
4. Repeat until the rebase completes

Checkpoint: `git branch --show-current` should show the target branch name, and `git status` should show a clean working tree.

## Step 3: Code Review (Resource Review)

Execute the **Resource Review** SOP process against the code on this branch.

> **Important:** Even if the Aone workitem already has historical review records, you MUST perform a fresh review.

Load the relevant review skill or follow the project's resource review SOP.

### If Review Fails

1. Summarize the failure reasons
2. Post the failure reasons as a comment on the Aone workitem via `coop_add_comment(workitemId, content)`
3. **Stop here** — validation is complete (failed)

### If Review Passes

Proceed to Step 4.

## Step 4: Locate and Verify Test Cases

1. Find the resource acceptance test file: `alicloud/resource_alicloud_<product>_<resource>_test.go`
2. Confirm test cases exist for this resource

### If No Test Cases Found

1. Post a comment on the Aone workitem: "找不到资源测试用例，验收失败"
2. **Stop here** — validation is complete (failed)

## Step 5: Run a Single Test Case First

Pick one specific test case and run it individually to verify the basic setup:

```bash
make test-resource-debug RESOURCE=alicloud_<product>_<resource> TESTCASE=<TestCaseName> LOGLEVEL=TRACE LOGFILE=vpc-test.log
```

Example:
```bash
make test-resource-debug RESOURCE=alicloud_vpc TESTCASE=TestAccAliCloudVpcVpc_basic9656 LOGLEVEL=TRACE LOGFILE=vpc-test.log
```

### If the Test Fails

1. Analyze the failure root cause
2. Fix the issue
3. **Record** the error details and fix description
4. Post the error info and fix process as a comment on the Aone workitem
5. Re-run the test to confirm it passes
6. Repeat until this test case passes

### If the Test Passes

Proceed to Step 6.

## Step 6: Run All Test Cases for the Resource

Run the full test suite for this resource without specifying a single test case:

```bash
make test-resource-debug RESOURCE=alicloud_<product>_<resource>
```

Example:
```bash
make test-resource-debug RESOURCE=alicloud_vpc
```

### If Any Tests Fail

1. For each failing test case, repeat the Step 5 fix process:
   - Analyze the root cause
   - Fix the issue
   - Post error details and fix process to the Aone workitem
2. Re-run all tests after fixing
3. Repeat until **all** test cases pass

### If All Tests Pass

Proceed to Step 7.

## Step 7: Commit Changes

Squash all changes into **a single commit** using the project's commit tool:

```bash
make commit
```

> **Important:** All changes (code fixes, test adjustments, etc.) must be squashed into one commit. Do NOT use `git commit` directly — always use `make commit`.

## Step 8: Run CI Check

Run the full CI check to ensure all checks pass:

```bash
make ci-check
```

### If CI Check Fails

1. Read the error output carefully
2. Fix the reported issues
3. Re-run the relevant tests if code was changed (repeat Step 5/6 as needed)
4. Squash fixes into the same commit: `make commit`
5. Re-run `make ci-check`
6. Repeat until CI check passes

### If CI Check Passes

Post a comment on the Aone workitem: "所有测试用例通过，CI检查通过，验收成功"

## Important Notes

1. **Skip unit tests** — Only run resource acceptance tests (`TestAcc*`), ignore unit tests
2. **All failures must be resolved** — API errors, inventory issues, and quota problems must all be fixed; never skip a failing test
3. **Remove auto-generated comment** — Once you modify any resource code file, delete the first-line comment: `// Package alicloud. This file is generated automatically...`
4. **Always report to Aone** — Every significant action (review failure, test failure, test success) must be posted as a comment to the Aone workitem

## Troubleshooting

### Quota Exceeded

If a test fails due to resource quota limits, use the sweep command to clean up resources in that region:

```bash
make sweep REGION=cn-hangzhou RESOURCE=alicloud_<resource>
```

> **Note:** The `RESOURCE` in the sweep command should match the resource that exceeded the quota, which may NOT be the resource under test.

### Region Not Supported

If a test fails because the region doesn't support the resource, add a region-specific `PreCheck` in the test case:

```go
// Reference: resource_alicloud_vpc_ipam_ipam_test.go
PreCheck: func() {
    testAccPreCheck(t)
    testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
},
```

## Acceptance Criteria

A release validation is **successful** only when:
1. Resource code review passes
2. **All** acceptance test cases pass
3. `make ci-check` passes with exit code 0

### Exception

If the failure is due to a **dependent resource type not existing** (dependency resource unavailable), the validation task ends as "unable to continue" — this is NOT considered a pass or fail, just an early termination.

## Example Workflow

### Validating a resource release for SLS Alert

1. **Fetch Aone workitem** — Get workitem details and extract branch `feature/SLS-Alert-update-1767947524953`
2. **Checkout branch** — `git checkout feature/SLS-Alert-update-1767947524953`
3. **Code review** — Run resource review SOP on the changed files
4. **Review passed** — Proceed to testing
5. **Run single test** — `make test-resource-debug RESOURCE=alicloud_sls_alert TESTCASE=TestAccAliCloudSlsAlert_basic0001`
6. **Single test passed** — Run all tests
7. **Run all tests** — `make test-resource-debug RESOURCE=alicloud_sls_alert`
8. **2 tests failed** — Fix issues, post error details to Aone
9. **Re-run all tests** — All pass
10. **Commit** — `make commit` to squash all changes into one commit
11. **CI check** — `make ci-check` passes
12. **Post success** — Comment on Aone: "所有测试用例通过，CI检查通过，验收成功"
