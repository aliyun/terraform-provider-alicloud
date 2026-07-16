---
name: terraformer-resource-dev
description: Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer, including unsupported resources, incomplete discovery, incorrect or multipart Import IDs, parent-scoped listing, pagination defects, endpoint failures, or invalid generated Terraform state or HCL.
---

# Terraformer Resource Development

## Core model

Treat a Terraformer resource as a discovery adapter, not a second Terraform Provider resource implementation:

```text
InitResources
  -> enumerate remote objects and emit Provider-compatible Import IDs
  -> ProviderWrapper.Refresh uses the installed Alicloud Provider Import/Read
  -> ConvertTFstate produces Terraform state and HCL
```

Diagnose failures at the correct layer: discovery, Import ID, Provider Read, or state/HCL conversion. Do not copy Provider CRUD logic into `InitResources`.

## Start every task

1. Resolve the repository with `bash bootstrap/workspace.sh dir terraformer`; resolve Provider evidence with `bash bootstrap/workspace.sh dir terraform_provider`.
2. Preserve dirty files in the Terraformer checkout and create an isolated worktree before modifying tracked files.
3. Use the existing Aone claim/bookend flow. Assign implementation to `terraform-rd` and acceptance verification to `terraform-qa`.
4. Read [references/alicloud-resource-development.md](references/alicloud-resource-development.md) before choosing an API or writing code.
5. Classify the request as a new resource or a repair. For a repair, change only files required by the demonstrated root cause and add a regression test.

## Evidence order

Use this order and record the decisive evidence:

1. Terraform Provider Resource source.
2. Provider Import documentation and Import acceptance tests.
3. Provider Data Source source for List/filter/pagination behavior only.
4. Provider service/client implementation.
5. Terraformer resources with the same discovery pattern.
6. OpenAPI metadata or official API documentation.
7. Read-only live API/export results when credentials and existing resources are available.

The Provider's `d.SetId(...)`, `ParseResourceId(...)`, Import docs, and Import tests define the Import ID. Do not infer it from names or Data Source arguments.

## Select one discovery pattern

Choose exactly one primary `InitResources` pattern:

- **A. Direct full List:** the List API enumerates resources without parent scope.
- **B. One List returns every composite-ID segment:** one response includes all parent and child ID pieces.
- **C. Parent-child traversal:** the child List API requires a parent ID, so enumerate parents and then children; reset pagination for each parent.
- **D. Complete enumeration is unavailable:** use an existing explicit scope/filter input or report the unsupported boundary.

A multipart Import ID does not imply pattern C. A Data Source may require a parent ID because its caller supplies a scope; Terraformer must discover that parent only when the child List API requires it.

## Change only applicable files

- Always add or repair `providers/alicloud/resource_alicloud_<name>.go`.
- Update `providers/alicloud/alicloud_provider.go` only when registration in `SupportedResourceByProduct` or the global-resource list is required.
- Add client/service or endpoint support only when the current product client cannot issue the required API call.
- Add resource-level tests that lock Import ID construction, pagination, empty results, and error propagation.
- Do not modify Terraform Provider code as part of a Terraformer task; split Provider defects into `provider-resource-dev` work.
- Do not produce or infer resource relationships. Read the unified relationship artifact and consume only an explicit matching declaration.

## Validation gates

Run target checks before broad checks:

1. Verify `gofmt` reports no target files.
2. Run the resource regression test and `go test ./providers/alicloud`.
3. Build the binary to `/tmp/terraformer` so the repository stays clean.
4. Confirm the resource is visible through the Terraformer CLI registration path.
5. Run or record `go test ./...`; compare failures with the baseline instead of hiding existing unrelated failures.
6. When an account and an existing resource are available, perform a read-only export, inspect state/HCL, run `terraform validate`, and run `terraform plan -refresh-only`.

If live validation is unavailable, report "static validation only" and list the missing acceptance evidence. Never create cloud resources merely to make a Terraformer discovery check possible unless the user explicitly authorizes it.

## Delivery

Keep the worktree after opening a CR/MR, link it to Aone immediately, and do not merge or release. Report the selected discovery pattern, Import ID evidence, files changed, tests run, existing baseline failures, and any live-validation gap.
