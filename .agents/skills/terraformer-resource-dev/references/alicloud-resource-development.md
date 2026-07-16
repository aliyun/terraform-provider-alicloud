# Alicloud Terraformer resource development

## Contents

1. Runtime architecture
2. Source-of-truth checklist
3. InitResources discovery patterns
4. Multipart Import IDs
5. Pagination and errors
6. File selection
7. Tests and validation
8. Common mistakes

## 1. Runtime architecture

`Generator.InitResources()` loads the Alicloud client, calls one or more read-only APIs, converts each discovered object into `terraformutils.Resource`, and appends it to `g.Resources`. `ProviderWrapper.Refresh` normally seeds prior state with that ID and calls the installed Provider's `ReadResource`; the implementation also contains an `ImportResourceState` fallback path. `ConvertTFstate` converts the returned Provider state to state and HCL.

Keep `InitResources` limited to discovery and Provider-compatible IDs. Do not reproduce Create, Update, Delete, schema flattening, or drift logic from the Provider.

## 2. Source-of-truth checklist

Read sources in this order:

1. Provider Resource: find `d.SetId(...)`, every `ParseResourceId(...)`, the Importer, and Read lookup parameters.
2. Import docs/tests: confirm segment order, delimiter, and import round trip.
3. Provider Data Source: reuse only the List API choice, filters, response path, and pagination semantics.
4. Provider service/client: confirm product endpoint, API version, RPC/ROA style, retryable errors, and response normalization.
5. Terraformer same-pattern resources: reuse repository conventions, not identity assumptions.
6. OpenAPI: verify request/response fields when Provider code is indirect or generated.
7. Live read-only call: validate only when credentials and an existing resource are available.

When sources conflict, Provider Import/Read behavior wins for the ID contract. Record the conflict rather than guessing.

## 3. InitResources discovery patterns

### A. Direct full List with a single-field Import ID

Use when one List API enumerates all resources without a parent identifier and each item exposes the Provider's single ID field. Paginate until the API's explicit completion signal, or until a short page when no stronger signal exists.

### B. One List returns every multipart-ID segment

Use when one response item contains every segment required by the Provider Import ID. Preserve the Provider-defined order and delimiter. Do not add a parent List merely because the ID has multiple segments.

### C. Parent-child traversal

Use only when the child List API requires parent scope and Terraformer must enumerate the whole account/region scope:

1. List all parents with complete pagination.
2. For each parent, create a fresh child request.
3. Reset pagination for every parent.
4. List every child page.
5. Join parent and child segments once, at the leaf, using the Provider contract.
6. Return errors with parent ID and page/token context; never silently skip one parent.

A Data Source may require the parent ID because the Terraform caller supplies a query scope. Terraformer cannot impose that Data Source input on a full export; it discovers parents only for this pattern.

The following is pseudocode for the loop shape, not a copy-ready SDK call:

```go
for _, parentID := range parentIDs {
    nextToken := ""

    for {
        children, returnedNextToken, err := listChildren(parentID, nextToken, pageSize)
        if err != nil {
            return nil, fmt.Errorf("list children for parent %s: %w", parentID, err)
        }
        for _, child := range children {
            importID, err := buildProviderImportID(parentID, child.ID)
            if err != nil {
                return nil, err
            }
            ids = append(ids, importID)
        }
        if returnedNextToken == "" {
            break
        }
        nextToken = returnedNextToken
    }
}
```

This example is token-only. For token pagination, stop when the returned next token is empty regardless of page length. For page-number pagination, increment the page number and stop using the API's explicit total/page metadata or a short page when no stronger signal exists. Do not combine token and page-number contracts unless the API actually defines both.

### D. Complete enumeration is unavailable

Use when the service offers only exact lookup, the parent cannot be enumerated, or permissions make account-wide discovery impossible. Reuse an existing Terraformer scope/filter mechanism when it can express the missing input. Otherwise stop and report the limitation; do not claim complete support or guess IDs.

## 4. Multipart Import IDs

The only valid evidence for segment count, order, and delimiter is the Provider Resource's `d.SetId(...)`, its `ParseResourceId(...)` calls, Import docs, and Import tests.

A multipart Import ID does not by itself require parent traversal. All segments may already be present in one List response (pattern B), or earlier segments may require parent discovery (pattern C).

Implementation rules:

- Carry parent, child, attachment, or account segments as separate variables while traversing.
- Validate every required segment before joining.
- Join exactly once when creating the leaf `terraformutils.Resource` ID.
- Do not trim, encode, reorder, or change delimiters without Provider evidence.
- Test the normal ID, missing segment, order, delimiter, and special-character boundary.

## 5. Pagination and errors

- Prefer `NextToken`, `TotalCount`, `IsTruncated`, or an equivalent explicit signal.
- When using returned item count, compare it with the exact page-size value sent in the request.
- Reset pagination for every parent; initialize page number/token inside the parent loop.
- Cover empty first page, short last page, exactly full last page, and multiple pages.
- Include action, resource type, parent ID, page number, or token in wrapped errors.
- Treat permission, endpoint, decode, and single-parent failures as errors, not empty results.
- Follow the repository's retry helpers and product client conventions; do not invent a second retry framework.

## 6. File selection

| File | Change when |
|---|---|
| `providers/alicloud/resource_alicloud_<name>.go` | Always for a new resource; normally for a repair |
| `providers/alicloud/alicloud_provider.go` | Registration or global-resource classification is missing |
| Product client/service files | No existing client can issue the API call |
| Endpoint configuration | The current endpoint resolution is proven insufficient |
| Resource `_test.go` | Lock ID, pagination, empty-result, and error behavior |
| Unified relationship consumer | The shared artifact explicitly declares this resource |

Do not produce or infer connections from Provider schema, Data Source arguments, or API field names. The unified producer owns relationship semantics.

If the unified artifact has no matching declaration, leave the relationship consumer unchanged and record the gap. That absence does not block core discovery and Import ID support unless relationship delivery is itself an explicit acceptance requirement.

Do not modify `cmd`, module entrypoints, README, Provider source, or unrelated shared code unless repository evidence proves the resource cannot work without that change.

## 7. Tests and validation

Use TDD for repairs: demonstrate the current failure, add the smallest regression, then implement the fix.

Static gates:

```bash
RESOURCE_FILE=providers/alicloud/resource_alicloud_example.go
gofmt -l "$RESOURCE_FILE"
go test ./providers/alicloud
go build -o /tmp/terraformer .
```

Confirm registration through the CLI's supported-resource listing or equivalent code path. Run or record `go test ./...`; the current repository has existing unrelated failures, so compare the broad result with the captured baseline while requiring all target-package checks to pass.

When live read-only validation is possible:

1. Export only the target product/resource.
2. Compare discovered count and IDs with the API response.
3. Inspect generated state and HCL.
4. Run `terraform init` and `terraform validate` in the generated directory.
5. Run `terraform plan -refresh-only` and investigate any read/import drift.

When credentials or an existing resource are unavailable, report static validation only and list the unverified live steps.

## 8. Common mistakes

| Mistake | Correct action |
|---|---|
| Treating every multipart ID as parent-child discovery | Select pattern B when one response already contains all segments |
| Copying a Data Source's required parent argument | Enumerate parents only when the child List API requires it |
| Initializing page number outside the parent loop | Reset pagination for every parent |
| Requesting one page size and testing termination with another | Use the same page-size variable for request and termination |
| Guessing an Import ID from API primary keys | Read Provider `d.SetId(...)`, `ParseResourceId(...)`, and Import evidence |
| Editing connection maps by inspection | Read only an explicit declaration from the unified relationship artifact |
| Treating `go test ./...` baseline failures as success or as a new regression | Report the baseline delta and require target tests to pass |
| Creating cloud resources for convenience | Use existing resources or obtain explicit authorization |
