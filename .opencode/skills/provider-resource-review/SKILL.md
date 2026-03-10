---
name: provider-resource-review
description: Code review SOP for Terraform Provider resource changes. Covers Schema validation, CRUD function correctness, ID construction, JSON encoding, test quality, and common bug patterns. Used standalone or referenced by provider-resource-acceptance.
metadata:
  version: "1.0.0"
  domain: terraform-provider
  triggers: code review, resource review, review SOP, 代码审查, 资源审查, review resource changes
---

# Provider Resource Review

Structured code review SOP for Terraform Provider resource changes. This skill defines **what to check** and **how to catch common bugs** in resource implementations.

## Applicable Scenarios

- Reviewing a feature branch that adds/modifies a Terraform Provider resource
- Called as a sub-step by `provider-resource-acceptance` during acceptance workflows
- Standalone review before submitting a PR

## Step 1: Identify Changed Files

Determine which files were changed on the current branch:

```bash
git diff --name-only master..HEAD
```

Expected file categories:
- `alicloud/resource_alicloud_<product>_<resource>.go` — Resource implementation (Schema + CRUD)
- `alicloud/resource_alicloud_<product>_<resource>_test.go` — Test cases
- `alicloud/service_alicloud_<product>[_v2].go` — Service layer (Describe/List methods)
- `website/docs/r/<product>_<resource>.html.markdown` — Documentation

## Step 2: Schema Review

Open the resource file and review the Schema definition.

### Checklist

1. **Auto-generated comment removed** — If the file was modified, the first line `// Package alicloud. This file is generated automatically...` MUST be deleted
2. **Attribute tags correctness**:
   - `Optional` vs `Required` — matches API requirements
   - `Computed` — set only when the API returns a value that may differ from user input. Do NOT set `Computed` by default
   - `ForceNew` — set when the attribute cannot be updated (requires resource recreation)
   - `Sensitive` — set for passwords, secrets, keys
3. **Type correctness** — `TypeString`, `TypeInt`, `TypeBool`, `TypeList`, `TypeSet`, `TypeMap` match actual data types
4. **Validation functions** — `ValidateFunc` present for attributes with constrained values (e.g., `StringInSlice`)
5. **Nested sub-structs** — `TypeList` with `MaxItems: 1` + `Elem: &schema.Resource{Schema: ...}` for object-typed attributes
6. **ConflictsWith / ExactlyOneOf / AtLeastOneOf** — mutual exclusion constraints are correct
7. **Default values** — `Default` is only set when the API has a documented default

### Common Schema Issues

| Issue | Symptom | Fix |
|---|---|---|
| Missing `Computed` | Existing tests show diff on unset attributes | Add `Computed: true` |
| Missing `ForceNew` | Update API returns error for this field | Add `ForceNew: true` |
| Wrong `Type` | Type assertion panics at runtime | Change to correct type |
| `ForceNew` on parent but not children | Nested attributes don't trigger recreation | Add `ForceNew` to child attributes OR remove from parent |

## Step 3: Create Function Review

### Resource ID Construction

The `d.SetId()` call determines how the resource is identified. Review for correctness.

**Correct Patterns:**

| Pattern | When to Use | Example |
|---|---|---|
| Single value from response | Resource has a unique server-generated ID | `d.SetId(fmt.Sprint(response["InstanceId"]))` |
| Composite from response + request | Child resource under a parent | `d.SetId(fmt.Sprintf("%v:%v", request["ParentId"], response["ChildId"]))` |
| Composite from `d.Get()` | Attachment/binding resource with no server-generated ID | `d.SetId(fmt.Sprintf("%v:%v", d.Get("parent_id"), d.Get("child_id")))` |
| Via `jsonpath.Get` from response | ID is nested in response JSON | `id, _ := jsonpath.Get("$.Resource.Id", response)` then `d.SetId(fmt.Sprint(id))` |

**Bug Pattern: `jsonpath.Get` on JSON strings**

When the API wraps a list/object as a JSON *string* (not a native object), `jsonpath.Get` returns `nil` because it cannot traverse a string.

```go
// ❌ BUG: PolicyBindingList is a JSON string like "[{...}]", not a traversable array
sourceType, _ := jsonpath.Get("PolicyBindingList[0].SourceType", response)
// sourceType == nil

// ✅ FIX: Use d.Get() for values already known from input
sourceType := d.Get("source_type")
```

**How to detect:** If the API field is documented as type `String` but contains JSON content, `jsonpath.Get` will NOT work on it.

### Request Parameter Construction

Review how request parameters are built.

**`GetOk` vs `GetOkExists` Rule:**

```go
// ✅ Use GetOkExists for boolean fields (and integer fields where 0 is valid)
if v, ok := d.GetOkExists("disabled"); ok {
    request["Disabled"] = v
}

// ✅ Use GetOk for string and complex fields
if v, ok := d.GetOk("description"); ok {
    request["Description"] = v
}
```

**Why:** `GetOk` returns `ok=false` when the value is Go's zero-value (`false` for bool, `0` for int). This silently drops valid user-specified `false`/`0` values. `GetOkExists` distinguishes "not set" from "set to zero-value".

**`sjson.Set` for JSON String Parameters:**

Some APIs expect parameters as JSON strings. The provider uses `sjson.Set` to build them:

```go
jsonString := convertObjectToJsonString(request)
jsonString, _ = sjson.Set(jsonString, "Items.0.Name", d.Get("name"))
jsonString, _ = sjson.Set(jsonString, "Items.0.Value", d.Get("value"))
_ = json.Unmarshal([]byte(jsonString), &request)
```

**Bug Pattern: `sjson.Set` creates Go arrays, not JSON strings**

When the API expects a parameter to be a JSON-encoded string like `"[\"value\"]"`, using `sjson.Set` creates a native Go array `["value"]` instead.

```go
// ❌ BUG: Creates native array, but API expects JSON string
jsonString, _ = sjson.Set(jsonString, "DataSourceIds.0", parts[2])
// Result: DataSourceIds = ["id123"]  (Go array)
// API expects: DataSourceIds = "[\"id123\"]"  (JSON string)

// ✅ FIX: Use string concatenation for JSON-string parameters
request["DataSourceIds"] = "[\"" + parts[2] + "\"]"
```

**How to detect:** Check the API documentation — if the parameter type is `String` and the example value contains escaped quotes like `"[\"...\"]"`, it's a JSON-string parameter.

### Array Type Assertions

```go
// ❌ WRONG: Direct type assertion on arrays
items := v.(*schema.Set)     // panics if type doesn't match
items := v.([]interface{})   // panics if type doesn't match

// ✅ CORRECT: Use provider helper
items := convertToInterfaceArray(v)
```

## Step 4: Read Function Review

### Checklist

1. **All schema attributes are set** — Every attribute in the Schema must have a corresponding `d.Set()` call in Read
2. **Not-found handling** — If the resource doesn't exist, call `d.SetId("")` and return `nil` (not an error)
3. **Nested object mapping** — For `TypeList` sub-structs, build a `[]map[string]interface{}` and set it
4. **Type consistency** — Values set via `d.Set()` must match the schema `Type`

### Common Read Issues

| Issue | Symptom | Fix |
|---|---|---|
| Missing `d.Set()` for new attribute | Attribute always shows as empty | Add `d.Set("attr", objectRaw["ApiField"])` |
| Wrong API field name | Attribute reads wrong value | Check API response structure |
| Missing not-found check | Destroy fails with error | Add `NotFoundError(err)` check |

## Step 5: Update Function Review

### Checklist

1. **Update vs ForceNew** — Attributes with `ForceNew: true` must NOT appear in the Update function
2. **Partial update support** — If the API supports updating individual fields, use `d.HasChange("field")` guards
3. **Read-after-update** — Update function should call the Read function at the end to refresh state

### Common Update Issues

| Issue | Symptom | Fix |
|---|---|---|
| Updating a ForceNew field | API returns error on update | Remove from Update, add `ForceNew: true` to schema |
| Missing `HasChange` guard | Unchanged fields sent to API | Wrap in `if d.HasChange("field")` |
| No Read after Update | State drift until next plan | Add `return resourceRead(d, meta)` at end |

## Step 6: Delete Function Review

### Checklist

1. **Resource ID parsing** — If composite ID, verify `strings.Split` with correct separator and index
2. **Not-found tolerance** — Delete should NOT error if resource is already gone
3. **Parameter encoding** — Same JSON-string parameter issues as Create (see Step 3)

## Step 7: Service File Review

### Checklist

1. **Describe method exists** — `Describe<Resource>(id)` is implemented
2. **ID parsing** — Composite IDs are split correctly
3. **Pagination** — If the API returns paginated results, pagination is handled
4. **Error mapping** — API "not found" errors are converted to `NotFoundError`

## Step 8: Test File Review

### Naming Conventions

```go
// Standard test name format
rand := acctest.RandIntRange(10000, 99999)
name := fmt.Sprintf("tf-testacc%s<resource_type>%d", defaultRegionToTest, rand)
```

**Test Suffixes:**
- `_twin` — Parallel/duplicate test with slightly different configuration
- `_raw` — Tests raw/alternative API configurations or edge cases

### Dependency Function Conventions

```go
// ✅ CORRECT: Use var.name for resource names
func AliCloudXxxBasicDependence1234(name string) string {
    return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
    vpc_name   = var.name          // ✅ References variable
    cidr_block = "172.16.0.0/12"
}
`, name)
}
```

```go
// ❌ ANTI-PATTERN: Hardcoded resource names
resource "alicloud_vpc" "default" {
    vpc_name = "tf-test-hardcoded-name"   // ❌ Should be var.name
}
```

**Why `var.name` matters:** Hardcoded names cause collisions when tests run in parallel, and make cleanup via sweep impossible.

### Test Coverage Checklist

1. **All new attributes tested** — Each new attribute appears in at least one test step
2. **Update tested** — If the attribute supports update, there's a step that changes its value
3. **Import tested** — `ImportState: true` step exists
4. **Region constraints** — If the resource is region-limited, `testAccPreCheckWithRegions` is used
5. **No hardcoded names** — All resource names use `var.name`

## Step 9: Documentation Review

### Checklist

1. **All new attributes documented** — Each new schema attribute has a corresponding entry in the docs
2. **Correct tags** — `(Required)`, `(Optional)`, `(Computed)`, `(ForceNew)` match the schema
3. **Valid values listed** — If `ValidateFunc` exists, document the allowed values
4. **Available-since version** — New attributes have `Available since vX.Y.Z`
5. **Example updated** — If new attributes are required, the example includes them
6. **Subcategory correct** — The front-matter `subcategory` matches the product

## Review Outcome

### PASS Criteria

All of the following must be true:
1. No bugs found in CRUD functions (or all bugs have been fixed)
2. Schema is correctly defined
3. Test coverage is adequate
4. Documentation is accurate

### FAIL Criteria

Any of the following triggers a failure:
1. Resource ID construction will produce incorrect IDs at runtime
2. `GetOk` used for boolean fields (will silently drop `false` values)
3. JSON encoding mismatch (Go array vs JSON string)
4. Missing `d.Set()` in Read for new attributes
5. No test coverage for new attributes
6. Critical documentation errors (wrong valid values, missing required attributes)

### Common Bug Pattern Summary

| # | Bug Pattern | Detection Method | Typical Fix |
|---|---|---|---|
| 1 | `jsonpath.Get` on JSON string fields | API docs say type=String but value is JSON | Use `d.Get()` or parse the JSON string first |
| 2 | `sjson.Set` for JSON-string parameters | API expects `"[\"...\"]"` format | Use string concatenation instead |
| 3 | `GetOk` on boolean/integer fields | Grep for `GetOk("bool_field")` | Change to `GetOkExists` |
| 4 | Hardcoded test resource names | Grep for `_name = "tf-` in test files | Change to `var.name` |
| 5 | Missing auto-gen comment removal | First line starts with `// Package alicloud. This file is generated` | Delete the line |
| 6 | Direct array type assertion | `v.(*schema.Set)` or `v.([]interface{})` | Use `convertToInterfaceArray(v)` |
| 7 | Ignored `jsonpath.Get` errors | `value, _ := jsonpath.Get(...)` with no nil check | Add nil/error check before using value |
