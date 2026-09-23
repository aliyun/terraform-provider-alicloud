---
name: provider-resource-review
description: Code review SOP for Terraform Provider resource changes. Covers delete implementation, code-doc consistency, conditional logic documentation, doc quality, test quality, and common bug patterns. Used standalone or referenced by provider-resource-acceptance.
metadata:
  version: "3.2.0"
  domain: terraform-provider
  triggers: code review, resource review, review SOP, review resource changes
outputs:
  result: "PASS | FAIL"
  report: "Markdown review report with ✅/❌/⚠️ per check item"
---

# Provider Resource Review

Review code, documentation, and tests for Terraform Provider resources. **The goal is to find issues and generate a review report, NOT to modify code.**

## Input Sources

- GitHub PR link → use `gh` to retrieve change information
- Resource name → locate files in the current Provider codebase
- Current branch → `git diff --name-only master..HEAD`

## Check Items

### 1. Delete Implementation

Verify the Delete function makes a **real API call**. If deletion only does `d.SetId("")` without calling a delete API → ⚠️ warning, request that a delete API be added later.

### 2. Code-Documentation Consistency

Check each field for consistency between the code Schema definition and documentation:
- Field type (`TypeString` / `TypeInt` / `TypeBool` / `TypeList` etc.)
- Constraints (`Required` / `Optional` / `Computed` / `ForceNew`)
- Valid values (`ValidateFunc` values vs values listed in documentation)

**Array field documentation type rule:**

For `TypeList`/`TypeSet` fields, the documentation label depends on `Elem`, not the code type:
- `Elem: &schema.Resource{...}` (Object) → doc MUST say **Set**
- `Elem: &schema.Schema{...}` (primitive) → doc MUST say **List**

Mismatch → ❌ error.

### 3. Conditional Logic Documentation

If the code contains logic conditional on attribute values, e.g. `if attr == "xxx" { yyy }`, verify the documentation describes this behavior. Undocumented → ❌ error.

### 4. Documentation Completeness

- Parameter descriptions must not be vague (e.g. only saying "The xxx of the resource" with no real information)
- Links in the form `[xxx](~~198289~~)` (bare numeric links) must not appear → critical error
- If an attribute description lost key information compared to the previous version (e.g. "when this value must be provided"), warn and request review

### 5. Test Case Quality

- Check for **hardcoded IDs** in tests (hardcoded resource IDs, Account IDs, etc.) → non-compliant, review fails
- Resource names should use `var.name`, not hardcoded strings

### 6. Code Bug Patterns

**`GetOk` on boolean/integer fields**

`GetOk` returns `ok=false` for Go zero values (`false`, `0`), silently dropping user-set values. TypeBool / TypeInt fields must use `GetOkExists`.

```go
// ❌ Drops false/0
if v, ok := d.GetOk("disabled"); ok { ... }

// ✅
if v, ok := d.GetOkExists("disabled"); ok { ... }
```

**Direct array type assertions**

Direct type assertions panic on type mismatch. Use the provider helper instead.

```go
// ❌ Panics
items := v.(*schema.Set)
items := v.([]interface{})

// ✅
items := convertToInterfaceArray(v)
```

## Important Notes

- API parameter names differing from resource attribute names (e.g. `request["scheduleTime"] = d.Get("recurrence")`) **is not necessarily a bug**. Flag a warning, but do not judge it as an error.
- This skill's goal is to **find issues**, not to modify code or documentation.

## Acceptance Criteria

After completing all checks, perform a **second-pass review**: re-verify every ❌ and ⚠️ in the report for accuracy. Once confirmed, output the final conclusion.

## Report Format

```
## Review Result: ✅ All Checks Passed / ❌ Issues Found

### Check Item Details

✅ Check 1: Delete Implementation
Details: Delete function calls DeleteInstance API
Evidence: alicloud/resource_alicloud_xxx.go:230

❌ Check 2: Conditional Logic Documentation
Details: Code branches on instance_type value, but documentation does not describe this
Evidence: alicloud/resource_alicloud_xxx.go:156
Suggestion: Document the special behavior when instance_type is "ecs.n1.small"

⚠️ Check 3: Parameter Description Clarity
Details: vpc_id description is too vague
Evidence: website/docs/r/xxx.html.markdown:45
Suggestion: Add specific usage details

### Issue Summary

**Critical Issues (must fix):**
1. Description — Location: file:line — Suggestion: fix approach

**Improvement Suggestions (optional):**
1. Description — Location: file:line
```
