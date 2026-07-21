---
name: provider-add-attribute
description: Add new attributes to an existing Terraform Provider resource. Covers the full workflow from code changes (Schema/Create/Read/Update), testing, documentation, to submitting a PR.
metadata:
  version: "2.0.0"
  domain: terraform-provider
  triggers: add attribute, add field, new attribute, new field, resource attribute
---

# Add Attribute to Provider Resource

Add one or more new attributes to an existing Terraform Provider resource.

## Requirement Sources

Requirements may come from:
- User's direct description
- An Aone workitem link
- A GitLab Code Review link

If the user provides a link, use the `link-info-extractor` skill to extract requirement details first:
> `task(category="quick", load_skills=["link-info-extractor"], ...)`

## Step 1: Prepare Development Environment

```bash
cd <provider_repo_path>
git checkout master
git pull --rebase alicloud master
git checkout -b feat/<briefDescription>
```

Checkpoint: `git branch --show-current` should show the new branch name.

## Step 2: Locate Target Files

Find these four files:
- `alicloud/resource_alicloud_<product>_<resource>.go` — Schema + CRUD
- `alicloud/resource_alicloud_<product>_<resource>_test.go` — Tests
- `alicloud/service_alicloud_<product>[_v2].go` — Describe method (data source for Read)
- `website/docs/r/<product>_<resource>.html.markdown` — Documentation

Find an existing attribute with the **same type and behavior** as a reference template (e.g., same `TypeBool` + `Optional` + `Computed` + `ForceNew` combination).

Note: Attribute tags must be set based on actual requirements:
- `Computed` is NOT set by default. If tests show the new attribute causes a diff when unset (typically in existing test cases), add `Computed: true`.
- `ForceNew` should be set to `true` when the attribute does not support updates.

## Step 3: Modify Code

### Rules

1. **Remove the auto-generated comment on line 1** once you modify a resource file: `// Package alicloud. This file is generated automatically...`
2. **Naming convention**: Terraform uses `snake_case` (e.g., `delete_on_release`), API uses `PascalCase` (e.g., `DeleteOnRelease`)
3. **Boolean attributes MUST use `GetOkExists`**: `GetOk` cannot distinguish `false` from unset (Go zero-value problem), causing `false` values to be silently dropped
4. **String / Integer attributes use `GetOk`**
5. **Do NOT use `v.(*schema.Set)` or `v.([]interface{})` for array type assertion** — use `convertToInterfaceArray` instead
6. **New attributes MUST have tests**. If the attribute is updatable, the update action MUST be tested
7. **Document the available-since version** for new attributes, e.g., `* \`tags\` - (Optional, Map, Available since v1.55.3) ...`. Use the upcoming version from CHANGELOG.md

## Step 4: Verify

```bash
# 1. Quick CI check
make ci-check-quick

# 2. Full CI check (includes example tests and integration tests)
make ci-check

# 3. If step 2 has test failures, debug with specific resource and test case:
make test-resource-debug RESOURCE=alicloud_vpc TESTCASE=TestAccAliCloudVPC_enableIpv6 LOGLEVEL=TRACE LOGFILE=vpc-test.log
```

## Step 5: Commit Code

Squash all changes into **a single commit**:

```bash
git status && git diff
make commit
git push origin feat/<brief_description>
```

Commit message format:
```
feat(<resource_name>): add <attribute_name> attribute

Closes: <requirement_url>
```

## Important Notes

1. Skip unit tests — only run resource acceptance tests
2. All test failures (API errors, inventory issues, quota problems) must be resolved, never skipped
3. All changes should be squashed into a single git commit
4. Before adjusting API call parameters, use `aliyun <Product> <Action> help` to check API documentation

## Acceptance Criteria

1. Requirement is correctly implemented
2. New test cases pass (Create + Import + Update all PASS)
3. All existing test cases for the resource still pass
4. ci-check-quick and ci-check pass
