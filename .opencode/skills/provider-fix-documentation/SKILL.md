---
name: provider-fix-documentation
description: Fix documentation issues in Terraform Provider resources. Covers attribute value corrections, description updates, example fixes, and documentation consistency checks.
metadata:
  version: "1.0.0"
  domain: terraform-provider
  triggers: fix documentation, update docs, correct attribute values, documentation error, wrong description
---

# Fix Provider Documentation

Fix documentation issues in Terraform Provider resources, including incorrect attribute values, outdated descriptions, missing information, or inconsistencies between code and documentation.

## Requirement Sources

Requirements may come from:
- User's direct description
- An Aone workitem link
- A GitLab Code Review link
- User-reported documentation errors

If the user provides a link, use the `link-info-extractor` skill to extract requirement details first:
> `task(category="quick", load_skills=["link-info-extractor"], ...)`

## Step 1: Identify the Issue

Common documentation issues include:

1. **Incorrect attribute values** — Valid values in docs don't match schema validation
2. **Outdated descriptions** — API behavior changed but docs weren't updated
3. **Missing attributes** — New attributes added to code but not documented
4. **Wrong examples** — Example code doesn't match current API or schema
5. **Inconsistent terminology** — Different terms used for same concept across docs
6. **Deprecated attributes** — Old attributes still documented without deprecation notice

## Step 2: Locate Target Files

Find the documentation file:
- `website/docs/r/<product>_<resource>.html.markdown` — Resource documentation
- `website/docs/d/<product>_<resource>.html.markdown` — Data source documentation

Find the corresponding code file for verification:
- `alicloud/resource_alicloud_<product>_<resource>.go` — Schema definition
- `alicloud/data_source_alicloud_<product>_<resource>.go` — Data source schema

## Step 3: Verify Against Code

### For Attribute Value Issues

1. **Find the schema definition** in the Go file:
   ```go
   "attribute_name": {
       Type:         schema.TypeString,
       Optional:     true,
       ValidateFunc: StringInSlice([]string{"Value1", "Value2"}, false),
   }
   ```

2. **Check the validation function** — The values in `StringInSlice` are the actual valid values

3. **Check conversion functions** — Look for mapping between Terraform values and API values:
   ```go
   // Request conversion (Terraform → API)
   if value == "Value1" {
       request.ApiValue = "api_value_1"
   }
   
   // Response conversion (API → Terraform)
   if response.ApiValue == "api_value_1" {
       terraformValue = "Value1"
   }
   ```

4. **Compare with documentation** — The docs should list the Terraform-side values (what users configure), NOT the API-side values

### For Description Issues

1. Check the schema `Description` field (if present)
2. Check API documentation via Alibaba Cloud OpenAPI Explorer
3. Verify against actual behavior in test files

## Step 4: Fix Documentation

### Rules

1. **Use Terraform-side values** — Document the values users configure in Terraform, not internal API values
2. **Match schema exactly** — Valid values must match `ValidateFunc` in schema
3. **Clear mapping notes** — If there's a mapping between Terraform and API values, document it clearly
4. **Consistent formatting** — Follow existing documentation style:
   ```markdown
   * `attribute_name` - (Optional, Computed) Description. Supported values:
     - `Value1`: Description of value 1
     - `Value2`: Description of value 2
   ```
5. **Mark deprecated attributes** — Add deprecation notice with version and replacement:
   ```markdown
   * `old_attribute` - (Optional, Deprecated since v1.261.0) Description. Use `new_attribute` instead.
   ```
6. **Update available-since version** — For new attributes, add version info:
   ```markdown
   * `new_attribute` - (Optional, Available since v1.267.0) Description.
   ```

### Common Fixes

#### Fix 1: Incorrect Valid Values

**Before:**
```markdown
* `payment_type` - (Optional) Billing method. Supported values:
  - `prepaid`: Subscription
  - `postpaid`: Pay-as-you-go
```

**After:**
```markdown
* `payment_type` - (Optional) Billing method. Supported values:
  - `PayAsYouGo`: Pay-as-you-go
  - `Subscription`: Subscription
```

#### Fix 2: Missing Deprecation Notice

**Before:**
```markdown
* `instance_charge_type` - (Optional) Instance charge type.
```

**After:**
```markdown
* `instance_charge_type` - (Optional, Deprecated since v1.261.0) Instance charge type. Use `payment_type` instead.
```

#### Fix 3: Incomplete Value Descriptions

**Before:**
```markdown
* `status` - (Optional) Instance status. Valid values: Active, Inactive.
```

**After:**
```markdown
* `status` - (Optional) Instance status. Supported values:
  - `Active`: Instance is active and running
  - `Inactive`: Instance is stopped
```

## Step 5: Verify

1. **Check consistency** — Ensure all similar attributes follow the same documentation pattern
2. **Verify values** — Cross-reference with schema `ValidateFunc` in code
3. **Check examples** — If examples use the attribute, ensure they use correct values
4. **Search for related mentions** — Grep for the attribute name to find all mentions:
   ```bash
   grep -r "attribute_name" website/docs/
   ```

## Step 6: Commit Changes

```bash
git status && git diff
make commit
git push origin fix/<brief_description>
```

Commit message format:
```
docs(<resource_name>): fix <attribute_name> documentation

- Correct valid values from <old_values> to <new_values>
- Update description to match schema definition

Closes: <requirement_url>
```

## Step 7: Run CI Checks (MANDATORY)

**After committing, you MUST run CI checks and ensure all pass:**

```bash
# Run quick CI check (required for all documentation changes)
make ci-check-quick
```

**Checkpoint:**
- ✅ `make ci-check-quick` exits with code 0
- ✅ No errors or warnings related to your changes
- ✅ Markdown linting passes
- ✅ Documentation format validation passes

**If CI checks fail:**
1. Read the error message carefully
2. Fix the reported issues
3. Amend the commit: `git commit --amend`
4. Re-run `make ci-check-quick`
5. Repeat until all checks pass

**DO NOT** skip this step or push code that fails CI checks.

## Important Notes

1. **Terraform values vs API values** — Users configure Terraform values, not API values. The provider handles mapping internally.
2. **Check conversion functions** — If values are mapped, ensure docs show Terraform-side values
3. **Schema is source of truth** — When in doubt, trust the schema `ValidateFunc` over existing docs
4. **Test files as reference** — Acceptance tests show actual valid values in use
5. **Don't change code** — This skill is for documentation fixes only. If code is wrong, use `provider-add-attribute` skill

## Acceptance Criteria

1. Documentation matches schema definition exactly
2. Valid values are correct and complete
3. Descriptions are clear and accurate
4. Consistent formatting with rest of documentation
5. Deprecation notices added where needed
6. All mentions of the attribute are updated consistently
7. **`make ci-check-quick` passes with exit code 0** (MANDATORY)

## Example Workflow

### Issue: User reports documentation error for payment_type

1. **Read the workitem** to understand the reported issue
2. **Find the schema** in `resource_alicloud_elasticsearch_instance.go`:
   ```go
   "payment_type": {
       Type:         schema.TypeString,
       Optional:     true,
       ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
   }
   ```
3. **Find the docs** in `website/docs/r/elasticsearch_instance.html.markdown`
4. **Compare** — Docs say `prepaid`/`postpaid`, schema says `PayAsYouGo`/`Subscription`
5. **Fix docs** — Update to show correct Terraform values
6. **Verify** — Check no other mentions need updating
7. **Git Commit** — Single commit with clear message
8. **CI Check** — Run `make ci-check-quick` and ensure all checks pass ✅
