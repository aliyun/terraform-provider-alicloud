---
name: provider-fix-regression
description: Automatically fetch failed regression test cases for a specific Alibaba Cloud product from ACube, analyze the root cause, fix the Provider resource code or test case, and submit a PR. Supports filtering by product / resource / error type / date range.
metadata:
  version: "1.0.0"
  domain: terraform-provider
  triggers: fix regression test, 修复回归测试, auto fix failed testcase, concourse build analysis, 修失败用例, regression repair, 自动修复回归
---

# Provider Fix Regression Tests

Fetch failed regression cases from ACube, analyze, fix the Provider code or test, and submit a PR.

## Data Source

ACube exposes regression analysis records (`concourse_build_analysis` table) via HTTP API. Each record contains: `caseName`, `errorType`, `conclusion`, and the full `logContext` of the failed test.

## API

```
GET https://pre-acube.aliyun-inc.com/api/v1/resource_type/info/query_build_analysis
```

Production host: `https://acube.aliyun-inc.com` (omit if unknown, ask user to confirm).

| Param | Required | Notes |
|---|---|---|
| `namespace` | Y | Product name, e.g. `ecs`, `vpc`, `rds` |
| `resourceCode` | N | Specific resource, omit to get all resources under the product |
| `errorType` | N | Filter by error type, default: all; use `OTHER` / `FAIL`-like types to get failures only |
| `startDate` | N | `yyyy-MM-dd`, default: last 7 days |
| `endDate` | N | `yyyy-MM-dd` |
| `pageNum` | N | default `0` |
| `pageSize` | N | default `20`, max `100` |

**Error type values** (from `ConcourseErrorTypeEnum`):

| value | meaning |
|---|---|
| `PASS` | passed (skip) |
| `SKIP` | skipped (skip) |
| `QUOTA_EXCEEDED` | quota problem — usually NOT a code bug, run `make sweep` |
| `ACCOUNT_AUTHORIZE` | auth problem — infra issue, not code |
| `RESOURCE_NOT_FOUND` | dependency resource missing |
| `REQUIRED_ARGUMENT_ABSENT` | required arg missing in test config |
| `REQUEST_PARAMETER_INVALID` | bad request params |
| `ATTRIBUTE_NOT_FOUND` | attribute missing from Read response |
| `ATTRIBUTE_NOT_EXPECTED` | attribute value mismatch in assertion |
| `INVALID_OPERATION` | resource state wrong for operation |
| `RESOURCE_NOT_DESTROYED` | cleanup failure |
| `OTHER` | uncategorized |

## Step 1: Obtain Scope

Determine from the user request:

- **Product** (`namespace`) — required. If missing, ask.
- **Resource** (`resourceCode`) — optional, narrows scope.
- **Time range** — default last 7 days.

If the user says "fix all failures for ecs", use `namespace=ecs` without `resourceCode`.

## Step 2: Fetch Failed Cases

```bash
# Fetch failures only (exclude PASS/SKIP)
curl -sG "https://pre-acube.aliyun-inc.com/api/v1/resource_type/info/query_build_analysis" \
  --data-urlencode "namespace=<product>" \
  --data-urlencode "resourceCode=<resource>" \
  --data-urlencode "pageSize=100" | jq .
```

Since `errorType` filters by exact value, to get all failures in one call omit `errorType` and client-side filter out `PASS` / `SKIP`.

**Group results by `caseName`**. Same case may fail multiple times across runs — use the most recent record per case (highest `id` or latest `gmtCreate`).

**Bucket by `errorType`** to prioritize:

1. **Skip-category (no code fix needed)** — `QUOTA_EXCEEDED`, `ACCOUNT_AUTHORIZE`, `PASS`, `SKIP`. Document, do not modify.
2. **Likely Provider bug** — `ATTRIBUTE_NOT_FOUND`, `ATTRIBUTE_NOT_EXPECTED`, `RESOURCE_NOT_DESTROYED`, `INVALID_OPERATION`.
3. **Likely test case bug** — `REQUIRED_ARGUMENT_ABSENT`, `REQUEST_PARAMETER_INVALID`, `RESOURCE_NOT_FOUND` (test's dependency config wrong).
4. **Needs deeper log analysis** — `OTHER`.

## Step 3: Prepare Development Environment

```bash
cd <terraform-provider-alicloud path>
git checkout master
git pull --rebase alicloud master
git checkout -b fix/<product>-<resource>-regression-$(date +%Y%m%d)
```

Use ONE branch for all failures within the same resource. If fixing multiple resources, use separate branches.

## Step 4: Locate Target Files

For each failing `caseName` (e.g. `TestAccAliCloudEcsInstance_basic`):

- Test file: `alicloud/resource_alicloud_<resource>_test.go`
- Resource: `alicloud/resource_alicloud_<resource>.go`
- Service: `alicloud/service_alicloud_<product>[_v2].go`
- Doc (if attribute-related): `website/docs/r/<product>_<resource>.html.markdown`

Extract resource name from caseName: `TestAccAliCloud<Camel>_<suffix>` → resource is `<snake_case of Camel>`.

## Step 5: Root-Cause Analysis

Parse `logContext` of the failing case. Key patterns to extract:

| Pattern in log | Likely cause | Where to fix |
|---|---|---|
| `Attribute '<name>' not found` | Read method missing `d.Set` call, or API response path changed | `resource_*.go` Read func |
| `Attribute '<name>' expected ... got ...` | Type coercion bug, zero-value drop, or API returned different format | Read func + possibly Schema |
| `Missing required argument "<name>" is required` | Test config missing required field | `*_test.go` config template |
| `<API>Failed: ... InvalidParameter ...` | Wrong request payload, e.g. `GetOk` on bool dropping `false` | Create/Update func |
| `the resource <addr> was not destroyed` | Delete API not wired up or dependency order wrong | Delete func |
| `cannot be found` during Read | Read returns object but not handling `NotFound` — should `d.SetId(""); return nil` | Read func |

If pattern unclear, **print the log snippet to the user with a suggested fix plan before modifying code**.

## Step 6: Apply Fix

Use sibling skills for the specific fix type:

| Failure type | Load skill |
|---|---|
| Attribute missing / wrong in docs | `provider-fix-documentation` |
| Need to add/change an attribute | `provider-add-attribute` |
| General code review needed before fix | `provider-resource-review` |

Common fix patterns (apply directly if simple):

### Fix A: `GetOk` drop on bool/int

```go
// ❌
if v, ok := d.GetOk("enabled"); ok { ... }
// ✅
if v, ok := d.GetOkExists("enabled"); ok { ... }
```

### Fix B: Missing `d.Set` in Read

Check every Schema field — each non-`ForceNew`, non-write-only field needs a matching `d.Set(..., object["..."])` in Read.

### Fix C: Wrong array type assertion

```go
// ❌ Panics
items := v.(*schema.Set).List()
items := v.([]interface{})
// ✅
items := convertToInterfaceArray(v)
```

### Fix D: Test config missing required arg

Update the `config` string in the `*_test.go` file to include the required field. Prefer referencing via `${var.name}` / `${alicloud_<dep>.default.id}`.

### Fix E: NotFound handling in Read

```go
object, err := service.DescribeXxx(d.Id())
if err != nil {
    if NotFoundError(err) {
        log.Printf("[DEBUG] Resource %s not found, removing from state", d.Id())
        d.SetId("")
        return nil
    }
    return WrapError(err)
}
```

## Step 7: Verify Fix Locally

For each fixed `caseName`:

```bash
make test-resource-debug RESOURCE=alicloud_<resource> TESTCASE=<TestName> LOGLEVEL=TRACE LOGFILE=<resource>-fix.log
```

All targeted cases must PASS. Then run the full resource test set:

```bash
make test-resource-debug RESOURCE=alicloud_<resource>
```

Other failures unrelated to our fix → document, do not modify (acceptable pre-existing state).

### 3-strike rule

If 3 consecutive fix attempts on the **same case** fail, `git stash`, report full context (log + attempted diffs) to the user, STOP.

## Step 8: Commit & CI

```bash
git status && git diff
make commit
make ci-check
```

Squash to single commit:

```bash
git log --oneline master..HEAD
# if >1 commit:
git reset --soft HEAD~N && git commit -m "<original message>"
```

Commit message format:

```
fix(<resource>): fix regression test cases

- <caseName_1>: <brief root cause> → <fix summary>
- <caseName_2>: <brief root cause> → <fix summary>

Refs: concourse_build_analysis#<id_1>,<id_2>
```

## Step 9: Report

Present to the user:

- **Fixed cases** (case name + root cause + diff file list)
- **Skipped cases** (case name + why skipped: quota / auth / pre-existing)
- **Unresolved cases** (case name + attempted fixes + why blocked)
- **Branch name + CI status**

## Important Notes

1. **Handle infra-category errors** (`QUOTA_EXCEEDED`, `ACCOUNT_AUTHORIZE`) — these are NOT code bugs. For `QUOTA_EXCEEDED` (e.g. `QuotaExceeded.Vpc`), first clean up dangling resources in the affected region before re-running tests:
   ```bash
   # Clean up VPC quota in the affected region (e.g. cn-hangzhou)
   make sweep REGION=cn-hangzhou RESOURCE=alicloud_vpc
   # Or sweep all resources
   make sweep REGION=cn-hangzhou
   # Then re-run the failing cases
   make test-resource-debug RESOURCE=alicloud_<resource> TESTCASE=<TestName>
   ```
   If sweep doesn't resolve the quota issue, escalate via ops ticket — do NOT modify code to workaround quota limits.
2. **One PR per resource** — do not batch unrelated resources into one branch.
3. **Always run `make ci-check`** before reporting success. `ci-check-quick` is only for docs-only changes.
4. **Do NOT modify** auto-generated resource files' first-line comment without a real content change.
5. The ACube API returns the `env` field as the run date (`yyyy-MM-dd`), not an environment name. Use it for time-based filtering.

## Acceptance Criteria

1. All in-scope failures addressed (fixed OR explicitly documented as out-of-scope).
2. `make test-resource-debug RESOURCE=alicloud_<X>` — all previously-fixed cases PASS.
3. `make ci-check` passes.
4. Single squashed commit per resource.
5. Report delivered to user listing each case's outcome.
